package named

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/xaque208/znet/pkg/events"
	grpc "google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Named struct {
	UnimplementedNamedServer
	services.Service

	logger log.Logger

	cfg *Config

	sch  *events.Scheduler
	conn *grpc.ClientConn

	// lights *lights.Lights
}

func New(cfg Config, logger log.Logger, conn *grpc.ClientConn) (*Named, error) {
	n := &Named{
		cfg:    &cfg,
		conn:   conn,
		logger: log.With(logger, "timer", "named"),
		sch:    events.NewScheduler(logger),
	}

	n.Service = services.NewBasicService(n.starting, n.running, n.stopping)

	return n, nil
}

func (t *Named) Observe(ctx context.Context, req *NamedTimeStamp) (*Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("unable to handle nil request")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("unable to handle request with empty name")
	}

	// return &Empty{}, t.lights.NamedTimerHandler(ctx, req.Name)

	return nil, nil
}

func (t *Named) Schedule(ctx context.Context, req *NamedTimeStamp) (*Empty, error) {

	err := t.sch.Set(req.Time.AsTime(), req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set named schedule")
	}

	return nil, nil
}

func (n *Named) starting(ctx context.Context) error {
	if len(n.cfg.Events) == 0 && len(n.cfg.RepeatEvents) == 0 {
		return fmt.Errorf("no Events or RepeatEvents config")
	}

	return nil
}

func (n *Named) running(ctx context.Context) error {
	err := n.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to Connect named")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (n *Named) stopping(_ error) error {
	return nil
}

//connect implements evenets.Producer
func (n *Named) Connect(ctx context.Context) error {
	n.logger.Log("msg", "starting eventProducer")

	go func() {
		err := n.scheduler(ctx)
		if err != nil {
			level.Error(n.logger).Log("msg", "scheduler failed",
				"err", err,
			)
		}
	}()

	return nil
}

func (n *Named) scheduler(ctx context.Context) error {
	namedClient := NewNamedClient(n.conn)

	go func() {
		for {
			for _, repeatEvent := range n.cfg.RepeatEvents {
				times := n.sch.TimesForName(repeatEvent.Produce)
				if len(times) == 0 {
					err := n.scheduleRepeatEvents(n.sch, repeatEvent)
					if err != nil {
						level.Error(n.logger).Log("msg", "failed to schedule repeat events",
							"repeatEvent", repeatEvent.Produce,
							"err", err,
						)
					}
				}
			}

			for _, event := range n.cfg.Events {
				if len(n.sch.TimesForName(event.Produce)) == 0 {
					err := n.scheduleEvents(event)
					if err != nil {
						level.Error(n.logger).Log("msg", "failed to schedul events",
							"event", event.Produce,
							"err", err,
						)
					}
				}
			}

			n.sch.Report()

			names := n.sch.WaitForNext()

			if len(names) == 0 {
				continue
			}

			for _, name := range names {
				_, err := namedClient.Observe(ctx, &NamedTimeStamp{Name: name, Time: timestamppb.Now()})
				if err != nil {
					level.Error(n.logger).Log("msg", "failed to Observe", "name", name, "err", err)
				}

				n.sch.Step()
			}
		}
	}()

	level.Debug(n.logger).Log("msg", "timer scheduler started",
		"repeated_events", len(n.cfg.RepeatEvents),
		"events", len(n.cfg.Events),
	)

	<-ctx.Done()
	level.Debug(n.logger).Log("msg", "scheduler done")

	return nil
}

func (n *Named) scheduleRepeatEvents(scheduledEvents *events.Scheduler, v RepeatEventConfig) error {

	// Stop calculating events beyond this time.
	end := time.Now().Add(time.Duration(n.cfg.FutureLimit) * time.Second)

	next := time.Now()
	for {
		next = next.Add(time.Duration(v.Every.Seconds) * time.Second)

		if next.Before(end) {
			err := scheduledEvents.Set(next, v.Produce)
			if err != nil {
				return err
			}
			continue
		}

		if next.After(end) {
			break
		}
	}

	return nil
}
func (n *Named) scheduleEvents(v EventConfig) error {
	loc, err := time.LoadLocation(n.cfg.TimeZone)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("15:04:05", v.Time, loc)
	if err != nil {
		return err
	}

	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)

	weekDayMatch := func(days []string) bool {
		for _, d := range days {
			if now.Weekday().String() == d {
				return true
			}
		}

		return false
	}(v.Days)

	if !weekDayMatch {
		return nil
	}

	timeRemaining := time.Until(d)

	if timeRemaining > 0 {
		err = n.sch.Set(d, v.Produce)
		if err != nil {
			return err
		}
	}

	return nil
}

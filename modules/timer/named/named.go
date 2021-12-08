package named

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/xaque208/znet/pkg/events"
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

func (t *Named) starting(ctx context.Context) error {
	if len(t.cfg.Events) == 0 && len(t.cfg.RepeatEvents) == 0 {
		return fmt.Errorf("no Events or RepeatEvents config")
	}

	return nil
}

func (t *Named) running(ctx context.Context) error {
	err := t.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to Connect named")
	}

	<-ctx.Done()
	return nil
}

func (t *Named) stopping(_ error) error {
	return nil
}

// Connect implements events.Producer
func (t *Named) Connect(ctx context.Context) error {
	_ = t.logger.Log("msg", "starting eventProducer")

	go func() {
		err := t.scheduler(ctx)
		if err != nil {
			_ = level.Error(t.logger).Log("msg", "scheduler failed",
				"err", err,
			)
		}
	}()

	return nil
}

func (t *Named) scheduler(ctx context.Context) error {
	namedClient := NewNamedClient(t.conn)

	go func() {
		for {
			for _, repeatEvent := range t.cfg.RepeatEvents {
				times := t.sch.TimesForName(repeatEvent.Produce)
				if len(times) == 0 {
					err := t.scheduleRepeatEvents(t.sch, repeatEvent)
					if err != nil {
						_ = level.Error(t.logger).Log("msg", "failed to schedule repeat events",
							"repeatEvent", repeatEvent.Produce,
							"err", err,
						)
					}
				}
			}

			for _, event := range t.cfg.Events {
				if len(t.sch.TimesForName(event.Produce)) == 0 {
					err := t.scheduleEvents(event)
					if err != nil {
						_ = level.Error(t.logger).Log("msg", "failed to schedul events",
							"event", event.Produce,
							"err", err,
						)
					}
				}
			}

			t.sch.Report()

			names := t.sch.WaitForNext()

			if len(names) == 0 {
				continue
			}

			for _, name := range names {
				_, err := namedClient.Observe(ctx, &NamedTimeStamp{Name: name, Time: timestamppb.Now()})
				if err != nil {
					_ = level.Error(t.logger).Log("msg", "failed to Observe", "name", name, "err", err)
				}

				t.sch.Step()
			}
		}
	}()

	_ = level.Debug(t.logger).Log("msg", "timer scheduler started",
		"repeated_events", len(t.cfg.RepeatEvents),
		"events", len(t.cfg.Events),
	)

	<-ctx.Done()
	_ = level.Debug(t.logger).Log("msg", "scheduler done")

	return nil
}

func (t *Named) scheduleRepeatEvents(scheduledEvents *events.Scheduler, v RepeatEventConfig) error {

	// Stop calculating events beyond this time.
	end := time.Now().Add(time.Duration(t.cfg.FutureLimit) * time.Second)

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
func (t *Named) scheduleEvents(v EventConfig) error {
	loc, err := time.LoadLocation(t.cfg.TimeZone)
	if err != nil {
		return err
	}

	timestamp, err := time.ParseInLocation("15:04:05", v.Time, loc)
	if err != nil {
		return err
	}

	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), 0, loc)

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
		err = t.sch.Set(d, v.Produce)
		if err != nil {
			return err
		}
	}

	return nil
}

package timer

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/pkg/events"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	conn   *grpc.ClientConn
	config Config
	ctx    context.Context
	cancel func()
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config Config) events.Producer {
	ctx, cancel := context.WithCancel(context.Background())

	var producer events.Producer = &EventProducer{
		conn:   conn,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	return producer
}

// Start initializes the producer.
func (e *EventProducer) Start() error {
	log.Info("starting timer eventProducer")

	go func() {
		err := e.scheduler()
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	e.cancel()
	return nil
}

func (e *EventProducer) scheduleEvents(scheduledEvents *events.Scheduler, v EventConfig) error {
	loc, err := time.LoadLocation(e.config.TimeZone)
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
		log.Tracef("skipping non-weekday match")
		return nil
	}

	timeRemaining := time.Until(d)

	if timeRemaining > 0 {
		err = scheduledEvents.Set(d, v.Produce)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EventProducer) scheduleRepeatEvents(scheduledEvents *events.Scheduler, v RepeatEventConfig) error {

	// Stop calculating events beyond this time.
	end := time.Now().Add(time.Duration(e.config.FutureLimit) * time.Second)

	next := time.Now()
	for {
		next = next.Add(time.Duration(v.Every.Seconds) * time.Second)

		log.WithFields(log.Fields{
			"next": next,
			"end":  end,
		}).Trace("timer")

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

func (e *EventProducer) scheduler() error {
	sch := events.NewScheduler()

	log.WithFields(log.Fields{
		"event_count": len(sch.All()),
	}).Debug("timer scheduler started")

	go func() {
		for {
			for _, repeatEvent := range e.config.RepeatEvents {
				times := sch.TimesForName(repeatEvent.Produce)
				if len(times) == 0 {
					err := e.scheduleRepeatEvents(sch, repeatEvent)
					if err != nil {
						log.Error(err)
					}
				}
			}

			for _, event := range e.config.Events {
				times := sch.TimesForName(event.Produce)
				if len(times) == 0 {
					err := e.scheduleEvents(sch, event)
					if err != nil {
						log.Error(err)
					}

					if len(sch.All()) == 0 {
						dur := 10 * time.Minute
						log.Debugf("no timer names found after reschedule, retry in %s seconds", dur)
						time.Sleep(dur)
					}

					// continue?

				}
			}

			sch.Report()

			names := sch.WaitForNext()

			if len(names) == 0 {
				continue
			}

			for _, n := range names {
				now := time.Now()

				ev := NamedTimer{
					Name: n,
					Time: &now,
				}

				err := events.ProduceEvent(e.conn, ev)
				if err != nil {
					log.Error(err)
				}

				sch.Step()
			}
		}
	}()

	<-e.ctx.Done()
	log.Debug("scheduler dying")

	return nil
}

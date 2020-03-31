package timer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	conn    *grpc.ClientConn
	config  Config
	diechan chan bool
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config Config) events.Producer {
	var producer events.Producer = &EventProducer{
		conn:   conn,
		config: config,
	}

	return producer
}

// Start initializes the producer.
func (e *EventProducer) Start() error {
	log.Info("starting timer eventProducer")
	e.diechan = make(chan bool)

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
	e.diechan <- true
	close(e.diechan)

	return nil
}

func (e *EventProducer) scheduleEvents(scheduledEvents *events.Scheduler) error {

	for _, v := range e.config.Events {

		loc, err := time.LoadLocation(e.config.TimeZone)
		if err != nil {
			log.Error(err)
			continue
		}

		t, err := time.ParseInLocation("15:04:05", v.Time, loc)
		if err != nil {
			log.Error(err)
			continue
		}

		now := time.Now()
		d := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)
		timeRemaining := time.Until(d)

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
			continue
		}

		if timeRemaining <= 0 {
			log.Tracef("skipping past event %s", v.Produce)
			continue
		}

		log.Infof("future event: %+v", d)

		if timeRemaining > 0 {
			scheduledEvents.Set(d, v.Produce)
		}

	}

	return nil
}

func (e *EventProducer) scheduleRepeatEvents(scheduledEvents *events.Scheduler) error {

	// Stop calculating events beyond this time.
	end := time.Now().Add(time.Duration(e.config.FutureLimit) * time.Second)
	scheduledEvents.Set(end, "ReloadConfig")

	for _, v := range e.config.RepeatEvents {
		next := time.Now()
		for {
			next = next.Add(time.Duration(v.Every.Seconds) * time.Second)

			log.Tracef("Repeat evert is: %+v", v.Every.Seconds)
			log.Tracef("Next is: %+v", next)
			log.Tracef("End is: %+v", end)
			log.Tracef("next.Before(end) is: %+v", next.Before(end))

			// TODO make the map handling here simpler.  Perhaps use the Scheduler interface
			// e.Schedule.Set(next, v.Produce)

			if next.Before(end) {
				scheduledEvents.Set(next, v.Produce)
				continue
			}

			break
		}
	}

	return nil
}

func (e *EventProducer) scheduler() error {
	log.Debug("timer scheduler started")

	sch := events.NewScheduler()

	err := e.scheduleEvents(sch)
	if err != nil {
		log.Error(err)
	}

	err = e.scheduleRepeatEvents(sch)
	if err != nil {
		log.Error(err)
	}

	log.Infof("%d timer events scheduled", len(sch.All()))

	otherchan := make(chan bool, 1)

	go func() {
		for {
			names := sch.WaitForNext()

			for _, n := range names {
				now := time.Now()

				ev := NamedTimer{
					Name: n,
					Time: &now,
				}

				err := e.Produce(ev)
				if err != nil {
					log.Error(err)
				}

				sch.Step()
			}
		}
	}()

	for {
		select {
		case <-otherchan:

		case <-e.diechan:
			log.Debugf("scheduler dying")
			return nil
		}
	}

}

// Produce implements the events.Producer interface.  Match the supported event
// types to know which event to notice, and then send notice of the event to
// the RPC server.
func (e *EventProducer) Produce(ev interface{}) error {
	// Create the RPC client
	ec := pb.NewEventsClient(e.conn)
	t := reflect.TypeOf(ev).String()

	var req *pb.Event

	switch t {
	case "timer.ExpiredTimer":
		x := ev.(ExpiredTimer)
		req = x.Make()
	case "timer.NamedTimer":
		x := ev.(NamedTimer)
		req = x.Make(x.Name)
	default:
		return fmt.Errorf("unhandled event type: %T", ev)
	}

	log.Tracef("producing RPC event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

package timer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
	"google.golang.org/grpc"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	conn   *grpc.ClientConn
	Config Config
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config Config) events.Producer {
	var producer events.Producer = &EventProducer{
		conn:   conn,
		Config: config,
	}

	SpawnReloader(producer, config)
	SpawnProducers(producer, config)

	return producer
}

// SpawnReloader creates a Go routine that waits on a timer for the next
// midnight to arrive, which schedules the next timers, and the next reloader.
func SpawnReloader(producer events.Producer, config Config) {
	now := time.Now()
	tomorrowNow := now.Add(time.Hour * 24)

	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		log.Error(err)
	}

	nextMidnight := time.Date(tomorrowNow.Year(), tomorrowNow.Month(), tomorrowNow.Day(), 0, 0, 0, 0, loc)
	timeRemaining := time.Until(nextMidnight)

	log.Debug("spawning timer reloader")
	go func(timeRemaining time.Duration, producer events.Producer, config Config) {
		t := time.NewTimer(timeRemaining)
		<-t.C

		SpawnReloader(producer, config)
		SpawnProducers(producer, config)
	}(timeRemaining, producer, config)
}

// SpawnProducers creates go routines that wait for a period of time and then
// produce an event.  Only events in the future for the current day are
// scheduled.  This is called daily by SpawnReloader at midnight.
func SpawnProducers(producer events.Producer, config Config) {
	log.Debug("spawning timers")

	for _, e := range config.Events {
		loc, err := time.LoadLocation(config.TimeZone)
		if err != nil {
			log.Error(err)
			continue
		}

		t, err := time.ParseInLocation("15:04:05", e.Time, loc)
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
		}(e.Days)

		if !weekDayMatch {
			log.Tracef("skipping non-weekday match")
			continue
		}

		if timeRemaining <= 0 {
			log.Tracef("skipping past event %s", e.Produce)
			continue
		}

		if timeRemaining > 0 {
			t := time.Now()
			ev := NamedTimer{
				Name: e.Produce,
				Time: &t,
			}

			log.Debugf("starting timer %s ending at %+s", e.Produce, d)

			go func(timeRemaining time.Duration, producer events.Producer, event NamedTimer) {
				t := time.NewTimer(timeRemaining)
				<-t.C

				err := producer.Produce(ev)
				if err != nil {
					log.Error(err)
				}
			}(timeRemaining, producer, ev)
		}
	}
}

// Produce satisies the events.Producer interface.  Match the supported event
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
		return fmt.Errorf("Unhandled event type: %T", ev)
	}

	log.Tracef("producing event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

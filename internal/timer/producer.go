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

var eventNames []string

// EventProducer implements events.Producer with an attached GRPC connection.
type EventProducer struct {
	conn *grpc.ClientConn
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config TimerConfig) events.Producer {
	var producer events.Producer = &EventProducer{
		conn: conn,
	}

	return producer
}

// SpawnProducers creates go routines that wait for a period of time and then produce an event.
func SpawnProducers(producer events.Producer, config TimerConfig) {

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
			log.Debugf("skipping non-weekday match")
			continue
		}

		if timeRemaining <= 0 {
			log.Debugf("skipping past event %s", e.Produce)
			continue
		}

		if timeRemaining > 0 {
			log.Debugf("timer %s ending at %+s", e.Produce, d)

			t := time.Now()

			ev := NamedTimer{
				Name: e.Produce,
				Time: &t,
			}

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

func (e *EventProducer) Names() []string {
	return eventNames
}

func (e *EventProducer) Produce(ev interface{}) error {
	ec := pb.NewEventsClient(e.conn)
	t := reflect.TypeOf(ev).String()

	var req *pb.Event

	switch t {
	case "timer.TimerExpired":
		x := ev.(TimerExpired)
		req = x.Make()
	case "timer.NamedTimer":
		x := ev.(NamedTimer)
		req = x.Make(x.Name)
	default:
		return fmt.Errorf("Unhandled event type: %T", ev)
	}

	log.Debugf("producing event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

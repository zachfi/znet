package astro

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
)

// EventProducer implements events.Producer with an attached GRPC connection.
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
	log.Info("starting astro eventProducer")

	e.diechan = make(chan bool)
	err := e.scheduler()
	if err != nil {
		log.Error(err)
	}

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	e.diechan <- true
	close(e.diechan)

	return nil
}

func (e *EventProducer) scheduleEvents(sch *events.Scheduler) error {
	clientConf := api.Config{
		Address: e.config.MetricsURL,
	}

	client, err := api.NewClient(clientConf)
	if err != nil {
		log.Error(err)
	}

	for _, l := range e.config.Locations {
		sunriseTime := queryForTime(client, fmt.Sprintf("owm_sunrise_time{location=\"%s\"}", l))
		sunsetTime := queryForTime(client, fmt.Sprintf("owm_sunset_time{location=\"%s\"}", l))

		sch.Set(sunriseTime, "Sunrise")
		sch.Set(sunsetTime, "Sunset")

		preSunset := sunsetTime.Add(-1 * time.Hour)

		sch.Set(preSunset, "PreSunset")
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

	log.Infof("%d astro events scheduled", len(sch.All()))

	otherchan := make(chan bool, 1)

	go func() {
		for {
			log.Info("astro")
			names := sch.WaitForNext()

			for _, n := range names {
				now := time.Now()

				ev := SolarEvent{
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
	case "astro.SolarEvent":
		x := ev.(SolarEvent)
		req = x.Make()
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

func queryForTime(client api.Client, query string) time.Time {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		log.Errorf("Error querying Prometheus: %+v\n", err)
	}

	if len(warnings) > 0 {
		log.Warn(warnings)
	}

	switch {
	case result.Type() == model.ValVector:
		vectorVal := result.(model.Vector)
		for _, elem := range vectorVal {
			i, err := strconv.ParseInt(elem.Value.String(), 10, 64)
			if err != nil {
				log.Error(err)
				continue
			}

			return time.Unix(i, 0)
		}
	}

	return time.Time{}
}

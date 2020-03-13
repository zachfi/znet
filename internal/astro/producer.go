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
	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
	"google.golang.org/grpc"
)

// EventProducer implements events.Producer with an attached GRPC connection.
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

	log.Debug("spawning astro reloader")
	go func(timeRemaining time.Duration, producer events.Producer, config Config) {
		t := time.NewTimer(timeRemaining)
		<-t.C

		SpawnReloader(producer, config)
		SpawnProducers(producer, config)
	}(timeRemaining, producer, config)
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
	case "astro.AstroEvent":
		x := ev.(AstroEvent)
		req = x.Make()
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

func SpawnProducers(producer events.Producer, config Config) {

	clientConf := api.Config{
		Address: config.MetricsURL,
	}

	client, err := api.NewClient(clientConf)
	if err != nil {
		log.Error(err)
	}

	for _, l := range config.Locations {
		sunriseTime := queryForTime(client, fmt.Sprintf("owm_sunrise_time{location=\"%s\"}", l))
		sunsetTime := queryForTime(client, fmt.Sprintf("owm_sunset_time{location=\"%s\"}", l))

		timeUntilSunrise := time.Until(sunriseTime)
		timeUntilSunset := time.Until(sunsetTime)

		// The function to create the RPC event.
		f := func(timeRemaining time.Duration, producer events.Producer, event AstroEvent) {
			t := time.NewTimer(timeRemaining)
			<-t.C

			now := time.Now()
			event.Time = &now

			err := producer.Produce(event)
			if err != nil {
				log.Error(err)
			}
		}

		if timeUntilSunrise > 0 {
			ev := AstroEvent{
				Name: "Sunrise",
			}

			log.Debugf("starting timer Sunrise at %s", sunsetTime)
			go f(timeUntilSunrise, producer, ev)
		}

		if timeUntilSunset > 0 {
			ev := AstroEvent{
				Name: "Sunset",
			}

			log.Debugf("starting timer Sunset at %s", sunsetTime)
			go f(timeUntilSunset, producer, ev)
		}

	}

}

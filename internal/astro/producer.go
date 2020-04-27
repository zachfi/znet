package astro

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/pkg/events"
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

		log.Tracef("astro found sunriseTime: %+v", sunriseTime)
		log.Tracef("astro found sunsetTime: %+v", sunsetTime)

		// Schedule tomorrow's sunrise based on today
		if time.Since(sunriseTime) > 0 {
			sunriseTime = sunriseTime.Add(24 * time.Hour)
		}

		err := sch.Set(sunriseTime, "Sunrise")
		if err != nil {
			log.Error(err)
		}

		// Schedule tomorrow's sunset based on today
		if time.Since(sunsetTime) > 0 {
			sunsetTime = sunsetTime.Add(24 * time.Hour)
		}

		err = sch.Set(sunsetTime, "Sunset")
		if err != nil {
			log.Error(err)
		}

		preSunset := sunsetTime.Add(-75 * time.Minute)

		err = sch.Set(preSunset, "PreSunset")
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (e *EventProducer) scheduler() error {
	log.Debug("astro scheduler started")

	sch := events.NewScheduler()

	err := e.scheduleEvents(sch)
	if err != nil {
		log.Error(err)
	}

	log.Infof("%d astro events scheduled: %+v", len(sch.All()), sch.All())

	otherchan := make(chan bool, 1)

	go func() {
		for {
			names := sch.WaitForNext()

			if len(names) == 0 {
				dur := 1 * time.Hour
				log.Debugf("no astro names, retry in %s", dur)
				time.Sleep(dur)
				continue
			}

			for _, n := range names {
				now := time.Now()

				ev := SolarEvent{
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

	for {
		select {
		case <-otherchan:

		case <-e.diechan:
			log.Debugf("scheduler dying")
			return nil
		}
	}
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

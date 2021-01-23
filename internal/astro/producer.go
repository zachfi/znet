package astro

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/events"
)

// EventProducer implements events.Producer with an attached GRPC connection.
type EventProducer struct {
	sync.Mutex
	conn   *grpc.ClientConn
	config *config.AstroConfig
}

// NewProducer receives a config to build a new EventProducer.
func NewProducer(cfg *config.AstroConfig) (events.Producer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("unable to create new astro producer from nil config")
	}

	var producer events.Producer = &EventProducer{
		config: cfg,
	}

	return producer, nil
}

// Connect starts the producer
func (e *EventProducer) Connect(ctx context.Context, conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("unable to connext with nil gRPC connection")
	}

	e.Lock()
	defer e.Unlock()

	if e.conn != nil {
		log.Warnf("replacing non-nil gRPC client connection")
		log.Debug("closing old connection")

		err := e.conn.Close()
		if err != nil {
			log.Error(err)
		}
	}

	e.conn = conn

	log.Info("starting astro eventProducer")

	go func() {
		err := e.scheduler(ctx)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (e *EventProducer) scheduleEvents(ctx context.Context, sch *events.Scheduler) error {
	clientConf := api.Config{
		Address: e.config.MetricsURL,
	}

	client, err := api.NewClient(clientConf)
	if err != nil {
		log.Error(err)
	}

	for _, l := range e.config.Locations {
		sunriseTime := queryForTime(ctx, client, fmt.Sprintf("owm_sunrise_time{location=\"%s\"}", l))
		sunsetTime := queryForTime(ctx, client, fmt.Sprintf("owm_sunset_time{location=\"%s\"}", l))

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

func (e *EventProducer) scheduler(ctx context.Context) error {
	sch := events.NewScheduler()

	err := e.scheduleEvents(ctx, sch)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"event_count": len(sch.All()),
	}).Debug("astro scheduler started")

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
				astroClient := NewAstroClient(e.conn)

				switch n {
				case "Sunrise":
					_, err := astroClient.Sunrise(ctx, &Empty{})
					if err != nil {
						log.Error(err)
					}
				case "Sunset":
					_, err := astroClient.Sunset(ctx, &Empty{})
					if err != nil {
						log.Error(err)
					}
				case "PreSunset":
					_, err := astroClient.PreSunset(ctx, &Empty{})
					if err != nil {
						log.Error(err)
					}
				default:
					log.Warnf("unknown astro event name: %s", n)
				}

				sch.Step()

			}
		}
	}()

	<-ctx.Done()
	log.Debugf("scheduler dying")

	return nil
}

func queryForTime(c context.Context, client api.Client, query string) time.Time {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
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

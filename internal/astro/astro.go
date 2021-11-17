package astro

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	grpc "google.golang.org/grpc"

	"github.com/xaque208/znet/modules/lights"
	"github.com/xaque208/znet/pkg/events"
)

type Astro struct {
	UnimplementedAstroServer
	services.Service

	logger log.Logger

	cfg *Config

	lights *lights.Lights
	sch    *events.Scheduler
	conn   *grpc.ClientConn
}

func New(cfg Config, logger log.Logger, conn *grpc.ClientConn) (*Astro, error) {
	a := &Astro{
		cfg:    &cfg,
		conn:   conn,
		logger: log.With(logger, "timer", "astro"),
		sch:    events.NewScheduler(logger),
	}

	a.Service = services.NewBasicService(a.starting, a.running, a.stopping)

	return a, nil
}

// Sunrise implements AstroServer.
func (a *Astro) Sunrise(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "Sunrise")
}

// Sunset implements AstroServer.
func (a *Astro) Sunset(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "SunSet")
}

// PreSunset implements AstroServer.
func (a *Astro) PreSunset(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "PreSunset")
}

func (a *Astro) starting(ctx context.Context) error {
	return nil
}

func (a *Astro) running(ctx context.Context) error {
	err := a.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to Connect astro")
	}

	<-ctx.Done()
	return nil
}

func (a *Astro) stopping(_ error) error {
	return nil
}

// Connect implements events.Producer
func (a *Astro) Connect(ctx context.Context) error {
	_ = a.logger.Log("msg", "starting eventProducer")

	go func() {
		err := a.scheduler(ctx)
		if err != nil {
			_ = level.Error(a.logger).Log("msg", "scheduler failed",
				"err", err,
			)
		}
	}()

	return nil
}

func (a *Astro) scheduleEvents(ctx context.Context, sch *events.Scheduler) error {
	clientConf := api.Config{
		Address: a.cfg.MetricsURL,
	}

	client, err := api.NewClient(clientConf)
	if err != nil {
		_ = level.Error(a.logger).Log("msg", "failed to create a prometheus api client",
			"err", err,
		)
	}

	for _, l := range a.cfg.Locations {
		sunriseTime, err := a.queryForTime(ctx, client, fmt.Sprintf("owm_sunrise_time{location=\"%s\"}", l))
		if err != nil {
			return errors.Wrap(err, "query failed")
		}

		sunsetTime, err := a.queryForTime(ctx, client, fmt.Sprintf("owm_sunset_time{location=\"%s\"}", l))
		if err != nil {
			return errors.Wrap(err, "query failed")
		}

		// Schedule tomorrow's sunrise based on today
		if time.Since(sunriseTime) > 0 {
			sunriseTime = sunriseTime.Add(24 * time.Hour)
		}

		err = sch.Set(sunriseTime, "Sunrise")
		if err != nil {
			_ = level.Error(a.logger).Log("msg", "failed to set Sunrise",
				"time", sunriseTime,
				"err", err,
			)
		}

		// Schedule tomorrow's sunset based on today
		if time.Since(sunsetTime) > 0 {
			sunsetTime = sunsetTime.Add(24 * time.Hour)
		}

		err = sch.Set(sunsetTime, "Sunset")
		if err != nil {
			_ = level.Error(a.logger).Log("msg", "failed to set Sunset",
				"time", sunsetTime,
				"err", err,
			)
		}

		preSunset := sunsetTime.Add(-75 * time.Minute)

		err = sch.Set(preSunset, "PreSunset")
		if err != nil {
			_ = level.Error(a.logger).Log("msg", "failed to set PreSunset",
				"time", preSunset,
				"err", err,
			)
		}
	}

	return nil
}

func (a *Astro) scheduler(ctx context.Context) error {
	err := a.scheduleEvents(ctx, a.sch)
	if err != nil {
		_ = level.Error(a.logger).Log("msg", "failed to schedule events",
			"err", err,
		)
	}

	astroClient := NewAstroClient(a.conn)

	go func() {
		for {
			names := a.sch.WaitForNext()

			if len(names) == 0 {
				dur := 1 * time.Hour

				_ = level.Debug(a.logger).Log("msg", "no astro names",
					"retry", dur,
				)
				time.Sleep(dur)
				continue
			}

			for _, n := range names {

				switch n {
				case "Sunrise":
					_, err := astroClient.Sunrise(ctx, &Empty{})
					if err != nil {
						_ = level.Error(a.logger).Log("msg", "failed to call Sunrise",
							"err", err,
						)
					}
				case "Sunset":
					_, err := astroClient.Sunset(ctx, &Empty{})
					if err != nil {
						_ = level.Error(a.logger).Log("msg", "failed to call Sunset",
							"err", err,
						)
					}
				case "PreSunset":
					_, err := astroClient.PreSunset(ctx, &Empty{})
					if err != nil {
						_ = level.Error(a.logger).Log("msg", "failed to call PreSunset",
							"err", err,
						)
					}
				default:
					_ = level.Warn(a.logger).Log("msg", "unknown event name",
						"event", n,
					)
				}

				a.sch.Step()
			}
		}
	}()

	<-ctx.Done()
	_ = level.Debug(a.logger).Log("msg", "scheduler done")

	return nil
}

func (a *Astro) queryForTime(c context.Context, client api.Client, query string) (time.Time, error) {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		return time.Time{}, errors.Wrap(err, "prometheus API call failed")
	}

	if len(warnings) > 0 {
		_ = level.Warn(a.logger).Log("msg", "warnings from prometheus",
			"warnings", warnings,
		)
	}

	if result != nil {
		switch {
		case result.Type() == model.ValVector:
			vectorVal := result.(model.Vector)
			for _, elem := range vectorVal {
				i, err := strconv.ParseInt(elem.Value.String(), 10, 64)
				if err != nil {
					_ = level.Error(a.logger).Log("msg", "failed to parse int",
						"err", err,
					)
					continue
				}

				return time.Unix(i, 0), nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("failed to query for time")
}

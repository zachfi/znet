package timer

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/modules/timer/named"

	"github.com/grafana/dskit/services"
)

type Timer struct {
	services.Service
	cfg *Config

	logger log.Logger

	// Manager for subservices
	subservices        *services.Manager
	subservicesWatcher *services.FailureWatcher

	// gRPC services.
	Astro astro.AstroServer
	Named named.NamedServer
}

func New(cfg Config, logger log.Logger, conn *grpc.ClientConn) (*Timer, error) {

	var err error
	subservices := []services.Service(nil)

	logger = log.With(logger, "module", "timer")

	a, err := astro.New(cfg.Astro, logger, conn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Astro")
	}
	subservices = append(subservices, a)

	n, err := named.New(cfg.Named, logger, conn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Named")
	}
	subservices = append(subservices, n)

	t := &Timer{
		cfg:    &cfg,
		logger: logger,
	}

	// GRPC services to attach to our caller's server.
	t.Named = n
	t.Astro = a

	t.subservices, err = services.NewManager(subservices...)
	if err != nil {
		return nil, fmt.Errorf("failed to create subservices %w", err)
	}

	t.subservicesWatcher = services.NewFailureWatcher()
	t.subservicesWatcher.WatchManager(t.subservices)

	t.Service = services.NewBasicService(t.starting, t.running, t.stopping)
	return t, nil
}

func (t *Timer) starting(ctx context.Context) error {
	err := services.StartManagerAndAwaitHealthy(ctx, t.subservices)
	if err != nil {
		return errors.Wrap(err, "failed to start subservices")
	}
	return nil
}

func (t *Timer) running(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-t.subservicesWatcher.Chan():
		return errors.Wrap(err, "timer subservices failed")
	}
}

func (t *Timer) stopping(_ error) error {
	return nil
}

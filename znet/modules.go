package znet

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/internal/timer/named"
	"github.com/xaque208/znet/modules/harvester"
)

type module int

const (
	Server string = "server"

	// Harvester string = "harvester"
	Timer string = "timer"

	// Telemetry string = "telemetry"

	// TODO currently we are using openweathermap_exporter as a source of data
	// for astro data when sending events to the server.  Perhaps it makes more
	// sense to include the exporter as a module, and then we can either
	// reference internal values, or provide RPC modules for requesting the
	// information necessary.  Without this, there is a requried coordination
	// between the exporter (another project), and the znet.
	// Weather string = "telemetry"

	All string = "all"
)

func (z *Znet) setupModuleManager() error {
	mm := modules.NewManager(z.logger)
	mm.RegisterModule(Server, z.initServer, modules.UserInvisibleModule)
	// mm.RegisterModule(Harvester, z.initHarvester)
	mm.RegisterModule(Timer, z.initTimer)
	mm.RegisterModule(All, nil)

	deps := map[string][]string{
		// Server:       nil,

		// Timer: {Server, Lights},
		// IOT: {Server},
		// Inventory: {Server},
		// Lights: {Server},
		// Astro: {Server},
		// Telemetry: {Server, IOT, Lights},

		// TODO
		// gitwatch:       nil,
		// build:       nil,
		// agent:       nil,

		// Harvester: {Server},
		Timer: {Server},
		All:   {Server, Timer},
	}

	for mod, targets := range deps {
		if err := mm.AddDependency(mod, targets...); err != nil {
			return err
		}
	}

	z.ModuleManager = mm

	return nil
}

func (z *Znet) initTimer() (services.Service, error) {

	conn := comms.SlimRPCClient(z.cfg.RPC.ServerAddress)

	t, err := timer.New(z.cfg.Timer, z.logger, conn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to init timer")
	}

	astro.RegisterAstroServer(z.Server.GRPC, t.Astro)
	named.RegisterNamedServer(z.Server.GRPC, t.Named)

	z.timer = t
	return t, nil
}

func (z *Znet) initHarvester() (services.Service, error) {
	h, err := harvester.New(z.cfg.Harvester)
	if err != nil {
		return nil, fmt.Errorf("failed to create harvester")
	}

	z.harvester = h
	return h, nil
}

func (z *Znet) initServer() (services.Service, error) {
	z.cfg.Server.MetricsNamespace = metricsNamespace
	z.cfg.Server.ExcludeRequestInLog = true

	// cortex.DisableSignalHandling(&t.cfg.Server)

	server, err := server.New(z.cfg.Server)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create server")
	}

	servicesToWaitFor := func() []services.Service {
		svs := []services.Service(nil)
		for m, s := range z.serviceMap {
			// Server should not wait for itself.
			if m != Server {
				svs = append(svs, s)
			}
		}
		return svs
	}

	z.Server = server

	serverDone := make(chan error, 1)

	runFn := func(ctx context.Context) error {
		go func() {
			defer close(serverDone)
			serverDone <- server.Run()
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-serverDone:
			if err != nil {
				return err
			}
			return fmt.Errorf("server stopped unexpectedly")
		}
	}

	stoppingFn := func(_ error) error {
		// wait until all modules are done, and then shutdown server.
		for _, s := range servicesToWaitFor() {
			_ = s.AwaitTerminated(context.Background())
		}

		// shutdown HTTP and gRPC servers (this also unblocks Run)
		server.Shutdown()

		// if not closed yet, wait until server stops.
		<-serverDone
		level.Info(z.logger).Log("msg", "server stopped")
		return nil
	}

	return services.NewBasicService(nil, runFn, stoppingFn), nil
}

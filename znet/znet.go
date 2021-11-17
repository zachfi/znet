package znet

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	"github.com/weaveworks/common/signals"

	"github.com/xaque208/znet/modules/harvester"
	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/modules/lights"
	"github.com/xaque208/znet/modules/telemetry"
	"github.com/xaque208/znet/modules/timer"
	"github.com/xaque208/znet/pkg/util"
)

const metricsNamespace = "znet"

// Znet is the core object for this project.  It keeps track of the data,
// configuration and flow control for starting the server process.
type Znet struct {
	cfg Config

	// ConfigDir   string
	// Data        netconfig.Data
	// Environment map[string]string

	Server *server.Server

	logger log.Logger

	// Modules.
	// server *server.Server

	telemetry *telemetry.Telemetry
	harvester *harvester.Harvester
	timer     *timer.Timer

	inventory *inventory.Server
	lights    *lights.Lights

	ModuleManager *modules.Manager
	serviceMap    map[string]services.Service
}

// New creates and returns a new Znet object.
func New(cfg Config) (*Znet, error) {
	z := &Znet{
		cfg: cfg,
	}

	z.logger = util.NewLogger()

	if z.cfg.Target == "" {
		z.cfg.Target = All
	}

	// var err error
	// var environment map[string]string
	//
	// l.Log("msg", "loading config", "file", file)
	//
	// cfg, err := loadConfig(file)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to load config file %s: %w", file, err)
	// }
	//
	// if cfg.Environments != nil && cfg.Vault != nil {
	// 	e, err := getEnvironmentConfig(*cfg.Environments, "common")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get environment config: %w", err)
	// 	}
	//
	// 	environment, err = LoadEnvironment(cfg.Vault, e)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to load environment: %w", err)
	// 	}
	// } else {
	// 	level.Debug(l).Log("missing vault/environment config")
	// }
	//
	// z := Znet{
	// 	cfg:         cfg,
	// 	Environment: environment,
	// }

	if err := z.setupModuleManager(); err != nil {
		return nil, errors.Wrap(err, "failed to setup module manager")
	}

	return z, nil
}

func (z *Znet) Run() error {
	serviceMap, err := z.ModuleManager.InitModuleServices(z.cfg.Target)
	if err != nil {
		return fmt.Errorf("failed to init module services %w", err)
	}
	z.serviceMap = serviceMap

	servs := []services.Service(nil)
	for _, s := range serviceMap {
		servs = append(servs, s)
	}

	sm, err := services.NewManager(servs...)
	if err != nil {
		return fmt.Errorf("failed to start service manager %w", err)
	}

	// Listen for events from this manager, and log them.
	healthy := func() { _ = level.Info(z.logger).Log("msg", "zNet started") }
	stopped := func() { _ = level.Info(z.logger).Log("msg", "zNet stopped") }
	serviceFailed := func(service services.Service) {
		// if any service fails, stop everything
		sm.StopAsync()

		// let's find out which module failed
		for m, s := range serviceMap {
			if s == service {
				if service.FailureCase() == modules.ErrStopProcess {
					_ = level.Info(z.logger).Log("msg", "received stop signal via return error", "module", m, "err", service.FailureCase())
				} else {
					_ = level.Error(z.logger).Log("msg", "module failed", "module", m, "err", service.FailureCase())
				}
				return
			}
		}

		_ = level.Error(z.logger).Log("msg", "module failed", "module", "unknown", "err", service.FailureCase())
	}
	sm.AddListener(services.NewManagerListener(healthy, stopped, serviceFailed))

	// Setup signal handler. If signal arrives, we stop the manager, which stops all the services.
	handler := signals.NewHandler(z.Server.Log)
	go func() {
		handler.Loop()
		sm.StopAsync()
	}()

	// Start all services. This can really only fail if some service is already
	// in other state than New, which should not be the case.
	err = sm.StartAsync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start service manager %w", err)
	}

	return sm.AwaitStopped(context.Background())
}

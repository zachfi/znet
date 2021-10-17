package znet

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/modules"

	"github.com/xaque208/znet/modules/harvester"
	"github.com/xaque208/znet/modules/server"
	"github.com/xaque208/znet/pkg/netconfig"
)

const metricsNamespace = "znet"

// Znet is the core object for this project.  It keeps track of the data,
// configuration and flow control for starting the server process.
type Znet struct {
	cfg *Config

	ConfigDir   string
	Data        netconfig.Data
	Environment map[string]string

	logger log.Logger

	harvester *harvester.Harvester
	server    *server.Server
}

// NewZnet creates and returns a new Znet object.
func NewZnet(file string, l log.Logger) (*Znet, error) {
	var err error
	var environment map[string]string

	l.Log("msg", "loading config", "file", file)

	cfg, err := loadConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", file, err)
	}

	if cfg.Environments != nil && cfg.Vault != nil {
		e, err := getEnvironmentConfig(*cfg.Environments, "common")
		if err != nil {
			return nil, fmt.Errorf("failed to get environment config: %w", err)
		}

		environment, err = LoadEnvironment(cfg.Vault, e)
		if err != nil {
			return nil, fmt.Errorf("failed to load environment: %w", err)
		}
	} else {
		level.Debug(l).Log("missing vault/environment config")
	}

	z := Znet{
		cfg:         cfg,
		Environment: environment,
	}

	return &z, nil
}

func (z *Znet) Run() {
	mm := modules.NewManager(z.logger)
	mm.RegisterModule(Harvester.String(), z.initHarvester)
	mm.RegisterModule(Server.String(), z.initServer)
}

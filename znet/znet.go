package znet

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/netconfig"
)

// Znet is the core object for this project.  It keeps track of the data,
// configuration and flow control for starting the server process.
type Znet struct {
	ConfigDir   string
	Config      *config.Config
	Data        netconfig.Data
	Environment map[string]string
}

// NewZnet creates and returns a new Znet object.
func NewZnet(file string) (*Znet, error) {
	var err error
	var environment map[string]string

	cfg, err := config.LoadConfig(file)
	if err != nil {
		return &Znet{}, fmt.Errorf("failed to load config file %s: %s", file, err)
	}

	if cfg.Environments != nil && cfg.Vault != nil {
		e, err := getEnvironmentConfig(*cfg.Environments, "common")
		if err != nil {
			log.Error(err)
		}

		environment, err = LoadEnvironment(cfg.Vault, e)
		if err != nil {
			log.Errorf("failed to load environment: %s", err)
		}
	} else {
		log.Debug("missing vault/environment config")
	}

	z := Znet{
		Config:      cfg,
		Environment: environment,
	}

	return &z, nil
}

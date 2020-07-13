package znet

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/pkg/netconfig"
)

// Znet is the core object for this project.  It keeps track of the data,
// configuration and flow control for starting the server process.
type Znet struct {
	ConfigDir   string
	Config      Config
	Data        netconfig.Data
	Environment map[string]string
	Lights      *lights.Lights
}

// NewZnet creates and returns a new Znet object.
func NewZnet(file string) (*Znet, error) {
	var err error
	var environment map[string]string

	config, err := loadConfig(file)
	if err != nil {
		return &Znet{}, fmt.Errorf("failed to load config file %s: %s", file, err)
	}

	if config.Environments != nil && config.Vault != nil {
		e, err := GetEnvironmentConfig(*config.Environments, "common")
		if err != nil {
			log.Error(err)
		}

		environment, err = LoadEnvironment(*config.Vault, e)
		if err != nil {
			log.Errorf("failed to load environment: %s", err)
		}
	} else {
		log.Debug("missing vault/environment config")
	}

	var light *lights.Lights
	if config.Lights != nil {
		light = lights.NewLights(*config.Lights)
	} else {
		log.Debug("missing lights config")
	}

	z := Znet{
		Config:      config,
		Environment: environment,
		Lights:      light,
	}

	return &z, nil
}

// LoadConfig receives a file path for a configuration to load.
func loadConfig(file string) (Config, error) {
	filename, _ := filepath.Abs(file)
	log.Debugf("loading config from: %s", filename)
	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Stop the znet connections
func (z *Znet) Stop() {
}

// Subscriptions is yet to be used, but conforms to the interface for
// generating consumers of named events.
func (z *Znet) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	s.Subscribe("ExecutionResult", z.executionResultHandler)

	return s.Table
}

func (z *Znet) executionResultHandler(name string, payload events.Payload) error {
	var x agent.ExecutionResult

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	executionExitStatus.With(prometheus.Labels{
		"command": x.Command,
	}).Set(float64(x.ExitCode))

	executionDuration.With(prometheus.Labels{
		"command": x.Command,
	}).Set(float64(x.Duration))

	return nil
}

package znet

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	"gopkg.in/yaml.v2"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/modules/harvester"
	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/modules/lights"
	"github.com/xaque208/znet/modules/telemetry"
	"github.com/xaque208/znet/modules/timer"
	"github.com/xaque208/znet/pkg/iot"
)

type Config struct {
	Target       string `yaml:"target"`
	OtelEndpoint string `yaml:"otel_endpoint"`

	// Environments []config.EnvironmentConfig `yaml:"environments,omitempty"`
	// Vault        config.VaultConfig         `yaml:"vault,omitempty"`

	// modules
	Server    server.Config    `yaml:"server,omitempty"`
	Harvester harvester.Config `yaml:"harvester"`
	Inventory inventory.Config `yaml:"inventory"`
	IOT       iot.Config       `yaml:"iot"`
	Lights    lights.Config    `yaml:"lights"`
	Telemetry telemetry.Config `yaml:"telemetry"`
	Timer     timer.Config     `yaml:"timer"`

	RPC config.RPCConfig `yaml:"rpc,omitempty"`
}

// LoadConfig receives a file path for a configuration to load.
func LoadConfig(file string) (Config, error) {
	filename, _ := filepath.Abs(file)

	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return config, errors.Wrap(err, "failed to load yaml file")
	}

	return config, nil
}

// loadYamlFile unmarshals a YAML file into the received interface{} or returns an error.
func loadYamlFile(filename string, d interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RegisterFlagsAndApplyDefaults(prefix string, f *flag.FlagSet) {
	c.Target = All
	f.StringVar(&c.Target, "target", All, "target module")
	f.StringVar(&c.OtelEndpoint, "otel_endpoint", "", "otel endpoint, eg: tempo:4317")
}

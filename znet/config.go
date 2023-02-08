package znet

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	"gopkg.in/yaml.v2"

	ztrace "github.com/zachfi/zkit/pkg/tracing"

	"github.com/zachfi/znet/internal/config"
	"github.com/zachfi/znet/modules/harvester"
	"github.com/zachfi/znet/modules/inventory"
	"github.com/zachfi/znet/modules/lights"
	"github.com/zachfi/znet/modules/telemetry"
	"github.com/zachfi/znet/modules/timer"
	"github.com/zachfi/znet/pkg/iot"
)

type Config struct {
	Target string `yaml:"target"`

	Tracing ztrace.Config `yaml:"tracing,omitempty"`

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
	yamlFile, err := os.ReadFile(filename)
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
	c.Tracing.RegisterFlagsAndApplyDefaults("tracing", f)
}

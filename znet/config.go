package znet

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/modules/harvester"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Target string `yaml:"target"`

	// Environments []config.EnvironmentConfig `yaml:"environments,omitempty"`
	// Vault        config.VaultConfig         `yaml:"vault,omitempty"`

	// modules
	Server    server.Config    `yaml:"server,omitempty"`
	Harvester harvester.Config `yaml:"harvester"`
	Timer     timer.Config     `yaml:"timer"`

	RPC *config.RPCConfig `yaml:"rpc,omitempty"`
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

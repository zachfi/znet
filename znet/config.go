package znet

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/modules/harvester"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Environments *[]config.EnvironmentConfig `yaml:"environments,omitempty"`
	Vault        *config.VaultConfig         `yaml:"vault,omitempty"`

	// modules
	Harvester harvester.Config `yaml:"harvester"`

	RPC *config.RPCConfig `yaml:"rpc,omitempty"`
}

// loadConfig receives a file path for a configuration to load.
func loadConfig(file string) (*Config, error) {
	filename, _ := filepath.Abs(file)

	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load yaml file: %w", err)
	}

	return &config, nil
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

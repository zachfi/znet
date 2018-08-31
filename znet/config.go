package znet

import (
	"io/ioutil"
	"path/filepath"

	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Rooms    []Room `yaml:"rooms"`
	Endpoint string `yaml:"endpoint"`
}

type Room struct {
	Name string `yaml:"name"`
	IDs  []int  `yaml:"ids"`
}

func LoadConfig(file string) (*Config, error) {
	filename, _ := filepath.Abs(file)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Error(err)
		return &config, err
	}

	return &config, nil
}

func (c *Config) Room(name string) (Room, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return Room{}, errors.New(fmt.Sprintf("Room %s not found in config", name))
}

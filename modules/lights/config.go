package lights

import "fmt"

// Config is the configuration for Lights
type Config struct {
	Rooms       []LightsRoom `yaml:"rooms"`
	Hue         *HueConfig   `yaml:"hue,omitempty"`
	PartyColors []string     `yaml:"party_colors,omitempty"`
}

// LightsRoom is a collection of device entries.
type LightsRoom struct {
	Name   string `yaml:"name"`
	IDs    []int  `yaml:"ids"`
	HueIDs []int  `yaml:"hue"`

	// Names of events
	On     []string `yaml:"turn_on"`
	Off    []string `yaml:"turn_off"`
	Dim    []string `yaml:"dim"`
	Alert  []string `yaml:"alert"`
	Toggle []string `yaml:"toggle"`
}

// HueConfig is the configuration for Philips Hue.
type HueConfig struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
}

// Room return the Room object for a room given by name.
func (c *Config) Room(name string) (LightsRoom, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return LightsRoom{}, fmt.Errorf("room %s not found in config", name)
}

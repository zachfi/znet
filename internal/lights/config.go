package lights

import "fmt"

type LightsConfig struct {
	Rooms []Room      `yaml:"rooms"`
	Hue   HueConfig   `yaml:"hue,omitempty"`
	RFToy RFToyConfig `yaml:"rftoy,omitempty"`
}

type Room struct {
	Name   string `yaml:"name"`
	IDs    []int  `yaml:"ids"`
	HueIDs []int  `yaml:"hue"`

	// Names of events
	On  []string `yaml:"turn_on"`
	Off []string `yaml:"turn_off"`
	Dim []string `yaml:"dim"`
}

type HueConfig struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
}

type RFToyConfig struct {
	Endpoint string `yaml:"endpoint,omitempty"`
}

func (c *LightsConfig) Room(name string) (Room, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return Room{}, fmt.Errorf("Room %s not found in config", name)
}

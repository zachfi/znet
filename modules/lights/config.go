package lights

import (
	"fmt"
)

// Config is the configuration for Lights
type Config struct {
	Rooms       []Room     `yaml:"rooms"`
	Hue         *HueConfig `yaml:"hue,omitempty"`
	PartyColors []string   `yaml:"party_colors,omitempty"`
	TimeZone    string     `yaml:"timezone"`
}

// Room is a collection of device entries.
type Room struct {
	Name   string `yaml:"name"`
	IDs    []int  `yaml:"ids"`
	HueIDs []int  `yaml:"hue"`

	States []StateSpec `yaml:"states"`
}

type StateSpec struct {
	State      ZoneState         `yaml:"state"`
	Brightness *Brightness       `yaml:"brightness,omitempty"`
	ColorTemp  *ColorTemperature `yaml:"color_temp,omitempty"`
	Event      string            `yaml:"event"`
}

// Implements the Unmarshaler interface of the yaml pkg.
func (s *StateSpec) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var fm map[string]string
	if err := unmarshal(&fm); err != nil {
		return err
	}

	for k, v := range fm {
		if k == "event" {
			s.Event = v
		}

		if k == "state" {
			if val, ok := ZoneState_value[v]; ok {
				s.State = ZoneState(val)
			} else {
				return fmt.Errorf("cannot unmarshal '%s' into %T", v, s.State)
			}
		}

		if k == "brightness" {
			if val, ok := Brightness_value[v]; ok {
				b := Brightness(val)
				s.Brightness = &b
			} else {
				return fmt.Errorf("cannot unmarshal '%s' into %T", v, s.State)
			}
		}

		if k == "color_temp" {
			if val, ok := ColorTemperature_value[v]; ok {
				b := ColorTemperature(val)
				s.ColorTemp = &b
			} else {
				return fmt.Errorf("cannot unmarshal '%s' into %T", v, s.State)
			}
		}
	}

	return nil
}

// HueConfig is the configuration for Philips Hue.
type HueConfig struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
}

// Room return the Room object for a room given by name.
func (c *Config) Room(name string) (Room, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return Room{}, fmt.Errorf("room %s not found in config", name)
}

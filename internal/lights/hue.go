package lights

import (
	"fmt"
	"strings"

	"github.com/amimof/huego"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/config"
)

type hueLight struct {
	config *config.LightsConfig
	hue    *huego.Bridge
}

func NewHueLight(cfg *config.LightsConfig) (Handler, error) {
	if cfg.Hue == nil {
		return nil, fmt.Errorf("unable to create new hue light with nil config")
	}

	h := hueLight{
		config: cfg,
		hue:    huego.New(cfg.Hue.Endpoint, cfg.Hue.User),
	}

	return h, nil
}

// On turns off the Hue Light for a room.
func (l hueLight) On(groupName string) error {
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)
		var light *huego.Light

		light, err = l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"name": light.Name,
			}).Debug("turning on light")

			err = light.On()
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		log.WithFields(log.Fields{
			"group": g.Name,
		}).Debug("turning on light group")

		err = g.On()
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// Off turns off the Hue Light for a room.
func (l hueLight) Off(groupName string) error {
	// try the light by group first
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)
		var light *huego.Light

		// then try to get just the light
		light, err = l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"name": light.Name,
			}).Debug("turning off light")

			err = light.Off()
			if err != nil {
				log.Error(err)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"group": g.Name,
		}).Debug("turning off light group")

		err = g.Off()
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// Dim modifies the brightness of a light group.
func (l hueLight) Dim(groupName string, brightness int32) error {
	room, err := l.config.Room(groupName)
	if err != nil {
		log.Error(err)
	}

	groups, err := l.hue.GetGroups()
	if err != nil {
		log.Error(err)
	}

	for _, g := range groups {
		for _, i := range room.HueIDs {
			if g.ID == i {
				log.WithFields(log.Fields{
					"name":  g.Name,
					"state": g.State,
				}).Debug("setting group brightness")

				err := g.Bri(uint8(brightness))
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}

// Alert blinks all hueLight in the given light group.
func (l hueLight) Alert(groupName string) error {
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)

		// then try to get just the light
		light, err := l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"name": light.Name,
			}).Debug("alerting light")

			err := light.Alert("select")
			if err != nil {
				log.Error(err)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"group": g.Name,
		}).Debug("alerting light group")

		err := g.Alert("select")
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// Toggle a Hue light in the given light group.
func (l hueLight) Toggle(groupName string) error {
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)

		// then try to get just the light
		light, err := l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"name": light.Name,
			}).Debug("alerting light")

			if light.IsOn() {
				return light.Off()
			} else {
				return light.On()
			}
		}

	} else {
		log.WithFields(log.Fields{
			"group": g.Name,
		}).Debug("alerting light group")

		if g.IsOn() {
			return g.Off()
		} else {
			return g.On()
		}
	}

	return nil
}

// SetColor applies the color to the light group.
func (l hueLight) SetColor(groupName string, hex string) error {
	return nil
}

// RandomColor applies a color selected at random to the light group.
func (l hueLight) RandomColor(groupName string, hex []string) error {
	return nil
}

// GetGroup calls the Hue bridge and looks for a group, that, when normalized,
// matches the name received.
func (l hueLight) getGroup(groupName string) (*huego.Group, error) {
	groups, err := l.hue.GetGroups()
	if err != nil {
		log.Error(err)
	}

	log.Tracef("found HUE groups: %+v", groups)

	for _, g := range groups {
		flatName := strings.ToLower(strings.ReplaceAll(g.Name, " ", "_"))

		if groupName == flatName {
			return &g, nil
		}
	}

	return &huego.Group{}, fmt.Errorf("group %s not found", groupName)
}

// getLight calls the Hue bridge and looks for a light, that, when normalized,
// matches the name received.
func (l hueLight) getLight(lightName string) (*huego.Light, error) {
	lights, err := l.hue.GetLights()
	if err != nil {
		log.Error(err)
	}

	log.Tracef("lights: %+v", lights)

	for _, g := range lights {
		flatName := strings.ToLower(strings.ReplaceAll(g.Name, " ", "_"))

		if lightName == flatName {
			return &g, nil
		}

	}

	return &huego.Light{}, fmt.Errorf("light %s not found", lightName)
}

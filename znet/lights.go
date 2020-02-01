package znet

import (
	"fmt"
	"strings"

	"github.com/amimof/huego"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/rftoy/rftoy"
)

// Lights holds the information necessary to communicate with lighting equipment.
type Lights struct {
	RFToy  *rftoy.RFToy
	HUE    *huego.Bridge
	config LightsConfig
}

// NewLights creates and returns a new Lights object based on the received configuration.
func NewLights(config LightsConfig) *Lights {
	return &Lights{
		HUE:    huego.New(config.Hue.Endpoint, config.Hue.User),
		RFToy:  &rftoy.RFToy{Address: config.RFToy.Endpoint},
		config: config,
	}
}

// GetLight calls the Hue bridge and looks for a light, that, when normalized, matches the name received.
func (l *Lights) GetLight(lightName string) (*huego.Light, error) {
	lights, err := l.HUE.GetLights()
	if err != nil {
		log.Error(err)
	}

	for _, g := range lights {
		flatName := strings.ToLower(strings.ReplaceAll(g.Name, " ", "_"))

		if lightName == flatName {
			return &g, nil
		}

	}

	return &huego.Light{}, fmt.Errorf("Light %s not found", lightName)
}

// GetGroup calls the Hue bridge and looks for a group, that, when normalized, matches the name received.
func (l *Lights) GetGroup(groupName string) (*huego.Group, error) {
	groups, err := l.HUE.GetGroups()
	if err != nil {
		log.Error(err)
	}

	for _, g := range groups {
		flatName := strings.ToLower(strings.ReplaceAll(g.Name, " ", "_"))

		if groupName == flatName {
			return &g, nil
		}
	}

	return &huego.Group{}, fmt.Errorf("Group %s not found", groupName)
}

// On turns off the Hue lights for a room.
func (l *Lights) On(groupName string) {
	room, err := l.config.Room(groupName)
	if err != nil {
		log.Error(err)
	}

	g, err := l.GetGroup(groupName)
	if err != nil {
		log.Error(err)

		light, err := l.GetLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.Debugf("Turning on light %s", light.Name)
			light.On()
		}

	} else {
		log.Debugf("Turning on light group %s", g.Name)
		g.On()
	}

	if len(room.IDs) > 0 {
		log.Debugf("Turning on rftoy lights: %+v", room.IDs)
		for _, i := range room.IDs {
			err := l.RFToy.On(i)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

// Off turns off the Hue lights for a room.
func (l *Lights) Off(groupName string) {
	room, err := l.config.Room(groupName)
	if err != nil {
		log.Error(err)
	}

	g, err := l.GetGroup(groupName)
	if err != nil {
		log.Error(err)

		light, err := l.GetLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.Debugf("Turning off light %s", light.Name)
			light.Off()
		}

	} else {
		log.Debugf("Turning off light group %s", g.Name)
		g.Off()
	}

	if len(room.IDs) > 0 {
		log.Debugf("Turning off rftoy lights: %+v", room.IDs)
		for _, i := range room.IDs {
			err := l.RFToy.Off(i)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

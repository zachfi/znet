package lights

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amimof/huego"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/timer"
)

// Lights holds the information necessary to communicate with lighting
// equipment, and the configuration to add a bit of context.
type Lights struct {
	RFToy  *rftoy.RFToy
	HUE    *huego.Bridge
	config Config
}

// NewLights creates and returns a new Lights object based on the received
// configuration.
func NewLights(config Config) *Lights {
	return &Lights{
		HUE:    huego.New(config.Hue.Endpoint, config.Hue.User),
		RFToy:  &rftoy.RFToy{Address: config.RFToy.Endpoint},
		config: config,
	}
}

// Subscriptions returns the data for mapping event names with functions.
func (l *Lights) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	s.Subscribe("PreSunset", l.eventHandler)
	s.Subscribe("Sunset", l.eventHandler)
	s.Subscribe("Sunrise", l.eventHandler)
	s.Subscribe("TimerExpired", l.eventHandler)
	s.Subscribe("NamedTimer", l.eventHandler)

	return s.Table
}

func (l *Lights) eventHandler(payload events.Payload) error {
	log.Debugf("Lights.eventHandler: %+v", string(payload))

	var e timer.NamedTimer

	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", e, err)
	}

	for _, room := range l.config.Rooms {
		for _, o := range room.On {
			if o == e.Name {
				l.On(room.Name)
			}
		}

		for _, o := range room.Off {
			if o == e.Name {
				l.Off(room.Name)
			}
		}

		for _, o := range room.Dim {
			if o == e.Name {
				l.Dim(room.Name, 100)
			}
		}

		for _, o := range room.Alert {
			if o == e.Name {
				l.Alert(room.Name, "alert")
			}
		}

	}

	return nil
}

// getLight calls the Hue bridge and looks for a light, that, when normalized,
// matches the name received.
func (l *Lights) getLight(lightName string) (*huego.Light, error) {
	lights, err := l.HUE.GetLights()
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

// GetGroup calls the Hue bridge and looks for a group, that, when normalized,
// matches the name received.
func (l *Lights) getGroup(groupName string) (*huego.Group, error) {
	groups, err := l.HUE.GetGroups()
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

// On turns off the Hue lights for a room.
func (l *Lights) On(groupName string) {
	room, err := l.config.Room(groupName)
	if err != nil {
		log.Error(err)
	}

	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)

		light, err := l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.Debugf("turning on light %s", light.Name)
			light.On()
		}

	} else {
		log.Debugf("turning on light group %s", g.Name)
		g.On()
	}

	if len(room.IDs) > 0 {
		log.Debugf("turning on rftoy lights: %+v", room.IDs)
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

	// try the light by group first
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)

		// then try to get just the light
		light, err := l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.Debugf("turning off light %s", light.Name)
			light.Off()
		}

	} else {
		log.Debugf("turning off light group %s", g.Name)
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

func (l *Lights) Dim(groupName string, brightness int32) {
	room, err := l.config.Room(groupName)
	if err != nil {
		log.Error(err)
	}

	groups, err := l.HUE.GetGroups()
	if err != nil {
		log.Error(err)
	}

	for _, g := range groups {
		for _, i := range room.HueIDs {
			if g.ID == i {
				log.Debugf("Setting brightness for group %s: %+v", g.Name, g.State)
				err := g.Bri(uint8(brightness))
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

func (l *Lights) Alert(groupName string, alertName string) {
	g, err := l.getGroup(groupName)
	if err != nil {
		log.Error(err)

		// then try to get just the light
		light, err := l.getLight(groupName)
		if err != nil {
			log.Error(err)
		} else {
			log.Debugf("alert light %s", light.Name)

			err := light.Alert("select")
			if err != nil {
				log.Error(err)
			}
		}

	} else {
		log.Debugf("alert light group %s", g.Name)
		err := g.Alert("select")
		if err != nil {
			log.Error(err)
		}
	}
}

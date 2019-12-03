package znet

import (
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

// On turns off the Hue lights for a room.
func (l *Lights) On(roomName string) {
	room, err := l.config.Room(roomName)
	if err != nil {
		log.Error(err)
		return
	}

	groups, err := l.HUE.GetGroups()
	if err != nil {
		log.Error(err)
	}
	log.Warnf("Groups: %+v", groups)

	for _, g := range groups {
		for _, i := range room.HueIDs {
			if g.ID == i {
				log.Debugf("Turning on %d lights in HUE group %s", len(g.Lights), g.Name)
				err := g.On()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	log.Debugf("Turning on rftoy lights: %+v", room.IDs)
	for _, i := range room.IDs {
		err := l.RFToy.On(i)
		if err != nil {
			log.Error(err)
		}
	}

}

// Off turns off the Hue lights for a room.
func (l *Lights) Off(roomName string) {
	room, err := l.config.Room(roomName)
	if err != nil {
		log.Error(err)
		return
	}

	groups, err := l.HUE.GetGroups()
	if err != nil {
		log.Error(err)
	}

	for _, g := range groups {
		for _, i := range room.HueIDs {
			if g.ID == i {
				log.Debugf("Turning off %d lights in HUE group %s", len(g.Lights), g.Name)
				err := g.Off()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	log.Debugf("Turning off rftoy lights: %+v", room.IDs)
	for _, i := range room.IDs {
		err := l.RFToy.Off(i)
		if err != nil {
			log.Error(err)
		}
	}

}

func (l *Lights) Status() []huego.Light {
	lights, err := l.HUE.GetLights()
	if err != nil {
		log.Error(err)
	}

	return lights
}

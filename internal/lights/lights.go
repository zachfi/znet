package lights

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/amimof/huego"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mpvl/unique"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/rftoy/rftoy"

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/pkg/iot"
	"github.com/xaque208/znet/rpc"
)

// Lights holds the information necessary to communicate with lighting
// equipment, and the configuration to add a bit of context.
type Lights struct {
	config   Config
	Handlers []Light
}

// NewLights creates and returns a new Lights object based on the received
// configuration.
func NewLights(config Config, inventoryClient rpc.InventoryClient, mqttClient mqtt.Client) *Lights {
	l := Lights{
		config: config,
	}

	hue := hueLight{
		config: config,
		hue:    huego.New(config.Hue.Endpoint, config.Hue.User),
	}

	zigbee := zigbeeLight{
		config:          config,
		inventoryClient: inventoryClient,
		mqttClient:      mqttClient,
	}

	rftoy := rftoyLight{
		endpoint: &rftoy.RFToy{Address: config.RFToy.Endpoint},
	}

	l.Handlers = []Light{
		hue,
		zigbee,
		rftoy,
	}

	return &l
}

// Subscriptions returns the data for mapping event names with functions.
func (l *Lights) Subscriptions() *events.Subscriptions {
	s := events.NewSubscriptions()

	eventNames := []string{
		"NamedTimer",
		"Click",
	}

	for _, e := range eventNames {
		switch e {
		case "SolarEvent":
			s.Subscribe(e, l.solarEventHandler)
		case "NamedTimer":
			s.Subscribe(e, l.namedTimerHandler)

			// f := &timer.EventFilter{}
			// f.Name = append(f.Name, iot.EventNames...)
			// s.Filter(e, f)
		case "Click":
			s.Subscribe(e, l.clickHandler)
		}
	}

	return s
}

func (l *Lights) solarEventHandler(name string, payload events.Payload) error {
	var e astro.SolarEvent

	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", e, err)
	}

	names := l.configuredEventNames()

	configuredEvent := func(name string, names []string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}(e.Name, names)

	if !configuredEvent {
		log.WithFields(log.Fields{
			"name":            e.Name,
			"configuredNames": names,
		}).Debug("unhandled lighting SolarEvent name")

		return nil
	}

	l.setRoomForEvent(e.Name)

	return nil
}

func (l *Lights) clickHandler(name string, payload events.Payload) error {
	log.Tracef("Lights.clickHandler: %s : %+v", name, string(payload))

	var e iot.Click

	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", e, err)
	}

	for _, room := range l.config.Rooms {
		alert := false
		dim := false
		off := false
		on := false
		toggle := false

		if room.Name == e.Zone {
			log.WithFields(log.Fields{
				"room_name": room.Name,
				"name":      e.Zone,
			}).Trace("using room")

			switch e.Count {
			case "single":
				toggle = true
			case "double":
				on = true
			case "triple":
				off = true
			case "long":
				dim = true
			case "many":
				alert = true
			default:
				log.Warnf("unknown click event: %s", e)
			}

			for _, h := range l.Handlers {
				if toggle {
					err := h.Toggle(room.Name)
					if err != nil {
						log.Error(err)
					}
				}

				if off {
					err := h.Off(room.Name)
					if err != nil {
						log.Error(err)
					}
				}

				if on {
					err := h.On(room.Name)
					if err != nil {
						log.Error(err)
					}
				}

				if alert {
					err := h.Alert(room.Name)
					if err != nil {
						log.Error(err)
					}
				}

				if dim {
					err := h.Dim(room.Name, 100)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

	return nil
}

// configuredEventNames is the collection of events that are configured in the lighting config.  This results in a
func (l *Lights) configuredEventNames() []string {

	names := []string{}

	for _, r := range l.config.Rooms {
		names = append(names, r.On...)
		names = append(names, r.Off...)
		names = append(names, r.Alert...)
		names = append(names, r.Toggle...)
		names = append(names, r.Dim...)
	}

	sort.Strings(names)
	unique.Strings(&names)

	return names
}

func (l *Lights) namedTimerHandler(name string, payload events.Payload) error {
	var e timer.NamedTimer

	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", e, err)
	}

	names := l.configuredEventNames()

	configuredEvent := func(name string, names []string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}(e.Name, names)

	if !configuredEvent {
		log.WithFields(log.Fields{
			"name":            e.Name,
			"configuredNames": names,
		}).Debug("unhandled lighting NamedTimer name")

		return nil
	}

	l.setRoomForEvent(e.Name)

	return nil
}

func (l *Lights) setRoomForEvent(name string) {
	for _, room := range l.config.Rooms {
		for _, o := range room.On {
			if o == name {
				for _, h := range l.Handlers {
					h.On(room.Name)
				}
			}
		}

		for _, o := range room.Off {
			if o == name {
				for _, h := range l.Handlers {
					h.Off(room.Name)
				}
			}
		}

		for _, o := range room.Dim {
			if o == name {
				for _, h := range l.Handlers {
					h.Dim(room.Name, 100)
				}
			}
		}

		for _, o := range room.Alert {
			if o == name {
				for _, h := range l.Handlers {
					h.Alert(room.Name)
				}
			}
		}
	}
}

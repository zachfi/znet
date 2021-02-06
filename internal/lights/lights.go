package lights

import (
	context "context"
	"fmt"
	"sort"
	sync "sync"

	"github.com/mpvl/unique"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/iot"
)

// Lights holds the information necessary to communicate with lighting
// equipment, and the configuration to add a bit of context.
type Lights struct {
	sync.Mutex
	config   *config.LightsConfig
	handlers []Handler
}

// NewLights creates and returns a new Lights object based on the received
// configuration.
// func NewLights(cfg *config.LightsConfig, inventoryClient rpc.InventoryClient, mqttClient mqtt.Client) *Lights {
func NewLights(cfg *config.Config) (*Lights, error) {
	if cfg == nil {
		return nil, fmt.Errorf("unable to create new Lights from nil config")
	}

	return &Lights{
		config: cfg.Lights,
	}, nil
}

// AddHandler is the way to register a new Handler with the Lights server.
func (l *Lights) AddHandler(h Handler) {
	l.Lock()
	defer l.Unlock()

	l.handlers = append(l.handlers, h)
}

func (l *Lights) ActionHandler(action *iot.Action) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := l.getRoom(action.Zone)
	if room == nil {
		return fmt.Errorf("no room named %s was found in config", action.Zone)
	}

	log.WithFields(log.Fields{
		"room_name": room.Name,
		"zone":      action.Zone,
		"device":    action.Device,
		"event":     action.Event,
	}).Trace("room action")

	request := &LightGroupRequest{
		Brightness: 254,
		Color:      "#ffffff",
		Colors:     l.config.PartyColors,
		Name:       room.Name,
	}

	switch action.Event {
	case "single":
		_, err := l.Toggle(ctx, request)
		return err
	case "double":
		_, err := l.On(ctx, request)
		if err != nil {
			return err
		}

		_, err = l.Dim(ctx, request)
		if err != nil {
			return err
		}

		_, err = l.SetColor(ctx, request)
		return err
	case "triple":
		_, err := l.Off(ctx, request)
		return err
	case "quadruple":
		_, err := l.RandomColor(ctx, request)
		return err
	case "hold", "release":
		request.Brightness = 110
		_, err := l.Dim(ctx, request)
		return err
	case "many":
		_, err := l.Alert(ctx, request)
		return err
	default:
		log.Debugf("unknown action event: %s", action.Event)
	}

	return nil
}

func (l *Lights) getRoom(name string) *config.LightsRoom {
	if l.config == nil {
		return nil
	}

	for _, room := range l.config.Rooms {
		if room.Name == name {
			return &room
		}
	}

	return nil
}

// configuredEventNames is a collection of events that are configured in the
// lighting config.  These event names determin all th epossible event names
// that will be responded to.
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

func (l *Lights) NamedTimerHandler(e string) error {
	names := l.configuredEventNames()

	configuredEvent := func(name string, names []string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}(e, names)

	if !configuredEvent {
		log.WithFields(log.Fields{
			"name":            e,
			"configuredNames": names,
		}).Debug("unhandled lighting NamedTimer name")

		return nil
	}

	l.SetRoomForEvent(e)

	return nil
}

func (l *Lights) SetRoomForEvent(name string) {
	for _, room := range l.config.Rooms {
		for _, o := range room.On {
			if o == name {
				for _, h := range l.handlers {
					err := h.On(room.Name)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}

		for _, o := range room.Off {
			if o == name {
				for _, h := range l.handlers {
					err := h.Off(room.Name)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}

		for _, o := range room.Dim {
			if o == name {
				for _, h := range l.handlers {
					err := h.Dim(room.Name, 110)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}

		for _, o := range room.Alert {
			if o == name {
				for _, h := range l.handlers {
					err := h.Alert(room.Name)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

// Alert calls Alert() on each handler.
func (l *Lights) Alert(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	for _, h := range l.handlers {
		err := h.Alert(req.Name)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

// Dim calls Dim() on each handler.
func (l *Lights) Dim(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	for _, h := range l.handlers {
		err := h.Dim(req.Name, req.Brightness)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

// Off calls Off() on each handler.
func (l *Lights) Off(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	for _, h := range l.handlers {
		err := h.Off(req.Name)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

// On calls On() on each handler.
func (l *Lights) On(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	for _, h := range l.handlers {
		err := h.On(req.Name)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

func (l *Lights) RandomColor(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	if len(req.Colors) == 0 {
		return nil, fmt.Errorf("request contained no colors to select from")
	}
	for _, h := range l.handlers {
		err := h.RandomColor(req.Name, req.Colors)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

func (l *Lights) SetColor(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	if req.Color == "" {
		return nil, fmt.Errorf("request missing color spec")
	}

	for _, h := range l.handlers {
		err := h.SetColor(req.Name, req.Color)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

func (l *Lights) Toggle(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	for _, h := range l.handlers {
		err := h.Toggle(req.Name)
		if err != nil {
			log.Error(err)
		}
	}

	return &LightResponse{}, nil
}

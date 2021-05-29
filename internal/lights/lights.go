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
	UnimplementedLightsServer
	sync.Mutex
	config   *config.LightsConfig
	handlers []Handler
	zones    *Zones
}

var defaultColorPool = []string{"#006c7f", "#e32636", "#b0bf1a"}

// NewLights creates and returns a new Lights object based on the received
// configuration.
func NewLights(cfg *config.LightsConfig) (*Lights, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &Lights{
		config: cfg,
		zones:  &Zones{},
	}, nil
}

// AddHandler is used to register the received Handler.
func (l *Lights) AddHandler(h Handler) {
	l.Lock()
	defer l.Unlock()

	l.handlers = append(l.handlers, h)
}

// ActionHandler is called when an action is requested against a light group.
// The action speciefies the a button press and a room to give enough context
// for how to change the behavior of the lights in response to the action.
func (l *Lights) ActionHandler(action *iot.Action) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	room := l.getRoom(action.Zone)
	if room == nil {
		return ErrRoomNotFound
	}

	log.WithFields(log.Fields{
		"room_name": room.Name,
		"zone":      action.Zone,
		"device":    action.Device,
		"event":     action.Event,
	}).Debug("room action")

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
		return fmt.Errorf("%s: %w", action.Event, ErrUnknownActionEvent)
	}
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
// lighting config.  These event names determin all the epossible event names
// that will be responded to.
func (l *Lights) configuredEventNames() ([]string, error) {
	names := []string{}

	if l.config == nil {
		return nil, ErrNilConfig
	}

	if l.config.Rooms == nil || len(l.config.Rooms) == 0 {
		return nil, ErrNoRoomsConfigured
	}

	for _, r := range l.config.Rooms {
		names = append(names, r.On...)
		names = append(names, r.Off...)
		names = append(names, r.Alert...)
		names = append(names, r.Toggle...)
		names = append(names, r.Dim...)
	}

	sort.Strings(names)
	unique.Strings(&names)

	return names, nil
}

func (l *Lights) NamedTimerHandler(ctx context.Context, e string) error {
	names, err := l.configuredEventNames()
	if err != nil {
		return err
	}

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

	return l.SetRoomForEvent(ctx, e)
}

func (l *Lights) SetRoomForEvent(ctx context.Context, event string) error {
	for _, room := range l.config.Rooms {
		for _, o := range room.On {
			if o == event {
				req := &LightGroupRequest{Name: room.Name}
				_, err := l.On(ctx, req)
				return err
			}
		}

		for _, o := range room.Off {
			if o == event {
				req := &LightGroupRequest{Name: room.Name}
				_, err := l.Off(ctx, req)
				return err
			}
		}

		for _, o := range room.Dim {
			if o == event {
				req := &LightGroupRequest{Name: room.Name, Brightness: 110}
				_, err := l.Dim(ctx, req)
				return err
			}
		}

		for _, o := range room.Alert {
			if o == event {
				req := &LightGroupRequest{Name: room.Name}
				_, err := l.Alert(ctx, req)
				return err
			}
		}
	}

	return nil
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
	z := l.zones.GetZone(req.Name)

	err := z.Dim(req.Brightness)
	if err != nil {
		return nil, err
	}

	err = z.Handle(req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

// Off calls Off() on each handler.
func (l *Lights) Off(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	z := l.zones.GetZone(req.Name)

	err := z.Off()
	if err != nil {
		return nil, err
	}

	err = z.Handle(req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

// On calls On() on each handler.
func (l *Lights) On(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	z := l.zones.GetZone(req.Name)

	err := z.On()
	if err != nil {
		return nil, err
	}

	err = z.Handle(req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

func (l *Lights) RandomColor(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	var colors []string

	if len(req.Colors) == 0 {
		log.Debug("using default colors")
		colors = defaultColorPool
	} else {
		colors = req.Colors
	}

	z := l.zones.GetZone(req.Name)
	err := z.RandomColor(colors)
	if err != nil {
		return nil, err
	}

	err = z.Handle(req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

func (l *Lights) SetColor(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	if req.Color == "" {
		return nil, fmt.Errorf("request missing color spec")
	}

	z := l.zones.GetZone(req.Name)
	err := z.SetColor(req.Color)
	if err != nil {
		return nil, err
	}

	err = z.Handle(req.Name, l.handlers...)
	if err != nil {
		return nil, err
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

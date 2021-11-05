package lights

import (
	"context"
	"fmt"
	"sort"
	sync "sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/mpvl/unique"
	"github.com/opentracing/opentracing-go"

	"github.com/xaque208/znet/pkg/iot"
)

const (
	brightnessLow  = 100
	brightnessHigh = 254
)

// Lights holds the information necessary to communicate with lighting
// equipment, and the configuration to add a bit of context.
type Lights struct {
	UnimplementedLightsServer

	services.Service
	cfg *Config

	logger log.Logger

	sync.Mutex
	handlers []Handler
	zones    *Zones
}

var defaultColorPool = []string{"#006c7f", "#e32636", "#b0bf1a"}

// NewLights creates and returns a new Lights object based on the received
// configuration.
func New(cfg Config, logger log.Logger) (*Lights, error) {
	l := &Lights{
		cfg:    &cfg,
		logger: log.With(logger, "module", "lights"),
		zones:  &Zones{},
	}

	l.Service = services.NewBasicService(l.starting, l.running, l.stopping)

	return l, nil
}

func (l *Lights) starting(ctx context.Context) error {
	return nil
}

func (l *Lights) running(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (l *Lights) stopping(_ error) error {
	return nil
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
func (l *Lights) ActionHandler(ctx context.Context, action *iot.Action) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Lights.ActionHandler")
	defer span.Finish()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	room := l.getRoom(action.Zone)
	if room == nil {
		return fmt.Errorf("%w: %s", ErrRoomNotFound, action.Zone)
	}

	level.Debug(l.logger).Log("msg", "room action",
		"room_name", room.Name,
		"zone", action.Zone,
		"device", action.Device,
		"event", action.Event,
	)

	request := &LightGroupRequest{
		Brightness: brightnessHigh,
		Color:      "#ffffff",
		Colors:     l.cfg.PartyColors,
		Name:       room.Name,
	}

	switch action.Event {
	case "single", "press":
		_, err := l.Toggle(ctx, request)
		return err
	case "on", "double", "tap", "rotate_right", "slide":
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
	case "off", "triple":
		_, err := l.Off(ctx, request)
		return err
	case "quadruple", "flip90", "flip180", "fall":
		_, err := l.RandomColor(ctx, request)
		return err
	case "hold", "release", "rotate_left":
		request.Brightness = brightnessLow
		_, err := l.Dim(ctx, request)
		return err
	case "many":
		_, err := l.Alert(ctx, request)
		return err
	case "wakeup": // do nothing
		return nil
	default:
		return fmt.Errorf("%s: %w", action.Event, ErrUnknownActionEvent)
	}
}

func (l *Lights) getRoom(name string) *LightsRoom {
	if l.cfg == nil {
		return nil
	}

	for _, room := range l.cfg.Rooms {
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

	if l.cfg == nil {
		return nil, ErrNilConfig
	}

	if l.cfg.Rooms == nil || len(l.cfg.Rooms) == 0 {
		return nil, ErrNoRoomsConfigured
	}

	for _, r := range l.cfg.Rooms {
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
		return ErrUnhandledEventName
	}

	return l.SetRoomForEvent(ctx, e)
}

func (l *Lights) SetRoomForEvent(ctx context.Context, event string) error {
	for _, room := range l.cfg.Rooms {
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
		err := h.Alert(ctx, req.Name)
		if err != nil {
			level.Error(l.logger).Log("err", err.Error())
		}
	}

	return &LightResponse{}, nil
}

// Dim calls Dim() on each handler.
func (l *Lights) Dim(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	z := l.zones.GetZone(req.Name)

	err := z.Dim(ctx, req.Brightness)
	if err != nil {
		return nil, err
	}

	err = z.Handle(ctx, req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

// Off calls Off() on each handler.
func (l *Lights) Off(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	z := l.zones.GetZone(req.Name)

	err := z.Off(ctx)
	if err != nil {
		return nil, err
	}

	err = z.Handle(ctx, req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

// On calls On() on each handler.
func (l *Lights) On(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	z := l.zones.GetZone(req.Name)

	err := z.On(ctx)
	if err != nil {
		return nil, err
	}

	err = z.Handle(ctx, req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

func (l *Lights) RandomColor(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	var colors []string

	if len(req.Colors) == 0 {
		level.Debug(l.logger).Log("msg", "using default colors")
		colors = defaultColorPool
	} else {
		colors = req.Colors
	}

	z := l.zones.GetZone(req.Name)
	err := z.RandomColor(ctx, colors)
	if err != nil {
		return nil, err
	}

	err = z.Handle(ctx, req.Name, l.handlers...)
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
	err := z.SetColor(ctx, req.Color)
	if err != nil {
		return nil, err
	}

	err = z.Handle(ctx, req.Name, l.handlers...)
	if err != nil {
		return nil, err
	}

	return &LightResponse{}, nil
}

func (l *Lights) Toggle(ctx context.Context, req *LightGroupRequest) (*LightResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Lights.Toggle")
	defer span.Finish()

	for _, h := range l.handlers {
		handlerSpan, handlerCtx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%T.Toggle()", h))
		err := h.Toggle(handlerCtx, req.Name)
		if err != nil {
			level.Error(l.logger).Log("err", err.Error())
		}
		handlerSpan.Finish()
	}

	return &LightResponse{}, nil
}

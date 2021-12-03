package lights

import (
	"context"
	"sort"
	"strings"
	sync "sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/mpvl/unique"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/xaque208/znet/pkg/iot"
)

const (
	brightnessLow  = 100
	brightnessHigh = 254

	nightTemp     = 500
	eveningTemp   = 400
	lateafternoon = 300
	day           = 200
	morningTemp   = 100
)

// Lights holds the information necessary to communicate with lighting
// equipment, and the configuration to add a bit of context.
type Lights struct {
	UnimplementedLightsServer

	services.Service
	cfg *Config

	logger log.Logger

	sync.Mutex
	handlers           []Handler
	colorTempScheduler ColorTempSchedulerFunc
	zones              *Zones
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

func (l *Lights) SetColorTempScheduler(c ColorTempSchedulerFunc) {
	l.Lock()
	defer l.Unlock()

	l.colorTempScheduler = c
}

// ActionHandler is called when an action is requested against a light group.
// The action speciefies the a button press and a room to give enough context
// for how to change the behavior of the lights in response to the action.
func (l *Lights) ActionHandler(ctx context.Context, action *iot.Action) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Lights.ActionHandler")
	defer span.Finish()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	z := l.zones.GetZone(action.Zone)
	z.SetHandlers(l.handlers...)

	_ = level.Debug(l.logger).Log("msg", "room action",
		"name", z.Name(),
		"zone", action.Zone,
		"device", action.Device,
		"event", action.Event,
	)

	switch action.Event {
	case "single", "press":
		return z.Toggle(ctx)
	case "on", "double", "tap", "rotate_right", "slide":
		if err := z.On(ctx); err != nil {
			return err
		}

		return z.Dim(ctx, brightnessHigh)
	case "off", "triple":
		return z.Off(ctx)
	case "quadruple", "flip90", "flip180", "fall":
		return z.RandomColor(ctx, z.colorPool)
	case "hold", "release", "rotate_left":
		return z.Dim(ctx, brightnessLow)
	case "many":
		return z.Alert(ctx)
	case "wakeup": // do nothing
		return nil
	default:
		return errors.Wrap(ErrUnknownActionEvent, action.Event)
	}
}

func (l *Lights) getRoom(name string) *Room {
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
// lighting config.  These event names determin all the possible event names
// that will be responded to.
func (l *Lights) configuredEventNames() ([]string, error) {
	names := []string{}

	if l.cfg == nil {
		return nil, ErrNilConfig
	}

	if l.cfg.Rooms == nil || len(l.cfg.Rooms) == 0 {
		return nil, ErrNoRoomsConfigured
	}

	for _, z := range l.cfg.Rooms {
		for _, s := range z.States {
			names = append(names, s.Event)
		}
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
	for _, zone := range l.cfg.Rooms {
		z := l.zones.GetZone(zone.Name)
		z.SetHandlers(l.handlers...)

		for _, s := range zone.States {
			if !strings.EqualFold(event, s.Event) {
				continue
			}

			return z.SetState(ctx, s.State)
		}
	}

	return nil
}

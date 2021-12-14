package lights

import (
	"context"
	"sort"
	"strings"
	sync "sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/mpvl/unique"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"

	"github.com/xaque208/znet/pkg/iot"
)

const (
	eveningTemp       = 500
	lateafternoonTemp = 400
	dayTemp           = 300
	morningTemp       = 200
	firstlightTemp    = 100

	nightVisionColor = `#FF00FF`
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
var defaultColorTemperatureMap = map[ColorTemperature]int32{
	ColorTemperature_FIRSTLIGHT:    firstlightTemp,
	ColorTemperature_MORNING:       morningTemp,
	ColorTemperature_DAY:           dayTemp,
	ColorTemperature_LATEAFTERNOON: lateafternoonTemp,
	ColorTemperature_EVENING:       eveningTemp,
}
var defaultBrightnessMap = map[Brightness]int32{
	Brightness_FULL: 254,
	Brightness_DIM:  100,
	Brightness_LOW:  90,
}
var defaultScheduleDuration = time.Minute * 10

// NewLights creates and returns a new Lights object based on the received
// configuration.
func New(cfg Config, logger log.Logger) (*Lights, error) {
	l := &Lights{
		cfg:    &cfg,
		logger: log.With(logger, "module", "lights"),
		zones:  &Zones{},
	}

	if len(l.cfg.PartyColors) == 0 {
		l.cfg.PartyColors = defaultColorPool
	}

	l.Service = services.NewBasicService(l.starting, l.running, l.stopping)

	return l, nil
}

func (l *Lights) starting(ctx context.Context) error {
	return nil
}

func (l *Lights) running(ctx context.Context) error {
	l.runColorTempScheduler(ctx)
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

func (l *Lights) runColorTempScheduler(ctx context.Context) {
	ticker := time.NewTicker(defaultScheduleDuration)

	go func(ctx context.Context) {
		zones := l.zones.GetZones()
		for _, room := range l.cfg.Rooms {
			z := l.zones.GetZone(room.Name)
			z.SetHandlers(l.handlers...)
		}
		update := func() {
			for _, z := range zones {
				temp := l.colorTempScheduler().MostRecent()
				_ = z.SetColorTemperature(ctx, temp)
			}
		}

		update()

		for {
			select {
			case <-ticker.C:
				update()
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

// ActionHandler is called when an action is requested against a light group.
// The action speciefies the a button press and a room to give enough context
// for how to change the behavior of the lights in response to the action.
func (l *Lights) ActionHandler(ctx context.Context, action *iot.Action) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	defer span.End()

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
		err := z.SetBrightness(ctx, Brightness_FULL)
		if err != nil {
			return err
		}

		return z.On(ctx)
		// if err := z.On(ctx); err != nil {
		// 	return err
		// }

		// return nil
	case "off", "triple":
		return z.Off(ctx)
	case "quadruple", "flip90", "flip180", "fall":
		return z.RandomColor(ctx, l.cfg.PartyColors)
	case "hold", "rotate_left":
		err := z.SetBrightness(ctx, Brightness_DIM)
		if err != nil {
			return err
		}

		return z.On(ctx)
	case "many":
		return z.Alert(ctx)
	case "wakeup", "release": // do nothing
		return nil
	default:
		return errors.Wrap(ErrUnknownActionEvent, action.Event)
	}
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
	span := trace.SpanFromContext(ctx)
	defer span.End()

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

// SetRoomForEvent is used to handle an event based on the room configuation.
func (l *Lights) SetRoomForEvent(ctx context.Context, event string) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	for _, zone := range l.cfg.Rooms {
		z := l.zones.GetZone(zone.Name)
		z.SetHandlers(l.handlers...)

		for _, s := range zone.States {
			if !strings.EqualFold(event, s.Event) {
				continue
			}

			if s.Brightness != nil {
				err := z.SetBrightness(ctx, *s.Brightness)
				if err != nil {
					return err
				}
			}

			if s.ColorTemp != nil {
				err := z.SetColorTemperature(ctx, *s.ColorTemp)
				if err != nil {
					return err
				}
			}

			return z.SetState(ctx, s.State)
		}
	}

	return nil
}

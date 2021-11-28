package lights

import (
	"context"
	"fmt"
	sync "sync"

	"github.com/opentracing/opentracing-go"
)

func NewZone(name string) *Zone {
	z := &Zone{}
	z.lock = new(sync.Mutex)
	z.SetName(name)

	return z
}

type Zone struct {
	lock *sync.Mutex

	name string

	brightness int32
	colorPool  []string
	color      string
	handlers   []Handler
	state      ZoneState
}

func (z *Zone) SetName(name string) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.name = name
}

func (z *Zone) Name() string {
	return z.name
}

func (z *Zone) SetHandlers(handlers ...Handler) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.handlers = handlers
}

func (z *Zone) Dim(ctx context.Context, brightness int32) error {
	z.brightness = brightness

	return z.SetState(ctx, ZoneState_DIM)
}

func (z *Zone) Off(ctx context.Context) error {
	return z.SetState(ctx, ZoneState_OFF)
}

func (z *Zone) On(ctx context.Context) error {
	return z.SetState(ctx, ZoneState_ON)
}

func (z *Zone) SetColor(ctx context.Context, color string) error {
	z.color = color
	return z.SetState(ctx, ZoneState_COLOR)
}

func (z *Zone) RandomColor(ctx context.Context, colors []string) error {
	z.colorPool = colors
	return z.SetState(ctx, ZoneState_RANDOMCOLOR)
}

func (z *Zone) SetState(ctx context.Context, state ZoneState) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zone.SetState")
	defer span.Finish()

	z.lock.Lock()
	defer z.lock.Unlock()

	z.state = state

	return z.handle(ctx)
}

func (z *Zone) handle(ctx context.Context) error {
	if z.name == "" {
		return fmt.Errorf("unable to handle unnamed zone")
	}

	if len(z.handlers) == 0 {
		return fmt.Errorf("no handlers for zone")
	}

	return z.Handle(ctx, z.name, z.handlers...)
}

func (z *Zone) Handle(ctx context.Context, name string, handlers ...Handler) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "zone.Handle")
	defer span.Finish()

	switch z.state {
	case ZoneState_ON:
		return handleOn(ctx, name, handlers...)
	case ZoneState_OFF:
		return handleOff(ctx, name, handlers...)
	case ZoneState_COLOR:
		for _, h := range handlers {
			err := h.SetColor(ctx, name, z.color)
			if err != nil {
				return fmt.Errorf("%s color: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_RANDOMCOLOR:
		for _, h := range handlers {
			err := h.RandomColor(ctx, name, z.colorPool)
			if err != nil {
				return fmt.Errorf("%s random color: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_DIM:
		for _, h := range handlers {
			err := h.Dim(ctx, name, z.brightness)
			if err != nil {
				return fmt.Errorf("%s dim: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_NIGHTVISION:
		for _, h := range handlers {
			err := h.SetTemp(ctx, name, nightTemp)
			if err != nil {
				return fmt.Errorf("%s night: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_EVENINGVISION:
		for _, h := range handlers {
			err := h.SetTemp(ctx, name, eveningTemp)
			if err != nil {
				return fmt.Errorf("%s evening: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_MORNINGVISION:
		for _, h := range handlers {
			err := h.SetTemp(ctx, name, morningTemp)
			if err != nil {
				return fmt.Errorf("%s morning: %w", name, ErrHandlerFailed)
			}
		}
	}

	return nil
}

type Zones struct {
	lock   *sync.Mutex
	state  map[string]*Zone
	states []*Zone
}

func (z *Zones) GetZone(name string) *Zone {
	if z.lock == nil {
		z.lock = new(sync.Mutex)
	}

	for _, zone := range z.states {
		if zone.Name() == name {
			return zone
		}
	}

	if zone, ok := z.state[name]; ok {
		return zone
	}

	if len(z.states) == 0 {
		z.states = make([]*Zone, 0)
	}

	z.lock.Lock()
	defer z.lock.Unlock()

	zone := NewZone(name)
	z.states = append(z.states, zone)
	return zone
}

func handleOn(ctx context.Context, name string, handlers ...Handler) error {
	for _, h := range handlers {
		err := h.On(ctx, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleOff(ctx context.Context, name string, handlers ...Handler) error {
	for _, h := range handlers {
		err := h.Off(ctx, name)
		if err != nil {
			return err
		}
	}

	return nil
}

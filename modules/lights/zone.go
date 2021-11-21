package lights

import (
	"context"
	"fmt"
	sync "sync"

	"github.com/opentracing/opentracing-go"
)

func NewZone() *Zone {
	z := &Zone{}
	z.lock = new(sync.Mutex)

	return z
}

type Zone struct {
	lock       *sync.Mutex
	state      ZoneState
	brightness int32
	color      string
	colorPool  []string
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

	// TODO
	// return z.Handle()

	return nil
}

func (z *Zone) Handle(ctx context.Context, name string, handlers ...Handler) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "zone.Handle")
	defer span.Finish()

	switch z.state {
	case ZoneState_ON:
		for _, h := range handlers {
			err := h.On(ctx, name)
			if err != nil {
				return fmt.Errorf("%s on: %w", name, ErrHandlerFailed)
			}
		}
	case ZoneState_OFF:
		for _, h := range handlers {
			err := h.Off(ctx, name)
			if err != nil {
				return fmt.Errorf("%s on: %w", name, ErrHandlerFailed)
			}
		}
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
	lock  *sync.Mutex
	state map[string]*Zone
}

func (z *Zones) GetZone(name string) *Zone {
	if z.lock == nil {
		z.lock = new(sync.Mutex)
	}

	if zone, ok := z.state[name]; ok {
		// zone.Name = name
		return zone
	}

	if len(z.state) == 0 {
		z.state = make(map[string]*Zone)
	}

	z.lock.Lock()
	defer z.lock.Unlock()

	z.state[name] = NewZone()

	return z.state[name]
}

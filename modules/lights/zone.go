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
	state      zoneState
	brightness int32
	color      string
	colorPool  []string
}

func (z *Zone) Dim(ctx context.Context, brightness int32) error {
	z.brightness = brightness
	return z.SetState(ctx, Dim)
}

func (z *Zone) Off(ctx context.Context) error {
	return z.SetState(ctx, Off)
}

func (z *Zone) On(ctx context.Context) error {
	return z.SetState(ctx, On)
}

func (z *Zone) RandomColor(ctx context.Context, colors []string) error {
	z.colorPool = colors
	return z.SetState(ctx, RandomColor)
}

func (z *Zone) SetColor(ctx context.Context, color string) error {
	z.color = color
	return z.SetState(ctx, Color)
}

func (z *Zone) SetState(ctx context.Context, state zoneState) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zone.SetState")
	defer span.Finish()

	z.lock.Lock()
	defer z.lock.Unlock()

	z.state = state

	return nil
}

func (z *Zone) Handle(ctx context.Context, name string, handlers ...Handler) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "zone.Handle")
	defer span.Finish()

	switch z.state {
	case Off:
		for _, h := range handlers {
			err := h.Off(ctx, name)
			if err != nil {
				return fmt.Errorf("%s off: %w", name, ErrHandlerFailed)
			}
		}
	case On:
		for _, h := range handlers {
			err := h.On(ctx, name)
			if err != nil {
				return fmt.Errorf("%s on: %w", name, ErrHandlerFailed)
			}
		}
	case Color:
		for _, h := range handlers {
			err := h.SetColor(ctx, name, z.color)
			if err != nil {
				return fmt.Errorf("%s color: %w", name, ErrHandlerFailed)
			}
		}
	case RandomColor:
		for _, h := range handlers {
			err := h.RandomColor(ctx, name, z.colorPool)
			if err != nil {
				return fmt.Errorf("%s random color: %w", name, ErrHandlerFailed)
			}
		}
	case Dim:
		for _, h := range handlers {
			err := h.Dim(ctx, name, z.brightness)
			if err != nil {
				return fmt.Errorf("%s dim: %w", name, ErrHandlerFailed)
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

	z.lock.Lock()
	defer z.lock.Unlock()

	if zone, ok := z.state[name]; ok {
		// zone.Name = name
		return zone
	}

	if len(z.state) == 0 {
		z.state = make(map[string]*Zone)
	}

	// z.state[name] = &Zone{Name: name}
	z.state[name] = NewZone()

	return z.state[name]
}

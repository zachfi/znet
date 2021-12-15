package lights

import (
	"context"
	"fmt"
	sync "sync"

	"go.opentelemetry.io/otel/trace"
)

func NewZone(name string, handlers ...Handler) *Zone {
	z := &Zone{}
	z.lock = new(sync.Mutex)
	z.SetName(name)
	// default
	z.colorTemp = dayTemp

	z.brightnessMap = defaultBrightnessMap
	z.colorTempMap = defaultColorTemperatureMap

	return z
}

type Zone struct {
	lock *sync.Mutex

	name string

	brightness int32
	colorPool  []string
	color      string
	colorTemp  int32
	handlers   []Handler
	state      ZoneState

	colorTempMap  map[ColorTemperature]int32
	brightnessMap map[Brightness]int32
}

func (z *Zone) SetName(name string) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.name = name
}

func (z *Zone) SetBrightnessMap(m map[Brightness]int32) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.brightnessMap = m
}

func (z *Zone) SetColorTemperatureMap(m map[ColorTemperature]int32) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.colorTempMap = m
}

func (z *Zone) Name() string {
	return z.name
}

func (z *Zone) SetHandlers(handlers ...Handler) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.handlers = handlers
}

func (z *Zone) SetColorTemperature(ctx context.Context, colorTemp ColorTemperature) error {
	z.colorTemp = z.colorTempMap[colorTemp]
	return nil
}

func (z *Zone) SetBrightness(ctx context.Context, brightness Brightness) error {
	z.brightness = z.brightnessMap[brightness]
	return nil
}

func (z *Zone) Off(ctx context.Context) error {
	return z.SetState(ctx, ZoneState_OFF)
}

func (z *Zone) On(ctx context.Context) error {
	return z.SetState(ctx, ZoneState_ON)
}

func (z *Zone) Toggle(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.Toggle(ctx, z.Name())
		if err != nil {
			return fmt.Errorf("%s toggle: %w", z.name, ErrHandlerFailed)
		}
	}

	return nil
}

func (z *Zone) Alert(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.Alert(ctx, z.Name())
		if err != nil {
			return fmt.Errorf("%s alert: %w", z.name, ErrHandlerFailed)
		}
	}

	return nil
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

	span := trace.SpanFromContext(ctx)
	defer span.End()

	z.lock.Lock()
	defer z.lock.Unlock()

	z.state = state

	return z.flush(ctx)
}

func (z *Zone) flush(ctx context.Context) error {
	if z.name == "" {
		return fmt.Errorf("unable to handle unnamed zone")
	}

	if len(z.handlers) == 0 {
		return fmt.Errorf("no handlers for zone")
	}

	return z.Flush(ctx)
}

// Flush handles pushing the current state out to each of the hnadlers.
func (z *Zone) Flush(ctx context.Context) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	switch z.state {
	case ZoneState_ON:
		return z.handleOn(ctx)
	case ZoneState_OFF:
		return z.handleOff(ctx)
	case ZoneState_COLOR:
		for _, h := range z.handlers {
			err := h.SetColor(ctx, z.name, z.color)
			if err != nil {
				return fmt.Errorf("%s color: %w", z.name, ErrHandlerFailed)
			}
		}
	case ZoneState_RANDOMCOLOR:
		for _, h := range z.handlers {
			err := h.RandomColor(ctx, z.name, z.colorPool)
			if err != nil {
				return fmt.Errorf("%s random color: %w", z.name, ErrHandlerFailed)
			}
		}
	case ZoneState_NIGHTVISION:
		z.color = nightVisionColor
		return z.handleColor(ctx)
	case ZoneState_EVENINGVISION:
		z.colorTemp = eveningTemp
		return z.handleColorTemperature(ctx)
	case ZoneState_MORNINGVISION:
		z.colorTemp = morningTemp
		return z.handleColorTemperature(ctx)
	}

	return nil
}

type Zones struct {
	lock   *sync.Mutex
	states []*Zone
}

func (z *Zones) GetZones() []*Zone {
	return z.states
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

	if len(z.states) == 0 {
		z.states = make([]*Zone, 0)
	}

	z.lock.Lock()
	defer z.lock.Unlock()

	zone := NewZone(name)
	z.states = append(z.states, zone)
	return zone
}

// handleOn takes care of the behavior when the light is set to On.  This
// includes brightness and color temperature.  The color hue of the light is
// handled by ZoneState_COLOR.
func (z *Zone) handleOn(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.On(ctx, z.name)
		if err != nil {
			return err
		}
	}

	err := z.handleBrightness(ctx)
	if err != nil {
		return err
	}

	return z.handleColorTemperature(ctx)

}

func (z *Zone) handleOff(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.Off(ctx, z.name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *Zone) handleColorTemperature(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.SetColorTemp(ctx, z.name, z.colorTemp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *Zone) handleBrightness(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.SetBrightness(ctx, z.name, z.brightness)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *Zone) handleColor(ctx context.Context) error {
	for _, h := range z.handlers {
		err := h.SetColor(ctx, z.name, z.color)
		if err != nil {
			return fmt.Errorf("%s color: %w", z.name, ErrHandlerFailed)
		}
	}

	return nil
}

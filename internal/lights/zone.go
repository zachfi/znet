package lights

import (
	"fmt"
	sync "sync"
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

func (z *Zone) Dim(brightness int32) error {
	z.brightness = brightness
	return z.SetState(Dim)
}

func (z *Zone) Off() error {
	return z.SetState(Off)
}

func (z *Zone) On() error {
	return z.SetState(On)
}

func (z *Zone) RandomColor(colors []string) error {
	z.colorPool = colors
	return z.SetState(RandomColor)
}

func (z *Zone) SetColor(color string) error {
	z.color = color
	return z.SetState(Color)
}

func (z *Zone) SetState(state zoneState) error {
	z.lock.Lock()
	defer z.lock.Unlock()

	z.state = state

	return nil
}

func (z *Zone) Handle(name string, handlers ...Handler) error {
	switch z.state {
	case Off:
		for _, h := range handlers {
			err := h.Off(name)
			if err != nil {

				return fmt.Errorf("%s off: %w", name, ErrHandlerFailed)
			}
		}
	case On:
		for _, h := range handlers {
			err := h.On(name)
			if err != nil {

				return fmt.Errorf("%s on: %w", name, ErrHandlerFailed)
			}
		}
	case Color:
		for _, h := range handlers {
			err := h.SetColor(name, z.color)
			if err != nil {

				return fmt.Errorf("%s color: %w", name, ErrHandlerFailed)
			}
		}
	case RandomColor:
		for _, h := range handlers {
			err := h.RandomColor(name, z.colorPool)
			if err != nil {

				return fmt.Errorf("%s random color: %w", name, ErrHandlerFailed)
			}
		}
	case Dim:
		for _, h := range handlers {
			err := h.Dim(name, z.brightness)
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

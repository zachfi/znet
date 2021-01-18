package astro

import (
	"context"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/lights"
)

// Astro exposes the RPC methods to communicate astrological events to the system.
type Astro struct {
	lights *lights.Lights
}

func NewServer(cfg *config.Config) (*Astro, error) {
	l, err := lights.NewLights(cfg)
	if err != nil {
		return nil, err
	}

	astroServer := &Astro{
		lights: l,
	}

	return astroServer, nil
}

func (a *Astro) Sunrise(ctx context.Context, req *Empty) (*Empty, error) {
	a.lights.SetRoomForEvent("Sunrise")
	return &Empty{}, nil
}

func (a *Astro) Sunset(ctx context.Context, req *Empty) (*Empty, error) {
	a.lights.SetRoomForEvent("SunSet")
	return &Empty{}, nil
}

func (a *Astro) PreSunset(ctx context.Context, req *Empty) (*Empty, error) {
	a.lights.SetRoomForEvent("PreSunset")
	return &Empty{}, nil
}

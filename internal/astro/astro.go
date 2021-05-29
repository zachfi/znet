package astro

import (
	"context"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/lights"
)

// Astro exposes the RPC methods to communicate astrological events to the system.
type Astro struct {
	UnimplementedAstroServer
	lights *lights.Lights
}

func NewAstro(cfg *config.Config, l *lights.Lights) (*Astro, error) {
	astroServer := &Astro{
		lights: l,
	}

	return astroServer, nil
}

func (a *Astro) Sunrise(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "Sunrise")
}

func (a *Astro) Sunset(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "SunSet")
}

func (a *Astro) PreSunset(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, a.lights.SetRoomForEvent(ctx, "PreSunset")
}

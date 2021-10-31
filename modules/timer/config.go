package timer

import (
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/timer/named"
)

type Config struct {
	Astro    astro.Config `yaml:"astro"`
	Named    named.Config `yaml:"named"`
	TimeZone string       `yaml:"timezone"`
}

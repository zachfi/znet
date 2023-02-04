package timer

import (
	"github.com/zachfi/znet/internal/astro"
	"github.com/zachfi/znet/modules/timer/named"
)

type Config struct {
	Astro    astro.Config `yaml:"astro"`
	Named    named.Config `yaml:"named"`
	TimeZone string       `yaml:"timezone" json:"timezone"`
}

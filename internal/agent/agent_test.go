package agent

import (
	"testing"
	"time"

	"github.com/tj/assert"
	"github.com/xaque208/znet/internal/gitwatch"
)

func TestPassFilter(t *testing.T) {
	config := Config{}

	filter := Filter{
		Names: []string{
			"znet",
		},
		URLs:     []string{},
		Branches: []string{},
	}

	ti := time.Now()

	x := gitwatch.NewTag{
		Time: &ti,
		Name: "znet",
		URL:  "https://github.com/xaque208/znet.git",
		Tag:  "v0.16.3",
	}

	ag := NewAgent(config, nil)

	assert.True(t, ag.passFilter(filter, x))
}

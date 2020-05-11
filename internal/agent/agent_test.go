package agent

import (
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/xaque208/znet/internal/gitwatch"
)

func TestPassFilter_Basic(t *testing.T) {
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

func TestPassFilter_Collection(t *testing.T) {
	config := Config{}

	filter := Filter{
		Names:    []string{},
		URLs:     []string{},
		Branches: []string{},
		Collections: []string{
			"go_project",
		},
	}

	ti := time.Now()

	x := gitwatch.NewTag{
		Time:       &ti,
		Name:       "znet",
		URL:        "https://github.com/xaque208/znet.git",
		Tag:        "v0.16.3",
		Collection: "go_project",
	}

	ag := NewAgent(config, nil)

	assert.True(t, ag.passFilter(filter, x))

	y := gitwatch.NewTag{
		Time:       &ti,
		Name:       "znet",
		URL:        "https://github.com/xaque208/znet.git",
		Tag:        "v0.16.3",
		Collection: "go_projects",
	}

	assert.False(t, ag.passFilter(filter, y))
}

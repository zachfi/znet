// +build unit

package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testEvent struct {
	Name string     `json:"name,omitempty"`
	Time *time.Time `json:"time,omitempty"`
}

type testFilter struct {
	Name string `json:"name,omitempty"`
}

func (n *testEvent) Filter(filter interface{}) bool {
	f := filter.(testFilter)
	return f.Name == n.Name
}

func TestFilter(t *testing.T) {
	now := time.Now()
	p := testEvent{
		Name: "Sunset",
		Time: &now,
	}

	f := testFilter{
		Name: "Sunset",
	}

	pass := p.Filter(f)
	require.True(t, pass)
}

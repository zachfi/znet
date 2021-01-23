package agent

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/events"
)

func TestAgent_Basic(t *testing.T) {
	config := &config.Config{}
	ag := NewAgent(config, nil)
	require.NotNil(t, ag)
}

func TestAgent_ConsumerInterface(t *testing.T) {
	var a events.Consumer = &Agent{}
	t.Logf("agent: %+v", a)
}

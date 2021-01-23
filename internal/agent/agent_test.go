package agent

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
)

func TestAgent_Basic(t *testing.T) {
	config := &config.Config{}
	ag := NewAgent(config, nil)
	require.NotNil(t, ag)
}

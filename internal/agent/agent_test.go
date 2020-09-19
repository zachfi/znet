package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgent_Basic(t *testing.T) {
	config := Config{}
	ag := NewAgent(config, nil)
	require.NotNil(t, ag)
}

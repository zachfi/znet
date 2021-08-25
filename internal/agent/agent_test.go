//go:build unit

package agent

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
)

func TestAgent_emptyConfig(t *testing.T) {
	config := &config.Config{}
	ag, err := NewAgent(config, nil)
	require.Error(t, err)
	require.Nil(t, ag)
}

func TestAgent_Config(t *testing.T) {
	c := &config.Config{}

	ag, err := NewAgent(c, nil)
	require.Error(t, err, "no TLS config")
	require.Nil(t, ag)

	c.TLS = &config.TLSConfig{}

	ag, err = NewAgent(c, nil)
	require.Error(t, err, "no Vault config")
	require.Nil(t, ag)

	c.Vault = &config.VaultConfig{}

	ag, err = NewAgent(c, nil)
	require.Error(t, err, "no RPC config")
	require.Nil(t, ag)

	c.RPC = &config.RPCConfig{}

	ag, err = NewAgent(c, nil)
	require.Error(t, err, "no Agent config")
	require.Nil(t, ag)

	c.Agent = &config.AgentConfig{}

	ag, err = NewAgent(c, nil)
	require.NoError(t, err)
	require.NotNil(t, ag)
}

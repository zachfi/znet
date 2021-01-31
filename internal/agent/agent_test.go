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

func TestAgent_IncompleteVaultConfig(t *testing.T) {
	config := &config.Config{
		TLS:   &config.TLSConfig{},
		Vault: &config.VaultConfig{},
		RPC:   &config.RPCConfig{},
		Agent: &config.AgentConfig{},
	}
	ag, err := NewAgent(config, nil)
	require.Error(t, err, "unsupported protocol scheme")
	require.Nil(t, ag)
}

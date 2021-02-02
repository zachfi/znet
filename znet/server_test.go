package znet

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaque208/znet/internal/config"
)

func TestNewServer_missingConfigs(t *testing.T) {

	cfg := &config.Config{}

	s, err := NewServer(cfg)
	require.Error(t, err, "unable to build znet Server with nil HTTP config")
	require.Nil(t, s)

	cfg.HTTP = &config.HTTPConfig{}

	s, err = NewServer(cfg)
	require.Error(t, err, "unable to build znet Server with nil RPC config")
	require.Nil(t, s)

	cfg.RPC = &config.RPCConfig{}

	s, err = NewServer(cfg)
	require.Error(t, err, "unable to build znet Server with nil Vault config")
	require.Nil(t, s)

	cfg.Vault = &config.VaultConfig{}

	s, err = NewServer(cfg)
	require.Error(t, err, "unable to build znet Server with nil TLS config")
	require.Nil(t, s)

	cfg.TLS = &config.TLSConfig{}

	s, err = NewServer(cfg)
	require.Error(t, err, "unable to summon vault token")
	require.Nil(t, s)

	// s, err = NewServer(cfg)
	// require.NoError(t, err)
	// require.NotNil(t, s)
}

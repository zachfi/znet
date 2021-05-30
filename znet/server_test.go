// +build unit

package znet

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/iot"
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
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestServer_Start(t *testing.T) {
	cfg := &config.Config{
		HTTP: &config.HTTPConfig{
			ListenAddress: ":0",
		},
		RPC: &config.RPCConfig{
			ListenAddress: ":0",
		},
		Vault: &config.VaultConfig{},
		TLS:   &config.TLSConfig{},
		Lights: &config.LightsConfig{
			Rooms: []config.LightsRoom{
				{
					Name: "zone1",
				},
			},
		},
	}

	s, err := NewServer(cfg)
	require.NoError(t, err)
	require.NotNil(t, s)

	cfgFile := "../testdata/config.yaml"
	z, err := NewZnet(cfgFile)
	require.NoError(t, err)
	require.NotNil(t, z)

	s.mqttClient = &iot.MockClient{}
	s.invClient = &inventory.MockInventory{}

	err = s.Start(z)
	require.NoError(t, err)

	err = s.Stop()
	require.NoError(t, err)

}

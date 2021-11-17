//go:build unit

package znet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewZnet(t *testing.T) {
	cfgFile := "../testdata/config.yaml"
	cfg, err := LoadConfig(cfgFile)
	require.NoError(t, err)
	z, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, z)
}

//go:build unit

package znet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewZnet(t *testing.T) {
	cfgFileNoExists := "../testdata/doesnotexist.yaml"
	z, err := NewZnet(cfgFileNoExists)
	require.Error(t, err, "failed to load config")
	require.Nil(t, z)

	cfgFile := "../testdata/config.yaml"
	z, err = NewZnet(cfgFile)
	require.NoError(t, err)
	require.NotNil(t, z)
}

//go:build unit

package znet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewZnet(t *testing.T) {
	cfgFile := "../testdata/config.yaml"
	z, err = NewZnet(cfgFile)
	require.NoError(t, err)
	require.NotNil(t, z)
}

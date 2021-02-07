package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLightsConfig(t *testing.T) {

	l := LightsConfig{}
	_, err := l.Room("name")
	require.Error(t, err)
}

// +build unit

package iot

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/require"
)

func TestMqttClientInterface(t *testing.T) {
	var client mqtt.Client = &MockClient{}
	require.NotNil(t, client)
}

func TestMqttToken_interface(t *testing.T) {
	var token mqtt.Token = &MockToken{}
	require.NotNil(t, token)
}

//go:build unit

package lights

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/iot"
)

func TestZigbeeLight_interface(t *testing.T) {
	var l Handler = zigbeeLight{}
	require.NotNil(t, l)
}

func TestZigbeeLight_New(t *testing.T) {
	invClient := &inventory.MockInventory{}
	mqttClient := &iot.MockClient{}

	l, err := NewZigbeeLight(&config.Config{}, mqttClient, invClient)
	require.NoError(t, err)
	require.NotNil(t, l)

	invClient.ListZigbeeDeviceResponse = []inventory.ZigbeeDevice{
		{
			IotZone: "group1",
			Name:    "testdevice1",
		},
	}

	require.NoError(t, l.On("group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"ON","transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Off("group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"OFF","transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Dim("group1", 123))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"brightness":123,"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Alert("group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"effect":"blink","transition":0.1}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.SetColor("group1", "#006c7f"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"color":{"hex":"#006c7f"},"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.RandomColor("group1", []string{"#006c7f"}))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"color":{"hex":"#006c7f"},"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Toggle("group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"TOGGLE","transition":0.5}`, mqttClient.LastPublishPayload)

}

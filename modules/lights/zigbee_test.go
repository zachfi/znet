//go:build unit

package lights

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/pkg/iot"
)

func TestZigbeeLight_interface(t *testing.T) {
	var l Handler = zigbeeLight{}
	require.NotNil(t, l)
}

func TestZigbeeLight_New(t *testing.T) {
	invClient := &inventory.MockInventory{}
	mqttClient := &iot.MockClient{}

	ctx := context.Background()

	l, err := NewZigbeeLight(Config{}, mqttClient, invClient)
	require.NoError(t, err)
	require.NotNil(t, l)

	invClient.ListZigbeeDeviceResponse = []inventory.ZigbeeDevice{
		{
			IotZone: "group1",
			Name:    "testdevice1",
			Vendor:  "Philips",
			Type:    "Router",
			ModelId: "LCA003",
		},
	}

	require.NoError(t, l.On(ctx, "group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"ON","transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Off(ctx, "group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"OFF","transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Dim(ctx, "group1", 123))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"brightness":123,"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Alert(ctx, "group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"effect":"blink","transition":0.1}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.SetColor(ctx, "group1", "#006c7f"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"color":{"hex":"#006c7f"},"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.RandomColor(ctx, "group1", []string{"#006c7f"}))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"color":{"hex":"#006c7f"},"transition":0.5}`, mqttClient.LastPublishPayload)

	require.NoError(t, l.Toggle(ctx, "group1"))
	require.Equal(t, "zigbee2mqtt/testdevice1/set", mqttClient.LastPublishTopic)
	require.Equal(t, `{"state":"TOGGLE","transition":0.5}`, mqttClient.LastPublishPayload)

}

func TestIsLightDevice(t *testing.T) {

	cases := []struct {
		d            inventory.ZigbeeDevice
		isLight      bool
		isColorLight bool
	}{
		{
			d: inventory.ZigbeeDevice{
				Vendor:  "Philips",
				Type:    "Router",
				ModelId: "LCA003",
			},
			isLight:      true,
			isColorLight: true,
		},
	}

	for _, tc := range cases {
		light := isLightDevice(&tc.d)
		require.Equal(t, tc.isLight, light)

		color := isColorLightDevice(&tc.d)
		require.Equal(t, tc.isColorLight, color)
	}

}

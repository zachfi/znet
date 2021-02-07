// +build unit

package telemetry

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/pkg/iot"
)

var zigbeeDeviceName string = "0x00158d0004238a81"

func TestNewServer(t *testing.T) {
	l := &lights.Lights{}
	h := &lights.MockLight{}
	l.AddHandler(h)

	s, err := NewServer(&inventory.MockInventory{}, l)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestReportIOTDevice_nilDiscovery(t *testing.T) {
	l := &lights.Lights{}
	h := &lights.MockLight{}
	l.AddHandler(h)

	s, err := NewServer(&inventory.MockInventory{}, l)
	require.NoError(t, err)
	require.NotNil(t, s)

	req := &inventory.IOTDevice{}

	response, err := s.ReportIOTDevice(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, response)
}

func TestReportIOTDevice_lights_handling(t *testing.T) {
	testCases := []struct {
		Handler *lights.MockLight
		Req     *inventory.IOTDevice
		Zone    string
	}{
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"double","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				OnCalls:       map[string]int{"dungeon": 1},
				SetColorCalls: map[string]int{"dungeon": 1},
				DimCalls:      map[string]int{"dungeon": 1},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"single","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				ToggleCalls: map[string]int{"dungeon": 1},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"hold","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				DimCalls: map[string]int{"dungeon": 1},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"quadruple","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				RandomColorCalls: map[string]int{"dungeon": 1},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"triple","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				OffCalls: map[string]int{"dungeon": 1},
			},
		},
	}

	for _, tc := range testCases {

		lightsConfig := &config.Config{
			Lights: &config.LightsConfig{
				PartyColors: []string{"#f33333"},
				Rooms: []config.LightsRoom{
					{
						Name: "dungeon",
					},
				},
			},
		}

		// l := &lights.Lights{}
		l, err := lights.NewLights(lightsConfig)
		require.NoError(t, err)
		require.NotNil(t, l)
		h := &lights.MockLight{}
		l.AddHandler(h)

		i := &inventory.MockInventory{}
		i.FetchZigbeeDeviceResponse = &inventory.ZigbeeDevice{
			Name:    zigbeeDeviceName,
			IotZone: "dungeon",
		}

		s, err := NewServer(i, l)
		require.NoError(t, err)
		require.NotNil(t, s)
		response, err := s.ReportIOTDevice(context.Background(), tc.Req)
		require.NoError(t, err)
		require.NotNil(t, response)

		// inventory
		require.Equal(t, 1, i.FetchZigbeeDeviceCalls[zigbeeDeviceName])
		require.Equal(t, 0, len(i.CreateZigbeeDeviceCalls))

		// lights handler
		require.Equal(t, tc.Handler, h)
	}

}

func TestReportIOTDevice_bridge_state(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	defer log.SetLevel(log.InfoLevel)

	testCases := []struct {
		// Handler *lights.MockLight
		Req *inventory.IOTDevice
		// Zone string
	}{
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  "bridge",
					Component: "zigbee2mqtt",
					Endpoint:  []string{"state"},
					Message:   []byte(`online`),
				},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  "bridge",
					Component: "zigbee2mqtt",
					Endpoint:  []string{"state"},
					Message:   []byte(`offline`),
				},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  "bridge",
					Component: "zigbee2mqtt",
					Endpoint:  []string{"log"},
					Message:   []byte(`{"message":"Update available for '0x001788010898e9c1'","meta":{"device":"0x001788010898e9c1","status":"available"},"type":"ota_update"}`),
				},
			},
		},
	}

	for _, tc := range testCases {

		lightsConfig := &config.Config{
			Lights: &config.LightsConfig{
				PartyColors: []string{"#f33333"},
				Rooms: []config.LightsRoom{
					{
						Name: "dungeon",
					},
				},
			},
		}

		// l := &lights.Lights{}
		l, err := lights.NewLights(lightsConfig)
		require.NoError(t, err)
		require.NotNil(t, l)

		h := &lights.MockLight{}
		l.AddHandler(h)

		i := &inventory.MockInventory{}
		// i.FetchZigbeeDeviceResponse = &inventory.ZigbeeDevice{
		// 	Name:    zigbeeDeviceName,
		// 	IotZone: "dungeon",
		// }

		iotServer, err := iot.NewServer(&iot.MockClient{})
		require.NoError(t, err)
		require.NotNil(t, iotServer)

		s, err := NewServer(i, l)
		require.NoError(t, err)
		require.NotNil(t, s)

		err = s.SetIOTServer(iotServer)
		require.NoError(t, err)

		response, err := s.ReportIOTDevice(context.Background(), tc.Req)
		require.NoError(t, err)
		require.NotNil(t, response)

		// inventory
		// require.Equal(t, 1, i.FetchZigbeeDeviceCalls[zigbeeDeviceName])
		// require.Equal(t, 0, len(i.CreateZigbeeDeviceCalls))

		// lights handler
		// require.Equal(t, tc.Handler, h)
	}

}

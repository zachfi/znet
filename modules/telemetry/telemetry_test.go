//go:build unit

package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/require"

	"github.com/zachfi/znet/modules/inventory"
	"github.com/zachfi/znet/modules/lights"
	"github.com/zachfi/znet/pkg/iot"
)

var zigbeeDeviceName = "0x00158d0004238a81"

func TestNewServer(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	l := &lights.Lights{}
	h := &lights.MockLight{}
	l.AddHandler(h)

	s, err := New(Config{}, logger, &inventory.MockInventory{}, l)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestReportIOTDevice_nilDiscovery(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	l := &lights.Lights{}
	h := &lights.MockLight{}
	l.AddHandler(h)

	s, err := New(Config{}, logger, &inventory.MockInventory{}, l)
	require.NoError(t, err)
	require.NotNil(t, s)

	req := &inventory.IOTDevice{}

	response, err := s.ReportIOTDevice(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, response)
}

func TestReportIOTDevice_lights_handling(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	testCases := map[string]struct {
		Handler *lights.MockLight
		Req     *inventory.IOTDevice
		Zone    string
	}{
		"double": {
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"double","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				OnCalls:            map[string]int{"dungeon": 1},
				SetColorTempCalls:  map[string]int{"dungeon": 1},
				SetBrightnessCalls: map[string]int{"dungeon": 1},
			},
		},
		"single": {
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
		"hold": {
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  zigbeeDeviceName,
					Component: "zigbee2mqtt",
					Message:   []byte(`{"action":"hold","battery":100,"linkquality":0,"voltage":3042}`),
				},
			},
			Zone: "dungeon",
			Handler: &lights.MockLight{
				OnCalls:            map[string]int{"dungeon": 1},
				SetBrightnessCalls: map[string]int{"dungeon": 1},
				SetColorTempCalls:  map[string]int{"dungeon": 1},
			},
		},
		"quadruple": {
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
		"triple": {
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

	for name, tc := range testCases {
		t.Logf("test: %s", name)

		lightsConfig := lights.Config{
			PartyColors: []string{"#f33333"},
			Rooms: []lights.Room{
				{
					Name: "dungeon",
				},
			},
		}

		// l := &lights.Lights{}
		l, err := lights.New(lightsConfig, logger)
		require.NoError(t, err)
		require.NotNil(t, l)
		h := &lights.MockLight{}
		l.AddHandler(h)

		i := &inventory.MockInventory{}
		i.FetchZigbeeDeviceResponse = &inventory.ZigbeeDevice{
			Name:    zigbeeDeviceName,
			IotZone: "dungeon",
		}

		s, err := New(Config{}, logger, i, l)
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
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	testCases := []struct {
		Req *inventory.IOTDevice
		Inv *inventory.MockInventory
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
					Message: []byte(`
					{
						"message":"Update available for '0x001788010898e9c1'",
						"meta":{"device":"0x001788010898e9c1","status":"available"},
						"type":"ota_update"
					}`),
				},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  "bridge",
					Component: "zigbee2mqtt",
					Endpoint:  []string{"logging"},
					Message: []byte(`
          {
            "level":"info",
            "message":"MQTT publish: topic 'zigbee2mqtt/0x00158d0006fafbeb',
              payload '{
              \"battery\":91,
              \"humidity\":29.3,
              \"linkquality\":86,
              \"pressure\":839,
              \"temperature\":21.27,
              \"voltage\":2985}'"
					}`),
				},
			},
		},
		{
			Req: &inventory.IOTDevice{
				DeviceDiscovery: &iot.DeviceDiscovery{
					ObjectId:  "bridge",
					Component: "zigbee2mqtt",
					Endpoint:  []string{"devices"},
					Message:   sampleDevices,
				},
			},
			Inv: &inventory.MockInventory{
				FetchZigbeeDeviceErr: fmt.Errorf("mock error"),
				FetchZigbeeDeviceCalls: map[string]int{
					"0x00158d0003960d06": 1,
					"0x00158d0004238a36": 1,
					"0x00158d0004238a81": 1,
					"0x00178801042131ca": 1,
					"0x0017880104215e6a": 1,
					"0x0017880104650857": 1,
					"0x00178801087fc8c8": 1,
					"0x001788010898e9c1": 1,
				},
				CreateZigbeeDeviceCalls: map[string]int{
					"0x00158d0003960d06": 1,
					"0x00158d0004238a36": 1,
					"0x00158d0004238a81": 1,
					"0x00178801042131ca": 1,
					"0x0017880104215e6a": 1,
					"0x0017880104650857": 1,
					"0x00178801087fc8c8": 1,
					"0x001788010898e9c1": 1,
				},
				UpdateZigbeeDeviceCalls: map[string]int{
					"0x00158d0003960d06": 1,
					"0x00158d0004238a36": 1,
					"0x00158d0004238a81": 1,
					"0x00178801042131ca": 1,
					"0x0017880104215e6a": 1,
					"0x0017880104650857": 1,
					"0x00178801087fc8c8": 1,
					"0x001788010898e9c1": 1,
				},
			},
		},
	}

	// nolint
	// zigbee2mqtt/bridge/log {"message":{"friendly_name":"0x00178801042131ca"},"type":"device_connected"}
	// zigbee2mqtt/bridge/log {"message":"interview_started","meta":{"friendly_name":"0x00178801042131ca"},"type":"pairing"}
	// zigbee2mqtt/bridge/log {"message":"announce","meta":{"friendly_name":"0x00178801042131ca"},"type":"device_announced"}
	// zigbee2mqtt/bridge/log {"message":"interview_successful","meta":{"description":"Hue white A60 bulb E27","friendly_name":"0x00178801042131ca","model":"9290011370","supported":true,"vendor":"Philips"},"type":"pairing"}
	// zigbee2mqtt/0x001788010898e9c1 {"brightness":254,"color":{"x":0.2061,"y":0.083},"color_temp":160,"linkquality":68,"state":"OFF","update_available":true}
	// zigbee2mqtt/bridge/log {"message":"Update available for '0x001788010898e9c1'","meta":{"device":"0x001788010898e9c1","status":"available"},"type":"ota_update"}
	// zigbee2mqtt/bridge/logging {"level":"info","message":"MQTT publish: topic 'zigbee2mqtt/0x00158d0006fafbeb', payload '{\"battery\":91,\"humidity\":29.3,\"linkquality\":86,\"pressure\":839,\"temperature\":21.27,\"voltage\":2985}'"}

	for _, tc := range testCases {

		lightsConfig := lights.Config{
			PartyColors: []string{"#f33333"},
			Rooms: []lights.Room{
				{
					Name: "dungeon",
				},
			},
		}

		// l := &lights.Lights{}
		l, err := lights.New(lightsConfig, logger)
		require.NoError(t, err)
		require.NotNil(t, l)

		h := &lights.MockLight{}
		l.AddHandler(h)

		i := &inventory.MockInventory{}

		if tc.Inv != nil {
			if tc.Inv.FetchZigbeeDeviceErr != nil {
				i.FetchZigbeeDeviceErr = tc.Inv.FetchZigbeeDeviceErr
			}
		}

		iotServer, err := iot.NewServer(&iot.MockClient{})
		require.NoError(t, err)
		require.NotNil(t, iotServer)

		s, err := New(Config{}, logger, i, l)
		require.NoError(t, err)
		require.NotNil(t, s)

		err = s.SetIOTServer(iotServer)
		require.NoError(t, err)

		response, err := s.ReportIOTDevice(context.Background(), tc.Req)
		require.NoError(t, err)
		require.NotNil(t, response)

		if tc.Inv != nil {
			require.Equal(t, tc.Inv, i)
		}
	}

}

var sampleDevices = []byte(`[
  {
    "dateCode": "20190608",
    "friendly_name": "Coordinator",
    "ieeeAddr": "0x00124b0014d91d6b",
    "lastSeen": 1612731063363,
    "networkAddress": 0,
    "softwareBuildID": "zStack12",
    "type": "Coordinator"
  },
  {
    "description": "MiJia wireless switch",
    "friendly_name": "0x00158d0004238a36",
    "ieeeAddr": "0x00158d0004238a36",
    "lastSeen": 1612729314959,
    "manufacturerID": 4151,
    "manufacturerName": "LUMI",
    "model": "WXKG01LM",
    "modelID": "lumi.sensor_switch",
    "networkAddress": 53291,
    "powerSource": "Battery",
    "type": "EndDevice",
    "vendor": "Xiaomi"
  },
  {
    "description": "MiJia wireless switch",
    "friendly_name": "0x00158d0004238a81",
    "ieeeAddr": "0x00158d0004238a81",
    "lastSeen": 1612730250297,
    "manufacturerID": 4151,
    "manufacturerName": "LUMI",
    "model": "WXKG01LM",
    "modelID": "lumi.sensor_switch",
    "networkAddress": 30828,
    "powerSource": "Battery",
    "type": "EndDevice",
    "vendor": "Xiaomi"
  },
  {
    "description": "MiJia wireless switch",
    "friendly_name": "0x00158d0003960d06",
    "ieeeAddr": "0x00158d0003960d06",
    "lastSeen": 1612729869704,
    "manufacturerID": 4151,
    "manufacturerName": "LUMI",
    "model": "WXKG01LM",
    "modelID": "lumi.sensor_switch",
    "networkAddress": 57158,
    "powerSource": "Battery",
    "type": "EndDevice",
    "vendor": "Xiaomi"
  },
  {
    "dateCode": "20191218",
    "description": "Hue white A60 bulb E27",
    "friendly_name": "0x0017880104215e6a",
    "hardwareVersion": 1,
    "ieeeAddr": "0x0017880104215e6a",
    "lastSeen": 1596902235381,
    "manufacturerID": 4107,
    "manufacturerName": "Philips",
    "model": "9290011370",
    "modelID": "LWB014",
    "networkAddress": 59215,
    "powerSource": "Mains (single phase)",
    "softwareBuildID": "1.50.2_r30933",
    "type": "Router",
    "vendor": "Philips"
  },
  {
    "dateCode": "20170311",
    "description": "Hue white A60 bulb E27",
    "friendly_name": "0x0017880104650857",
    "hardwareVersion": 1,
    "ieeeAddr": "0x0017880104650857",
    "lastSeen": 1612726431059,
    "manufacturerID": 4107,
    "manufacturerName": "Philips",
    "model": "9290011370",
    "modelID": "LWB014",
    "networkAddress": 56735,
    "powerSource": "Mains (single phase)",
    "softwareBuildID": "1.23.0_r20156",
    "type": "Router",
    "vendor": "Philips"
  },
  {
    "dateCode": "20200327",
    "description": "Hue white and color ambiance E26/E27",
    "friendly_name": "0x00178801087fc8c8",
    "hardwareVersion": 1,
    "ieeeAddr": "0x00178801087fc8c8",
    "lastSeen": 1612730510331,
    "manufacturerID": 4107,
    "manufacturerName": "Philips",
    "model": "9290022166",
    "modelID": "LCA003",
    "networkAddress": 42279,
    "powerSource": "Mains (single phase)",
    "softwareBuildID": "1.65.11_hB798F2B",
    "type": "Router",
    "vendor": "Philips"
  },
  {
    "dateCode": "20200327",
    "description": "Hue white and color ambiance E26/E27",
    "friendly_name": "0x001788010898e9c1",
    "hardwareVersion": 1,
    "ieeeAddr": "0x001788010898e9c1",
    "lastSeen": 1612730931271,
    "manufacturerID": 4107,
    "manufacturerName": "Philips",
    "model": "9290022166",
    "modelID": "LCA003",
    "networkAddress": 36588,
    "powerSource": "Mains (single phase)",
    "softwareBuildID": "1.65.11_hB798F2B",
    "type": "Router",
    "vendor": "Philips"
  },
  {
    "dateCode": "20191218",
    "description": "Hue white A60 bulb E27",
    "friendly_name": "0x00178801042131ca",
    "hardwareVersion": 1,
    "ieeeAddr": "0x00178801042131ca",
    "lastSeen": 1612730883251,
    "manufacturerID": 4107,
    "manufacturerName": "Philips",
    "model": "9290011370",
    "modelID": "LWB014",
    "networkAddress": 65099,
    "powerSource": "Mains (single phase)",
    "softwareBuildID": "1.50.2_r30933",
    "type": "Router",
    "vendor": "Philips"
  }
]`)

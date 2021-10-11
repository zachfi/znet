//go:build unit

package iot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZigbeeBridgeLog(t *testing.T) {

	cases := []struct {
		Message []byte
		Obj     ZigbeeBridgeLog
	}{
		{
			[]byte(`{
				"message":"interview_successful",
				"meta":{
					"description":"Hue white A60 bulb E27",
					"friendly_name":"0x0017880102be373d",
					"model":"9290011370",
					"supported":true,
					"vendor":"Philips"
				},"type":"pairing"
			}`),
			ZigbeeBridgeLog{
				Type:    "pairing",
				Message: "interview_successful",
				Meta: map[string]interface{}{
					"description":   "Hue white A60 bulb E27",
					"friendly_name": "0x0017880102be373d",
					"model":         "9290011370",
					"vendor":        "Philips",
					"supported":     true,
				},
			},
		},
		{
			[]byte(`{"type":"device_announced","message":"announce","meta":{"friendly_name":"0x0017880104650857"}}`),
			ZigbeeBridgeLog{
				Type:    "device_announced",
				Message: "announce",
				Meta: map[string]interface{}{
					"friendly_name": "0x0017880104650857",
				},
			},
		},
		{
			[]byte(`{
					"message":[
						{
							"dateCode":"20190608",
							"friendly_name":"Coordinator",
							"ieeeAddr":"0x00124b0014d91d6b",
							"lastSeen":1612129344843,
							"networkAddress":0,
							"softwareBuildID":"zStack12",
							"type":"Coordinator"
						},{
							"dateCode":"20200327",
							"description":"Hue white and color ambiance E26/E27",
							"friendly_name":"0x001788010898e9c1",
							"hardwareVersion":1,
							"ieeeAddr":"0x001788010898e9c1",
							"lastSeen":1612127953195,
							"manufacturerID":4107,
							"manufacturerName":"Philips",
							"model":"9290022166",
							"modelID":"LCA003",
							"networkAddress":36588,
							"powerSource":"Mains (single phase)",
							"softwareBuildID":"1.65.11_hB798F2B",
							"type":"Router",
							"vendor":"Philips"
						}
					],
					"type":"devices"
				}`),
			ZigbeeBridgeLog{
				Type: "devices",
				Message: ZigbeeBridgeMessageDevices{
					{
						IeeeAddr:        "0x00124b0014d91d6b",
						Type:            "Coordinator",
						FriendlyName:    "Coordinator",
						SoftwareBuildID: "zStack12",
						DateCode:        "20190608",
						LastSeen:        1612129344843,
					},
					{
						IeeeAddr:         "0x001788010898e9c1",
						Type:             "Router",
						NetworkAddress:   36588,
						FriendlyName:     "0x001788010898e9c1",
						SoftwareBuildID:  "1.65.11_hB798F2B",
						DateCode:         "20200327",
						LastSeen:         1612127953195,
						Model:            "9290022166",
						Vendor:           "Philips",
						Description:      "Hue white and color ambiance E26/E27",
						ManufacturerID:   4107,
						ManufacturerName: "Philips",
						PowerSource:      "Mains (single phase)",
						ModelID:          "LCA003",
						HardwareVersion:  1,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		obj := ZigbeeBridgeLog{}
		err := json.Unmarshal(tc.Message, &obj)
		require.NoError(t, err)
		require.Equal(t, tc.Obj, obj)
	}

}

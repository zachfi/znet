//go:build unit

package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/require"
)

func TTTestUpdateMessageFixtures(t *testing.T) {
	cfg := MQTTConfig{
		URL:      "tcp://localhost:1883",
		Topic:    "#",
		Username: "iot",
		Password: "xxx",
	}

	var onMessageReceived mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		topicPath, err := ParseTopicPath(msg.Topic())
		require.NoError(t, err)
		discovery := ParseDiscoveryMessage(topicPath, msg)
		_, err = ReadZigbeeMessage(discovery.ObjectId, discovery.Message, discovery.Endpoint...)
		require.NoError(t, err)

		e := strings.Join(discovery.Endpoint, "/")

		switch e {
		case "devices":
			err := os.WriteFile("../../testdata/devices.json", msg.Payload(), 0644)
			require.NoError(t, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mqttClient, err := NewMQTTClient(cfg, log.NewNopLogger())
	require.NoError(t, err)

	token := mqttClient.Subscribe(cfg.Topic, 0, onMessageReceived)
	token.Wait()
	require.NoError(t, token.Error())

	topic := "zigbee2mqtt/bridge/config/devices"
	token = mqttClient.Publish(topic, byte(0), false, "")
	token.Wait()
	require.NoError(t, token.Error())

	<-ctx.Done()
}

func TestZigbeeBridgeLog_devices(t *testing.T) {
	cases := []struct {
		JSONFile string
	}{
		{
			JSONFile: "../../testdata/devices.json",
		},
	}

	for _, tc := range cases {
		jsonFile, err := os.Open(tc.JSONFile)
		require.NoError(t, err)
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		require.NoError(t, err)
		obj := ZigbeeMessageBridgeDevices{}
		err = json.Unmarshal(byteValue, &obj)
		require.NoError(t, err)
		require.Greater(t, len(obj), 0)

		for _, d := range obj {
			x := ZigbeeDeviceType(d)
			require.Greater(t, x, Unknown,
				fmt.Sprintf("device: %+v", d),
			)
		}

	}
}

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
							"date_code":"20190608",
							"friendly_name":"Coordinator",
							"ieee_address":"0x00124b0014d91d6b",
							"networkAddress":0,
							"software_build_id":"zStack12",
							"type":"Coordinator"
						},{
							"date_code":"20200327",
              "definition": {
                "description":"Hue white and color ambiance E26/E27",
                "model":"9290022166",
                "vendor":"Philips"
              },
							"friendly_name":"0x001788010898e9c1",
							"hardwareVersion":1,
							"ieee_address":"0x001788010898e9c1",
							"lastSeen":1612127953195,
							"manufacturerID":4107,
							"manufacturerName":"Philips",
							"model_id":"LCA003",
							"network_address":36588,
							"power_source":"Mains (single phase)",
							"software_build_id":"1.65.11_hB798F2B",
							"type":"Router"
						}
					],
					"type":"devices"
				}`),
			ZigbeeBridgeLog{
				Type: "devices",
				Message: ZigbeeMessageBridgeDevices{
					{
						IeeeAddress:     "0x00124b0014d91d6b",
						Type:            "Coordinator",
						FriendlyName:    "Coordinator",
						SoftwareBuildID: "zStack12",
						DateCode:        "20190608",
						// LastSeen:        1612129344843,
					},
					{
						IeeeAddress:     "0x001788010898e9c1",
						Type:            "Router",
						NetworkAddress:  36588,
						FriendlyName:    "0x001788010898e9c1",
						SoftwareBuildID: "1.65.11_hB798F2B",
						DateCode:        "20200327",
						// LastSeen:         1612127953195,

						Definition: ZigbeeBridgeDeviceDefinition{
							Model:       "9290022166",
							Vendor:      "Philips",
							Description: "Hue white and color ambiance E26/E27",
						},
						PowerSource: "Mains (single phase)",
						ModelID:     "LCA003",
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

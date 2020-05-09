// +build unit

package iot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTopicPath(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Topic  string
		Should TopicPath
	}

	samples := []testStruct{
		{
			Topic: "stat/ff3e8a21dc12507d3a159b4792403f01/tempcoef",
			Should: TopicPath{
				Component: "stat",
				ObjectID:  "tempcoef",
				NodeID:    "ff3e8a21dc12507d3a159b4792403f01",
				Endpoint:  []string{},
			},
		},

		{
			Topic: "stat/ff3e8a21dc12507d3a159b4792403f01/water/tempcoef",
			Should: TopicPath{
				Component: "stat",
				NodeID:    "ff3e8a21dc12507d3a159b4792403f01",
				ObjectID:  "water",
				Endpoint:  []string{"tempcoef"},
			},
		},

		{
			Topic: "homeassistant/binary_sensor/garden/config",
			Should: TopicPath{
				Component: "homeassistant",
				ObjectID:  "binary_sensor",
				Endpoint:  []string{"garden", "config"},
			},
		},

		{
			Topic: "homeassistant/binary_sensor/garden/state",
			Should: TopicPath{
				Component: "homeassistant",
				ObjectID:  "binary_sensor",
				Endpoint:  []string{"garden", "state"},
			},
		},

		{
			Topic: "workgroup/92696ed2ae92b430f4e9447583936628/wifi/ssid",
			Should: TopicPath{
				Component: "workgroup",
				NodeID:    "92696ed2ae92b430f4e9447583936628",
				ObjectID:  "wifi",
				Endpoint:  []string{"ssid"},
			},
		},

		{
			Topic: "homeassistant/light/18c114ad3dec7c1d29bc888e4e748f89/led1/config",
			Should: TopicPath{
				DiscoveryPrefix: "homeassistant",
				Component:       "light",
				NodeID:          "18c114ad3dec7c1d29bc888e4e748f89",
				ObjectID:        "led1",
				Endpoint:        []string{"config"},
			},
		},

		// "workgroup/92696ed2ae92b430f4e9447583936628/wifi/bssid",

		// "stat/92696ed2ae92b430f4e9447583936628/tempcoef",
		// "stat/92696ed2ae92b430f4e9447583936628/water/tempcoef",
		// "stat/f51d958dbc60b0519d7e64f14cc733ab/tempcoef",
		// "stat/f51d958dbc60b0519d7e64f14cc733ab/water/tempcoef",
		// "stat/18c114ad3dec7c1d29bc888e4e748f89/led1/color",
		// "stat/18c114ad3dec7c1d29bc888e4e748f89/led1/power",
		// "stat/18c114ad3dec7c1d29bc888e4e748f89/led2/color",
		// "stat/18c114ad3dec7c1d29bc888e4e748f89/led2/power",
		// "workgroup/92696ed2ae92b430f4e9447583936628/wifi/ssid",
		// "workgroup/92696ed2ae92b430f4e9447583936628/wifi/bssid",
		// "workgroup/92696ed2ae92b430f4e9447583936628/wifi/rssi",
		// "workgroup/92696ed2ae92b430f4e9447583936628/wifi/ip",
		// "workgroup/92696ed2ae92b430f4e9447583936628/device",
		// "workgroup/92696ed2ae92b430f4e9447583936628/sketch",
		// "workgroup/92696ed2ae92b430f4e9447583936628/air/temperature",
		// "workgroup/92696ed2ae92b430f4e9447583936628/air/humidity",
		// "workgroup/92696ed2ae92b430f4e9447583936628/air/heatindex",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/air/temperature",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/air/humidity",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/air/heatindex",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/wifi/ssid",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/wifi/bssid",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/wifi/rssi",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/wifi/ip",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/device",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/sketch",
		// "workgroup/ff3e8a21dc12507d3a159b4792403f01/light",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/wifi/ssid",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/wifi/bssid",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/wifi/rssi",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/wifi/ip",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/device",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/sketch",
		// "workgroup/f51d958dbc60b0519d7e64f14cc733ab/light",
		// "workgroup/18c114ad3dec7c1d29bc888e4e748f89/wifi/ssid",
		// "workgroup/18c114ad3dec7c1d29bc888e4e748f89/wifi/bssid",
		// "workgroup/18c114ad3dec7c1d29bc888e4e748f89/wifi/rssi",
		// "workgroup/18c114ad3dec7c1d29bc888e4e748f89/wifi/ip",
		// "workgroup/18c114ad3dec7c1d29bc888e4e748f89/sketch",
		// "homeassistant/light/18c114ad3dec7c1d29bc888e4e748f89/led1/config",
		// "homeassistant/light/18c114ad3dec7c1d29bc888e4e748f89/led2/config",
	}

	for _, s := range samples {
		result, err := ParseTopicPath(s.Topic)
		assert.Nil(t, err)
		assert.Equal(t, s.Should.DiscoveryPrefix, result.DiscoveryPrefix)
		assert.Equal(t, s.Should.Component, result.Component)
		assert.Equal(t, s.Should.NodeID, result.NodeID)
		assert.Equal(t, s.Should.ObjectID, result.ObjectID)
		assert.Equal(t, s.Should.Endpoint, result.Endpoint)
	}

}

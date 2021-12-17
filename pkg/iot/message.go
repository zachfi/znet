package iot

import (
	"encoding/json"
)

type ZigbeeMessage struct {
	Action      *string `json:"action,omitempty"`
	Battery     *int    `json:"battery,omitempty"`
	Illuminance *int    `json:"illuminance,omitempty"`
	LinkQuality *int    `json:"linkquality,omitempty"`
	Occupancy   *bool   `json:"occupancy,omitempty"`
	Tamper      *bool   `json:"tamper,omitempty"`
	Temperature *int    `json:"temperature,omitempty"`
	Voltage     *int    `json:"voltage,omitempty"`
	WaterLeak   *bool   `json:"water_leak,omitempty"`
}

type ZigbeeBridgeState string

const (
	Offline ZigbeeBridgeState = "offline"
	Online  ZigbeeBridgeState = "online"
)

// ZigbeeBridgeLogMessage
// https://www.zigbee2mqtt.io/information/mqtt_topics_and_message_structure.html#zigbee2mqttbridgelog
// zigbee2mqtt/bridge/log
// {"type":"device_announced","message":"announce","meta":{"friendly_name":"0x0017880104650857"}}
type ZigbeeBridgeLog struct {
	Type    string                 `json:"type,omitempty"`
	Message interface{}            `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

func (z *ZigbeeBridgeLog) UnmarshalJSON(data []byte) error {

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	z.Type, _ = v["type"].(string)
	message := v["message"]

	switch z.Type {
	case "device_announced":
		z.Message = v["message"].(string)
		z.Meta = v["meta"].(map[string]interface{})
	case "devices":
		j, err := json.Marshal(message)
		if err != nil {
			return err
		}

		m := ZigbeeMessageBridgeDevices{}
		err = json.Unmarshal(j, &m)
		if err != nil {
			return err
		}

		z.Message = m
	case "ota_update":
		z.Meta = v["meta"].(map[string]interface{})
		z.Message = v["message"].(string)
	case "pairing":
		z.Meta = v["meta"].(map[string]interface{})
		z.Message = v["message"].(string)
	}

	return nil
}

type ZigbeeMessageBridgeDevices []ZigbeeBridgeDevice

type ZigbeeBridgeDevice struct {
	IeeeAddress    string `json:"ieee_address"`
	Type           string `json:"type"`
	NetworkAddress int    `json:"network_address"`
	Supported      bool   `json:"supported"`
	FriendlyName   string `json:"friendly_name"`
	Endpoints      struct {
		Num1 struct {
			Bindings             []interface{} `json:"bindings"`
			ConfiguredReportings []interface{} `json:"configured_reportings"`
			Clusters             struct {
				Input  []string      `json:"input"`
				Output []interface{} `json:"output"`
			} `json:"clusters"`
		} `json:"1"`
	} `json:"endpoints"`
	Definition  ZigbeeBridgeDeviceDefinition `json:"definition"`
	PowerSource string                       `json:"power_source"`
	DateCode    string                       `json:"date_code"`
	ModelID     string                       `json:"model_id"`
	Scenes      []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"scenes"`
	Interviewing       bool   `json:"interviewing"`
	InterviewCompleted bool   `json:"interview_completed"`
	SoftwareBuildID    string `json:"software_build_id,omitempty"`
}

type ZigbeeBridgeDeviceDefinition struct {
	Model       string `json:"model"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
}

type WifiMessage struct {
	BSSID string `json:"bssid,omitempty"`
	IP    string `json:"ip,omitempty"`
	RSSI  int    `json:"rssi,omitempty"`
	SSID  string `json:"ssid,omitempty"`
}

type AirMessage struct {
	Humidity    *float32 `json:"humidity,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	HeatIndex   *float32 `json:"heatindex,omitempty"`
	TempCoef    *float64 `json:"tempcoef,omitempty"`
}

type WaterMessage struct {
	Temperature *float32 `json:"temperature,omitempty"`
	TempCoef    *float64 `json:"tempcoef,omitempty"`
}

type LEDConfig struct {
	Schema       string   `json:"schema"`
	Brightness   bool     `json:"brightness"`
	Rgb          bool     `json:"rgb"`
	Effect       bool     `json:"effect"`
	EffectList   []string `json:"effect_list"`
	Name         string   `json:"name"`
	UniqueID     string   `json:"unique_id"`
	CommandTopic string   `json:"command_topic"`
	StateTopic   string   `json:"state_topic"`
	Device       struct {
		Identifiers  string     `json:"identifiers"`
		Manufacturer string     `json:"manufacturer"`
		Model        string     `json:"model"`
		Name         string     `json:"name"`
		SwVersion    string     `json:"sw_version"`
		Connections  [][]string `json:"connections"`
	} `json:"device"`
}

type LEDColor struct {
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
	Effect     string `json:"effect"`
	Color      struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	} `json:"color"`
}

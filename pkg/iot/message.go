package iot

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type ZigbeeMessage struct {
	Battery     int    `json:"battery,omitempty"`
	LinkQuality int    `json:"linkquality,omitempty"`
	Click       string `json:"click,omitempty"`
	Voltage     int    `json:"voltage,omitempty"`
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

	z.Meta = v["meta"].(map[string]interface{})

	switch z.Type {
	case "device_announced":
		z.Message = v["message"].(string)
	case "devices":
		j, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
		}

		m := ZigbeeBridgeMessageDevices{}
		err = json.Unmarshal(j, &m)
		if err != nil {
			return err
		}

		z.Message = m
	}

	return nil
}

type ZigbeeBridgeMessageDevices []struct {
	IeeeAddr         string `json:"ieeeAddr"`
	Type             string `json:"type"`
	NetworkAddress   int    `json:"networkAddress"`
	FriendlyName     string `json:"friendly_name"`
	SoftwareBuildID  string `json:"softwareBuildID,omitempty"`
	DateCode         string `json:"dateCode,omitempty"`
	LastSeen         int64  `json:"lastSeen"`
	Model            string `json:"model,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Description      string `json:"description,omitempty"`
	ManufacturerID   int    `json:"manufacturerID,omitempty"`
	ManufacturerName string `json:"manufacturerName,omitempty"`
	PowerSource      string `json:"powerSource,omitempty"`
	ModelID          string `json:"modelID,omitempty"`
	HardwareVersion  int    `json:"hardwareVersion,omitempty"`
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

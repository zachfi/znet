package lights

import (
	"encoding/json"
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/modules/inventory"
)

type zigbeeLight struct {
	cfg        *Config
	inv        inventory.Inventory
	mqttClient mqtt.Client
}

const defaultTransitionTime = 0.5

func NewZigbeeLight(cfg Config, mqttClient mqtt.Client, inv inventory.Inventory) (Handler, error) {
	return &zigbeeLight{
		cfg:        &cfg,
		inv:        inv,
		mqttClient: mqttClient,
	}, nil
}

func (l zigbeeLight) Toggle(groupName string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"state":      "TOGGLE",
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) Alert(groupName string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"effect":     "blink",
			"transition": 0.1,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}
func (l zigbeeLight) On(groupName string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"state":      "ON",
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}

func (l zigbeeLight) Off(groupName string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"state":      "OFF",
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}

func (l zigbeeLight) Dim(groupName string, brightness int32) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"brightness": brightness,
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) SetColor(groupName string, hex string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isColorLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"transition": defaultTransitionTime,
			"color": map[string]string{
				"hex": hex,
			},
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) RandomColor(groupName string, hex []string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for i := range devices {
		if !isColorLightDevice(devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"transition": defaultTransitionTime,
			"color": map[string]string{
				"hex": hex[rand.Intn(len(hex))],
			},
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func isLightDevice(z inventory.ZigbeeDevice) bool {
	switch z.Vendor {
	case "Philips":
		if z.Type == "Router" {
			return true
		}
	}

	return false
}

func isColorLightDevice(z inventory.ZigbeeDevice) bool {
	switch z.Vendor {
	case "Philips":
		if z.ModelId == "LCA003" {
			return true
		}
	}

	return false
}

package lights

import (
	"encoding/json"
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/iot"
)

type zigbeeLight struct {
	config     *config.LightsConfig
	inv        *inventory.Inventory
	mqttClient mqtt.Client
}

func NewZigbeeLight(cfg *config.Config) (Handler, error) {
	inv, err := inventory.NewInventory(cfg.LDAP)
	if err != nil {
		return nil, err
	}

	mqttClient, err := iot.NewMQTTClient(cfg.MQTT)
	if err != nil {
		return nil, err
	}

	return &zigbeeLight{
		config:     cfg.Lights,
		inv:        inv,
		mqttClient: mqttClient,
	}, nil
}

func (l zigbeeLight) Toggle(groupName string) error {
	devices, err := l.inv.ListZigbeeDevices()
	if err != nil {
		return err
	}

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"state":      "TOGGLE",
			"transition": 0.5,
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"state":      "ON",
			"transition": 0.5,
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"state":      "OFF",
			"transition": 0.5,
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"brightness": brightness,
			"transition": 0.5,
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"transition": 0.5,
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

	for _, d := range *devices {
		if d.IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]interface{}{
			"transition": 0.5,
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

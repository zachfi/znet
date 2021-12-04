package lights

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/modules/inventory"
)

type zigbeeLight struct {
	cfg        *Config
	inv        inventory.Inventory
	mqttClient mqtt.Client
}

const defaultTransitionTime = 0.5
const slowTransitionTime = 5

func NewZigbeeLight(cfg Config, mqttClient mqtt.Client, inv inventory.Inventory) (Handler, error) {
	return &zigbeeLight{
		cfg:        &cfg,
		inv:        inv,
		mqttClient: mqttClient,
	}, nil
}

func (l zigbeeLight) Toggle(ctx context.Context, groupName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.Toggle")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(&devices[i]) {
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

func (l zigbeeLight) Alert(ctx context.Context, groupName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.Alert")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(&devices[i]) {
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

func (l zigbeeLight) On(ctx context.Context, groupName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.On")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		if !isLightDevice(&devices[i]) {
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

func (l zigbeeLight) Off(ctx context.Context, groupName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.Off")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		deviceSpan, _ := opentracing.StartSpanFromContext(ctx, devices[i].Name)
		if !isLightDevice(&devices[i]) {
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
		deviceSpan.Finish()
	}
	return nil
}

func (l zigbeeLight) SetBrightness(ctx context.Context, groupName string, brightness int32) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.Dim")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		deviceSpan, _ := opentracing.StartSpanFromContext(ctx, devices[i].Name)
		if !isLightDevice(&devices[i]) {
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
		deviceSpan.Finish()
	}

	return nil
}

func (l zigbeeLight) SetColor(ctx context.Context, groupName string, hex string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.SetColor")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		deviceSpan, _ := opentracing.StartSpanFromContext(ctx, devices[i].Name)
		if !isColorLightDevice(&devices[i]) {
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
		deviceSpan.Finish()
	}

	return nil
}

func (l zigbeeLight) RandomColor(ctx context.Context, groupName string, hex []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.RandomColor")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		deviceSpan, _ := opentracing.StartSpanFromContext(ctx, devices[i].Name)
		if !isColorLightDevice(&devices[i]) {
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
		deviceSpan.Finish()
	}

	return nil
}

func (l zigbeeLight) SetColorTemp(ctx context.Context, groupName string, temp int32) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "zigbeeLight.SetTemp")
	defer span.Finish()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		deviceSpan, _ := opentracing.StartSpanFromContext(ctx, devices[i].Name)
		if !isColorLightDevice(&devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"transition": slowTransitionTime,
			"color_temp": temp,
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
		deviceSpan.Finish()
	}

	return nil
}

func isLightDevice(z *inventory.ZigbeeDevice) bool {
	switch z.Vendor {
	case "Philips":
		if z.Type == "Router" {
			return true
		}
	}

	return false
}

func isColorLightDevice(z *inventory.ZigbeeDevice) bool {
	switch z.Vendor {
	case "Philips":
		if z.ModelId == "LCA003" {
			return true
		}
	}

	return false
}

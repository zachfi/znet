package lights

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/pkg/iot"
)

type zigbeeLight struct {
	cfg        *Config
	inv        inventory.Inventory
	mqttClient mqtt.Client

	logger log.Logger
	tracer trace.Tracer
}

const defaultTransitionTime = 0.5

func NewZigbeeLight(cfg Config, mqttClient mqtt.Client, inv inventory.Inventory, logger log.Logger) (Handler, error) {
	return &zigbeeLight{
		cfg:        &cfg,
		inv:        inv,
		mqttClient: mqttClient,
		logger:     log.With(logger, "light", "zigbee"),
		tracer:     otel.Tracer("zigbeeLight"),
	}, nil
}

func (l zigbeeLight) Toggle(ctx context.Context, groupName string) error {
	ctx, span := l.tracer.Start(ctx, "Toggle")
	defer span.End()

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
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) Alert(ctx context.Context, groupName string) error {
	ctx, span := l.tracer.Start(ctx, "Alert")
	defer span.End()

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
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}

func (l zigbeeLight) On(ctx context.Context, groupName string) error {
	ctx, span := l.tracer.Start(ctx, "On")
	defer span.End()

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
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}

func (l zigbeeLight) Off(ctx context.Context, groupName string) error {
	ctx, span := l.tracer.Start(ctx, "Off")
	defer span.End()

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
			"state":      "OFF",
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}
	return nil
}

func (l zigbeeLight) SetBrightness(ctx context.Context, groupName string, brightness int32) error {
	ctx, span := l.tracer.Start(ctx, "SetBrightness")
	defer span.End()

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
			"brightness": brightness,
			"transition": defaultTransitionTime,
		}

		m, err := json.Marshal(message)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) SetColor(ctx context.Context, groupName string, hex string) error {
	ctx, span := l.tracer.Start(ctx, "SetColor")
	defer span.End()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
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
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) RandomColor(ctx context.Context, groupName string, hex []string) error {
	ctx, span := l.tracer.Start(ctx, "RandomColor")
	defer span.End()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
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
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func (l zigbeeLight) SetColorTemp(ctx context.Context, groupName string, temp int32) error {
	ctx, span := l.tracer.Start(ctx, "SetColorTemp")
	defer span.End()

	devices, err := l.inv.ListZigbeeDevices(ctx)
	if err != nil {
		return err
	}

	for i := range devices {
		if !isColorLightDevice(&devices[i]) {
			continue
		}

		if devices[i].IotZone != groupName {
			continue
		}

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", devices[i].Name)
		message := map[string]interface{}{
			"transition": defaultTransitionTime,
			"color_temp": temp,
		}

		m, err := json.Marshal(message)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			_ = level.Error(l.logger).Log("err", err.Error())
			continue
		}

		l.mqttClient.Publish(topic, byte(0), false, string(m))
	}

	return nil
}

func isLightDevice(z *inventory.ZigbeeDevice) bool {
	switch zigbeeDeviceType(z) {
	case iot.ColorLight, iot.BasicLight, iot.Relay:
		return true
	}

	return false
}

func isColorLightDevice(z *inventory.ZigbeeDevice) bool {
	switch zigbeeDeviceType(z) {
	case iot.ColorLight:
		return true
	}

	return false
}

func zigbeeDeviceType(z *inventory.ZigbeeDevice) iot.DeviceType {

	d := iot.ZigbeeBridgeDevice{}
	d.Definition.Vendor = z.GetVendor()
	d.Definition.Model = z.GetModel()
	d.ModelID = z.GetModelId()

	return iot.ZigbeeDeviceType(d)

}

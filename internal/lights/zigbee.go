package lights

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type zigbeeLight struct {
	config          *config.LightsConfig
	inventoryClient rpc.InventoryClient
	mqttClient      mqtt.Client
}

func (l zigbeeLight) Toggle(groupName string) error {
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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
	ctx := context.Background()

	rand.Seed(time.Now().Unix())

	stream, err := l.inventoryClient.ListZigbeeDevices(ctx, &rpc.Empty{})
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var d *rpc.ZigbeeDevice

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

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

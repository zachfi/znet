package lights

import (
	"context"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type zigbeeLight struct {
	config          Config
	inventoryClient rpc.InventoryClient
	mqttClient      mqtt.Client
}

func (l zigbeeLight) Toggle(groupName string) error {

	log.Debugf("toggle: %s", groupName)

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

		log.Debugf("match: %s", d)

		topic := fmt.Sprintf("zigbee2mqtt/%s/set", d.Name)
		message := map[string]string{
			"state": "TOGGLE",
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
	return nil
}
func (l zigbeeLight) On(groupName string) error {
	return nil
}
func (l zigbeeLight) Off(groupName string) error {
	return nil
}
func (l zigbeeLight) Dim(groupName string, brightness int32) error {
	return nil
}

package iot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type TopicPath struct {
	DiscoveryPrefix string
	Component       string
	NodeID          string
	ObjectID        string
	Endpoint        []string
}

func ParseTopicPath(topic string) (TopicPath, error) {
	// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config

	var tp TopicPath
	tp.Endpoint = make([]string, 0)

	nodeIDRegex := regexp.MustCompile(`^.*/([0-9a-z]{32})/.*$`)
	parts := strings.Split(topic, "/")

	match := nodeIDRegex.FindAllStringSubmatch(topic, -1)
	if len(match) == 1 {
		// A node ID has been found in the topic path.

		// Figure out which part matches the node ID.
		nodeIndex := func(parts []string) int {
			for i, p := range parts {
				if p == match[0][1] {
					return i
				}
			}

			return 0
		}(parts)

		// Determine if we have a discovery_prefix to set.
		if nodeIndex == 1 {
			tp.Component = parts[0]
			tp.NodeID = parts[nodeIndex]
		} else if nodeIndex == 2 {
			tp.DiscoveryPrefix = parts[0]
			tp.Component = parts[1]
			tp.NodeID = parts[nodeIndex]
		}

		if nodeIndex > 0 {
			next := nodeIndex + 1
			tp.ObjectID = parts[next]
			next++
			tp.Endpoint = parts[next:]
		}

	} else {
		// else a node ID is not matched in the topic path.
		tp.Component = parts[0]
		tp.ObjectID = parts[1]
		tp.Endpoint = parts[2:]
	}

	return tp, nil
}

// ReadZigbeeMessage implements the payload unmarshaling for zigbee2mqtt
// https://www.zigbee2mqtt.io/information/mqtt_topics_and_message_structure.html
func ReadZigbeeMessage(objectID string, payload []byte, endpoint ...string) (interface{}, error) {

	switch objectID {
	case "bridge":
		if len(endpoint) == 1 {
			// topic: zigbee2mqtt/bridge/log
			switch endpoint[0] {
			case "log":
				m := ZigbeeBridgeLog{}
				err := json.Unmarshal(payload, &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			case "state":
				m := ZigbeeBridgeState(string(payload))
				if m != "" {
					return m, nil
				}
			case "config", "logging":
				// do nothing for a config message
				return nil, nil
			case "devices":
				m := ZigbeeBridgeMessageDevices{}
				err := json.Unmarshal(payload, &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			}
		}
		return nil, fmt.Errorf("unhandled bridge endpoint: %s: %+v", endpoint, string(payload))
	default:
		if len(endpoint) == 0 {
			m := ZigbeeMessage{}
			err := json.Unmarshal(payload, &m)
			if err != nil {
				return nil, err
			}
			return m, nil
		}
	}

	return nil, nil
}

func ReadMessage(objectID string, payload []byte, endpoint ...string) (interface{}, error) {
	log.WithFields(log.Fields{
		"objectID": objectID,
		"endpoint": endpoint,
		"payload":  string(payload),
	}).Trace("ReadMessage()")

	switch objectID {
	case "wifi":
		m := WifiMessage{}
		err := json.Unmarshal(payload, &m)
		if err != nil {
			log.Error(err)
		}
		return m, nil
	case "air":
		m := AirMessage{}
		err := json.Unmarshal(payload, &m)
		if err != nil {
			log.Error(err)
		}
		return m, nil
	case "water":
		m := WaterMessage{}
		err := json.Unmarshal(payload, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	case "led":
		if len(endpoint) > 0 {
			if endpoint[0] == "config" {
				m := LEDConfig{}
				err := json.Unmarshal(payload, &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			} else if endpoint[0] == "color" {
				m := LEDColor{}
				err := json.Unmarshal(payload, &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			}
			return nil, fmt.Errorf("unhandled led endpoint: %s : %+v", endpoint, string(payload))
		}
	}

	return nil, nil
}

func NewMQTTClient(cfg MQTTConfig) (mqtt.Client, error) {
	var mqttClient mqtt.Client

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(cfg.URL)
	mqttOpts.SetCleanSession(true)

	if cfg.Username != "" && cfg.Password != "" {
		mqttOpts.Username = cfg.Username
		mqttOpts.Password = cfg.Password
	}

	mqttClient = mqtt.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	} else {
		log.WithFields(log.Fields{
			"url": cfg.URL,
		}).Debug("connected to MQTT")
	}

	return mqttClient, nil
}

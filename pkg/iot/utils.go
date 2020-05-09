package iot

import (
	"encoding/json"
	"regexp"
	"strings"

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

func ReadMessage(endpoint string, payload []byte) interface{} {

	switch endpoint {
	// case "air":
	case "wifi":
		m := WifiMessage{}
		err := json.Unmarshal(payload, &m)
		if err != nil {
			log.Error(err)
		}
		return m
	case "air":
		m := AirMessage{}
		err := json.Unmarshal(payload, &m)
		if err != nil {
			log.Error(err)
		}
		return m
	}

	return nil
}

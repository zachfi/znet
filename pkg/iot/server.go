package iot

import (
	context "context"
	"encoding/json"
	sync "sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Server is the struct to implement the IOTServer.
type Server struct {
	UnimplementedIOTServer
	mutex      sync.Mutex
	mqttClient mqtt.Client
}

// NewServer receives an mqttClient to return a new Server.
func NewServer(mqttClient mqtt.Client) (*Server, error) {
	return &Server{
		mqttClient: mqttClient,
		mutex:      sync.Mutex{},
	}, nil
}

// UpdateDevice implements the OTA update for zigbee devices.
//
// https://www.zigbee2mqtt.io/information/ota_updates.html#update-to-latest-firmware
func (s *Server) UpdateDevice(ctx context.Context, req *UpdateRequest) (*Empty, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	topic := "zigbee2mqtt/bridge/request/device/ota_update/update"
	message := map[string]interface{}{
		"id": req.Device,
	}

	m, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	s.mqttClient.Publish(topic, byte(0), false, string(m))

	time.Sleep(10 * time.Minute)

	return &Empty{}, nil
}

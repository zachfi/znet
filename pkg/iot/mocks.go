package iot

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/require"
)

func TestMqttToken_interface(t *testing.T) {
	var token mqtt.Token = &MockToken{}
	require.NotNil(t, token)
}

type MockToken struct {
	ready    bool
	err      error
	complete chan struct{}
}

func (mt *MockToken) Wait() bool {
	return mt.ready
}

func (mt *MockToken) WaitTimeout(time.Duration) bool {
	return mt.ready
}

func (mt *MockToken) Error() error {
	return mt.err
}

func (mt *MockToken) Done() <-chan struct{} {
	return mt.complete
}

func TestMqttClient_interface(t *testing.T) {
	var client mqtt.Client = &MockClient{}
	require.NotNil(t, client)
}

type MockClient struct {
	isConnected bool
	token       mqtt.Token
}

func (mc *MockClient) IsConnected() bool {
	return mc.isConnected
}

func (mc *MockClient) IsConnectionOpen() bool {
	return mc.isConnected
}

func (mc *MockClient) Connect() mqtt.Token {
	return mc.token
}

func (mc *MockClient) Disconnect(quiesce uint) {}

func (mc *MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return mc.token
}

func (mc *MockClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return mc.token
}

func (mc *MockClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return mc.token
}

func (mc *MockClient) Unsubscribe(topics ...string) mqtt.Token {
	return mc.token
}

func (mc *MockClient) AddRoute(topic string, handler mqtt.MessageHandler) {}

func (mc *MockClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

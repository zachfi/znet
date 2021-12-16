package iot

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func NewMQTTClient(cfg MQTTConfig, logger log.Logger) (mqtt.Client, error) {
	var mqttClient mqtt.Client

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(cfg.URL)
	mqttOpts.SetCleanSession(true)
	mqttOpts.SetAutoReconnect(true)
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetConnectRetryInterval(10 * time.Second)

	if cfg.Username != "" && cfg.Password != "" {
		mqttOpts.SetUsername(cfg.Username)
		mqttOpts.SetPassword(cfg.Password)
	}

	mqttClient = mqtt.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		_ = level.Error(logger).Log("err", token.Error())
	} else {
		_ = level.Debug(logger).Log("msg", "mqtt connected", "url", cfg.URL)
	}

	return mqttClient, nil
}

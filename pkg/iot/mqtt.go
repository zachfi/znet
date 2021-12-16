package iot

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func NewMQTTClient(cfg MQTTConfig, logger log.Logger) (mqtt.Client, error) {
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
		_ = level.Error(logger).Log("err", token.Error())
	} else {
		_ = level.Debug(logger).Log("msg", "mqtt connected", "url", cfg.URL)
	}

	return mqttClient, nil
}

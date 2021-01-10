package cmd

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xaque208/znet/internal/config"
)

func mqttConnect(cfg config.MQTTConfig) mqtt.Client {
	var mqttClient mqtt.Client

	viper.SetDefault("mqtt.url", "tcp://localhost:1883")

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

	return mqttClient
}

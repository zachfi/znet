package cmd

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xaque208/znet/znet"
)

func mqttConnect(config znet.MQTTConfig) mqtt.Client {
	var mqttClient mqtt.Client

	viper.SetDefault("mqtt.url", "tcp://localhost:1883")

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(config.URL)
	mqttOpts.SetCleanSession(true)

	if config.Username != "" && config.Password != "" {
		mqttOpts.Username = config.Username
		mqttOpts.Password = config.Password
	}

	mqttClient = mqtt.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	} else {
		log.WithFields(log.Fields{
			"url": config.URL,
		}).Debug("connected to MQTT")
	}

	return mqttClient
}

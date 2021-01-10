package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/iot"
	"github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var harvestCmd = &cobra.Command{
	Use:     "harvest",
	Short:   "Run an mqtt harvester",
	Long:    "Subscribe to MQTT topics and parse out the data",
	Example: "znet harves",
	Run:     runHarvest,
}

func init() {
	rootCmd.AddCommand(harvestCmd)
}

func runHarvest(cmd *cobra.Command, args []string) {
	formatter := log.TextFormatter{
		FullTimestamp: true,
	}

	log.SetFormatter(&formatter)
	if trace {
		log.SetLevel(log.TraceLevel)
	} else if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		done <- true
	}()

	telemetryClient := rpc.NewTelemetryClient(conn)

	var onMessageReceived mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		topicPath, err := iot.ParseTopicPath(msg.Topic())
		if err != nil {
			log.Error(errors.Wrap(err, "failed to to parse topic path"))
			return
		}

		discovery := &rpc.DeviceDiscovery{
			Component: topicPath.Component,
			NodeId:    topicPath.NodeID,
			ObjectId:  topicPath.ObjectID,
			Endpoint:  topicPath.Endpoint,
			Message:   msg.Payload(),
		}

		iotDevice := &rpc.IOTDevice{
			DeviceDiscovery: discovery,
		}

		_, err = telemetryClient.ReportIOTDevice(context.Background(), iotDevice)
		if err != nil {
			log.Error(err)
		}
	}

	viper.SetDefault("mqtt.url", "tcp://localhost:1883")
	viper.SetDefault("mqtt.topic", "#")

	mqttURL := viper.GetString("mqtt.url")
	mqttTopic := viper.GetString("mqtt.topic")
	mqttUsername := viper.GetString("mqtt.username")
	mqttPassword := viper.GetString("mqtt.password")

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(mqttURL)
	mqttOpts.SetCleanSession(true)
	mqttOpts.OnConnect = func(c mqtt.Client) {
		token := c.Subscribe(mqttTopic, 0, onMessageReceived)
		token.Wait()
		if token.Error() != nil {
			log.Error(token.Error())
		}
	}

	if mqttUsername != "" && mqttPassword != "" {
		mqttOpts.Username = mqttUsername
		mqttOpts.Password = mqttPassword
	}

	mqttClient := mqtt.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	} else {
		log.Debugf("connected to MQTT: %s", mqttURL)
	}

	log.Debugf("HTTP listening: %s", listenAddr)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		if token := mqttClient.Unsubscribe(mqttTopic); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		}

		mqttClient.Disconnect(250)

		done <- true
	}()

	<-done
}

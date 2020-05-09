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
	"google.golang.org/grpc"

	"github.com/xaque208/znet/pkg/iot"
	pb "github.com/xaque208/znet/rpc"
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

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	thingClient := pb.NewThingsClient(conn)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var onMessageReceived mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		topicPath, err := iot.ParseTopicPath(msg.Topic())
		if err != nil {
			log.Error(errors.Wrap(err, "failed to to parse topic path"))
			return
		}

		log.Infof("have topicPath: %+v", topicPath)

		m := &pb.DeviceDiscovery{
			Component: topicPath.Component,
			NodeID:    topicPath.NodeID,
			ObjectID:  topicPath.ObjectID,
			Endpoint:  topicPath.Endpoint,
			Message:   msg.Payload(),
		}

		_, err = thingClient.Notice(context.Background(), m)
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
		panic(token.Error())
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

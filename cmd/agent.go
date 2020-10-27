package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/pkg/iot"
	"github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   "Run a znet agent",
	Long:    "Subscribe to RPC events to perform actions",
	Example: "znet agent",
	Run:     runAgent,
}

func init() {
	rootCmd.AddCommand(agentCmd)
}

func runAgent(cmd *cobra.Command, args []string) {
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

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	var mqttClient mqtt.Client

	if z.Config.MQTT != nil {
		viper.SetDefault("mqtt.url", "tcp://localhost:1883")

		mqttURL := viper.GetString("mqtt.url")
		mqttUsername := viper.GetString("mqtt.username")
		mqttPassword := viper.GetString("mqtt.password")

		mqttOpts := mqtt.NewClientOptions()
		mqttOpts.AddBroker(mqttURL)
		mqttOpts.SetCleanSession(true)

		if mqttUsername != "" && mqttPassword != "" {
			mqttOpts.Username = mqttUsername
			mqttOpts.Password = mqttPassword
		}

		mqttClient = mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		} else {
			log.Debugf("connected to MQTT: %s", mqttURL)
		}
	}

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)

	if z.Config.Agent == nil {
		log.Fatal("unable to create agent with nil Agent configuration")
	}

	consumers := []events.Consumer{}

	var agentConsumer *agent.Agent
	if z.Config.Agent != nil {
		agentConsumer = agent.NewAgent(*z.Config.Agent, conn)
		consumers = append(consumers, agentConsumer)
	}

	eventNames := []string{}
	inventoryClient := rpc.NewInventoryClient(conn)

	if z.Config.Lights != nil && inventoryClient != nil && mqttClient != nil {

		lightsConsumer := lights.NewLights(*z.Config.Lights, inventoryClient, mqttClient)
		consumers = append(consumers, lightsConsumer)

		eventNames = append(eventNames, timer.EventNames...)
		eventNames = append(eventNames, astro.EventNames...)
		eventNames = append(eventNames, iot.EventNames...)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	machine, err := eventmachine.New(ctx, &consumers)
	if err != nil {
		log.Error(err)
	}

	client := rpc.NewEventsClient(conn)

	for _, exec := range z.Config.Agent.Executions {
		eventNames = append(eventNames, exec.Events...)
	}

	eventSub := &rpc.EventSub{
		EventNames: eventNames,
	}

	go machine.ReadStream(client, eventSub)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Debugf("caught signal: %s", sig.String())

		done <- true
	}()

	<-done

	log.Debug("closing RPC connection")
	err = conn.Close()
	if err != nil {
		log.Error(err)
	}

	log.Debug("stopping event machine")
	err = machine.Stop()
	if err != nil {
		log.Error(err)
	}
}

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
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/network"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
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
	initLogger()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatalf("failed to create znet: %s", err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	var mqttClient mqtt.Client

	if z.Config.MQTT != nil {
		mqttClient = mqttConnect(*z.Config.MQTT)
	}

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	if z.Config.Agent == nil {
		log.Fatal("unable to create agent with nil Agent configuration")
	}

	consumers := []events.Consumer{}

	var agentServer *agent.Agent
	if z.Config.Agent != nil {
		agentServer = agent.NewAgent(z.Config, conn)
		consumers = append(consumers, agentServer)
	}

	eventNames := []string{}
	inventoryClient := inventory.NewInventoryClient(conn)

	// Configure the Lights consumer if we have...
	// mqttClient for publishing messages
	// inventoryClient for device lookup and group selection
	if z.Config.Lights != nil && inventoryClient != nil && mqttClient != nil {
		lightsConsumer, err := lights.NewLights(z.Config)
		if err != nil {
			log.Fatal(err)
		}
		consumers = append(consumers, lightsConsumer)

		// The lightsConsumer responds to a bunch of events.
		eventNames = append(eventNames, timer.EventNames...)
		// eventNames = append(eventNames, astro.EventNames...)
		// eventNames = append(eventNames, iot.EventNames...)
	}

	if z.Config.Network != nil && inventoryClient != nil {
		networkConsumer := network.NewNetwork(z.Config.Network, conn)

		consumers = append(consumers, networkConsumer)
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

	err = agentServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Debugf("caught signal: %s", sig.String())

		done <- true
	}()

	<-done

	log.Debug("stopping event machine")
	err = machine.Stop()
	if err != nil {
		log.Error(err)
	}

	log.Debug("terminating RPC server")
	err = agentServer.Stop()
	if err != nil {
		log.Error(err)
	}

	log.Debug("closing RPC client connection")
	err = conn.Close()
	if err != nil {
		log.Error(err)
	}

}

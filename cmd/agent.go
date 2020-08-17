package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
	pb "github.com/xaque208/znet/rpc"
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

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)

	if z.Config.Agent == nil {
		log.Fatal("unable to create agent with nil Agent configuration")
	}

	ag := agent.NewAgent(*z.Config.Agent, conn)

	consumers := []events.Consumer{
		ag,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	machine, err := eventmachine.New(ctx, consumers)
	if err != nil {
		log.Error(err)
	}

	client := pb.NewEventsClient(conn)

	eventSub := &pb.EventSub{
		Name: ag.EventNames(),
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

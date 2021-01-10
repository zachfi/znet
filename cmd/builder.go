package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/builder"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var builderCmd = &cobra.Command{
	Use:     "builder",
	Short:   "Run a git repo builder",
	Long:    "Listen for NewCommit and NewTag events, and bulid repos using that event information",
	Example: "znet builder",
	Run:     runBuilder,
}

func init() {
	rootCmd.AddCommand(builderCmd)
}

func runBuilder(cmd *cobra.Command, args []string) {
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

	if z.Config.Builder == nil {
		log.Fatal("unable to create agent with nil Builder configuration")
	}

	x := builder.NewBuilder(conn, *z.Config.Builder)

	consumers := []events.Consumer{
		x,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	machine, err := eventmachine.New(ctx, &consumers)
	if err != nil {
		log.Error(err)
	}

	client := rpc.NewEventsClient(conn)

	eventNames := []string{}
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
		log.Warnf("caught signal: %s", sig.String())

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

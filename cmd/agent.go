package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/pkg/eventmachine"
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

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	ag := agent.NewAgent(z.Config.Agent, conn)

	consumers := []events.Consumer{
		ag,
	}

	machine, err := eventmachine.New(consumers)
	if err != nil {
		log.Error(err)
	}

	client := pb.NewEventsClient(conn)

	ctx, cancel := context.WithCancel(context.Background())

	eventSub := &pb.EventSub{
		Name: ag.EventNames(),
	}

	// Run the receiver forever.
	go func() {
		for {
			stream, receiverErr := client.SubscribeEvents(ctx, eventSub)
			if receiverErr == grpc.ErrClientConnClosing {
				return
			}

			if receiverErr != nil {
				log.Errorf("calling %+v.SubscribeEvents(_) = _, %+v", client, receiverErr)

				time.Sleep(10 * time.Second)
				continue
			}

			var ev *pb.Event

			ev, err = stream.Recv()
			if status.Code(err) != codes.OK {
				log.Errorf("stream.Recv() = %v, %v; want _, status.Code(err)=%v", ev, err, codes.OK)
				continue
			}

			if err != nil {
				log.Errorf("%v.SubscribeEvents(_) = _, %v", client, err)
				time.Sleep(10 * time.Second)
				continue
			}

			log.Tracef("received event: %+v", ev)

			evE := events.Event{
				Name:    ev.Name,
				Payload: ev.Payload,
			}

			machine.EventChannel <- evE
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

	<-done

	cancel()

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

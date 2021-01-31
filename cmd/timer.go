package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/znet"
)

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Run a timer",
	Long: `Run a timer daemon to send events to the gRPC server when the timers expire.

Several flavors of a timer are implemented.

* astro timers send an AstroEvent based on data read from openweathermap_exporter
* repeating timers send a NamedTimer event every interval
* static timers send a NamedTimer event at a specific time of day, on specific days
	`,
	Example: "znet timer -v --config timer.yaml",
	Run:     runTimer,
}

func init() {
	rootCmd.AddCommand(timerCmd)
}

func runTimer(cmd *cobra.Command, args []string) {
	initLogger()

	viper.SetDefault("timer.future_limit", 1000)

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.Timer.FutureLimit = viper.GetInt("timer.future_limit")
	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	if z.Config.Astro == nil {
		log.Fatal("unable to create EventProducer with nil Astro configuration")
	}

	producers := make([]events.Producer, 0)

	y, err := astro.NewProducer(z.Config.Astro)
	if err != nil {
		log.Error(err)
	} else {
		producers = append(producers, y)
	}

	if z.Config.Timer == nil {
		log.Fatal("unable to create EventProducer with nil Timer configuration")
	}

	if z.Config.Lights == nil {
		log.Fatal("unable to create EventProducer with nil Lights configuration")
	}

	x, err := timer.NewProducer(z.Config.Timer)
	if err != nil {
		log.Error(err)
	} else {
		producers = append(producers, x)
	}

	log.WithFields(log.Fields{
		"count": len(producers),
	}).Debug("starting producers")

	for _, p := range producers {
		err = p.Connect(ctx, conn)
		if err != nil {
			log.Error(err)
		}
	}

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
}

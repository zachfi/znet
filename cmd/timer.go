package cmd

import (
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
	Example: "znet timer -v --config ~/.timer.yaml",
	Run:     runTimer,
}

func init() {
	rootCmd.AddCommand(timerCmd)
}

func runTimer(cmd *cobra.Command, args []string) {
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

	log.Infof("%s starting", Version)

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

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	if z.Config.Astro == nil {
		log.Fatal("unable to create EventProducer with nil Astro configuration")
	}

	producers := make([]events.Producer, 0)

	y := astro.NewProducer(conn, *z.Config.Astro)
	producers = append(producers, y)

	if z.Config.Timer == nil {
		log.Fatal("unable to create EventProducer with nil Timer configuration")
	}

	x := timer.NewProducer(conn, *z.Config.Timer)
	producers = append(producers, x)

	for _, p := range producers {
		err = p.Start()
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

	for _, p := range producers {
		err = p.Stop()
		if err != nil {
			log.Error(err)
		}
	}
}

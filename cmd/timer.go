package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/znet"
)

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Run a timer",
	Run:   runTimer,
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

	viper.SetDefault("timer.future_limit", 1000)

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")
	z.Config.Timer.FutureLimit = viper.GetInt("timer.future_limit")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	xConn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err = xConn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	yConn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err = yConn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	y := astro.NewProducer(yConn, z.Config.Astro)
	err = y.Start()
	if err != nil {
		log.Error(err)
	}

	x := timer.NewProducer(xConn, z.Config.Timer)
	err = x.Start()
	if err != nil {
		log.Error(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())
		done <- true
	}()

	<-done

	err = y.Stop()
	if err != nil {
		log.Error(err)
	}

	err = x.Stop()
	if err != nil {
		log.Error(err)
	}
}

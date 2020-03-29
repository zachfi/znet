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
	if verbose {
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

	xConn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}
	defer xConn.Close()

	yConn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}
	defer yConn.Close()

	y := astro.NewProducer(yConn, z.Config.Astro)
	y.Start()

	x := timer.NewProducer(xConn, z.Config.Timer)
	x.Start()

	// log.Infof("x: %+v", x)
	// log.Infof("y: %+v", y)
	//

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())
		done <- true
	}()

	<-done
	y.Stop()
	x.Stop()
}

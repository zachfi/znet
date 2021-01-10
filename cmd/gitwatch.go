package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/znet"
)

var gitwatchCmd = &cobra.Command{
	Use:     "gitwatch",
	Short:   "Run a git watcher",
	Long:    "Watch git repos and emit events when changes are noticed.",
	Example: "znet gitwatch",
	Run:     runGitwatch,
}

func init() {
	rootCmd.AddCommand(gitwatchCmd)
}

func runGitwatch(cmd *cobra.Command, args []string) {
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

	if z.Config.GitWatch == nil {
		log.Fatal("unable to create agent with nil GitWatch configuration")
	}

	x := gitwatch.NewProducer(conn, *z.Config.GitWatch)
	err = x.Start()
	if err != nil {
		log.Error(err)
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

	err = x.Stop()
	if err != nil {
		log.Error(err)
	}

	log.Debug("closing RPC connection")
	err = conn.Close()
	if err != nil {
		log.Error(err)
	}
}

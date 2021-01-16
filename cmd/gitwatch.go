package cmd

import (
	"context"
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
	initLogger()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	if z.Config.GitWatch == nil {
		log.Fatal("unable to create agent with nil GitWatch configuration")
	}

	x := gitwatch.NewProducer(z.Config.GitWatch)

	err = x.Connect(ctx, conn)
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

	cancel()

	log.Debug("closing RPC connection")
	err = conn.Close()
	if err != nil {
		log.Error(err)
	}
}

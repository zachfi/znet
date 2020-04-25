// Copyright © 2018 Zach Leslie <code@zleslie.info>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/znet"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for commands/events/messages sent to the RPC server",
	Long: `

`,
	Example: "znet listen -v --trace",
	Run:     listen,
}

var listenAddr string
var rpcListenAddr string

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", ":9100", "Specify HTTP listen address")
	listenCmd.PersistentFlags().StringVarP(&rpcListenAddr, "rpc", "r", ":8800", "Specify RPC listen address")
}

func listen(cmd *cobra.Command, args []string) {
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

	// Handle environment variables
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("nats.url", nats.DefaultURL)
	viper.SetDefault("nats.topic", "things")

	viper.AutomaticEnv()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.Nats.URL = viper.GetString("nats.url")
	z.Config.Nats.Topic = viper.GetString("nats.topic")
	z.Config.RPC.ListenAddress = viper.GetString("rpc.listen")
	z.Config.HTTP.ListenAddress = viper.GetString("http.listen")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		done <- true
	}()

	consumers := []events.Consumer{
		z.Lights,
		z,
	}

	znetServer := znet.NewServer(z.Config.HTTP, z.Config.RPC, consumers)
	err = znetServer.Start(z)
	if err != nil {
		log.Error(err)
	}

	<-done

	err = znetServer.Stop()
	if err != nil {
		log.Error(err)
	}

	err = z.Stop()
	if err != nil {
		log.Error(err)
	}
}

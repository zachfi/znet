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
	"github.com/xaque208/znet/znet"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for commands/events/messages on the bus",
	Long: `Messages issued sent to the bus might be events, messages, or command requests.  The listen command here subscribes to a topic and handles what actions need to be taken.

`,
	Run: listen,
}

var listenAddr string

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", ":9100", "Specify listen address")

}

func listen(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle environment variables
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("nats.url", nats.DefaultURL)
	viper.SetDefault("nats.topic", "things")
	viper.SetDefault("redis.host", "localhost")

	viper.AutomaticEnv()

	z := znet.Znet{}
	z.LoadConfig(cfgFile)

	z.Config.Nats.URL = viper.GetString("nats.url")
	z.Config.Nats.Topic = viper.GetString("nats.topic")

	go func() {
		sig := <-sigs
		log.Warnf("Caught signal: %s", sig.String())
		done <- true
	}()

	z.Listen(listenAddr, done)
}

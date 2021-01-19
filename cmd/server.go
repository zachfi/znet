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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/znet"
)

// serverCmd represents the listen command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Listen for commands/events/messages sent to the RPC server",
	Long: `

`,
	Example: "znet server -v --trace",
	Run:     server,
}

var listenAddr string
var rpcListenAddr string

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", ":9100", "Specify HTTP listen address")
	serverCmd.PersistentFlags().StringVarP(&rpcListenAddr, "rpc", "r", ":8800", "Specify RPC listen address")
}

func server(cmd *cobra.Command, args []string) {
	initLogger()

	// Handle environment variables
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("http.listen_address", listenAddr)
	viper.SetDefault("rpc.listen_address", rpcListenAddr)

	if z.Config.RPC == nil {
		z.Config.RPC = &config.RPCConfig{}
	}

	z.Config.RPC.ListenAddress = viper.GetString("rpc.listen_address")

	if z.Config.HTTP == nil {
		z.Config.HTTP = &config.HTTPConfig{}
	}

	z.Config.HTTP.ListenAddress = viper.GetString("http.listen_address")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		done <- true
	}()

	server := znet.NewServer(z.Config)
	err = server.Start(z)
	if err != nil {
		log.Error(err)
	}

	<-done

	err = server.Stop()
	if err != nil {
		log.Error(err)
	}
}

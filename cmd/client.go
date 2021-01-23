// Copyright © 2021 Zach Leslie <code@zleslie.info>
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
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/znet"
)

// clientCmd represents the on command
var clientCmd = &cobra.Command{
	Use:     "client",
	Short:   "Issue calls to the RPC agent for testing and manual triggering.",
	Long:    "A test client for getting the agent RPC up and running",
	Example: "znet client",
	Run:     client,
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func client(cmd *cobra.Command, args []string) {
	initLogger()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Error(err)
	}

	// Connect to the agent locally using the agent_listen_address
	z.Config.RPC.ServerAddress = viper.GetString("rpc.agent_listen_address")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	defer func() {
		log.Debug("closing RPC client connection")
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	a := agent.NewNodeClient(conn)

	res, err := a.RunPuppetAgent(context.Background(), &agent.Empty{})
	if err != nil {
		log.Error(err)
	}

	if len(res.Error) > 0 {
		log.Errorf("error: %s", string(res.Error))
	}

}

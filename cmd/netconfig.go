// Copyright © 2020 Zach Leslie <code@zleslie.info>
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

	"github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xaque208/znet/pkg/netconfig"
	"github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var commit bool
var diff bool
var limit int
var confirm int

// netconfigCmd represents the netconfig command
var netconfigCmd = &cobra.Command{
	Use:     "netconfig",
	Short:   "Configure Junos Devices",
	Long:    "Configure the network",
	Example: "znet netconfig",
	Run:     runNetconfig,
}

func init() {
	rootCmd.AddCommand(netconfigCmd)

	netconfigCmd.PersistentFlags().BoolVarP(&commit, "commit", "", false, "Commit the configuration")
	netconfigCmd.PersistentFlags().BoolVarP(&diff, "diff", "d", false, "Show the rendered templates")
	netconfigCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Limit the number of devices to configure")
	netconfigCmd.PersistentFlags().IntVarP(&confirm, "confirm", "", 0, "Number of minutes at which the config will be rolled back")
}

func runNetconfig(cmd *cobra.Command, args []string) {
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

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	// Load the network data.
	configDir := viper.GetString("netconfig.configdir")

	inventoryClient := rpc.NewInventoryClient(conn)

	ctx := context.Background()
	stream, err := inventoryClient.ListNetworkHosts(ctx, &rpc.Empty{})
	if err != nil {
		log.Error(err)
	}

	hosts := []*rpc.NetworkHost{}

	for {
		var d *rpc.NetworkHost

		d, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				log.Error(err)
				break
			}
		}

		hosts = append(hosts, d)
	}

	if len(hosts) == 0 {
		log.Fatalf("zero hosts to configure")
	}

	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	nc, err := netconfig.NewNetConfig(configDir, hosts, auth, z.Environment)
	if err != nil {
		log.Fatal(err)
	}

	err = nc.ConfigureNetwork(commit, confirm, diff)
	if err != nil {
		log.Error(err)
	}

}

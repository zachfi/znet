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
	"io"
	"time"

	"github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/netconfig"
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

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	// Load the network data.
	configDir := viper.GetString("netconfig.configdir")

	inventoryClient := inventory.NewInventoryClient(conn)

	var stream inventory.Inventory_ListNetworkHostsClient

	ctx := context.Background()

	for {
		stream, err = inventoryClient.ListNetworkHosts(ctx, &inventory.Empty{})
		if err != nil {
			switch status.Code(err) {
			case codes.Unavailable:
				time.Sleep(3 * time.Second)
				continue
			default:
				log.Error(err)
			}
		}
		break
	}

	hosts := []*inventory.NetworkHost{}

	for {
		var d *inventory.NetworkHost

		d, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			switch status.Code(err) {
			case codes.OK:
				log.Info("ok")
				continue
			default:
				log.Error(err)
			}
		}

		hosts = append(hosts, d)
	}

	if len(hosts) == 0 {
		log.Fatalf("zero hosts to configure")
	}

	auth := &junos.AuthMethod{
		Username:   viper.GetString("network.junos.username"),
		PrivateKey: viper.GetString("network.junos.keyfile"),
	}

	log.Debugf("auth: %+v", auth)

	nc, err := netconfig.NewNetConfig(configDir, hosts, auth, z.Environment)
	if err != nil {
		log.Fatal(err)
	}

	err = nc.ConfigureNetwork(commit, confirm, diff)
	if err != nil {
		log.Error(err)
	}

}

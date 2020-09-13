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
	"context"
	"os"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pb "github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

// invCmd represents the inv command
var invCmd = &cobra.Command{
	Use:     "inv",
	Short:   "Report on inventory",
	Long:    "Run an inventory report",
	Example: "znet inv",
	Run:     runInv,
}

var adopt string

func init() {
	rootCmd.AddCommand(invCmd)

	invCmd.PersistentFlags().StringVarP(&rpcServer, "rpc", "r", ":8800", "Specify RPC server address")
	invCmd.Flags().BoolP("verbose", "v", false, "Raise verbosity")

	invCmd.PersistentFlags().StringVarP(&adopt, "adopt", "a", "", "Adopt an unknown host by MAC address")
}

func runInv(cmd *cobra.Command, args []string) {
	formatter := log.TextFormatter{
		DisableQuote:     true,
		DisableTimestamp: true,
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
		log.Error(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}

		z.Stop()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryClient := pb.NewInventoryClient(conn)

	resp, err := inventoryClient.Search(ctx, &pb.SearchRequest{})
	if err != nil {
		log.Error(err)
	}

	// resp, err := inventoryClient.ListNetworkHosts(context.Background(), &pb.Empty{})
	// if err != nil {
	// 	log.Error(err)
	// }

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Platform", "Type", "Description"})

	for _, h := range resp.Hosts {
		t.AppendRow([]interface{}{
			h.Name,
			h.Platform,
			h.Type,
			h.Description,
		})
	}

	t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()

	// t = table.NewWriter()
	// t.SetOutputMirror(os.Stdout)
	// t.AppendHeader(table.Row{"Name", "IP", "MAC"})

	// for _, h := range res.UnknownHosts {
	// 	t.AppendRow([]interface{}{
	// 		h.Name,
	// 		h.Ip,
	// 		h.Mac,
	// 	})
	// }

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	// t.Render()

	// if adopt != "" {
	// 	for _, h := range res.UnknownHosts {
	// 		if strings.EqualFold(h.Mac, adopt) {
	// 			x := inventory.UnknownHost{
	// 				Name:       h.Name,
	// 				MACAddress: h.Mac,
	// 				IP:         h.Ip,
	// 			}
	// 			z.Inventory.AdoptUnknownHost(x, "cn=new,ou=network,dc=znet")
	// 		}
	// 	}
	// }
	//

	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		// "DateCode",
		"Description",
		"Dn",
		"IotZone",
		"LastSeen",
		// "ManufacturerName",
		"Model",
		"ModelId",
		"Name",
		"PowerSource",
		"SoftwareBuildId",
		"Type",
		"Vendor",
	})

	for _, h := range resp.ZigbeeDevices {
		t.AppendRow([]interface{}{
			// h.DateCode,
			h.Description,
			h.Dn,
			h.IotZone,
			h.LastSeen,
			// h.ManufacturerName,
			h.Model,
			h.ModelId,
			h.Name,
			h.PowerSource,
			h.SoftwareBuildId,
			h.Type,
			h.Vendor,
		})
	}

	t.Render()

}

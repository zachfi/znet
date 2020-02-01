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
	"google.golang.org/grpc"
)

// lightsCmd represents the lights command
var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Send an on/off command to the bus for a given room's lights.",
	Long:  "",
	Run:   runLights,
}

var roomName string
var roomBrightness uint8

func init() {
	rootCmd.AddCommand(lightsCmd)

	lightsCmd.PersistentFlags().StringVarP(&roomName, "room", "r", "", "Specify a configured room")
	lightsCmd.PersistentFlags().Uint8VarP(&roomBrightness, "brightness", "b", 254, "Set the brightness of the room")
	lightsCmd.Flags().BoolP("verbose", "v", false, "Raise verbosity")
}

func runLights(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Error(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	lc := pb.NewLightsClient(conn)

	req := &pb.LightRequest{}

	res, err := lc.Status(context.Background(), req)
	if err != nil {
		log.Errorf("RPC Error: %s Response: %+v", err, res)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Id", "Type", "Brightness", "On"})

	for _, g := range res.Groups {
		t.AppendRow([]interface{}{
			g.Name,
			g.Id,
			g.Type,
			g.State.Brightness,
			g.State.On,
		})
	}

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()

}

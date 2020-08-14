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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pb "github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn all the lights in a configured room on/off",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: off,
}

func init() {
	lightsCmd.AddCommand(offCmd)
}

func off(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Error(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(z.Config.RPC.ServerAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	lc := pb.NewLightsClient(conn)

	req := &pb.LightGroup{
		Name: roomName,
	}

	res, err := lc.Off(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	log.Infof("RPC Response: %+v", res)

}

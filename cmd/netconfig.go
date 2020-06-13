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
	"sync"

	"github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/znet"
)

var commit bool
var show bool
var limit int
var confirm int

// invCmd represents the inv command
var netconfigCmd = &cobra.Command{
	Use:     "netconfig",
	Short:   "Configure Junos Devices",
	Long:    "Run an inventory report",
	Example: "znet inv",
	Run:     runNetconfig,
}

func init() {
	rootCmd.AddCommand(netconfigCmd)

	netconfigCmd.PersistentFlags().BoolVarP(&commit, "commit", "", false, "Commit the configuration")
	netconfigCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	netconfigCmd.PersistentFlags().BoolVarP(&show, "show", "s", false, "Show the rendered templates")
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

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	// Load the network data.
	configDir := viper.GetString("netconfig.configdir")
	z.ConfigDir = configDir
	z.LoadData(configDir)

	hosts, err := z.Inventory.NetworkHosts()
	if err != nil {
		log.Error(err)
	}

	if len(hosts) == 0 {
		log.Fatalf("no hosts.")
	}

	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	wg := sync.WaitGroup{}

	for _, host := range hosts {
		wg.Add(1)
		go func(h inventory.NetworkHost) {

			if h.Platform == "junos" {
				log.Debugf("configuring network host: %+v", h.HostName)

				err = z.ConfigureNetworkHost(&h, commit, confirm, auth, show)
				if err != nil {
					log.Error(err)
				}
			}

			wg.Done()
		}(host)
	}

	wg.Wait()
}

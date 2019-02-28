// Copyright © 2019 Zach Leslie <code@zleslie.info>
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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/znet/znet"
)

// netconfigCmd represents the on command
var netconfigCmd = &cobra.Command{
	Use:   "netconfig",
	Short: "Configure Junos Devices",
	Long:  ``,
	Run:   netconfig,
}

var commit bool

func init() {
	rootCmd.AddCommand(netconfigCmd)

	netconfigCmd.PersistentFlags().BoolVarP(&commit, "commit", "", false, "Commit the configuration")
}

func netconfig(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	viper.SetDefault("netconfig.configdir", "etc/")
	viper.AutomaticEnv()

	z := znet.Znet{}
	z.LoadConfig(cfgFile)

	configDir := viper.GetString("netconfig.configdir")
	z.ConfigDir = configDir

	z.LoadData(configDir)

	l, err := z.NewLDAPClient(z.Config.Ldap)
	if err != nil {
		log.Error(err)
	}
	defer l.Close()

	// zones := z.GetNetworkZones(l, z.Config.Ldap.BaseDN)
	// log.Warnf("Zones: %+v", zones)

	hosts := z.GetNetworkHosts(l, z.Config.Ldap.BaseDN)
	// log.Warnf("Hosts: %+v", hosts)

	for _, host := range hosts {
		z.ConfigureNetworkHost(&host, commit)
	}

}

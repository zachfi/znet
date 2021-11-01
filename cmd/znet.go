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
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/znet/znet"
)

var cfgFile string
var verbose bool
var trace bool

// Version is the version of the project
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "znet",
	Short: "zNet",
	Long: `zNet 

Run zNet.
`,
	Run: runZnet,
}

func runZnet(cmd *cobra.Command, args []string) {
	initLogger()

	logger := newLogger()

	z, err := znet.NewZnet(cfgFile, logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create Znet", "err", err)
		os.Exit(1)
	}

	z.Run()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.znet.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "", false, "Trace level verbosity")

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

	rootCmd.AddCommand(inventoryCommand)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".znet" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".znet")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("using config file: %s", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	}
}

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
	"time"

	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerLogger "github.com/uber/jaeger-client-go/log"

	"github.com/go-kit/log/level"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/znet/pkg/util"
	"github.com/xaque208/znet/znet"
)

var (
	cfgFile         string
	target          string
	tracingEndpoint string

	verbose bool
	trace   bool
)

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
	logger := util.NewLogger()

	// jaegerCfg, err := jaegerConfig.FromEnv()
	// if err != nil {
	// 	level.Error(logger).Log("msg", "failed to get jaegerCfg from environment", "err", err)
	// 	os.Exit(1)
	// }

	jaegerCfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   "http://10.42.0.44:14268/api/traces",
		},
	}

	tracer, closer, err := jaegerCfg.New(
		"znet",
		jaegerConfig.Logger(jaegerLogger.StdLogger),
	)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create new tracer", "err", err)
		os.Exit(1)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	cfg, err := znet.LoadConfig(cfgFile)
	if err != nil {
		level.Error(logger).Log("msg", "failed to load config file", "err", err)
		os.Exit(1)
	}

	z, err := znet.New(cfg)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create Znet", "err", err)
		os.Exit(1)
	}

	if err := z.Run(); err != nil {
		level.Error(logger).Log("msg", "error running zNet", "err", err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version

	logger := util.NewLogger()

	if err := rootCmd.Execute(); err != nil {
		level.Error(logger).Log("msg", "failed to execute", "err", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.znet.yaml)")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "all", "Run a specific module")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "", false, "Trace level verbosity")
	rootCmd.PersistentFlags().StringVarP(&tracingEndpoint, "tracing-endpoint", "", "", "Jaeger reporter endpoint URL")

	rootCmd.MarkFlagRequired("tracing-endpoint")

	rootCmd.AddCommand(inventoryCommand)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logger := util.NewLogger()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			level.Error(logger).Log("msg", "failed to get homedir", "err", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".znet" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".znet")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		level.Debug(logger).Log("msg", "using config file", "file", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	}
}

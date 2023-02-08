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

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/flagext"

	yaml "gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"

	t "github.com/zachfi/zkit/pkg/tracing"

	"github.com/zachfi/znet/pkg/util"
	"github.com/zachfi/znet/znet"
)

const appName = "znet"

// Version is set via build flag -ldflags -X main.Version
var (
	Version  string
	Branch   string
	Revision string
)

func init() {
	version.Version = Version
	version.Branch = Branch
	version.Revision = Revision
	prometheus.MustRegister(version.NewCollector(appName))
}

func main() {
	logger := util.NewLogger()

	cfg, err := loadConfig()
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to load config file", "err", err)
		os.Exit(1)
	}

	shutdownTracer, err := t.InstallOpenTelemetryTracer(&cfg.Tracing, logger, appName, Version)
	if err != nil {
		_ = level.Error(logger).Log("msg", "error initialising tracer", "err", err)
		os.Exit(1)
	}
	defer shutdownTracer()

	z, err := znet.New(*cfg)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to create Znet", "err", err)
		os.Exit(1)
	}

	if err := z.Run(); err != nil {
		_ = level.Error(logger).Log("msg", "error running zNet", "err", err)
		os.Exit(1)
	}
}

func loadConfig() (*znet.Config, error) {
	const (
		configFileOption = "config.file"
	)

	var (
		configFile string
	)

	args := os.Args[1:]
	config := &znet.Config{}

	// first get the config file
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.StringVar(&configFile, configFileOption, "", "")

	// Try to find -config.file & -config.expand-env flags. As Parsing stops on the first error, eg. unknown flag,
	// we simply try remaining parameters until we find config flag, or there are no params left.
	// (ContinueOnError just means that flag.Parse doesn't call panic or os.Exit, but it returns error, which we ignore)
	for len(args) > 0 {
		_ = fs.Parse(args)
		args = args[1:]
	}

	// load config defaults and register flags
	config.RegisterFlagsAndApplyDefaults("", flag.CommandLine)

	// overlay with config file if provided
	if configFile != "" {
		buff, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read configFile %s: %w", configFile, err)
		}

		err = yaml.UnmarshalStrict(buff, config)
		if err != nil {
			return nil, fmt.Errorf("failed to parse configFile %s: %w", configFile, err)
		}
	}

	// overlay with cli
	flagext.IgnoredFlag(flag.CommandLine, configFileOption, "Configuration file to load")
	flag.Parse()

	return config, nil
}

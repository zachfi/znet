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
	"github.com/spf13/cobra"
)

// lightsCmd represents the lights command
var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "",
	Long:  "",
	Run:   run,
}

var roomName string

func init() {
	rootCmd.AddCommand(lightsCmd)

	lightsCmd.PersistentFlags().StringVarP(&roomName, "room", "r", "", "Specify a configured room")
	lightsCmd.Flags().BoolP("verbose", "v", false, "Raise verbosity")
}

func run(cmd *cobra.Command, args []string) {
}

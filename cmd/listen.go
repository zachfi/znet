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
	"time"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/things/things"
	"github.com/xaque208/znet/znet"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for commands/events/messages on the bus",
	Long: `Messages issued sent to the bus might be events, messages, or command requests.  The listen command here subscribes to a topic and handles what actions need to be taken.

`,
	Run: listen,
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

func lightsHandler(command things.Command) {

	roomName := command.Arguments["room"]
	state := command.Arguments["state"]

	if state != "on" && state != "off" {
		log.Errorf("Unknown light state received %s", state)
	}

	z, err := znet.LoadConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	r := rftoy.RFToy{Address: z.Endpoint}

	room, err := z.Room(roomName.(string))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Turning %s room %s", state, room.Name)
	for _, sid := range room.IDs {
		if state == "on" {
			r.On(sid)
		} else if state == "off" {
			r.Off(sid)
		}
		time.Sleep(2 * time.Second)
	}

}

func messageHandler(messages chan things.Message) {
	for {
		select {
		case msg := <-messages:
			log.Debugf("new message: %+v", msg)

			for _, c := range msg.Commands {
				if c.Name == "lights" {
					go lightsHandler(c)
				} else {
					log.Warnf("Unknown command %s", c.Name)
				}
			}

		}
	}
}

func listen(cmd *cobra.Command, args []string) {

	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	viper.SetDefault("nats.url", nats.DefaultURL)
	viper.SetDefault("nats.topic", "things")

	url := viper.GetString("nats.url")
	topic := viper.GetString("nats.topic")

	server, err := things.NewServer(url, topic)
	if err != nil {
		log.Error(err)
	}
	defer server.Close()
	messages := make(chan things.Message)
	server.Listen(messages, messageHandler)

}

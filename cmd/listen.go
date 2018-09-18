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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	nats "github.com/nats-io/go-nats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/things/things"
	"github.com/xaque208/znet/arpwatch"
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

var (
	macAddress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mac",
		Help: "MAC Address",
	}, []string{"mac", "ip"})
)

var listenAddr string

const (
	macsList  = "macs"
	macsTable = "mac:*"
)

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", ":9100", "Specify listen address")

	prometheus.MustRegister(macAddress)
}

func listen(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	messages := make(chan things.Message)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	viper.SetDefault("nats.url", nats.DefaultURL)
	viper.SetDefault("nats.topic", "things")

	url := viper.GetString("nats.url")
	topic := viper.GetString("nats.topic")

	log.Debug("Pre-reqs met")

	server, err := things.NewServer(url, topic)
	if err != nil {
		log.Error(err)
	}

	redisClient := arpwatch.NewRedisClient()

	log.Info("Listening to nats")
	go server.Listen(messages, messageHandler)

	log.Debug("Starting arpwatch")
	go arpWatch(redisClient)

	log.Debugf("HTTP listening on %s", listenAddr)
	srv := httpListen(listenAddr)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		log.Info("Closing thing server")
		server.Close()

		log.Info("Disconnecting redis")
		redisClient.Close()

		log.Info("HTTP shutting down")
		srv.Shutdown(nil)

		done <- true
	}()

	<-done

}

func lightsHandler(command things.Command) {

	roomName := command.Arguments["room"]
	state := command.Arguments["state"]

	if state != "on" && state != "off" {
		log.Errorf("Unknown light state received %s", state)
	}

	z, err := znet.LoadConfig(cfgFile)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Using RFToy at %s", z.Endpoint)
	r := rftoy.RFToy{Address: z.Endpoint}

	room, err := z.Room(roomName.(string))
	if err != nil {
		log.Error(err)
		return
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
			log.Debugf("New message: %+v", msg)

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

func httpListen(listenAddress string) *http.Server {
	srv := &http.Server{Addr: listenAddress}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	return srv
}

func arpWatch(redisClient *redis.Client) {

	hosts := viper.GetStringSlice("junos.hosts")
	if len(hosts) == 0 {
		log.Error("List of hosts required")
		return
	}

	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	aw := arpwatch.ArpWatch{
		Hosts: hosts,
		Auth:  auth,
	}

	go aw.Update()

	go func() {
		for {
			select {
			default:
				data, err := redisClient.SMembers(macsList).Result()
				if err != nil {
					log.Error(err)
				}

				for _, i := range data {
					r, err := redisClient.HGetAll(fmt.Sprintf("mac:%s", i)).Result()
					if err != nil {
						log.Error(err)
					}

					if len(r) == 0 {
						log.Debugf("Empty data set for %s", i)
						break
					}

					macAddress.WithLabelValues(r["mac"], r["ip"]).Set(1)
				}

				// log.Debugf("Sleeping %d seconds", 30)
				time.Sleep(time.Second * 30)
			}
		}

	}()

}

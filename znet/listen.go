package znet

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/things/things"
	"github.com/xaque208/znet/arpwatch"
)

type Listener struct {
	config      *Config
	thingServer *things.Server
	redisClient *redis.Client
	httpServer  *http.Server
}

const (
	macsList  = "macs"
	macsTable = "mac:*"
)

var (
	macAddress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mac",
		Help: "MAC Address",
	}, []string{"mac", "ip"})
)

func NewListener(config *Config) (Listener, error) {
	znetListener := Listener{
		config: config,
	}

	var err error
	prometheus.MustRegister(macAddress)

	// Attach a things server
	znetListener.thingServer, err = things.NewServer(znetListener.config.nats.URL, znetListener.config.nats.Topic)
	if err != nil {
		return Listener{}, err
	}

	// Attach a redis client
	znetListener.redisClient, err = NewRedisClient(znetListener.config.redis.Host)

	return znetListener, nil
}

func (l *Listener) Listen(listenAddr string, ch chan bool) {
	log.Info("ZNET Listening")

	l.httpServer = httpListen(listenAddr)

	// messages := make(chan things.Message)
	// go l.thingServer.Listen(messages, messageHandler)
	//
	// log.Debug("Starting arpwatch")
	// go arpWatch(l.redisConfig)
	//
	// log.Debugf("HTTP listening on %s", listenAddr)
	// srv := httpListen(listenAddr)

	defer l.shutdown()
	<-ch
}

func (l *Listener) shutdown() {
	log.Info("ZNET Shutting Down")

	log.Info("closing redis connection")
	l.redisClient.Close()

	// log.Info("halting things server")
	// l.thingsServer.Close()
	// l.httpServer.Shutdown(nil)
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

func (l *Listener) lightsHandler(command things.Command) {

	roomName := command.Arguments["room"]
	state := command.Arguments["state"]

	if state != "on" && state != "off" {
		log.Errorf("Unknown light state received %s", state)
	}

	log.Debugf("Using RFToy at %s", l.config.Endpoint)
	r := rftoy.RFToy{Address: l.config.Endpoint}

	room, err := l.config.Room(roomName.(string))
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

func (l *Listener) messageHandler(messages chan things.Message) {
	for {
		select {
		case msg := <-messages:
			log.Debugf("New message: %+v", msg)

			for _, c := range msg.Commands {
				if c.Name == "lights" {
					go l.lightsHandler(c)
				} else {
					log.Warnf("Unknown command %s", c.Name)
				}
			}

		}
	}
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

	ticker := time.NewTicker(30 * time.Second)

	go func() {

		for {
			select {
			default:
				go aw.Update()

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

				<-ticker.C
			}
		}

	}()

}

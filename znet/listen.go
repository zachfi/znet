package znet

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/things/things"
	"github.com/xaque208/znet/arpwatch"
)

type Listener struct {
	Config      *Config
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

func (z *Znet) Listen(listenAddr string, ch chan bool) {
	var err error
	z.listener, err = NewListener(&z.Config)
	if err != nil {
		log.Fatal(err)
	}

	z.listener.Listen(listenAddr, ch)
}

func NewListener(config *Config) (*Listener, error) {
	l := &Listener{
		Config: config,
	}

	var err error
	prometheus.MustRegister(macAddress)

	// Attach a things server
	log.Debugf("Using nats %s#%s", l.Config.Nats.URL, l.Config.Nats.Topic)
	l.thingServer, err = things.NewServer(l.Config.Nats.URL, l.Config.Nats.Topic)
	if err != nil {
		return &Listener{}, err
	}

	// Attach a redis client
	l.redisClient, err = NewRedisClient(l.Config.Redis.Host)

	return l, nil
}

func (l *Listener) Listen(listenAddr string, ch chan bool) {
	log.Infof("Listening on %s", listenAddr)

	l.httpServer = httpListen(listenAddr)

	messages := make(chan things.Message)
	go l.messageHandler(messages)
	go l.thingServer.Listen(messages)

	// log.Debug("Starting arpwatch")
	// go arpWatch(l.redisClient)

	<-ch
	l.Shutdown()
}

func (l *Listener) Shutdown() {
	log.Info("ZNET Shutting Down")

	// log.Info("closing redis connection")
	// l.redisClient.Close()

	log.Info("halting Things server")
	l.thingServer.Close()

	log.Info("halting HTTP server")
	l.httpServer.Shutdown(nil)
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

	log.Debugf("Using RFToy at %s", l.Config.Endpoint)
	r := rftoy.RFToy{Address: l.Config.Endpoint}

	room, err := l.Config.Room(roomName.(string))
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

// messageHandler
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

	aw := arpwatch.NewArpWatch()

	ticker := time.NewTicker(30 * time.Second)

	go func() {

		for {
			select {
			default:
				aw.Update()

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

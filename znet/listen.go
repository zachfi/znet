package znet

import (
	"context"
	"net"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/things/things"
	pb "github.com/xaque208/znet/rpc"
	"google.golang.org/grpc"
)

// Listener is a znet server
type Listener struct {
	Config      *Config
	thingServer *things.Server
	redisClient *redis.Client
	httpServer  *http.Server
}

// Listen starts the znet listener
func (z *Znet) Listen(listenAddr string, ch chan bool) {
	var err error
	z.listener, err = NewListener(&z.Config)
	if err != nil {
		log.Fatal(err)
	}

	z.listenRPC()

	z.listener.Listen(listenAddr, ch)
}

func (z *Znet) listenRPC() {

	if z.Config.RPC.ListenAddress != "" {
		log.Debugf("Starting RPC listener on %s", z.Config.RPC.ListenAddress)

		inventoryServer := &inventoryServer{
			inventory: z.Inventory,
		}

		lightServer := &lightServer{
			lights: z.Lights,
		}

		go func() {
			lis, err := net.Listen("tcp", z.Config.RPC.ListenAddress)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()

			pb.RegisterInventoryServer(grpcServer, inventoryServer)
			pb.RegisterLightsServer(grpcServer, lightServer)

			err = grpcServer.Serve(lis)
			if err != nil {
				log.Error(err)
			}

		}()
	}

}

// NewListener builds a new Listener object from the received configuration.
func NewListener(config *Config) (*Listener, error) {
	l := &Listener{
		Config: config,
	}

	var err error

	// Attach a things server
	log.Debugf("Using nats %s#%s", l.Config.Nats.URL, l.Config.Nats.Topic)
	l.thingServer, err = things.NewServer(l.Config.Nats.URL, l.Config.Nats.Topic)
	if err != nil {
		return &Listener{}, err
	}

	// Attach a redis client
	l.redisClient, err = NewRedisClient(l.Config.Redis.Host)
	if err != nil {
		return &Listener{}, err
	}

	return l, nil
}

// Listen starts the http listener
func (l *Listener) Listen(listenAddr string, ch chan bool) {
	log.Infof("Listening on %s", listenAddr)

	l.httpServer = httpListen(listenAddr)

	messages := make(chan things.Message)
	go l.messageHandler(messages)
	// go l.thingServer.Listen(messages)

	<-ch
	l.Shutdown()
}

// Shutdown closes down the to the message bus and shuts down the HTTP server.
func (l *Listener) Shutdown() {
	log.Info("ZNET Shutting Down")

	// log.Info("closing redis connection")
	// l.redisClient.Close()

	log.Info("halting Things server")
	l.thingServer.Close()

	log.Info("halting HTTP server")
	err := l.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Error(err)
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

// messageHandler
func (l *Listener) messageHandler(messages chan things.Message) {
	for msg := range messages {
		log.Debugf("New message: %+v", msg)

		for _, c := range msg.Commands {
			log.Warnf("Unknown command %s", c.Name)
		}

	}
}

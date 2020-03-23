package znet

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/timer"
	pb "github.com/xaque208/znet/rpc"
	"google.golang.org/grpc"
)

// Listen starts the znet listener
func (z *Znet) Listen(listenAddr string, ch chan bool) {
	var err error
	z.listener, err = NewListener(&z.Config)
	if err != nil {
		log.Fatal(err)
	}

	consumers := []events.Consumer{
		z.Lights,
		z,
	}

	log.Tracef("consumers: %+v", consumers)

	z.EventChannel = make(chan events.Event)
	z.EventConsumers = make(map[string][]events.Handler)

	z.initEventConsumers(consumers)
	z.initEventConsumer()
	z.listenRPC()

	z.listener.Listen(listenAddr, ch)
}

// Subscriptions is yet to be used, but conforms to the interface for generating consumers of named events.
func (z *Znet) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()
	return s.Table
}

// listenRPC starts the RPC server and all the services.
func (z *Znet) listenRPC() {
	if z.Config.RPC.ListenAddress != "" {
		log.Infof("Starting RPC listener on %s", z.Config.RPC.ListenAddress)

		inventoryServer := &inventoryServer{
			inventory: z.Inventory,
		}

		lightServer := &lightServer{
			lights: z.Lights,
		}

		eventServer := &eventServer{
			eventNames: timer.EventNames,
			ch:         z.EventChannel,
		}

		eventServer.RegisterEvents(timer.EventNames)
		eventServer.RegisterEvents(astro.EventNames)

		go func() {
			lis, err := net.Listen("tcp", z.Config.RPC.ListenAddress)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			grpcServer := grpc.NewServer()

			pb.RegisterInventoryServer(grpcServer, inventoryServer)
			pb.RegisterLightsServer(grpcServer, lightServer)
			pb.RegisterEventsServer(grpcServer, eventServer)

			err = grpcServer.Serve(lis)
			if err != nil {
				log.Error(err)
			}
		}()
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

// initEventConsumers updates the EventConsumers map for each consumer, to
// append a handler for the discovered topic keys.
func (z *Znet) initEventConsumers(consumers []events.Consumer) {
	for _, e := range consumers {
		subs := e.Subscriptions()
		for k, handlers := range subs {
			for _, x := range handlers {
				z.EventConsumers[k] = append(z.EventConsumers[k], x)
			}
		}
	}
}

// initEventConsumer starts a routine that never ends to read from the
// EventChannel and execute the loaded handlers with the event Payload.
func (z *Znet) initEventConsumer() {
	go func(ch chan events.Event) {
		for e := range ch {
			log.Debugf("z.EventConsumers: %+v", z.EventConsumers)

			if handlers, ok := z.EventConsumers[e.Name]; ok {
				log.Infof("handling message %s", e.Name)
				for _, h := range handlers {
					h(e.Payload)
				}
			}
		}
	}(z.EventChannel)
}

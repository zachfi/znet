package znet

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/timer"
	pb "github.com/xaque208/znet/rpc"
)

// Listen starts the znet listener.  The Listener is responsible for starting
// up all the event handling threads, and then blocking on the final HTTP
// listener.
func (z *Znet) Listen(listenAddr string, ch chan bool) {
	var err error
	z.listener, err = NewListener(&z.Config)
	if err != nil {
		log.Fatal(err)
	}

	// Each events.Consumer here is later used
	consumers := []events.Consumer{
		z.Lights,
		z,
	}

	log.Tracef("%d event consumers", len(consumers))

	z.EventChannel = make(chan events.Event)
	z.EventConsumers = make(map[string][]events.Handler)

	z.initEventConsumers(consumers)
	z.initEventConsumer()
	z.listenRPC()

	z.listener.Listen(listenAddr, ch)
}

// Subscriptions is yet to be used, but conforms to the interface for
// generating consumers of named events.
func (z *Znet) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()
	return s.Table
}

// listenRPC starts the RPC server and all the services.
func (z *Znet) listenRPC() {
	if z.Config.RPC.ListenAddress != "" {
		log.Infof("starting RPC listener %s", z.Config.RPC.ListenAddress)

		inventoryServer := &inventoryServer{
			inventory: z.Inventory,
		}

		lightServer := &lightServer{
			lights: z.Lights,
		}

		eventServer := &eventServer{
			ch:         z.EventChannel,
			remoteChan: make(chan *pb.Event, 1000),
		}

		eventServer.RegisterEvents(timer.EventNames)
		eventServer.RegisterEvents(astro.EventNames)
		eventServer.RegisterEvents(gitwatch.EventNames)

		go func() {
			lis, err := net.Listen("tcp", z.Config.RPC.ListenAddress)
			if err != nil {
				log.Fatalf("failed to listen: %s", err)
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

// initEventConsumers updates the z.EventConsumers map for each received
// consumer in order to append a handler for the discovered topic keys.
func (z *Znet) initEventConsumers(consumers []events.Consumer) {
	for _, e := range consumers {
		subs := e.Subscriptions()
		for k, handlers := range subs {
			z.EventConsumers[k] = append(z.EventConsumers[k], handlers...)
		}
	}
}

// initEventConsumer starts a routine that never ends to read from
// z.EventChannel and execute the loaded handlers with the event Payload.
func (z *Znet) initEventConsumer() {
	go func(ch chan events.Event) {
		log.Debugf("total %d z.EventConsumers", len(z.EventConsumers))

		for e := range ch {
			if handlers, ok := z.EventConsumers[e.Name]; ok {
				log.Debugf("executing %d handlers for event %s", len(handlers), e.Name)
				log.Tracef("listener heard event %s: %s", e.Name, string(e.Payload))
				for _, h := range handlers {
					err := h(e.Payload)
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				log.Warnf("received event with no handlers: %+v", e.Name)
			}
		}
	}(z.EventChannel)
}

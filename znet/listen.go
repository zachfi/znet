package znet

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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

		// commandServer := &commandServer{}

		go func() {
			lis, err := net.Listen("tcp", z.Config.RPC.ListenAddress)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()

			pb.RegisterInventoryServer(grpcServer, inventoryServer)
			pb.RegisterLightsServer(grpcServer, lightServer)
			// pb.RegisterCommandServer(grpcServer, commandServer)

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

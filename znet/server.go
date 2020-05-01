package znet

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
	pb "github.com/xaque208/znet/rpc"
)

// Server is a znet Server.
type Server struct {
	eventMachine *eventmachine.EventMachine

	httpConfig HTTPConfig
	rpcConfig  RPCConfig

	httpServer *http.Server
	grpcServer *grpc.Server

	rpcEventServer *eventServer
}

var (
	executionExitStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "znet_execution_result",
		Help: "Stats on the received ExecutionResult RPC events",
	}, []string{"command"})

	executionDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "znet_execution_duration",
		Help: "Stats on the received ExecutionResult RPC events",
	}, []string{"command"})

	eventTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "znet_event_total",
		Help: "The total number of events that have been seen since start",
	}, []string{})
)

func init() {
	prometheus.MustRegister(
		executionExitStatus,
		executionDuration,
		eventTotal,
	)
}

// NewServer creates a new Server composed of the received information.
func NewServer(httpConfig HTTPConfig, rpcConfig RPCConfig, consumers []events.Consumer) *Server {

	eventMachine, err := eventmachine.New(consumers)
	if err != nil {
		log.Error(err)
	}

	var httpServer *http.Server
	if httpConfig.ListenAddress != "" {
		httpServer = &http.Server{Addr: httpConfig.ListenAddress}
	}

	s := &Server{
		eventMachine: eventMachine,

		rpcEventServer: &eventServer{ch: eventMachine.EventChannel},

		httpConfig: httpConfig,
		rpcConfig:  rpcConfig,

		httpServer: httpServer,
		grpcServer: grpc.NewServer(),
	}

	return s
}

// Start is used to launch the server routines.
func (s *Server) Start(z *Znet) error {
	http.Handle("/metrics", promhttp.Handler())

	if s.httpConfig.ListenAddress != "" {
		log.Infof("starting HTTP listener %s", s.httpConfig.ListenAddress)

		go func() {
			if err := s.httpServer.ListenAndServe(); err != nil {
				log.Error(err)
			}
		}()
	}

	if s.rpcConfig.ListenAddress != "" {
		log.Infof("starting RPC listener %s", s.rpcConfig.ListenAddress)

		rpcInventoryServer := &inventoryServer{
			inventory: z.Inventory,
		}

		rpcLightServer := &lightServer{
			lights: z.Lights,
		}

		pb.RegisterInventoryServer(s.grpcServer, rpcInventoryServer)
		pb.RegisterLightsServer(s.grpcServer, rpcLightServer)
		pb.RegisterEventsServer(s.grpcServer, s.rpcEventServer)

		s.rpcEventServer.RegisterEvents(agent.EventNames)
		s.rpcEventServer.RegisterEvents(astro.EventNames)
		s.rpcEventServer.RegisterEvents(continuous.EventNames)
		s.rpcEventServer.RegisterEvents(gitwatch.EventNames)
		s.rpcEventServer.RegisterEvents(timer.EventNames)

		go func() {
			lis, err := net.Listen("tcp", s.rpcConfig.ListenAddress)
			if err != nil {
				log.Errorf("failed to listen: %s", err)
			}

			err = s.grpcServer.Serve(lis)
			if err != nil {
				log.Error(err)
			}
		}()
	}

	return nil
}

// Stop is used to close up the connections and channels.
func (s *Server) Stop() error {
	errs := []error{}

	var err error

	s.grpcServer.Stop()

	err = s.httpServer.Shutdown(context.Background())
	if err != nil {
		errs = append(errs, err)
	}

	err = s.eventMachine.Stop()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors while shutting down: %s", errs)
	}

	return nil
}

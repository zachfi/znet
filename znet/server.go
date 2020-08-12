package znet

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/events"
	pb "github.com/xaque208/znet/rpc"
)

// Server is a znet Server.
type Server struct {
	cancel         func()
	ctx            context.Context
	eventMachine   *eventmachine.EventMachine
	grpcServer     *grpc.Server
	httpConfig     *HTTPConfig
	httpServer     *http.Server
	ldapConfig     *inventory.LDAPConfig
	rpcConfig      *RPCConfig
	rpcEventServer *eventServer
}

type statusCheckHandler struct {
	server *Server
}

func init() {
	prometheus.MustRegister(
		eventMachineConsumers,
		eventMachineHandlers,
		eventTotal,
		executionDuration,
		executionExitStatus,

		airHeatindex,
		airHumidity,
		airTemperature,
		tempCoef,
		waterTempCoef,
		waterTemperature,
		thingWireless,

		// rpc
		rpcEventServerEventCount,
		rpcEventServerSubscriberCount,

		telemetryIOTUnhandledReport,
		telemetryIOTReport,
		telemetryIOTBatteryPercent,
		telemetryIOTLinkQuality,
	)
}

// NewServer creates a new Server composed of the received information.
func NewServer(config Config, consumers []events.Consumer) *Server {
	ctx := context.Background()

	eventMachine, err := eventmachine.New(ctx, consumers)
	if err != nil {
		log.Error(err)
	}

	var httpServer *http.Server
	if config.HTTP.ListenAddress != "" {

		httpServer = &http.Server{Addr: config.HTTP.ListenAddress}
	}

	roots, err := CABundle(config.Vault)
	if err != nil {
		log.Error(err)
	}

	c, err := newCertify(config.Vault, config.TLS)
	if err != nil {
		log.Error(err)
	}

	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
		ClientCAs:      roots,
		// RootCAs:        cp,
		ClientAuth: tls.RequireAndVerifyClientCert,
		// ClientAuth:           tls.VerifyClientCertIfGiven,
	}

	if config.HTTP == nil || config.RPC == nil {
		log.Errorf("unable to build znet Server with nil HTTPConfig or RPCConfig")
	}

	s := &Server{
		ctx:          ctx,
		eventMachine: eventMachine,

		rpcEventServer: &eventServer{
			eventMachineChannel: eventMachine.EventChannel,
			ctx:                 ctx,
		},

		httpConfig: config.HTTP,
		rpcConfig:  config.RPC,
		ldapConfig: config.LDAP,

		httpServer: httpServer,

		grpcServer: grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig))),
	}

	return s
}

func (s *statusCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := make(map[string]interface{})

	status["status"] = "healthy"

	grpcInfo := s.server.grpcServer.GetServiceInfo()

	status["grpcServices"] = len(grpcInfo)

	payload, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Error(err)
	}

	fmt.Fprint(w, string(payload))
}

// Start is used to launch the server routines.
func (s *Server) Start(z *Znet) error {
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/status/check", &statusCheckHandler{server: s})

	if s.httpConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"listen_address": s.httpConfig.ListenAddress,
		}).Info("starting HTTP listener")

		go func() {
			if err := s.httpServer.ListenAndServe(); err != nil {
				log.Error(err)
			}
		}()
	}

	if s.rpcConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"listen_address": s.rpcConfig.ListenAddress,
		}).Info("starting RPC listener")

		inv := inventory.NewInventory(*s.ldapConfig)

		// inventoryServer
		rpcInventoryServer := &inventoryServer{
			inventory: inv,
		}
		pb.RegisterInventoryServer(s.grpcServer, rpcInventoryServer)

		// lightServer
		rpcLightServer := &lightServer{
			lights: z.Lights,
		}
		pb.RegisterLightsServer(s.grpcServer, rpcLightServer)

		// telemetryServer
		rpcTelemetryServer := newTelemetryServer(inv)
		pb.RegisterTelemetryServer(s.grpcServer, rpcTelemetryServer)

		// Register and configure the rpcEventServer
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

	go func() {
		t := time.NewTicker(10 * time.Second)

		for range t.C {
			// export the eventMachine data
			eventMachineConsumers.WithLabelValues().Set(float64(len(s.eventMachine.EventConsumers)))

			for name, handlers := range s.eventMachine.EventConsumers {
				eventMachineHandlers.WithLabelValues(name).Set(float64(len(handlers)))
			}

			// export the event RPC server data
			subscriberCount, eventCount := s.rpcEventServer.Report()

			rpcEventServerSubscriberCount.WithLabelValues().Set(float64(subscriberCount))
			rpcEventServerEventCount.WithLabelValues().Set(float64(eventCount))
		}
	}()

	return nil
}

// Stop is used to close up the connections and channels.
func (s *Server) Stop() error {
	errs := []error{}
	var err error

	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

	err = s.httpServer.Shutdown(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	s.grpcServer.GracefulStop()

	cancel()
	s.cancel()

	if len(errs) > 0 {
		return fmt.Errorf("errors while shutting down: %s", errs)
	}

	return nil
}

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
	"github.com/xaque208/znet/pkg/iot"
	pb "github.com/xaque208/znet/rpc"
)

// Server is a znet Server.
type Server struct {
	cancel       func()
	ctx          context.Context
	eventMachine *eventmachine.EventMachine
	grpcServer   *grpc.Server
	httpConfig   *HTTPConfig
	httpServer   *http.Server
	ldapConfig   *inventory.LDAPConfig
	rpcConfig    *RPCConfig
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
		telemetryIOTBridgeState,
	)
}

// NewServer creates a new Server composed of the received information.
func NewServer(config Config, consumers []events.Consumer) *Server {
	ctx, cancel := context.WithCancel(context.Background())

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
		cancel:       cancel,
		eventMachine: eventMachine,

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
	log.WithFields(log.Fields{
		"rpc_listen":  s.rpcConfig.ListenAddress,
		"http_listen": s.httpConfig.ListenAddress,
	}).Debug("starting znetd")

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/status/check", &statusCheckHandler{server: s})

	if s.httpConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"listen_address": s.httpConfig.ListenAddress,
		}).Info("starting HTTP listener")

		go func() {
			if err := s.httpServer.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Error(err)
				}
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

		// lightServer to receive RPC calls for lighting changes directly.
		rpcLightServer := &lightServer{
			lights: z.Lights,
		}
		pb.RegisterLightsServer(s.grpcServer, rpcLightServer)

		// telemetryServer
		rpcTelemetryServer := newTelemetryServer(inv, s.eventMachine)
		pb.RegisterTelemetryServer(s.grpcServer, rpcTelemetryServer)

		// Register and configure the rpcEventServer
		rpcEventServer := &eventServer{
			eventMachineChannel: s.eventMachine.EventChannel,
			ctx:                 s.ctx,
		}

		pb.RegisterEventsServer(s.grpcServer, rpcEventServer)
		rpcEventServer.RegisterEvents(agent.EventNames)
		rpcEventServer.RegisterEvents(astro.EventNames)
		rpcEventServer.RegisterEvents(continuous.EventNames)
		rpcEventServer.RegisterEvents(gitwatch.EventNames)
		rpcEventServer.RegisterEvents(iot.EventNames)
		rpcEventServer.RegisterEvents(timer.EventNames)

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
		}
	}()

	return nil
}

// Stop is used to close up the connections and channels.
func (s *Server) Stop() error {
	errs := []error{}
	var err error

	log.Debug("stopping znetd")

	d := time.Now().Add(1500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(s.ctx, d)
	defer cancel()

	err = s.httpServer.Shutdown(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	select {
	case <-time.After(500 * time.Millisecond):
	case <-ctx.Done():
		log.Error(ctx.Err())
		err = ctx.Err()
		if err != nil {
			errs = append(errs, err)
		}
	}

	s.grpcServer.Stop()
	err = s.eventMachine.Stop()
	if err != nil {
		errs = append(errs, err)
	}
	s.cancel()

	if len(errs) > 0 {
		return fmt.Errorf("errors while shutting down: %s", errs)
	}

	return nil
}

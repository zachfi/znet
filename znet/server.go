package znet

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/iot"
	"github.com/xaque208/znet/rpc"
)

// Server is a znet Server.
type Server struct {
	cancel       func()
	ctx          context.Context
	eventMachine *eventmachine.EventMachine
	grpcServer   *grpc.Server
	httpConfig   *config.HTTPConfig
	httpServer   *http.Server
	ldapConfig   *inventory.LDAPConfig
	rpcConfig    *config.RPCConfig
}

type statusCheckHandler struct {
	server *Server
}

func init() {
	prometheus.MustRegister(
		eventTotal,
		executionDuration,
		executionExitStatus,

		airHeatindex,
		airHumidity,
		airTemperature,
		tempCoef,
		thingWireless,
		waterTempCoef,
		waterTemperature,

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
// func NewServer(config Config, consumers []events.Consumer) *Server {
func NewServer(config Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	eventMachine, err := eventmachine.New(ctx, nil)
	if err != nil {
		log.Error(err)
	}

	var httpServer *http.Server
	if config.HTTP.ListenAddress != "" {
		httpServer = &http.Server{Addr: config.HTTP.ListenAddress}
	}

	if config.HTTP == nil || config.RPC == nil {
		log.Errorf("unable to build znet Server with nil HTTPConfig or RPCConfig")
	}

	return &Server{
		ctx:          ctx,
		cancel:       cancel,
		eventMachine: eventMachine,

		httpConfig: config.HTTP,
		rpcConfig:  config.RPC,
		ldapConfig: config.LDAP,

		httpServer: httpServer,

		grpcServer: comms.StandardRPCServer(config.Vault, config.TLS),
	}
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
func (s *Server) startRPCListener() error {
	rpcInventoryServer := inventory.NewRPCServer(*s.ldapConfig)
	rpc.RegisterInventoryServer(s.grpcServer, rpcInventoryServer)

	// telemetryServer
	inv := inventory.NewInventory(*s.ldapConfig)
	rpcTelemetryServer := newTelemetryServer(inv, s.eventMachine)
	rpc.RegisterTelemetryServer(s.grpcServer, rpcTelemetryServer)

	// rpcEventServer
	rpcEventServer := &eventServer{eventMachine: s.eventMachine, ctx: s.ctx}

	rpc.RegisterEventsServer(s.grpcServer, rpcEventServer)
	rpcEventServer.RegisterEvents(agent.EventNames)
	rpcEventServer.RegisterEvents(astro.EventNames)
	rpcEventServer.RegisterEvents(continuous.EventNames)
	rpcEventServer.RegisterEvents(gitwatch.EventNames)
	rpcEventServer.RegisterEvents(iot.EventNames)
	rpcEventServer.RegisterEvents(timer.EventNames)

	go func() {
		lis, err := net.Listen("tcp", s.rpcConfig.ListenAddress)
		if err != nil {
			log.Errorf("rpc failed to listen: %s", err)
		}

		err = s.grpcServer.Serve(lis)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (s *Server) startHTTPListener() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Errorf("http failed to listen: %s", err)
			}
		}
	}()

	return nil
}

// Start is used to launch the server routines.
func (s *Server) Start(z *Znet) error {

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/status/check", &statusCheckHandler{server: s})
	http.HandleFunc("/alerts", s.alertsHandler)

	if s.httpConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"http_listen": s.httpConfig.ListenAddress,
		}).Info("starting HTTP listener")

		s.startHTTPListener()
	}

	if s.rpcConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"rpc_listen": s.rpcConfig.ListenAddress,
		}).Debug("starting RPC listener")

		s.startRPCListener()
	}

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

func (s *Server) alertsHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var m hookMessage
	if err := dec.Decode(&m); err != nil {
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "invalid request body", 400)
		return
	}

	log.Debugf("webhook alert: %+v", m)
}

type hookMessage struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []alert           `json:"alerts"`
}

type alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"EndsAt,omitempty"`
}

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

	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/telemetry"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/eventmachine"
)

// Server is a znet Server.
type Server struct {
	cancel       func()
	ctx          context.Context
	eventMachine *eventmachine.EventMachine
	grpcServer   *grpc.Server
	httpConfig   *config.HTTPConfig
	httpServer   *http.Server
	ldapConfig   *config.LDAPConfig
	rpcConfig    *config.RPCConfig
	config       *config.Config
}

type statusCheckHandler struct {
	server *Server
}

func init() {
	prometheus.MustRegister(
		eventTotal,
		executionDuration,
		executionExitStatus,

		tempCoef,
		waterTempCoef,
	)
}

// NewServer creates a new Server composed of the received information.
// func NewServer(config Config, consumers []events.Consumer) *Server {
func NewServer(cfg *config.Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	eventMachine, err := eventmachine.New(ctx, nil)
	if err != nil {
		log.Error(err)
	}

	var httpServer *http.Server
	if cfg.HTTP.ListenAddress != "" {
		httpServer = &http.Server{Addr: cfg.HTTP.ListenAddress}
	}

	if cfg.HTTP == nil || cfg.RPC == nil {
		log.Errorf("unable to build znet Server with nil HTTPConfig or RPCConfig")
	}

	return &Server{
		ctx:          ctx,
		cancel:       cancel,
		eventMachine: eventMachine,

		config:     cfg,
		httpConfig: cfg.HTTP,
		rpcConfig:  cfg.RPC,
		ldapConfig: cfg.LDAP,

		httpServer: httpServer,

		grpcServer: comms.StandardRPCServer(cfg.Vault, cfg.TLS),
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
	//
	// inventoryServer
	inventoryServer, err := inventory.NewServer(s.ldapConfig)
	if err != nil {
		return err
	}

	inventory.RegisterInventoryServer(s.grpcServer, inventoryServer)

	//
	// astroServer
	astroServer, err := astro.NewAstro(s.config)
	if err != nil {
		return err
	}

	astro.RegisterAstroServer(s.grpcServer, astroServer)

	//
	// lightsServer
	lightsServer, err := lights.NewLights(s.config)
	if err != nil {
		return err
	}

	lights.RegisterLightsServer(s.grpcServer, lightsServer)

	//
	// telemetryServer
	inv, err := inventory.NewInventory(s.ldapConfig)
	if err != nil {
		return err
	}

	telemetryServer, err := telemetry.NewServer(inv, lightsServer)
	if err != nil {
		return err
	}

	telemetry.RegisterTelemetryServer(s.grpcServer, telemetryServer)

	//
	// rpcEventServer
	// rpcEventServer := &eventServer{eventMachine: s.eventMachine, ctx: s.ctx}

	// rpc.RegisterEventsServer(s.grpcServer, rpcEventServer)
	// rpcEventServer.RegisterEvents(continuous.EventNames)
	// rpcEventServer.RegisterEvents(gitwatch.EventNames)
	// rpcEventServer.RegisterEvents(timer.EventNames)

	//
	// timerServer
	timerServer, err := timer.NewServer(lightsServer)
	if err != nil {
		return err
	}

	timer.RegisterTimerServer(s.grpcServer, timerServer)

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

		err := s.startHTTPListener()
		if err != nil {
			return err
		}
	}

	if s.rpcConfig.ListenAddress != "" {
		log.WithFields(log.Fields{
			"rpc_listen": s.rpcConfig.ListenAddress,
		}).Debug("starting RPC listener")

		err := s.startRPCListener()
		if err != nil {
			return err
		}
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

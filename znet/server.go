package znet

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	"github.com/xaque208/znet/pkg/iot"
)

// Server is a znet Server.
type Server struct {
	sync.Mutex
	config        *config.Config
	grpcServer    *grpc.Server
	httpServer    *http.Server
	NewHTTPServer comms.HTTPServerFunc
	NewRPCServer  comms.RPCServerFunc

	mqttClient mqtt.Client
	invClient  inventory.Inventory
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
func NewServer(cfg *config.Config) (*Server, error) {
	if cfg.HTTP == nil {
		return nil, fmt.Errorf("unable to build znet Server with nil HTTP config")
	}

	if cfg.RPC == nil {
		return nil, fmt.Errorf("unable to build znet Server with nil RPC config")
	}

	if cfg.Vault == nil {
		return nil, fmt.Errorf("unable to build znet Server with nil Vault config")
	}

	if cfg.TLS == nil {
		return nil, fmt.Errorf("unable to build znet Server with nil TLS config")
	}

	return &Server{
		config: cfg,

		NewRPCServer:  comms.StandardRPCServer,
		NewHTTPServer: comms.StandardHTTPServer,
	}, nil
}

func (s *Server) startRPCListener() error {
	s.Lock()
	defer s.Unlock()

	var err error

	if s.grpcServer == nil {
		grpcServer, serverErr := s.NewRPCServer(s.config)
		if serverErr != nil {
			return serverErr
		}
		s.grpcServer = grpcServer
	}

	if s.mqttClient == nil {
		// clients to be used by the servers
		s.mqttClient, err = iot.NewMQTTClient(s.config.MQTT)
		if err != nil {
			return err
		}
	}

	if s.invClient == nil {
		s.invClient, err = inventory.NewLDAPInventory(s.config.LDAP)
		if err != nil {
			return err
		}
	}

	//
	// inventoryServer
	inventoryServer, err := inventory.NewLDAPServer(s.invClient)
	if err != nil {
		return err
	}

	inventory.RegisterInventoryServer(s.grpcServer, inventoryServer)

	//
	// lightsServer
	lightsServer, err := lights.NewLights(s.config.Lights)
	if err != nil {
		return err
	}

	hue, err := lights.NewHueLight(s.config.Lights)
	if err != nil {
		log.Error(err)
	} else {
		lightsServer.AddHandler(hue)
	}

	zigbee, err := lights.NewZigbeeLight(s.config, s.mqttClient, s.invClient)
	if err != nil {
		log.Error(err)
	} else {
		lightsServer.AddHandler(zigbee)
	}

	lights.RegisterLightsServer(s.grpcServer, lightsServer)

	//
	// astroServer
	astroServer, err := astro.NewAstro(s.config, lightsServer)
	if err != nil {
		return err
	}

	astro.RegisterAstroServer(s.grpcServer, astroServer)

	// iotServer
	iotServer, err := iot.NewServer(s.mqttClient)
	if err != nil {
		return err
	}

	iot.RegisterIOTServer(s.grpcServer, iotServer)

	//
	// telemetryServer
	telemetryServer, err := telemetry.NewServer(s.invClient, lightsServer)
	if err != nil {
		return err
	}

	err = telemetryServer.SetIOTServer(iotServer)
	if err != nil {
		return err
	}

	telemetry.RegisterTelemetryServer(s.grpcServer, telemetryServer)

	//
	// timerServer
	timerServer, err := timer.NewServer(lightsServer)
	if err != nil {
		return err
	}

	timer.RegisterTimerServer(s.grpcServer, timerServer)

	go func() {
		lis, err := net.Listen("tcp", s.config.RPC.ListenAddress)
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
	s.Lock()
	defer s.Unlock()

	if s.httpServer == nil {
		httpServer, err := s.NewHTTPServer(s.config)
		if err != nil {
			return err
		}
		s.httpServer = httpServer
	}

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

	if s.config.HTTP.ListenAddress != "" {
		log.WithFields(log.Fields{
			"http_listen": s.config.HTTP.ListenAddress,
		}).Debug("starting HTTP")

		err := s.startHTTPListener()
		if err != nil {
			return err
		}
	}

	if s.config.RPC.ListenAddress != "" {
		log.WithFields(log.Fields{
			"rpc_listen": s.config.RPC.ListenAddress,
		}).Debug("starting RPC")

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
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	if s.httpServer != nil {
		err = s.httpServer.Shutdown(ctx)
		if err != nil {
			errs = append(errs, err)
		}
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

	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}

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

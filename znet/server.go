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
	eventMachine *eventmachine.EventMachine

	httpConfig *HTTPConfig
	rpcConfig  *RPCConfig
	ldapConfig *inventory.LDAPConfig

	httpServer *http.Server
	grpcServer *grpc.Server

	rpcEventServer *eventServer
}

type statusCheckHandler struct {
	server *Server
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
	}, []string{"event_name"})

	eventMachineConsumers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "znet_event_machine_consumers",
		Help: "The current number of event consumers",
	}, []string{})

	eventMachineHandlers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "znet_event_machine_handlers",
		Help: "The current number of event handlers per consumer",
	}, []string{"event_name"})

	// ciJobsRunning = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	// 	Name: "znet_ci_jobs_running",
	// 	Help: "Stats on running CI jobs",
	// }, []string{})

	airTemperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_air_temperature",
		Help: "Temperature",
	}, []string{"device"})

	airHumidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_air_humidity",
		Help: "humidity",
	}, []string{"device"})

	airHeatindex = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_air_heatindex",
		Help: "computed heat index",
	}, []string{"device"})

	waterTemperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_water_temperature",
		Help: "Water Temperature",
	}, []string{"device"})

	tempCoef = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_air_temperature_coef",
		Help: "Air Temperature Coefficient",
	}, []string{"device"})

	waterTempCoef = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_water_temperature_coef",
		Help: "Water Temperature Coefficient",
	}, []string{"device"})

	thingWireless = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_wireless",
		Help: "wireless information",
	}, []string{"device", "ssid", "bssid", "ip"})

	// rpc eventServer
	rpcEventServerSubscriberCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_eventserver_subscriber_count",
		Help: "The current number of rpc subscribers",
	}, []string{})

	rpcEventServerEventCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_eventserver_event_count",
		Help: "The current number of rpc events that are subscribed",
	}, []string{})

	// rpc telemetry
	telemetryIOTUnhandledReport = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetryunhandled_object_report",
		Help: "The total number of notice calls that include an unhandled object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTReport = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetryobject_report",
		Help: "The total number of notice calls for an object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTBatteryPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_battery_percent",
		Help: "The reported batter percentage remaining.",
	}, []string{"object_id", "component"})

	telemetryIOTLinkQuality = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_link_quality",
		Help: "The reported link quality",
	}, []string{"object_id", "component"})
)

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
	eventMachine, err := eventmachine.New(consumers)
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
		eventMachine: eventMachine,

		rpcEventServer: &eventServer{ch: eventMachine.EventChannel},

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
		log.Infof("starting HTTP listener %s", s.httpConfig.ListenAddress)

		go func() {
			if err := s.httpServer.ListenAndServe(); err != nil {
				log.Error(err)
			}
		}()
	}

	if s.rpcConfig.ListenAddress != "" {
		log.Infof("starting RPC listener %s", s.rpcConfig.ListenAddress)

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

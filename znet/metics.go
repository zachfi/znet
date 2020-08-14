package znet

import "github.com/prometheus/client_golang/prometheus"

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
		Name: "rpc_telemetry_unhandled_object_report",
		Help: "The total number of notice calls that include an unhandled object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTReport = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetry_object_report",
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

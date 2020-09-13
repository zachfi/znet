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

	eventRemoteSendErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "znet_remote_send_error_total",
		Help: "The total number of that returned an error when sent",
	}, []string{"event_name"})

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
)

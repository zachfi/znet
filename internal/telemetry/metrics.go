package telemetry

import "github.com/prometheus/client_golang/prometheus"

var (
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

	telemetryIOTBridgeState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_bridge_state",
		Help: "The reported bridge state",
	}, []string{})

	//
	waterTemperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_water_temperature",
		Help: "Water Temperature",
	}, []string{"device"})

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

	thingWireless = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "thing_wireless",
		Help: "wireless information",
	}, []string{"device", "ssid", "bssid", "ip"})
)

func init() {
	prometheus.MustRegister(
		telemetryIOTUnhandledReport,
		telemetryIOTReport,
		telemetryIOTBatteryPercent,
		telemetryIOTLinkQuality,
		telemetryIOTBridgeState,

		airHeatindex,
		airHumidity,
		airTemperature,

		thingWireless,

		waterTemperature,
	)
}

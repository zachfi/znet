package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	telemetryIOTUnhandledReport = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetry_unhandled_object_report",
		Help: "The total number of notice calls that include an unhandled object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTReport = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetry_object_report",
		Help: "The total number of notice calls for an object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTBatteryPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_battery_percent",
		Help: "The reported batter percentage remaining.",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTLinkQuality = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_link_quality",
		Help: "The reported link quality",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTBridgeState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_bridge_state",
		Help: "The reported bridge state",
	}, []string{})

	telemetryIOTOccupancy = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_occupancy",
		Help: "Occupancy binary",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTWaterLeak = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_water_leak",
		Help: "Water leak binary",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_temperature",
		Help: "Sensor Temperature(C)",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTTamper = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_tamper",
		Help: "Tamper binary",
	}, []string{"object_id", "component", "zone"})

	telemetryIOTIlluminance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_illuminance",
		Help: "Illuminance(LQI)",
	}, []string{"object_id", "component", "zone"})

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
		airHeatindex,
		airHumidity,
		airTemperature,

		thingWireless,

		waterTemperature,
	)
}

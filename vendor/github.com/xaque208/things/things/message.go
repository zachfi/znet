package things

type Message struct {
	// Commands []Command       `json:"commands,omitempty"`
	// Device string `json:"id,omitempty"`
	// Events   []Event         `json:"events,omitempty"`
	Topic string `json:"topic,omitempty"`
	// Sensors  []SensorReading `json:"sensors,omitempty"`

	// sketch stats
	Sketch string `json:"sketch,omitempty"`

	// wifi stats
	BSSID string `json:"bssid,omitempty"`
	IP    string `json:"ip,omitempty"`
	RSSI  int    `json:"rssi,omitempty"`
	SSID  string `json:"ssid,omitempty"`

	// air stats
	Humidity    float32 `json:"humidity,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	HeatIndex   float32 `json:"heatindex,omitempty"`
}

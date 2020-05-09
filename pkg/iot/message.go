package iot

type WifiMessage struct {
	BSSID string `json:"bssid,omitempty"`
	IP    string `json:"ip,omitempty"`
	RSSI  int    `json:"rssi,omitempty"`
	SSID  string `json:"ssid,omitempty"`
}

type AirMessage struct {
	Humidity    float32 `json:"humidity,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	HeatIndex   float32 `json:"heatindex,omitempty"`
	TempCoef    float64 `json:"tempcoef,omitempty"`
}

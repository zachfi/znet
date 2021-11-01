package iot

type Config struct {
	MQTT MQTTConfig `yaml:"mqtt,omitempty"`
}

type MQTTConfig struct {
	URL      string `yaml:"url,omitempty"`
	Topic    string `yaml:"topic,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

package config

// HTTPConfig is the configuration for the listening HTTP server.
type HTTPConfig struct {
	ListenAddress string `yaml:"listen_address,omitempty"`
}

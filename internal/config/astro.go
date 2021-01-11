package config

// AstroConfig is where to find the information.
type AstroConfig struct {
	MetricsURL string   `yaml:"metrics_url"`
	Locations  []string `yaml:"locations"`
	TimeZone   string   `yaml:"timezone"`
}

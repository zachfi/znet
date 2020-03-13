package astro

// Config is where to find the information.
type Config struct {
	MetricsURL string   `yaml:"metrics_url"`
	Locations  []string `yaml:"locations"`
	TimeZone   string   `yaml:"timezone"`
}

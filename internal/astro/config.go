package astro

// Config is where the information necessary to query Prometheus for
// the given metrics that we will keep an eye on.
type Config struct {
	MetricsURL string   `yaml:"metrics_url"`
	Locations  []string `yaml:"locations"`
}

package config

// NetworkConfig are the details needed for the network package.
type NetworkConfig struct {
	ScrapeInterval int         `yaml:"scrape_interval,omitempty"`
	Junos          JunosConfig `yaml:"junos,omitempty"`
}

// JunosConfig is the configuration for Junos devices.
type JunosConfig struct {
	Hosts      []string `yaml:"hosts,omitempty"`
	Username   string   `yaml:"username,omitempty"`
	PrivateKey string   `yaml:"private_key,omitempty"`
}

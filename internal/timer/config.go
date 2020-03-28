package timer

// Config is the information necessary for Timer to generate timers.
type Config struct {
	TimeZone string `yaml:"timezone"`
	Events   []struct {
		// Produce is the name of the event to emit.
		Produce string   `yaml:"produce"`
		Time    string   `yaml:"time"`
		Days    []string `yaml:"days"`
	} `yaml:"events"`
}
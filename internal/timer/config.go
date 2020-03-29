package timer

// Config is the information necessary for Timer to generate timers.
type Config struct {
	TimeZone       string `yaml:"timezone"`
	ReloadInterval int    `yaml:"reload_interval"`
	FutureLimit    int    `yaml:"future_limit"`
	RepeatEvents   []struct {
		// Produce is the name of the event to emit.
		Produce string `yaml:"produce"`
		Every   struct {
			Seconds int `yaml:"seconds"`
		} `yaml:"every"`
	} `yaml:"repeat_events"`
	Events []struct {
		// Produce is the name of the event to emit.
		Produce string   `yaml:"produce"`
		Time    string   `yaml:"time"`
		Days    []string `yaml:"days"`
	} `yaml:"events"`
}

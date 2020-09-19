package agent

// Config is the agent configuration.
type Config struct {
	Executions []Execution `yaml:"executions"`
}

// Execution is a single execution.
type Execution struct {
	// Args is the command arguments to pass for execution.
	Args []string `yaml:"args"`

	// Command is the name of the command to execute.
	Command string `yaml:"command"`

	Dir         string            `yaml:"dir"`
	Environment map[string]string `yaml:"environment"`

	// Events is the slice of names upon which to execute the given executions.
	Events []string `yaml:"events"`

	Filter map[string]interface{} `yaml:"filter"`
}

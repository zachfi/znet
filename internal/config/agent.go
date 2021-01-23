package config

// AgentConfig is the agent configuration.
type AgentConfig struct {
	Executions []AgentExecution `yaml:"executions"`
}

// AgentExecution is a single execution instruction.
type AgentExecution struct {
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

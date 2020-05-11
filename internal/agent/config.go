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

	Filter Filter `yaml:"filter"`
}

// Filter is a way of reducing when the executions fire, in the case of many
// repos.  This does mean the events are handed to all scubscribers ,and the
// subscriber is responsible for reducing the executions.
type Filter struct {
	Names       []string
	URLs        []string
	Branches    []string
	Collections []string
}

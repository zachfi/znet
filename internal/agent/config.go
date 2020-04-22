package agent

// Config is the agent configuration.
type Config struct {
	Executions []Execution
}

// Execution is a single execution.
type Execution struct {
	// Args is the command arguments to pass for execution.
	Args []string

	// Command is the name of the command to execute.
	Command string

	Dir         string
	Environment map[string]string

	// Events is the slice of names upon which to execute the given executions.
	Events []string

	// TODO
	Filter map[string]interface{}

	// TODO
	OnSuccess []string

	// TODO
	OnFailure []string
}

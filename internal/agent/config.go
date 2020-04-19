package agent

// Config is the agent configuration.
type Config struct {
	Executions []Execution
}

// Execution is a single execution.
type Execution struct {
	Args        []string
	Command     string
	Dir         string
	Environment map[string]string
	Event       string
	Filter      map[string]interface{}
}

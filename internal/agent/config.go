package agent

type Config struct {
	Executions []Execution
}

type Execution struct {
	Args        []string
	Command     string
	Dir         string
	Environment map[string]string
	Event       string
}

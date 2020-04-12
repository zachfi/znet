package agent

type Config struct {
	Executions []Execution
}

type Execution struct {
	Event       string
	Shell       string
	Environment map[string]string
}

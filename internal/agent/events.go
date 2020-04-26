package agent

import (
	"time"
)

var EventNames = []string{
	"ExecutionResult",
}

type ExecutionResult struct {
	Command  string
	Args     []string
	Dir      string
	Output   []byte
	Time     *time.Time
	ExitCode int
	Duration time.Duration
}

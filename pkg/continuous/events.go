package continuous

import "time"

var EventNames = []string{
	"BuildResult",
}

type BuildResult struct {
	Command  string
	Args     []string
	Dir      string
	Output   []byte
	Time     *time.Time
	ExitCode int
	Duration time.Duration
}

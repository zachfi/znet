package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/pkg/events"
)

type Agent struct {
	config Config
	conn   *grpc.ClientConn
}

func NewAgent(config Config, conn *grpc.ClientConn) *Agent {
	return &Agent{
		config: config,
		conn:   conn,
	}
}

func (a *Agent) EventNames() []string {
	var names []string

	for _, e := range a.config.Executions {
		for _, x := range e.Events {
			if x != "" {
				names = append(names, x)
			}
		}
	}

	log.Debugf("agent responding to %d event names: %+v", len(names), names)

	return names
}

// Subscriptions implements the events.Consumer interface
func (a *Agent) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	for _, e := range a.config.Executions {
		for _, x := range e.Events {
			switch x {
			case "NewCommit":
				s.Subscribe(x, a.newCommitHandler)
			case "NewTag":
				s.Subscribe(x, a.newTagHandler)
			default:
				log.Errorf("unhandled execution event %s", x)
			}
		}
	}

	log.Debugf("event subscriptions %+v", s.Table)

	return s.Table
}

func (a *Agent) newTagHandler(name string, payload events.Payload) error {

	var x gitwatch.NewTag

	err := json.Unmarshal(payload, &x)
	if err != nil {
		log.Errorf("failed to unmarshal %T: %s", x, err)
	}

	err = a.executeForEvent(x)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (a *Agent) newCommitHandler(name string, payload events.Payload) error {
	log.Debugf("Agent.newCommitHandler: %+v", string(payload))
	log.Debugf("Agent.newCommitHandler config: %+v", a.config)

	var x gitwatch.NewCommit

	err := json.Unmarshal(payload, &x)
	if err != nil {
		log.Errorf("failed to unmarshal %T: %s", x, err)
	}

	err = a.executeForEvent(x)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (a *Agent) passFilter(filter Filter, x interface{}) bool {
	var xName string
	var xURL string
	var xBranch string
	var xCollection string

	t := reflect.TypeOf(x).String()

	switch t {
	case "gitwatch.NewTag":
		xName = x.(gitwatch.NewTag).Name
		xURL = x.(gitwatch.NewTag).URL
		xCollection = x.(gitwatch.NewTag).Collection
	case "gitwatch.NewCommit":
		xName = x.(gitwatch.NewCommit).Name
		xURL = x.(gitwatch.NewCommit).URL
		xBranch = x.(gitwatch.NewCommit).Branch
		xCollection = x.(gitwatch.NewCommit).Collection
	}

	passName := func() bool {
		if len(filter.Names) == 0 {
			return true
		}

		for _, name := range filter.Names {
			if name == xName {
				return true
			}
		}

		return false
	}

	passURL := func() bool {
		if len(filter.URLs) == 0 {
			return true
		}

		for _, url := range filter.URLs {
			if url == xURL {
				return true
			}
		}
		return false
	}

	passBranch := func() bool {
		if len(filter.Branches) == 0 {
			return true
		}

		for _, branch := range filter.Branches {
			if branch == xBranch {
				return true
			}
		}
		return false
	}

	passCollection := func() bool {
		if len(filter.Collections) == 0 {
			return true
		}

		for _, collection := range filter.Collections {
			if collection == xCollection {
				return true
			}
		}
		return false
	}

	return passName() && passURL() && passBranch() && passCollection()
}

func (a *Agent) executeForEvent(x interface{}) error {
	log.Tracef("executeForEvent %+v", x)

	for _, execution := range a.config.Executions {

		if !a.passFilter(execution.Filter, x) {
			return fmt.Errorf("event did not pass filter: %+v, %+v", x, execution.Filter)
		}

		for _, xx := range execution.Events {
			if xx != "" {
				var args []string

				// Render the args as template strings, passing the current x interface.
				for _, v := range execution.Args {
					tmpl, err := template.New("env").Parse(v)
					if err != nil {
						log.Errorf("failed to parse template %s: %s", v, err)
					}

					var buf bytes.Buffer

					err = tmpl.Execute(&buf, x)
					if err != nil {
						log.Error(err)
					}

					args = append(args, buf.String())
				}

				cmd := exec.Command(execution.Command, args...)

				if execution.Dir != "" {
					cmd.Dir = execution.Dir
				}

				var env []string

				// Render the values of the environment variables as templates using the received event.
				for k, v := range execution.Environment {

					tmpl, err := template.New("env").Parse(v)
					if err != nil {
						log.Errorf("failed to parse template %s: %s", v, err)
					}

					var buf bytes.Buffer

					err = tmpl.Execute(&buf, x)
					if err != nil {
						log.Error(err)
					}

					env = append(env, fmt.Sprintf("%s=%s", k, buf.String()))
				}

				if len(env) > 0 {
					cmd.Env = append(os.Environ(), env...)
				}

				start := time.Now()
				// var out bytes.Buffer
				// cmd.Stdout = &out
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Errorf("command execution failed: %s", err)
				}

				now := time.Now()

				ev := ExecutionResult{
					Time:     &now,
					Command:  execution.Command,
					Args:     args,
					Dir:      execution.Dir,
					Output:   output,
					ExitCode: cmd.ProcessState.ExitCode(),
					Duration: time.Since(start),
				}

				err = events.ProduceEvent(a.conn, ev)
				if err != nil {
					log.Error(err)
				}

			}
		}
	}

	return nil
}

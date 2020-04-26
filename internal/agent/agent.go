package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/xaque208/znet/rpc"

	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/gitwatch"
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

	filterHasName := func() bool {
		for _, name := range filter.Names {
			if name == xName {
				return true
			}
		}
		return false
	}

	filterHasURL := func() bool {
		for _, url := range filter.URLs {
			if url == xURL {
				return true
			}
		}
		return false
	}

	filterHasBranch := func() bool {
		for _, branch := range filter.Branches {
			if branch == xBranch {
				return true
			}
		}
		return false
	}

	if len(filter.Names) > 0 {
		if !filterHasName() {
			return false
		}
	}

	if len(filter.URLs) > 0 {
		if !filterHasURL() {
			return false
		}
	}

	if len(filter.Branches) > 0 {
		if !filterHasBranch() {
			return false
		}
	}

	return true
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

				// Render the args as template strings
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
				}

				err = a.Produce(ev)
				if err != nil {
					log.Error(err)
				}

			}
		}
	}

	return nil
}

// Produce implements the events.Producer interface.  Match the supported event
// types to know which event to notice, and then send notice of the event to
// the RPC server.
func (a *Agent) Produce(ev interface{}) error {
	// Create the RPC client
	ec := pb.NewEventsClient(a.conn)
	t := reflect.TypeOf(ev).String()

	var req *pb.Event

	switch t {
	case "agent.ExecutionResult":
		x := ev.(ExecutionResult)
		req = events.MakeEvent(x)
	default:
		return fmt.Errorf("unhandled event type: %T", ev)
	}

	log.Tracef("agent producing RPC event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

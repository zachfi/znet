package agent

import (
	"bytes"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
)

type Agent struct {
	config Config
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
	}

}

func (a *Agent) EventNames() []string {
	var names []string

	for _, e := range a.config.Executions {
		if e.Event != "" {
			names = append(names, e.Event)
		}
	}

	log.Debugf("agent responding to %d event names: %+v", len(names), names)

	return names
}

// Subscriptions implements the events.Consumer interface
func (a *Agent) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	for _, e := range a.config.Executions {
		if e.Event != "" {
			s.Subscribe(e.Event, a.eventHandler)
		}
	}

	log.Debugf("event subscriptions %+v", s.Table)

	return s.Table
}

func (a *Agent) eventHandler(name string, payload events.Payload) error {
	log.Debugf("Agent.eventHandler: %+v", string(payload))
	log.Debugf("Agent.eventHandler config: %+v", a.config)

	for _, e := range a.config.Executions {
		if e.Event != "" {
			cmd := exec.Command(e.Command, e.Args...)

			if e.Dir != "" {
				cmd.Dir = e.Dir
			}

			var env []string

			for k, v := range e.Environment {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}

			// cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

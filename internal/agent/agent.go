package agent

import (
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

	return nil
}

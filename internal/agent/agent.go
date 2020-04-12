package agent

import (
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
)

type Agent struct {
	config Config
}

// Subscriptions implements the events.Consumer interface
func (a *Agent) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	s.Subscribe("NewTag", a.eventHandler)
	s.Subscribe("NewCommit", a.eventHandler)

	return s.Table
}

func (a *Agent) eventHandler(payload events.Payload) error {
	log.Debugf("Agent.eventHandler: %+v", string(payload))
	log.Debugf("Agent.eventHandler config: %+v", a.config)

	return nil
}

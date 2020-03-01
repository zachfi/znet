package lights

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/events"
)

// Handler is the EventHandler for astro events.
var Handler events.EventHandler = &LightsHandler{}

type LightsHandler struct {
	config events.HandlerConfig
}

func New() events.EventHandler {
	return &LightsHandler{}
}

func (h *LightsHandler) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	s.Subscribe("Sunrise", sunriseHandler())
	// s.Subscribe("Sunset", sunsetHandler())
	// s.Subscribe("Bedtime", bedtimeHandler())
	// s.Subscribe("Sleep", sleepHandler())

	return s.Table
}

func (h *LightsHandler) Config() (events.HandlerConfig, error) {
	return events.HandlerConfig{
		Name: "Lights",
		Events: []string{
			"Sunset",
			"Sunrise",
			"Bedtime",
			"Sleep",
		},
	}, nil
}

func (h *LightsHandler) Handler() events.Handler {
	return func(bits []byte) error {
		var t time.Time

		err := json.Unmarshal(bits, &t)
		if err != nil {
			return fmt.Errorf("failed to unmarshal %T: %s", t, err)
		}

		log.Infof("Got %s", t)

		return nil
	}
}

func sunriseHandler() events.Handler {
	return func(bits []byte) error {

		return nil
	}
}

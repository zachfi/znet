package timer

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/events"
)

type TimerHandler struct {
	config events.HandlerConfig
}

func (h *TimerHandler) Subscriptions() map[string][]func([]byte) error {
	subscriptions := make(map[string][]func([]byte) error)

	return subscriptions
}

func (h *TimerHandler) Handler() events.Handler {
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

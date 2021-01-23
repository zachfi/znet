package gitwatch

import (
	"testing"

	"github.com/xaque208/znet/pkg/events"
)

func TestProducer_interface(t *testing.T) {
	var a events.Producer = &EventProducer{}
	t.Logf("eventProducer: %+v", a)
}

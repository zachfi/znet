package znet

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/things/things"
)

// Listener is a znet server
type Listener struct {
	Config     *Config
	httpServer *http.Server
}

// NewListener builds a new Listener object from the received configuration.
func NewListener(config *Config) (*Listener, error) {
	l := &Listener{
		Config: config,
	}

	return l, nil
}

// Listen starts the http listener
func (l *Listener) Listen(listenAddr string, ch chan bool) {
	log.Infof("Listening on %s", listenAddr)
	l.httpServer = httpListen(listenAddr)

	messages := make(chan things.Message)
	go l.messageHandler(messages)
	// go l.thingServer.Listen(messages)

	<-ch
	l.Shutdown()
}

// Shutdown closes down the to the message bus and shuts down the HTTP server.
func (l *Listener) Shutdown() {
	log.Info("ZNET Shutting Down")

	log.Info("halting HTTP server")
	err := l.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Error(err)
	}
}

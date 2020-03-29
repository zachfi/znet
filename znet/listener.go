package znet

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
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

// Listen starts the http listener and blocks listening for any message on the
// given channel, shutting down the HTTP.
func (l *Listener) Listen(listenAddr string, ch chan bool) {
	log.Infof("starting znet listener %s", listenAddr)
	l.httpServer = httpListen(listenAddr)

	<-ch
	l.Shutdown()
}

// Shutdown closes down the to the message bus and shuts down the HTTP server.
func (l *Listener) Shutdown() {
	log.Info("znet HTTP listener shutting down")
	err := l.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Error(err)
	}
}

package things

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	nats "github.com/nats-io/go-nats"
)

type MessageHandler func(chan Message)

type Server struct {
	URL         string
	Topic       string
	Conn        nats.Conn
	EncodedConn nats.EncodedConn
}

func NewServer(url, topic string) (*Server, error) {
	server := Server{
		URL:   url,
		Topic: topic,
	}

	natsConn, err := nats.Connect(url)
	if err != nil {
		return &server, err
	}

	eConn, err := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	if err != nil {
		return &server, err
	}

	server.Conn = *natsConn
	server.EncodedConn = *eConn

	return &server, nil
}

func (s Server) Listen(messages chan Message, callback MessageHandler) error {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go callback(messages)
	s.EncodedConn.Subscribe(s.Topic, func(msg *Message) {
		messages <- *msg
	})

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())
		s.Close()

		done <- true
	}()

	<-done

	return nil
}

func (s Server) Close() {
	s.EncodedConn.Flush()

	if !s.Conn.IsClosed() {
		log.Debug("Closing nats connection")
		s.Conn.Close()
	}
}

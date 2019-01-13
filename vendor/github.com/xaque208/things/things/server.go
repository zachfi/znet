package things

import (
	"errors"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

type MessageHandler func(chan Message)

type Server struct {
	URL         string
	Topic       string
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
}

func NewServer(url, topic string) (*Server, error) {
	if topic == "" || url == "" {
		return &Server{}, errors.New("NATS URL and Topic are required")
	}

	server := &Server{
		URL:   url,
		Topic: topic,
	}

	log.Debugf("Things Server using nats: %+v", server)

	var err error

	server.Conn, err = nats.Connect(url)
	if err != nil {
		return &Server{}, err
	}

	server.EncodedConn, err = nats.NewEncodedConn(server.Conn, nats.JSON_ENCODER)
	if err != nil {
		return &Server{}, err
	}

	return server, nil
}

// Listen
func (s Server) Listen(messages chan Message) {
	s.EncodedConn.Subscribe(s.Topic, func(msg *Message) {
		log.Debugf("%+v", *msg)
		messages <- *msg
	})
}

func (s Server) Close() {
	s.EncodedConn.Flush()

	if !s.Conn.IsClosed() {
		log.Debug("Closing nats connection")
		s.EncodedConn.Close()
	}
}

package things

import (
	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
)

type Client struct {
	Conn        nats.Conn
	EncodedConn nats.EncodedConn
}

func NewClient(url, topic string) (*Client, error) {
	client := Client{}

	natsConn, err := nats.Connect(url)
	if err != nil {
		return &client, err
	}

	eConn, err := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	if err != nil {
		return &client, err
	}

	client.Conn = *natsConn
	client.EncodedConn = *eConn

	return &client, nil
}

func (c *Client) Close() {
	c.EncodedConn.Flush()

	if !c.Conn.IsClosed() {
		log.Debug("Closing nats connection")
		c.Conn.Close()
	}
}

package things

import (
	"errors"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
)

// Client represents a Thigns client, holding the decoding and connection to the NATS server.
type Client struct {
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
}

// NewClient creates and returns a new Things client.
func NewClient(url, topic string) (*Client, error) {
	if topic == "" || url == "" {
		return &Client{}, errors.New("NATS URL and Topic are required")
	}

	client := &Client{}
	log.Debugf("Things Server using nats: %+v", client)

	var err error

	client.Conn, err = nats.Connect(url)
	if err != nil {
		return &Client{}, err
	}

	client.EncodedConn, err = nats.NewEncodedConn(client.Conn, nats.JSON_ENCODER)
	if err != nil {
		return &Client{}, err
	}

	return client, nil
}

// Close closes the connection to the server.
func (c *Client) Close() {
	c.EncodedConn.Flush()

	if !c.Conn.IsClosed() {
		log.Debug("Closing nats connection")
		c.Conn.Close()
	}
}

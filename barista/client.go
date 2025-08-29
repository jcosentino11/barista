package barista

import (
	"net"
	"strconv"
)

type Client struct {
	Config ClientConfig
	writer NetworkWriter
	logger Logger
}

// TODO convert to interface
type ClientConfig struct {
	ServerHost string
	ServerPort int
}

func NewClient(config ClientConfig) Client {
	logger := NewConsoleLogger("server")
	return Client{
		Config: config,
		logger: logger,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("udp", net.JoinHostPort(c.Config.ServerHost, strconv.Itoa(c.Config.ServerPort)))
	if err != nil {
		return err
	}
	writer := NewNetworkWriter(conn)
	c.writer = &writer
	return nil
}

func (c *Client) Subscribe(topic string, messageCallback func(string, string)) error {
	packet := &SubscribePacket{
		Topic: topic,
	}
	// TODO register callback, wait for acknowledgement
	return c.writer.Write(packet.Bytes())
}

func (c *Client) Publish(topic string, message string) error {
	packet := &PublishPacket{
		Topic:   topic,
		Content: message,
	}
	return c.writer.Write(packet.Bytes())
}

// TODO make thread safe?
func (c *Client) Close() error {
	if c.writer != nil {
		return c.writer.Close()
	}
	return nil
}

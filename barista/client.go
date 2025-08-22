package barista

import (
	"net"
	"strconv"
)

type Client struct {
	Config ClientConfig
	writer PacketWriter
}

// TODO convert to interface
type ClientConfig struct {
	ServerHost string
	ServerPort int
}

func NewClient(config ClientConfig) Client {
	return Client{
		Config: config,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("udp", net.JoinHostPort(c.Config.ServerHost, strconv.Itoa(c.Config.ServerPort)))
	if err != nil {
		return err
	}
	writer := NewNetworkPacketWriter(conn)
	c.writer = &writer
	return nil
}

func (c *Client) Subscribe(topic string, messageCallback func(string, string)) error {
	packet := &SubscribePacket{
		Topic: topic,
	}
	// TODO register callback, wait for acknowledgement
	return c.writer.Write(packet)
}

func (c *Client) Publish(topic string, message string) error {
	packet := &PublishPacket{
		Topic:   topic,
		Content: message,
	}
	return c.writer.Write(packet)
}

// TODO make thread safe?
func (c *Client) Close() error {
	if c.writer != nil {
		return c.writer.Close()
	}
	return nil
}

package barista

import (
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	Port int
}

func NewClient() Client {
	return Client{
		Port: 8080,
	}
}

func (c *Client) sendDatagram(message string) error {
	conn, err := net.Dial("udp", net.JoinHostPort("localhost", strconv.Itoa(c.Port)))
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	return err
}

func (c *Client) Subscribe(topic string, messageCallback func(string, string)) error {
	err := c.sendDatagram("SUBSCRIBE " + topic)
	if err != nil {
		return fmt.Errorf("unable to subscribe to topic %s: %w", topic, err)
	}
	return nil
}

func (c *Client) Publish(topic string, message string) error {
	err := c.sendDatagram("PUBLISH " + topic + " " + message)
	if err != nil {
		return fmt.Errorf("unable to publish to topic %s: %w", topic, err)
	}
	return nil
}

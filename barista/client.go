package barista

import (
	"context"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	Config ClientConfig

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	conn      net.Conn
	processor PacketStreamProcessor
	writer    NetworkWriter

	logger Logger
}

// TODO convert to interface
type ClientConfig struct {
	ServerHost string
	ServerPort int
}

func NewClient(config ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	logger := NewConsoleLogger("client")
	client := Client{
		Config: config,
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
	return &client
}

func (c *Client) Start() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	c.conn = conn

	stream, err := NewNetworkPacketStreamFromConn(c.ctx, conn)
	if err != nil {
		return err
	}

	c.processor = NewPacketStreamProcessor(c.ctx, stream, c.handlePacket)

	writer := NewNetworkConnWriter(conn)
	c.writer = &writer

	return nil
}

func (c *Client) connect() (net.Conn, error) {
	conn, err := net.Dial("udp", net.JoinHostPort(c.Config.ServerHost, strconv.Itoa(c.Config.ServerPort)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Client) Stop() error {
	c.cancel()
	c.wg.Wait()
	return c.closeConn()
}

func (c *Client) closeConn() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (s *Client) handlePacket(packet PacketResult) error {
	s.logger.Printf("Received packet: %v\n", packet.Packet)
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

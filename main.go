package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	server := NewServer()

	err := server.Start()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	defer func() {
		server.Stop()
	}()

	messageCallback := func(topic string, message string) {
		fmt.Printf("Received message on topic '%s': %s\n", topic, message)
	}

	client := NewClient()

	err = client.Subscribe("topic", messageCallback)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	err = client.Publish("topic", "test")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	time.Sleep(100 * time.Millisecond)
}

type client struct {
	id     string
	topics map[string]bool
}

type Server struct {
	Port     int
	conn     *net.UDPConn
	clients  map[string]client
	shutdown chan bool
	parser   Parser
}

func NewServer() Server {
	return Server{
		Port:     8080,
		clients:  make(map[string]client),
		shutdown: make(chan bool),
		parser:   &DefaultParser{},
	}
}

func (s *Server) Start() error {
	if s.conn != nil {
		return fmt.Errorf("server is already listening at %s", s.conn.LocalAddr())
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.Port})
	if err != nil {
		return fmt.Errorf("unable to start server on port %d: %w", s.Port, err)
	}
	s.conn = conn

	go s.receiveDatagrams()

	return nil
}

func (s *Server) receiveDatagrams() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-s.shutdown:
			return
		default:
			n, addr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Error reading from UDP: %s\n", err)
				continue
			}
			packet, err := s.parser.Parse(buf[:n])
			if err != nil {
				fmt.Printf("packet parsing failed: %w", err)
				continue
			}
			fmt.Printf("Received from %s: %s\n", addr, packet)
		}
	}
}

func (s *Server) Stop() error {
	close(s.shutdown)
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

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

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
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

type Server struct {
	Port   int
	conn   *net.UDPConn
	parser Parser
	wg     sync.WaitGroup
}

func NewServer() Server {
	return Server{
		Port:   8080,
		parser: &DefaultParser{},
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

	s.wg.Add(1)
	go s.receiveDatagrams()

	return nil
}

func (s *Server) receiveDatagrams() {
	defer s.wg.Done()

	buf := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
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

func (s *Server) Stop() error {
	if s.conn != nil {
		s.conn.Close()
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Printf("shutdown complete\n")
	case <-time.After(2 * time.Second):
		fmt.Printf("shutdown timed out\n")
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

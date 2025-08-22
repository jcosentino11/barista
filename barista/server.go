package barista

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
)

const (
	ErrServerClosed        = "server is closed"
	ErrServerRunning       = "server is running"
	ErrServerContextClosed = "server context is closed"
)

type Server struct {
	Config ServerConfig
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type ServerConfig struct {
	Port int
}

func NewServer(config ServerConfig) Server {
	ctx, cancel := context.WithCancel(context.Background())
	return Server{
		Config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) Start() error {
	select {
	case <-s.ctx.Done():
		return errors.New(ErrServerClosed)
	default:
		s.wg.Add(1)
		go s.worker()
		return nil
	}
}

func (s *Server) worker() {
	defer s.wg.Done()

	reader, err := s.newPacketReader()
	if err != nil {
		fmt.Printf("unable to create reader: %s\n", err)
		return
	}

	packets := reader.Packets()

	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Printf("unable to close reader: %s\n", err)
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case packet, ok := <-packets:
			if !ok {
				fmt.Println("reader channel closed, exiting worker")
				return
			}
			if err := s.handlePacket(packet); err != nil {
				fmt.Printf("err handling packet: %s", err)
			}
		}
	}
}

func (s *Server) newPacketReader() (PacketReader, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.Config.Port})
	if err != nil {
		return nil, fmt.Errorf("unable to start server on port %d: %w", s.Config.Port, err)
	}

	reader := NewUdpPacketReader(conn)
	return &reader, nil
}

// TODO pass this in to server
func (s *Server) handlePacket(packet PacketResult) error {
	fmt.Printf("Received packet: %s\n", packet)
	return nil
}

func (s *Server) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

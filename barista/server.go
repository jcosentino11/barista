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
	Config    ServerConfig
	processor PacketProcessor
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	logger    Logger
}

type ServerConfig struct {
	Port int
}

func NewServer(config ServerConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	// TODO configurable, log to file
	logger := NewConsoleLogger("server")
	server := Server{
		Config: config,
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
	server.processor = NewPacketProcessor(server.newPacketReader, server.handlePacket)
	return &server
}

func (s *Server) Start() error {
	select {
	case <-s.ctx.Done():
		return errors.New(ErrServerClosed)
	default:
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.processor.ProcessPackets(s.ctx); err != nil {
				s.cancel()
			}
		}()
		return nil
	}
}

func (s *Server) newPacketReader() (PacketReader, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.Config.Port})
	if err != nil {
		return nil, fmt.Errorf("unable to start server on port %d: %w", s.Config.Port, err)
	}

	reader := NewUdpPacketReader(s.ctx, conn)
	return &reader, nil
}

// TODO pass this in to server
func (s *Server) handlePacket(packet PacketResult) error {
	s.logger.Printf("Received packet: %v\n", packet)
	return nil
}

func (s *Server) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

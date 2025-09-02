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

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	conn      net.Conn
	processor PacketStreamProcessor
	writer    NetworkWriter

	logger Logger
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
	return &server
}

func (s *Server) Start() error {
	select {
	case <-s.ctx.Done():
		return errors.New(ErrServerClosed)
	default:
		conn, err := s.listen()
		if err != nil {
			return err
		}
		s.conn = conn

		stream, err := NewNetworkPacketStreamFromConn(s.ctx, conn)
		if err != nil {
			return err
		}

		s.processor = NewPacketStreamProcessor(s.ctx, stream, s.handlePacket)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.processor.ProcessPackets(); err != nil {
				s.Stop()
			}
		}()

		return nil
	}
}

func (s *Server) listen() (net.Conn, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.Config.Port})
	if err != nil {
		return nil, fmt.Errorf("unable to start server on port %d: %w", s.Config.Port, err)
	}
	return conn, nil
}

func (s *Server) Stop() error {
	s.cancel()
	s.wg.Wait()
	return s.closeConn()
}

func (s *Server) closeConn() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// TODO pass this in to server
func (s *Server) handlePacket(packet PacketResult) error {
	s.logger.Printf("Received packet: %v\n", packet.Packet)
	return nil
}

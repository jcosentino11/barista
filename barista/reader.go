package barista

import (
	"context"
	"errors"
	"net"
	"sync"
)

const (
	ErrReaderClosed  = "reader is closed"
	ErrReaderRunning = "reader is running"
	ErrContextClosed = "context is closed"
)

type PacketResult struct {
	Packet Packet
	Err    error
}

type PacketReader interface {
	Packets() (<-chan PacketResult, error)
	Close() error
}

type UdpPacketReader struct {
	conn    *net.UDPConn
	parser  PacketParser
	packets chan PacketResult
	ctx     context.Context
	wg      sync.WaitGroup
	logger  Logger
}

func NewUdpPacketReader(ctx context.Context, conn *net.UDPConn) UdpPacketReader {
	logger := NewConsoleLogger("udp-packet-reader")
	logger.Verbose = true // TODO
	return UdpPacketReader{
		conn:   conn,
		parser: &DefaultParser{},
		// TODO set bounds, handle backpressure
		packets: make(chan PacketResult),
		ctx:     ctx,
		logger:  logger,
	}
}

func (r *UdpPacketReader) Packets() (<-chan PacketResult, error) {
	select {
	case <-r.ctx.Done():
		return nil, errors.New(ErrContextClosed)
	default:
		r.wg.Add(1)
		go r.readPackets()
		return r.packets, nil
	}
}

func (r *UdpPacketReader) readPackets() {
	defer r.wg.Done()

	buffer := make([]byte, 1024)
	for {
		select {
		case <-r.ctx.Done():
			r.logger.Verbosef("closed ctx detected\n")
			return
		default:
			bytesRead, _, err := r.conn.ReadFromUDP(buffer)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				r.logger.Printf("error reading from UDP: %s\n", err)
				continue
			}
			packet, err := r.parser.Parse(buffer[:bytesRead])
			r.packets <- PacketResult{
				Packet: packet,
				Err:    err,
			}
		}
	}
}

func (r *UdpPacketReader) Close() error {
	close(r.packets)
	r.logger.Verbosef("packets channel closed\n")

	if err := r.closeConnection(); err != nil {
		r.logger.Printf("failed to close connection: %s\n", err)
	}
	r.logger.Verbosef("connection closed\n")
	r.wg.Wait()
	return nil
}

func (r *UdpPacketReader) closeConnection() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

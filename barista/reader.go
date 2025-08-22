package barista

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
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
	Packets() <-chan PacketResult
	Start() error
	Close() error
}

type UdpPacketReader struct {
	conn    *net.UDPConn
	parser  PacketParser
	packets chan PacketResult
	closed  atomic.Bool
	once    sync.Once
	wg      sync.WaitGroup
	logger  Logger
}

func NewUdpPacketReader(conn *net.UDPConn) UdpPacketReader {
	logger := NewConsoleLogger("udppacketreader")
	return UdpPacketReader{
		conn:   conn,
		parser: &DefaultParser{},
		// TODO set bounds, handle backpressure
		packets: make(chan PacketResult),
		logger:  logger,
	}
}

func (r *UdpPacketReader) Packets() <-chan PacketResult {
	return r.packets
}

// TODO use context to manage lifecycle?
func (r *UdpPacketReader) Start() error {
	if r.closed.Load() {
		return errors.New(ErrReaderClosed)
	}

	r.once.Do(r.readPacketsAsync)

	return nil
}

func (r *UdpPacketReader) Close() error {
	if !r.closed.CompareAndSwap(false, true) {
		return nil
	}

	if err := r.closeConnection(); err != nil {
		return err
	}

	r.wg.Wait() // TODO consider timeout

	close(r.packets)

	return nil
}

func (r *UdpPacketReader) closeConnection() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// TODO feels awkward to have this as it's own method
func (r *UdpPacketReader) readPacketsAsync() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.readPackets()
	}()
}

func (r *UdpPacketReader) readPackets() {
	buffer := make([]byte, 1024)
	for {
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

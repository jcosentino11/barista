package barista

import (
	"context"
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

type PacketStream interface {
	Stream() (<-chan PacketResult, error)
	Close() error
}

type NetworkPacketStream struct {
	networkReader NetworkReader
	parser        PacketParser
	packets       chan PacketResult
	ctx           context.Context
	wg            sync.WaitGroup
	logger        Logger
}

func NewNetworkPacketStreamFromConn(ctx context.Context, conn net.Conn) (PacketStream, error) {
	reader := NewDefaultNetworkReader(conn)
	stream := NewNetworkPacketStream(ctx, &reader)
	return &stream, nil
}

func NewNetworkPacketStream(ctx context.Context, networkReader NetworkReader) NetworkPacketStream {
	logger := NewConsoleLogger("network-packet-stream")
	logger.Verbose = false // TODO

	parser := NewDefaultParser()

	return NetworkPacketStream{
		networkReader: networkReader,
		parser:        &parser,
		// TODO set bounds, handle backpressure
		packets: make(chan PacketResult),
		ctx:     ctx,
		logger:  logger,
	}
}

func (r *NetworkPacketStream) Stream() (<-chan PacketResult, error) {
	r.wg.Add(1)
	go r.readPackets()
	return r.packets, nil
}

func (r *NetworkPacketStream) readPackets() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			// TODO handle cleanup better
			r.logger.Verbosef("closed ctx detected\n")
			return
		default:
			buf, err := r.networkReader.Bytes()
			if err != nil {
				continue
			}
			packet, err := r.parser.Parse(buf)
			r.packets <- PacketResult{
				Packet: packet,
				Err:    err,
			}
		}
	}
}

func (r *NetworkPacketStream) Close() error {
	close(r.packets)
	r.wg.Wait()
	return nil
}

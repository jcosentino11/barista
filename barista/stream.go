package barista

import (
	"context"
	"errors"
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

func NewNetworkPacketStream(ctx context.Context, networkReader NetworkReader) NetworkPacketStream {
	logger := NewConsoleLogger("default-packet-reader")
	logger.Verbose = true // TODO
	return NetworkPacketStream{
		networkReader: networkReader,
		parser:        &DefaultParser{},
		// TODO set bounds, handle backpressure
		packets: make(chan PacketResult),
		ctx:     ctx,
		logger:  logger,
	}
}

func (r *NetworkPacketStream) Stream() (<-chan PacketResult, error) {
	select {
	case <-r.ctx.Done():
		return nil, errors.New(ErrContextClosed)
	default:
		r.wg.Add(1)
		go r.readPackets()
		return r.packets, nil
	}
}

func (r *NetworkPacketStream) readPackets() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Verbosef("closed ctx detected\n")
			return
		default:
			buf, err := r.networkReader.Bytes()
			if err != nil {
				r.logger.Printf("error reading from UDP: %s\n", err)
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
	if err := r.networkReader.Close(); err != nil {
		return err
	}
	r.wg.Wait()
	return nil
}

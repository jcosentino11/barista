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

type PacketReader interface {
	Packets() (<-chan PacketResult, error)
	Close() error
}

type DefaultPacketReader struct {
	networkReader NetworkReader
	parser        PacketParser
	packets       chan PacketResult
	ctx           context.Context
	wg            sync.WaitGroup
	logger        Logger
}

func NewDefaultPacketReader(ctx context.Context, networkReader NetworkReader) DefaultPacketReader {
	logger := NewConsoleLogger("default-packet-reader")
	logger.Verbose = true // TODO
	return DefaultPacketReader{
		networkReader: networkReader,
		parser:        &DefaultParser{},
		// TODO set bounds, handle backpressure
		packets: make(chan PacketResult),
		ctx:     ctx,
		logger:  logger,
	}
}

func (r *DefaultPacketReader) Packets() (<-chan PacketResult, error) {
	select {
	case <-r.ctx.Done():
		return nil, errors.New(ErrContextClosed)
	default:
		r.wg.Add(1)
		go r.readPackets()
		return r.packets, nil
	}
}

func (r *DefaultPacketReader) readPackets() {
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

func (r *DefaultPacketReader) Close() error {
	close(r.packets)
	if err := r.networkReader.Close(); err != nil {
		return err
	}
	r.wg.Wait()
	return nil
}

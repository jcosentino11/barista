package barista

import "context"

type PacketStreamProcessor struct {
	ctx     context.Context
	stream  PacketStream
	handler func(PacketResult) error
	logger  Logger
}

func NewPacketStreamProcessor(
	ctx context.Context,
	stream PacketStream,
	handler func(PacketResult) error) PacketStreamProcessor {

	logger := NewConsoleLogger("packet-processor")
	logger.Verbose = false // TODO
	return PacketStreamProcessor{
		ctx:     ctx,
		stream:  stream,
		handler: handler,
		logger:  logger,
	}
}

func (p *PacketStreamProcessor) ProcessPackets() error {
	defer p.stream.Close()

	packets, err := p.stream.Stream()
	if err != nil {
		p.logger.Printf("unable to get reader packets: %s\n", err)
		return err
	}

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Verbosef("closed ctx detected\n")
			return nil
		case packet, ok := <-packets:
			if !ok {
				p.logger.Printf("reader channel closed, exiting worker\n")
				return nil
			}
			if err := p.handler(packet); err != nil {
				p.logger.Printf("err handling packet: %s\n", err)
			}
		}
	}
}

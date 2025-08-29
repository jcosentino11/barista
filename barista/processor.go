package barista

import "context"

type PacketProcessor struct {
	newReader func() (PacketReader, error)
	handler   func(PacketResult) error
	logger    Logger
}

func NewPacketProcessor(
	newReader func() (PacketReader, error),
	handler func(PacketResult) error) PacketProcessor {

	logger := NewConsoleLogger("packet-processor")
	logger.Verbose = false // TODO
	return PacketProcessor{
		newReader: newReader,
		logger:    logger,
		handler:   handler,
	}
}

func (w *PacketProcessor) ProcessPackets(ctx context.Context) error {
	reader, err := w.newReader()
	if err != nil {
		w.logger.Printf("unable to create reader: %s\n", err)
		return err
	}

	defer reader.Close()

	packets, err := reader.Packets()
	if err != nil {
		w.logger.Printf("unable to get reader packets: %s\n", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Verbosef("closed ctx detected\n")
			return nil
		case packet, ok := <-packets:
			if !ok {
				w.logger.Printf("reader channel closed, exiting worker\n")
				return nil
			}
			if err := w.handler(packet); err != nil {
				w.logger.Printf("err handling packet: %s\n", err)
			}
		}
	}
}

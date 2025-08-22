package barista

import (
	"net"
)

type PacketWriter interface {
	Write(Packet) error
	Close() error
}

type NetworkPacketWriter struct {
	conn net.Conn
}

func NewNetworkPacketWriter(conn net.Conn) NetworkPacketWriter {
	return NetworkPacketWriter{
		conn: conn,
	}
}

func (w *NetworkPacketWriter) Write(packet Packet) error {
	// TODO handle partial writes?
	_, err := w.conn.Write(packet.Bytes())
	return err
}

// TODO make thread safe?
func (w *NetworkPacketWriter) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

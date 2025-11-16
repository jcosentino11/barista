package barista

import (
	"net"
)

type NetworkWriter interface {
	Write(buf []byte) error
}

type NetworkConnWriter struct {
	conn net.Conn
}

func NewNetworkConnWriter(conn net.Conn) NetworkConnWriter {
	return NetworkConnWriter{
		conn: conn,
	}
}

func (w *NetworkConnWriter) Write(buf []byte) error {
	// TODO handle partial writes?
	_, err := w.conn.Write(buf)
	return err
}

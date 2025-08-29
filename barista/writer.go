package barista

import (
	"net"
)

type NetworkWriter interface {
	Write(buf []byte) error
	Close() error
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

// TODO make thread safe?
func (w *NetworkConnWriter) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

package barista

import (
	"net"
	"time"
)

type NetworkReader interface {
	Bytes() ([]byte, error)
}

type DefaultNetworkReader struct {
	conn net.Conn
	buf  []byte
}

func NewDefaultNetworkReader(conn net.Conn) DefaultNetworkReader {
	return DefaultNetworkReader{
		conn: conn,
		buf:  make([]byte, 1024),
	}
}

func (r *DefaultNetworkReader) Bytes() ([]byte, error) {
	r.conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // TODO configurable
	bytesRead, err := r.conn.Read(r.buf)
	if err != nil {
		return nil, err
	}
	return r.buf[:bytesRead], nil
}

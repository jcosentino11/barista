package barista

import (
	"net"
)

type NetworkReader interface {
	Bytes() ([]byte, error)
	Close() error
}

type UdpNetworkReader struct {
	conn *net.UDPConn
	buf  []byte
}

func NewUdpNetworkReader(conn *net.UDPConn) UdpNetworkReader {
	return UdpNetworkReader{
		conn: conn,
		buf:  make([]byte, 1024),
	}
}

func (r *UdpNetworkReader) Bytes() ([]byte, error) {
	bytesRead, _, err := r.conn.ReadFromUDP(r.buf)
	if err != nil {
		return nil, err
	}
	return r.buf[:bytesRead], nil
}

func (r *UdpNetworkReader) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

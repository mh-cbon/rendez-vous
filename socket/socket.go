package socket

import (
	"net"
	"sync"
)

// FromAddr is a ctor
func FromAddr(address string) (Socket, error) {
	protocol := "udp"

	udpAddr, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		return nil, err
	}
	logger.Info("listening on ", udpAddr.String())
	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		return nil, err
	}
	return FromConn(conn), nil
}

// FromConn is a ctor
func FromConn(conn net.PacketConn) Socket {
	t := &Tx{
		UDP: UDP{conn},
		l:   sync.Mutex{},
	}
	t.init()
	return t
}

// Socket with transaction support
type Socket interface {
	Listen(queryHandler TxHandler) error
	Query(data []byte, remote net.Addr, h ResponseHandler) error
	Reply(data []byte, remote net.Addr, txID uint16) error
	Close() error
	Conn() net.PacketConn
}

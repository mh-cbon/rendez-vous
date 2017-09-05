package socket

import (
	"net"
	"sync"
)

// Server handles queries/responses
type Server struct {
	Socket
	Handler TxHandler
}

// ListenAndServe the queries/responses
func (s *Server) ListenAndServe() error {
	return s.Socket.Listen(s.Handler)
}

// Handle set the queries/responses handler
func (s *Server) Handle(h TxHandler) *Server {
	s.Handler = h
	return s
}

// Close the server and the underlying socket.
func (s *Server) Close() error { return s.Socket.Close() }

// Socket with transaction support
type Socket interface {
	Listen(queryHandler TxHandler) error
	Query(data []byte, remote net.Addr, h ResponseHandler) error
	Reply(data []byte, remote net.Addr, txID uint16) error
	Close() error
	Conn() net.PacketConn
}

// FromAddr is a ctor
func FromAddr(address string) (*Server, error) {
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
func FromConn(conn net.PacketConn) *Server {
	t := &Tx{
		UDP: UDP{conn},
		l:   sync.Mutex{},
	}
	t.init()
	return &Server{
		Socket: t,
	}
}

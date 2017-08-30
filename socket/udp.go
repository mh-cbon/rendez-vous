package socket

import (
	"fmt"
	"io"
	"log"
	"net"
)

// UDP is a struct
type UDP struct {
	conn *net.UDPConn
}

// FromAddr is a ctor
func FromAddr(address string) (*UDP, error) {
	protocol := "udp"

	udpAddr, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		return nil, err
	}
	log.Println("listening on ", udpAddr.String())
	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		return nil, err
	}
	return &UDP{conn: conn}, nil
}

// FromConn is a ctor
func FromConn(conn *net.UDPConn) *UDP { return &UDP{conn: conn} }

// Handler handle incoming messages
type Handler func(data []byte, remote Socket) error

// Socket writes to a remote
type Socket struct {
	remote *net.UDPAddr
	conn   *net.UDPConn
}

func (s Socket) Write(data []byte) error {
	_, err := s.conn.WriteTo(data, s.remote)
	return err
}

//Addr of the remote
func (s Socket) Addr() string {
	return s.remote.String()
}

//Close the socket
func (u *UDP) Close() error {
	return u.conn.Close()
}

// Listen invoke process when a new message income
func (u *UDP) Listen(h Handler) error {

	conn := u.conn
	var b [0x10000]byte
	for {
		n, addr, readErr := conn.ReadFromUDP(b[:])
		if readErr == nil && n == len(b) {
			readErr = fmt.Errorf("received packet exceeds buffer size %q", len(b))
		}

		if readErr != nil {
			if x, ok := readErr.(*net.OpError); ok && x.Temporary() == false {
				return io.EOF
			}
			log.Printf("read error: %#v\n", readErr)
			continue

		} else if h != nil {
			x := make([]byte, n)
			copy(x, b[:n])
			go func(remote *net.UDPAddr, data []byte) {
				s := Socket{remote: remote, conn: conn}
				if err := h(data, s); err != nil {
					log.Println("handling error:", err)
				}
			}(addr, x)
		}
	}
}

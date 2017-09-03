package socket

import (
	"fmt"
	"io"
	"net"

	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("rendez-vous")

// UDP is a struct
type UDP struct {
	conn net.PacketConn
}

//Close the socket
func (u *UDP) Close() error {
	return u.conn.Close()
}

//Conn of the underlying
func (u *UDP) Conn() net.PacketConn {
	return u.conn
}

// Handler handle incoming messages
type Handler func(data []byte, remote net.Addr) error

//Conn of the underlying
func (u *UDP) Write(data []byte, remote net.Addr) (int, error) {
	return u.conn.WriteTo(data, remote)
}

// Listen invoke process when a new message income
func (u *UDP) Listen(h Handler) error {

	conn := u.conn
	var b [0x10000]byte
	for {
		n, addr, readErr := conn.ReadFrom(b[:])
		if readErr == nil && n == len(b) {
			readErr = fmt.Errorf("received packet exceeds buffer size %q", len(b))
		}

		if readErr != nil {
			if x, ok := readErr.(*net.OpError); ok && x.Temporary() == false {
				return io.EOF
			}
			logger.Errorf("read error: %#v\n", readErr)
			continue

		} else if h != nil {
			x := make([]byte, n)
			copy(x, b[:n])
			go func(remote net.Addr, data []byte) {
				if err := h(data, remote); err != nil {
					logger.Error("handling error:", err)
				}
			}(addr, x)
		}
	}
}

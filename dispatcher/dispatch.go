package dispatcher

import (
	"bytes"
	"fmt"
	"net"

	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("rendez-vous")

// Dispatch multiplex many packetconn over one packetcon.
// dispatched conn emits packets with a protocol prefix.
type Dispatch struct {
	src      net.PacketConn
	dispatch map[string]*dispatchedPacketConn
	close    chan bool
}

// New dispatcher onto given conn
func New(src net.PacketConn) *Dispatch {
	d := &Dispatch{
		src,
		map[string]*dispatchedPacketConn{},
		make(chan bool),
	}
	go d.loop()
	return d
}

// UDP builds a dispatcher from given addr
func UDP(address string) (*Dispatch, error) {
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
	return New(conn), nil
}

// Proto creates a dispatched conn with given protocol id
func (d *Dispatch) Proto(p int) *dispatchedPacketConn {
	return d.New(fmt.Sprintf("%v", p))
}

// New creates a dispatched conn with given protocol name
func (d *Dispatch) New(name string) *dispatchedPacketConn {
	if x, ok := d.dispatch[name]; ok {
		return x
	}
	d.dispatch[name] = &dispatchedPacketConn{
		src:       d.src,
		name:      name,
		closed:    false,
		readWaits: make(chan *pendingRead),
	}
	return d.dispatch[name]
}

func (d *Dispatch) loop() {
	in := make([]byte, 2048)
	for {
		n, addr, err := d.src.ReadFrom(in)
		select {
		case <-d.close:
			return
		default:
			//keep going
		}
		if n > 0 {
			for name, dispatched := range d.dispatch {
				if dispatched.IsClosed() {
					delete(d.dispatch, name)
					continue
				}
				u := []byte(fmt.Sprintf("%v:%v", len(name), name))
				if n >= len(u) && bytes.Equal(u, in[:len(u)]) {
					out := in[len(u):]
					select {
					case p := <-dispatched.readWaits:
						p.err = err
						p.n = n - len(u)
						p.addr = addr
						m := len(p.b)
						if m > len(out) {
							m = len(out)
						}
						copy(p.b, out[:m])
						p.done <- true
					default:
						// do not block.
					}
				}
			}
		}
	}
}

type pendingRead struct {
	b    []byte
	n    int
	addr net.Addr
	err  error
	done chan bool
}

// // Accept ...
// func (d *dispatchedPacketConn) Accept() (net.Conn, error) {
// 	d.ReadFrom(b)
// }

// type dispatchedConn struct {
// 	*dispatchedPacketConn
// 	t net.Adr
// }
//
// // Read ...
// func (d *dispatchedPacketConn) Read(b []byte) (n int, err error) {
// 	n,_,err := d.ReadFrom(b)
// 	return n, err
// }
//
// // Write ...
// func (d *dispatchedPacketConn) Write(b []byte) (n int, err error) {
// 	return d.dispathed.WriteTo(b, d.t)
// }

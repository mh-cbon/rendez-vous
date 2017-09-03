package dispatcher

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type dispatchedPacketConn struct {
	name      string
	src       net.PacketConn
	closed    bool
	l         sync.Mutex
	readWaits chan *pendingRead
}

// Pending read requests to honor
func (d *dispatchedPacketConn) Pending() chan *pendingRead {
	return d.readWaits
}

// ReadFrom ...
func (d *dispatchedPacketConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	x := &pendingRead{b: b, done: make(chan bool)}
	d.readWaits <- x
	<-x.done
	return x.n, x.addr, x.err
}

// WriteTo ...
func (d *dispatchedPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	u := []byte(fmt.Sprintf("%v:%v", len(d.name), d.name))
	b = append(u, b...)
	n, err = d.src.WriteTo(b, addr)
	if n > len(u) {
		n -= len(u)
	}
	return n, err
}

func (d *dispatchedPacketConn) IsClosed() bool {
	d.l.Lock()
	defer d.l.Unlock()
	return d.closed
}

// Close ...
func (d *dispatchedPacketConn) Close() error {
	d.l.Lock()
	defer d.l.Unlock()
	d.closed = true
	return nil
}

// LocalAddr ...
func (d *dispatchedPacketConn) LocalAddr() net.Addr {
	return d.src.LocalAddr()
}

// SetDeadline ...
func (d *dispatchedPacketConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline ...
func (d *dispatchedPacketConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline ...
func (d *dispatchedPacketConn) SetWriteDeadline(t time.Time) error {
	return nil
}

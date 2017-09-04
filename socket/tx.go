package socket

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

// Tx ...
type Tx struct {
	UDP
	transactions map[uint16]pendingQuery
	id           uint16
	l            sync.Mutex
	closed       chan bool
}

//Close the socket
func (t *Tx) init() {
	t.l.Lock()
	t.transactions = map[uint16]pendingQuery{}
	t.id = 0
	t.l.Unlock()
}

//Close the socket
func (t *Tx) Close() error {
	if t.closed != nil {
		t.closed <- true
	}
	t.init()
	return t.UDP.Close()
}

//Conn of the underlying
func (t *Tx) Conn() net.PacketConn {
	return t.UDP.Conn()
}

type pendingQuery struct {
	h      ResponseHandler
	create time.Time
}

func (t *Tx) loop() {
	for {
		select {
		case <-t.closed:
			return
		case <-time.After(time.Millisecond * 100):
			t.l.Lock()
			for index, p := range t.transactions {
				if p.create.Add(time.Second * 5).Before(time.Now()) {
					t.transactions[index].h(nil, true)
					delete(t.transactions, index)
				}
			}
			t.l.Unlock()
		}
	}
}

func (t *Tx) makeID() uint16 {
	t.id++
	//todo: find a better way, 10k is maybe not that much.
	if t.id > 10000 {
		t.id = 0
	}
	return t.id
}

// Query a remote
func (t *Tx) Query(data []byte, remote net.Addr, h ResponseHandler) error {
	t.l.Lock()
	txID := t.makeID()
	t.transactions[txID] = pendingQuery{h, time.Now()}
	t.l.Unlock()
	b := make([]byte, binary.MaxVarintLen16)
	binary.LittleEndian.PutUint16(b, txID)
	data = append(b, data...)
	data = append([]byte("q"), data...)
	_, err := t.UDP.Write(data, remote)
	//todo: handle _
	return err
}

// Reply to a remote
func (t *Tx) Reply(data []byte, remote net.Addr, txID uint16) error {
	b := make([]byte, binary.MaxVarintLen16)
	binary.LittleEndian.PutUint16(b, txID)
	data = append(b, data...)
	data = append([]byte("r"), data...)
	_, err := t.UDP.Write(data, remote)
	//todo: handle _
	return err
}

// TxHandler handle tx messages
type TxHandler func(remote net.Addr, data []byte, reply ResponseWriter) error

// ResponseWriter writes response
type ResponseWriter func(data []byte) error

// ResponseHandler respond to a query
type ResponseHandler func(data []byte, timedout bool) error

// Listen ...
func (t *Tx) Listen(queryHandler TxHandler) error {
	t.closed = make(chan bool)
	go t.loop()
	return t.UDP.Listen(func(data []byte, remote net.Addr) error {
		if len(data) < 1 {
			return fmt.Errorf("data too small")
		}
		kind := string(data[0])
		data = data[1:]
		txID := binary.LittleEndian.Uint16(data[:binary.MaxVarintLen16])
		data = data[binary.MaxVarintLen16:]
		if kind == "q" {
			if queryHandler == nil {
				return nil
			}
			return queryHandler(remote, data, func(data []byte) error {
				return t.Reply(data, remote, txID)
			})

		} else if kind == "r" {
			t.l.Lock()
			if handler, ok := t.transactions[txID]; ok {
				delete(t.transactions, txID)
				t.l.Unlock()
				return handler.h(data, false)
			}
			t.l.Unlock()
			return fmt.Errorf("transaction id not found: %v", txID)
		}
		return fmt.Errorf("wrong message remote:%v kind:%v data.len:%v", remote.String(), kind, len(data))
	})
}

package store

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"
)

// NewRegistrations thread safe store
func NewRegistrations(s *Registrations) *TSRegistrations {
	if s == nil {
		s = &Registrations{}
	}
	return &TSRegistrations{store: s, m: &sync.Mutex{}}
}

// Registrations of peers
type Registrations struct {
	Peers []*Peer
}

// Peer is an address and a pbk
type Peer struct {
	Create     time.Time
	Address    net.Addr
	Pbk        []byte
	Sign       []byte
	Value      string
	PortStatus int
}

// Add a peer (remote+pbk)
func (s *Registrations) Add(address net.Addr, pbk []byte, sign []byte, value string) {
	p := &Peer{time.Now(), address, make([]byte, len(pbk)), make([]byte, len(sign)), value, 0}
	copy(p.Pbk, pbk)
	copy(p.Sign, sign)
	s.Peers = append(s.Peers, p)
}

// AddUpdate a peer (remote+pbk)
func (s *Registrations) AddUpdate(address net.Addr, pbk []byte, sign []byte, value string) error {
	p := s.GetByPbk(pbk)
	if p != nil {
		if p.Address.String() != address.String() {
			return errors.New("Mismatching address")
		}
		p.Create = time.Now()
	} else {
		s.Add(address, pbk, sign, value)
	}
	return nil
}

// GetByAddr find a peer by its addresss
func (s *Registrations) GetByAddr(address string) *Peer {
	for _, p := range s.Peers {
		if p.Address.String() == address {
			return p
		}
	}
	return nil
}

// GetByPbk find a peer by its Pbk
func (s *Registrations) GetByPbk(pbk []byte) *Peer {
	for _, p := range s.Peers {
		if bytes.Equal(p.Pbk, pbk) {
			return p
		}
	}
	return nil
}

// RemoveByAddr removes all peer with given address
func (s *Registrations) RemoveByAddr(address string) bool {
	ok := false
	for i, p := range s.Peers {
		if p.Address.String() == address {
			s.Peers = append(s.Peers[:i], s.Peers[i+1:]...)
			ok = true
		}
	}
	return ok
}

// RemoveByPbk removes all peer with given pbk
func (s *Registrations) RemoveByPbk(pbk []byte) bool {
	ok := false
	for i, p := range s.Peers {
		if bytes.Equal(p.Pbk, pbk) {
			s.Peers = append(s.Peers[:i], s.Peers[i+1:]...)
			ok = true
		}
	}
	return ok
}

// HasAddr return true if a peer has given address
func (s *Registrations) HasAddr(address string) bool {
	for _, p := range s.Peers {
		if p.Address.String() == address {
			return true
		}
	}
	return false
}

// HasPbk return true if a peer has given pbk
func (s *Registrations) HasPbk(pbk []byte) bool {
	for _, p := range s.Peers {
		if bytes.Equal(p.Pbk, pbk) {
			return true
		}
	}
	return false
}

// Select some peers
func (s *Registrations) Select(start, limit int) []*Peer {
	var ret []*Peer
	if start < len(s.Peers) {
		max := start + limit
		if max > len(s.Peers) {
			max = len(s.Peers)
		}
		ret = s.Peers[start:max]
	}
	return ret
}

// SetPortStatus of a peer
func (s *Registrations) SetPortStatus(address string, status int) bool {
	for _, p := range s.Peers {
		if p.Address.String() == address {
			p.PortStatus = status
			return true
		}
	}
	return false
}

// TSRegistrations is a TS Registrations
type TSRegistrations struct {
	store *Registrations
	m     *sync.Mutex
}

// Transact ....
func (s *TSRegistrations) Transact(t func(s *Registrations)) {
	s.m.Lock()
	defer s.m.Unlock()
	t(s.store)
}

// Add a peer (remote+pbk)
func (s *TSRegistrations) Add(address net.Addr, pbk, sign []byte, value string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.store.Add(address, pbk, sign, value)
}

// AddUpdate a peer (remote+pbk)
func (s *TSRegistrations) AddUpdate(address net.Addr, pbk, sign []byte, value string) error {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.AddUpdate(address, pbk, sign, value)
}

// GetByAddr find a peer by its addresss
func (s *TSRegistrations) GetByAddr(address string) *Peer {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.GetByAddr(address)
}

// GetByPbk find a peer by its Pbk
func (s *TSRegistrations) GetByPbk(pbk []byte) *Peer {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.GetByPbk(pbk)
}

// RemoveByAddr removes all peer with given address
func (s *TSRegistrations) RemoveByAddr(address string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.RemoveByAddr(address)
}

// RemoveByPbk removes all peer with given pbk
func (s *TSRegistrations) RemoveByPbk(pbk []byte) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.RemoveByPbk(pbk)
}

// HasAddr return true if a peer has given address
func (s *TSRegistrations) HasAddr(address string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.HasAddr(address)
}

// HasPbk return true if a peer has given pbk
func (s *TSRegistrations) HasPbk(pbk []byte) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.HasPbk(pbk)
}

// Select some peers
func (s *TSRegistrations) Select(start, limit int) []*Peer {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.Select(start, limit)
}

// SetPortStatus of a peer
func (s *TSRegistrations) SetPortStatus(address string, status int) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.SetPortStatus(address, status)
}

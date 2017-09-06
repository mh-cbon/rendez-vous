package store

import (
	"bytes"
	"net"
	"sync"
	"time"
)

// Registrations of peers
type Registrations struct {
	Peers []Peer
}

// Peer is an address and a pbk
type Peer struct {
	Address net.Addr
	Pbk     []byte
	create  time.Time
}

// Add a peer (remote+pbk)
func (s *Registrations) Add(address net.Addr, pbk []byte) {
	p := Peer{address, make([]byte, len(pbk)), time.Now()}
	copy(p.Pbk, pbk)
	s.Peers = append(s.Peers, p)
}

// GetByAddr find a peer by its addresss
func (s *Registrations) GetByAddr(address string) *Peer {
	for _, p := range s.Peers {
		if p.Address.String() == address {
			return &p
		}
	}
	return nil
}

// GetByPbk find a peer by its Pbk
func (s *Registrations) GetByPbk(pbk []byte) *Peer {
	for _, p := range s.Peers {
		if bytes.Equal(p.Pbk, pbk) {
			return &p
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

// TSRegistrations is a TS Registrations
type TSRegistrations struct {
	store *Registrations
	m     *sync.Mutex
}

// New thread safe store
func New(s *Registrations) *TSRegistrations {
	if s == nil {
		s = &Registrations{}
	}
	return &TSRegistrations{store: s, m: &sync.Mutex{}}
}

// Add a peer (remote+pbk)
func (s *TSRegistrations) Add(address net.Addr, pbk []byte) {
	s.m.Lock()
	defer s.m.Unlock()
	s.store.Add(address, pbk)
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

package store

import (
	"bytes"
	"log"
	"sync"
)

// Store of peers
type Store struct {
	Peers []Peer
}

// Peer is an address and a pbk
type Peer struct {
	Address string
	Pbk     []byte
}

// Add a peer (remote+pbk)
func (s *Store) Add(address string, pbk []byte) {
	p := Peer{Address: address, Pbk: make([]byte, len(pbk))}
	copy(p.Pbk, pbk)
	s.Peers = append(s.Peers, p)
}

// GetByAddr find a peer by its addresss
func (s *Store) GetByAddr(address string) *Peer {
	for _, p := range s.Peers {
		if p.Address == address {
			return &p
		}
	}
	return nil
}

// GetByPbk find a peer by its Pbk
func (s *Store) GetByPbk(pbk []byte) *Peer {
	for _, p := range s.Peers {
		log.Printf("%x <> %x\n", p.Pbk, pbk)
		if bytes.Equal(p.Pbk, pbk) {
			return &p
		}
	}
	return nil
}

// RemoveByAddr removes all peer with given address
func (s *Store) RemoveByAddr(address string) bool {
	ok := false
	for i, p := range s.Peers {
		if p.Address == address {
			s.Peers = append(s.Peers[:i], s.Peers[i+1:]...)
			ok = true
		}
	}
	return ok
}

// RemoveByPbk removes all peer with given pbk
func (s *Store) RemoveByPbk(pbk []byte) bool {
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
func (s *Store) HasAddr(address string) bool {
	for _, p := range s.Peers {
		if p.Address == address {
			return true
		}
	}
	return false
}

// HasPbk return true if a peer has given pbk
func (s *Store) HasPbk(pbk []byte) bool {
	for _, p := range s.Peers {
		if bytes.Equal(p.Pbk, pbk) {
			return true
		}
	}
	return false
}

// TSStore is a TS Store
type TSStore struct {
	store *Store
	m     *sync.Mutex
}

// New thread safe store
func New(s *Store) *TSStore {
	if s == nil {
		s = &Store{}
	}
	return &TSStore{store: s, m: &sync.Mutex{}}
}

// Add a peer (remote+pbk)
func (s *TSStore) Add(address string, pbk []byte) {
	s.m.Lock()
	defer s.m.Unlock()
	s.store.Add(address, pbk)
}

// GetByAddr find a peer by its addresss
func (s *TSStore) GetByAddr(address string) *Peer {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.GetByAddr(address)
}

// GetByPbk find a peer by its Pbk
func (s *TSStore) GetByPbk(pbk []byte) *Peer {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.GetByPbk(pbk)
}

// RemoveByAddr removes all peer with given address
func (s *TSStore) RemoveByAddr(address string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.RemoveByAddr(address)
}

// RemoveByPbk removes all peer with given pbk
func (s *TSStore) RemoveByPbk(pbk []byte) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.RemoveByPbk(pbk)
}

// HasAddr return true if a peer has given address
func (s *TSStore) HasAddr(address string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.HasAddr(address)
}

// HasPbk return true if a peer has given pbk
func (s *TSStore) HasPbk(pbk []byte) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store.HasPbk(pbk)
}

package client

import (
	"fmt"
	"net"
	"sync"
)

// TokenStore of tokens received after sending "get_peers"/"get" query.
type TokenStore struct {
	tokens map[string]string
} //todos: add a limit on maximum number of tokens.

// NewTokenStore creates a new TokenStore
func NewTokenStore() *TokenStore {
	return &TokenStore{tokens: map[string]string{}}
}

//SetToken for given remote
func (s *TokenStore) SetToken(token string, remote *net.UDPAddr) error {
	id := remote.String()
	s.tokens[id] = token
	return nil
}

//RmByAddr for given remote.
func (s *TokenStore) RmByAddr(remote *net.UDPAddr) error {
	id := remote.String()
	if _, ok := s.tokens[id]; ok {
		delete(s.tokens, id)
		return nil
	}
	return fmt.Errorf("remote address not found: %v", id)
}

//RmByToken for given token
func (s *TokenStore) RmByToken(token string) error {
	for id, t := range s.tokens {
		if t == token {
			delete(s.tokens, id)
		}
	}
	return nil
}

//GetToken for a remote.
func (s *TokenStore) GetToken(remote net.UDPAddr) string {
	id := remote.String()
	if token, ok := s.tokens[id]; ok {
		return token
	}
	return ""
}

//Clear the storage.
func (s *TokenStore) Clear() {
	s.tokens = map[string]string{}
}

// TSTokenStore tokens for a given address
type TSTokenStore struct {
	store *TokenStore
	mu    *sync.RWMutex
}

// NewTSTokenStore creates a new TS store
func NewTSTokenStore() *TSTokenStore {
	return &TSTokenStore{
		store: NewTokenStore(),
		mu:    &sync.RWMutex{},
	}
}

//SetToken for given remote
func (s *TSTokenStore) SetToken(token string, remote *net.UDPAddr) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.SetToken(token, remote)
}

//RmByAddr for given remote.
func (s *TSTokenStore) RmByAddr(remote *net.UDPAddr) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.RmByAddr(remote)
}

//RmByToken for given token
func (s *TSTokenStore) RmByToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.RmByToken(token)
}

//GetToken for a remote.
func (s *TSTokenStore) GetToken(remote net.UDPAddr) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.GetToken(remote)
}

//Clear the storage.
func (s *TSTokenStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store.Clear()
}

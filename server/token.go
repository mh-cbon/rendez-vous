package server

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"net"
	"time"

	"github.com/bradfitz/iter"
)

// NewDefaultTokenServer returns a pre configured token server.
func NewDefaultTokenServer(secret []byte) *TokenServer {
	return NewTokenServer(secret, 5*time.Minute, 2, time.Now)
}

// Server manages creation and validation of tokens issued to querying nodes.
type TokenServer struct {
	secret           []byte
	interval         time.Duration
	maxIntervalDelta int
	timeNow          func() time.Time
}

// NewServer returns a token server.
func NewTokenServer(secret []byte, interval time.Duration, maxIntervalDelta int, timeNow func() time.Time) *TokenServer {
	if timeNow == nil {
		timeNow = time.Now
	}
	ret := &TokenServer{
		interval:         interval,
		maxIntervalDelta: maxIntervalDelta,
		timeNow:          timeNow,
	}
	ret.SetSecret(secret)
	return ret
}

// SetSecret to create tokens.
func (s *TokenServer) SetSecret(secret []byte) {
	if secret == nil {
		secret = make([]byte, 20)
		rand.Read(secret)
	}
	if len(secret) > 20 {
		secret = secret[:20]
	}
	s.secret = make([]byte, len(secret))
	copy(s.secret, secret)
}

// CreateToken for given addr.
func (s TokenServer) CreateToken(addr *net.UDPAddr) string {
	return s.createToken(addr, s.timeNow())
}

func (s TokenServer) createToken(addr *net.UDPAddr, t time.Time) string {
	h := sha1.New()
	ip := addr.IP.To16()
	if len(ip) != 16 {
		panic(ip)
	}
	h.Write(ip)
	ti := t.UnixNano() / int64(s.interval)
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(ti))
	h.Write(b[:])
	h.Write(s.secret)
	return string(h.Sum(nil))
}

// ValidToken for given address.
func (s *TokenServer) ValidToken(token string, addr *net.UDPAddr) bool {
	t := s.timeNow()
	for range iter.N(s.maxIntervalDelta + 1) {
		if s.createToken(addr, t) == token {
			return true
		}
		t = t.Add(-s.interval)
	}
	return false
}

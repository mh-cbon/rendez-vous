package utils

import (
	"net"
	"net/http"

	"github.com/anacrolix/utp"
)

// UDP listen for given address
func UDP(address string) (*net.UDPConn, error) {
	protocol := "udp"

	udpAddr, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//UTPDialer http transport on utp
func UTPDialer(dialAddr string) *http.Transport {
	return &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			s, err := utp.NewSocket("udp", ":8082")
			if err != nil {
				return nil, err
			}
			defer s.Close()
			return s.DialTimeout(dialAddr, 0)
			// return utp.Dial(dialAddr)
		},
	}
}

//ServeHTTPFromListener returns an http server serving on the given net.Listener
func ServeHTTPFromListener(l net.Listener, s *http.Server) *HTTPu {
	return &HTTPu{Server: s, ln: l}
}

// HTTPu embeds http server and serves to given conn
type HTTPu struct {
	*http.Server
	ln net.Listener
}

// ListenAndServe the http requests
func (h HTTPu) ListenAndServe() error {
	return h.Server.Serve(h.ln)
}

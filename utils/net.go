package utils

import (
	"net"
	"net/http"
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

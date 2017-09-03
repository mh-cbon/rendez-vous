package utils

import "net"

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

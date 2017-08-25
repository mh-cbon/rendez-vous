package client

import (
	"encoding/json"
	"net"

	"github.com/pkg/errors"

	"github.com/mh-cbon/rendez-vous/model"
)

// FromAddr is a ctor
func FromAddr(address string) (*Client, error) {
	protocol := "udp"

	udpAddr, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// FromConn is a ctor
func FromConn(conn *net.UDPConn) *Client { return &Client{conn: conn} }

// Client to speak with a rendez-vous server
type Client struct {
	conn *net.UDPConn
}

func (c *Client) query(remote *net.UDPAddr, q model.Message) (model.Message, error) {
	var ret model.Message
	data, err := json.Marshal(q)
	if err != nil {
		return ret, errors.WithMessage(err, "query marshal")
	}
	if _, err2 := c.conn.WriteToUDP(data, remote); err2 != nil {
		return ret, errors.WithMessage(err, "query write")
	}
	res := make([]byte, 1000)
	n, err := c.conn.Read(res)
	if err != nil {
		return ret, errors.WithMessage(err, "response read")
	}
	res = res[:n]
	if err := json.Unmarshal(res, &ret); err != nil {
		return ret, errors.WithMessage(err, "response unmarshal")
	}
	return ret, nil
}

// Conn returns the underlying udp conn
func (c *Client) Conn() *net.UDPConn {
	return c.conn
}

// Ping remote
func (c *Client) Ping(remote *net.UDPAddr) (model.Message, error) {
	m := model.Message{
		Type:  "q",
		Query: model.Ping,
	}
	return c.query(remote, m)
}

// Find peer for given pbk
func (c *Client) Find(remote *net.UDPAddr, pbk []byte) (model.Message, error) {
	m := model.Message{
		Type:  "q",
		Query: model.Find,
		Pbk:   pbk,
	}
	return c.query(remote, m)
}

// Register yourself
func (c *Client) Register(remote *net.UDPAddr, pbk []byte, sign []byte, value string) (model.Message, error) {
	m := model.Message{
		Type:  "q",
		Query: model.Register,
		Pbk:   pbk,
		Sign:  sign,
		Value: value,
	}
	return c.query(remote, m)
}

// Unregister yourself
func (c *Client) Unregister(remote *net.UDPAddr, pbk []byte) (model.Message, error) {
	m := model.Message{
		Type:  "q",
		Query: model.Unregister,
		Pbk:   pbk,
	}
	return c.query(remote, m)
}

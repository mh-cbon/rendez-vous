// Package client implements a client to query a rendez-vous server.
// It also provides a client-server implementation for p2p communication.
package client

import (
	"encoding/hex"
	"net"

	logging "github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

var logger = logging.MustGetLogger("rendez-vous")

// FromSocket ...
func FromSocket(s socket.Socket) *Client {
	return &Client{model.JSONClient{Socket: s}}
}

// Client to speak with a rendez-vous server
type Client struct {
	s model.MessageQuerier
}

func (c *Client) query(remote string, q model.Message) (model.Message, error) {
	var ret model.Message
	addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return ret, err
	}
	w := make(chan error)
	queryErr := c.s.Query(q, addr, func(res model.Message, timedout bool) error {
		var replyErr error
		if timedout {
			replyErr = errors.New("query has timedout")
		} else {
			ret = res
		}
		w <- replyErr
		return replyErr
	})
	if queryErr == nil {
		queryErr = <-w
	}
	return ret, queryErr
}

// Ping remote
func (c *Client) Ping(remote string) (model.Message, error) {
	m := model.Message{
		Query: model.Ping,
	}
	return c.query(remote, m)
}

// Find peer for given pbk
func (c *Client) Find(remote string, id *identity.PublicIdentity) (model.Message, error) {
	bPbk, err := hex.DecodeString(id.Pbk)
	if err != nil {
		return model.Message{}, err
	}
	m := model.Message{
		Query: model.Find,
		Pbk:   bPbk,
		Value: id.Value,
	}
	return c.query(remote, m)
}

// Register yourself
func (c *Client) Register(remote string, id *identity.Identity) (model.Message, error) {
	bPbk, err := hex.DecodeString(id.Pbk)
	if err != nil {
		return model.Message{}, err
	}
	bSign, err2 := hex.DecodeString(id.Sign)
	if err2 != nil {
		return model.Message{}, err2
	}
	m := model.Message{
		Query: model.Register,
		Pbk:   bPbk,
		Sign:  bSign,
		Value: id.Value,
	}
	return c.query(remote, m)
}

// Unregister yourself
func (c *Client) Unregister(remote string, id *identity.Identity) (model.Message, error) {
	bPbk, err := hex.DecodeString(id.Pbk)
	if err != nil {
		return model.Message{}, err
	}
	bSign, err2 := hex.DecodeString(id.Sign)
	if err2 != nil {
		return model.Message{}, err2
	}
	m := model.Message{
		Query: model.Unregister,
		Pbk:   bPbk,
		Sign:  bSign,
		Value: id.Value,
	}
	return c.query(remote, m)
}

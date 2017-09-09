// Package client implements a client to query a rendez-vous server.
// It also provides a client-server implementation for p2p communication.
package client

import (
	"net"

	logging "github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

var logger = logging.MustGetLogger("rendez-vous")

// New ...
func New(encoder MessageEncoder) *Client {
	return &Client{encoder: encoder}
}

// JSON ...
func JSON(s socket.Socket) *JSONClient {
	return &JSONClient{s}
}

// Bencode ...
func Bencode(s socket.Socket) *BencodeClient {
	return &BencodeClient{s}
}

// Proto ...
func Proto(s socket.Socket) *ProtoClient {
	return &ProtoClient{s}
}

// MessageResponseHandler handles  query's response
type MessageResponseHandler func(data model.Message, timedout bool) error

// MessageEncoder handles encoded query/response
type MessageEncoder interface {
	Query(q model.Message, remote net.Addr, h MessageResponseHandler) error
}

// Client to speak with a rendez-vous server
type Client struct {
	encoder MessageEncoder
}

func (c *Client) query(remote string, q model.Message) (model.Message, error) {
	var ret model.Message
	addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return ret, err
	}
	w := make(chan error)
	queryErr := c.encoder.Query(q, addr, func(res model.Message, timedout bool) error {
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

// ReqKnock help
func (c *Client) ReqKnock(remote string, id *identity.PublicIdentity, token string) (model.Message, error) {
	m := model.Message{
		Query: model.ReqKnock,
		Pbk:   id.BPbk,
		Token: token,
	}
	return c.query(remote, m)
}

// Knock send
func (c *Client) Knock(remote string, knockToken string) (model.Message, error) {
	m := model.Message{
		Query: model.Knock,
		Token: knockToken,
	}
	return c.query(remote, m)
}

// DoKnock help
func (c *Client) DoKnock(remote string, knockAddress string, knockToken string) (model.Message, error) {
	m := model.Message{
		Query: model.DoKnock,
		Data:  knockAddress,
		Token: knockToken,
	}
	return c.query(remote, m)
}

// Find peer for given pbk
func (c *Client) Find(remote string, id *identity.PublicIdentity) (model.Message, error) {
	m := model.Message{
		Query: model.Find,
		Pbk:   id.BPbk,
		Value: id.Value,
	}
	return c.query(remote, m)
}

// Register yourself
func (c *Client) Register(remote string, id *identity.Identity) (model.Message, error) {
	m := model.Message{
		Query: model.Register,
		Pbk:   id.BPbk,
		Sign:  id.BSign,
		Value: id.Value,
	}
	return c.query(remote, m)
}

// Unregister yourself
func (c *Client) Unregister(remote string, id *identity.Identity) (model.Message, error) {
	unregister, err := id.Derive("unregister")
	if err != nil {
		return model.Message{}, err
	}
	m := model.Message{
		Query: model.Unregister,
		Pbk:   unregister.BPbk,
		Sign:  unregister.BSign,
		Value: unregister.Value,
	}
	return c.query(remote, m)
}

// List some peers
func (c *Client) List(remote string, start, limit int) ([]*model.Peer, error) {
	m := model.Message{
		Query: model.List,
		Start: int32(start),
		Limit: int32(limit),
	}
	res, err := c.query(remote, m)
	if err != nil {
		return nil, err
	}
	return res.Peers, nil
}

// TestPort sends testport query
func (c *Client) TestPort(remote string, token string) (model.Message, error) {
	m := model.Message{
		Query: model.TestPort,
		Token: token,
	}
	return c.query(remote, m)
}

// PortTest anwers TestPort
func (c *Client) PortTest(remote string, token string) (model.Message, error) {
	m := model.Message{
		Query:   model.PortTest,
		Token:   token,
		Address: remote,
	}
	return c.query(remote, m)
}

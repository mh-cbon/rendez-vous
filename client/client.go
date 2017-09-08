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
	"github.com/mh-cbon/rendez-vous/store"
)

var logger = logging.MustGetLogger("rendez-vous")

// New ...
func New(encoder MessageEncoder, knocks *store.TSPendingKnocks) *Client {
	return &Client{encoder: encoder, knocks: knocks}
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
	knocks  *store.TSPendingKnocks
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
func (c *Client) ReqKnock(remote string, id *identity.PublicIdentity) (string, error) {
	bPbk, err := hex.DecodeString(id.Pbk)
	if err != nil {
		return "", err
	}
	knock := c.knocks.Add("")
	defer c.knocks.Rm(knock)
	m := model.Message{
		Query: model.ReqKnock,
		Pbk:   bPbk,
		Data:  knock.ID,
	}
	f, err2 := c.query(remote, m)
	if err2 == nil {
		for i := 0; i < 5; i++ {
			var res string
			res, err2 = knock.Run(func() error {
				go c.Knock(f.Data, knock.ID)
				return nil
			})
			if err2 == nil {
				return res, err2
			}
		}
	}
	return "", err2
}

// Knock send
func (c *Client) Knock(remote string, knockToken string) (model.Message, error) {
	m := model.Message{
		Query: model.Knock,
		Data:  knockToken,
	}
	return c.query(remote, m)
}

// DoKnock help
func (c *Client) DoKnock(remote string, knockAddress string, knockToken string) (model.Message, error) {
	m := model.Message{
		Query: model.DoKnock,
		Data:  knockAddress,
		Value: knockToken,
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

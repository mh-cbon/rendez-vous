package client

import (
	"encoding/hex"
	"encoding/json"
	"net"

	logging "github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

var logger = logging.MustGetLogger("rendez-vous")

// FromSocket ...
func FromSocket(s socket.Socket) Client {
	return Client{s}
}

// Client to speak with a rendez-vous server
type Client struct {
	s socket.Socket
}

//Listen ...
func (c *Client) Listen() error {
	return c.s.Listen(func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		return nil
	})
}

func (c *Client) query(remote string, q model.Message) (model.Message, error) {
	var ret model.Message
	addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return ret, err
	}
	data, err := json.Marshal(q)
	if err != nil {
		return ret, errors.WithMessage(err, "query marshal")
	}
	w := make(chan error)
	queryErr := c.s.Query(data, addr, func(data []byte, timedout bool) error {
		var replyErr error
		if timedout {
			replyErr = errors.New("query has timedout")
		} else if replyErr = json.Unmarshal(data, &ret); replyErr != nil {
			replyErr = errors.WithMessage(replyErr, "response unmarshal")
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
		// Type:  "q",
		Query: model.Ping,
	}
	return c.query(remote, m)
}

// Find peer for given pbk
func (c *Client) Find(remote string, pbk string) (model.Message, error) {
	bPbk, err := hex.DecodeString(pbk)
	if err != nil {
		return model.Message{}, err
	}
	m := model.Message{
		// Type:  "q",
		Query: model.Find,
		Pbk:   bPbk,
	}
	return c.query(remote, m)
}

// Register yourself
func (c *Client) Register(remote string, pbk string, sign string, value string) (model.Message, error) {
	bPbk, err := hex.DecodeString(pbk)
	if err != nil {
		return model.Message{}, err
	}
	bSign, err2 := hex.DecodeString(sign)
	if err2 != nil {
		return model.Message{}, err2
	}
	m := model.Message{
		// Type:  "q",
		Query: model.Register,
		Pbk:   bPbk,
		Sign:  bSign,
		Value: value,
	}
	return c.query(remote, m)
}

// Unregister yourself
func (c *Client) Unregister(remote string, pbk string) (model.Message, error) {
	bPbk, err := hex.DecodeString(pbk)
	if err != nil {
		return model.Message{}, err
	}
	m := model.Message{
		// Type:  "q",
		Query: model.Unregister,
		Pbk:   bPbk,
	}
	return c.query(remote, m)
}

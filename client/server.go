package client

import (
	"encoding/json"
	"net"

	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/pkg/errors"
)

//handleQuery ...
func (c *Client) handleQuery(remote net.Addr, data []byte, writer socket.ResponseWriter) error {

	var v model.Message
	err := json.Unmarshal(data, &v)
	if err != nil {
		return errors.WithMessage(err, "json unmarshal")
	}

	logger.Info(remote.String(), "<-", v)

	var res *model.Message

	switch v.Query {

	case model.Services:
		res = model.ReplyOk(remote, "")

	}

	if res != nil {
		b, err := json.Marshal(*res)
		if err != nil {
			return errors.WithMessage(err, "json marshal")
		}
		return writer(b)
	}
	return nil
}

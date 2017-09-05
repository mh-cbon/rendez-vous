package client

import (
	"net"

	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

//HandleQuery handles p2p communication.
func HandleQuery() socket.TxHandler {
	return model.ProtoHandler(func(remote net.Addr, v model.Message, writer model.MessageResponseWriter) error {

		var res *model.Message

		switch v.Query {

		case model.Ping:
			res = model.ReplyOk(remote, "")

		case model.Services:
			res = model.ReplyOk(remote, "")

		}

		return writer(*res)
	})
}

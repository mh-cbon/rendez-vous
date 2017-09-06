package client

import (
	"log"
	"net"

	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

//HandleQuery handles p2p communication.
func HandleQuery(c *Client) socket.TxHandler {
	return model.JSONHandler(func(remote net.Addr, v model.Message, writer model.MessageResponseWriter) error {
		var res *model.Message

		switch v.Query {

		case model.Ping:
			res = model.ReplyOk(remote, "")

		case model.Services:
			res = model.ReplyOk(remote, "")

		case model.DoKnock:
			addrToKnock := v.Data
			knockToken := v.Value

			knock := c.knocks.Add(addrToKnock, knockToken)
			okAddr, err3 := knock.Run(c)
			log.Println(okAddr, err3)
			c.knocks.Rm(knock)

		case model.Knock:
			c.knocks.Resolve(remote.String(), v.Data)
		}

		if res != nil {
			return writer(remote, *res)
		}
		return nil
	})
}

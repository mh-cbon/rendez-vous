package client

import (
	"fmt"
	"net"

	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
)

//HandleQuery handles p2p communication.
func HandleQuery(client *Client) socket.TxHandler {
	return model.JSONHandler(func(remote net.Addr, v model.Message, writer model.MessageResponseWriter) error {

		var res *model.Message

		switch v.Query {

		case model.Ping:
			res = model.ReplyOk(remote, "")

		case model.Services:
			res = model.ReplyOk(remote, "")

		case model.DoKnock:
			_, err := client.Ping(res.Data)
			if err != nil {
				for i := 0; i < 5; i++ {
					_, err = client.Ping(res.Data)
					if err == nil {
						break
					}
				}
			}
			if err != nil {
				return fmt.Errorf("knock failure: %v", err.Error())
			}
		}

		return writer(remote, *res)
	})
}

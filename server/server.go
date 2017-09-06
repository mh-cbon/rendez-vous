// Package server runs a rendez-vous meeting point server.
// Its a server onto which clients can announce/find services.
package server

import (
	"net"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
	logging "github.com/op/go-logging"
)

var (
	missingPbk   = 301
	missingSign  = 302
	invalidValue = 303
	invalidSign  = 304
	wrongQuery   = 305
	wrongPbk     = 306

	notFound = 404
)

var logger = logging.MustGetLogger("rendez-vous")

//HandleQuery ...
func HandleQuery(c *client.Client, registrations *store.TSRegistrations) socket.TxHandler {
	if registrations == nil {
		registrations = store.New(nil)
	}
	return model.JSONHandler(func(remote net.Addr, v model.Message, writer model.MessageResponseWriter) error {
		var res *model.Message

		switch v.Query {

		case model.Ping:
			res = model.ReplyOk(remote, "")

		case model.Knock:
			if len(v.Pbk) == 0 {
				res = model.ReplyError(remote, missingPbk)

			} else if len(v.Pbk) != 32 {
				res = model.ReplyError(remote, wrongPbk)

			} else {
				peer := registrations.GetByPbk(v.Pbk)
				if peer != nil {
					go c.DoKnock(peer.Address.String(), remote.String())
					res = model.ReplyOk(remote, peer.Address.String())
				}
			}

		case model.Register:
			//todo: rendez-vous server should implement a write token

			if len(v.Pbk) == 0 {
				res = model.ReplyError(remote, missingPbk)

			} else if len(v.Pbk) != 32 {
				res = model.ReplyError(remote, wrongPbk)

			} else if len(v.Sign) == 0 {
				res = model.ReplyError(remote, missingSign)

			} else if len(v.Value) > 100 {
				res = model.ReplyError(remote, invalidValue)

			} else if ed25519.Verify(v.Pbk, []byte(v.Value), v.Sign) == false {
				res = model.ReplyError(remote, invalidSign)

			} else {
				go func() {
					registrations.RemoveByAddr(remote.String())
					registrations.Add(remote, v.Pbk)
				}()
				res = model.ReplyOk(remote, "")
			}

		case model.Unregister:
			//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
			if len(v.Pbk) == 0 {
				res = model.ReplyError(remote, missingPbk)

			} else {
				addr := remote.String() //is it a safe value ?
				go registrations.RemoveByAddr(addr)
				res = model.ReplyOk(remote, "")
			}

		case model.Find:
			if len(v.Pbk) == 0 {
				res = model.ReplyError(remote, missingPbk)

			} else if peer := registrations.GetByPbk(v.Pbk); peer != nil {
				res = model.ReplyOk(remote, peer.Address.String())

			} else {
				res = model.ReplyError(remote, notFound)
			}

		case model.Join:
			//todo: Join the swarm
		case model.Leave:
			//todo: leave the swarm
		}

		return writer(remote, *res)
	})
}

// Package server runs a rendez-vous meeting point server.
// Its a server onto which clients can announce/find services.
package server

import (
	"net"

	"github.com/mh-cbon/dht/ed25519"
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
func HandleQuery(registrations *store.TSRegistrations) socket.TxHandler {
	if registrations == nil {
		registrations = store.New(nil)
	}
	return model.ProtoHandler(func(remote net.Addr, v model.Message, writer model.MessageResponseWriter) error {
		var res *model.Message

		switch v.Query {

		case model.Ping:
			res = model.ReplyOk(remote, "")

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
				addr := remote.String() //is it a safe value ?
				go func() {
					registrations.RemoveByAddr(addr)
					registrations.Add(addr, v.Pbk)
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
				res = model.ReplyOk(remote, peer.Address)

			} else {
				res = model.ReplyError(remote, notFound)
			}

		case model.Join:
			//todo: Join the swarm
		case model.Leave:
			//todo: leave the swarm
		}

		return writer(*res)
	})
}

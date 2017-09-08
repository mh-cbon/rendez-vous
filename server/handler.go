// Package server runs a rendez-vous meeting point server.
// Its a server onto which clients can announce/find services.
package server

import (
	"errors"
	"log"
	"net"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/model"
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
	missingData  = 307
	wrongData    = 308

	notFound = 404
)

var logger = logging.MustGetLogger("rendez-vous")

// MessageQueryHandler handles json requests
type MessageQueryHandler func(remote net.Addr, m model.Message, reply MessageResponseWriter) error

// MessageResponseWriter writes response
type MessageResponseWriter func(remote net.Addr, m *model.Message) error

func OneOf(many ...MessageQueryHandler) MessageQueryHandler {
	return func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		for _, one := range many {
			err := one(remote, m, reply)
			if err == NotHandled {
				continue
			}
			return err
		}
		return errors.New("not handlded " + m.Query)
	}
}

var NotHandled = errors.New("not handlded")

func QueryHandler(query string, handler MessageQueryHandler) MessageQueryHandler {
	return func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		if m.Query == query {
			return handler(remote, m, func(remote net.Addr, m *model.Message) error {
				if m == nil {
					return nil
				}
				return reply(remote, m)
			})
		}
		return NotHandled
	}
}

func HandlePing() MessageQueryHandler {
	return QueryHandler(model.Ping, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		return reply(remote, model.ReplyOk(remote, ""))
	})
}

func HandleRegister(registrations *store.TSRegistrations) MessageQueryHandler {
	return QueryHandler(model.Register, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		//todo: rendez-vous server should implement a write token
		if len(m.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else if len(m.Pbk) != 32 {
			res = model.ReplyError(remote, wrongPbk)

		} else if len(m.Sign) == 0 {
			res = model.ReplyError(remote, missingSign)

		} else if len(m.Value) > 100 {
			res = model.ReplyError(remote, invalidValue)

		} else if ed25519.Verify(m.Pbk, []byte(m.Value), m.Sign) == false {
			res = model.ReplyError(remote, invalidSign)

		} else {
			go func() {
				registrations.RemoveByAddr(remote.String())
				registrations.RemoveByPbk(m.Pbk)
				registrations.Add(remote, m.Pbk)
			}()
			res = model.ReplyOk(remote, "")
		}
		return reply(remote, res)
	})
}

func HandleUnregister(registrations *store.TSRegistrations) MessageQueryHandler {
	return QueryHandler(model.Unregister, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
		if len(m.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else {
			go registrations.RemoveByAddr(remote.String())
			res = model.ReplyOk(remote, "")
		}
		return reply(remote, res)
	})
}

func HandleFind(registrations *store.TSRegistrations) MessageQueryHandler {
	return QueryHandler(model.Find, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		if len(m.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else if peer := registrations.GetByPbk(m.Pbk); peer != nil {
			res = model.ReplyOk(remote, peer.Address.String())

		} else {
			res = model.ReplyError(remote, notFound)
		}
		return reply(remote, res)
	})
}

func HandleRequestKnock(registrations *store.TSRegistrations, c *client.Client) MessageQueryHandler {
	return QueryHandler(model.ReqKnock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		if len(m.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else if len(m.Pbk) != 32 {
			res = model.ReplyError(remote, wrongPbk)

		} else if len(m.Data) == 0 {
			res = model.ReplyError(remote, missingData)

		} else if len(m.Data) > 100 {
			res = model.ReplyError(remote, wrongData)

		} else {
			peer := registrations.GetByPbk(m.Pbk)
			if peer != nil {
				go c.DoKnock(peer.Address.String(), remote.String(), m.Data)
				res = model.ReplyOk(remote, peer.Address.String())
			} else {
				res = model.ReplyError(remote, notFound)
			}
		}
		return reply(remote, res)
	})
}

func HandleDoKnock(c *client.Client, knocks *store.TSPendingKnocks) MessageQueryHandler {
	return QueryHandler(model.DoKnock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		//todo: protect from undesired usage.
		addrToKnock := m.Data
		knockToken := m.Value
		knock := knocks.Add(knockToken)
		go func() {
			defer knocks.Rm(knock)
			for i := 0; i < 5; i++ {
				_, err := knock.Run(func() error {
					_, err := c.Knock(addrToKnock, knock.ID)
					return err
				})
				if err == nil {
					break
				}
			}
		}()
		return nil
	})
}

func HandleKnock(knocks *store.TSPendingKnocks) MessageQueryHandler {
	return QueryHandler(model.Knock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		log.Println("knock q: ", remote.String(), m.Data)
		if knocks.Resolve(remote.String(), m.Data) {
			res = model.ReplyOk(remote, m.Data)
			log.Println("knock success")
		} else {
			log.Println("knock fail")
		}
		if res == nil {
			return nil
		}
		return reply(remote, res)
	})
}

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
	missingPbk   int32 = 301
	missingSign  int32 = 302
	invalidValue int32 = 303
	invalidSign  int32 = 304
	wrongQuery   int32 = 305
	wrongPbk     int32 = 306
	missingData  int32 = 307
	wrongData    int32 = 308
	notFound     int32 = 404
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
			log.Printf("registration failed %x\n", m.Pbk)
			log.Printf("registration failed %x\n", m.Sign)
			log.Printf("registration failed %v\n", m.Value)

		} else {
			go func() {
				registrations.AddUpdate(remote, m.Pbk, m.Sign, m.Value)
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

		} else if ed25519.Verify(m.Pbk, []byte(m.Value), m.Sign) == false {
			res = model.ReplyError(remote, invalidSign)

		} else {
			go registrations.RemoveByPbk(m.Pbk)
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
			res.PortStatus = int32(peer.PortStatus)
			res.Value = peer.Value
			res.Pbk = peer.Pbk
			res.Sign = peer.Sign

		} else {
			res = model.ReplyError(remote, notFound)
		}
		return reply(remote, res)
	})
}

func HandleList(registrations *store.TSRegistrations) MessageQueryHandler {
	return QueryHandler(model.List, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		found := registrations.Select(int(m.Start), int(m.Limit))
		res = model.ReplyOk(remote, "")
		res.Peers = []*model.Peer{}
		for _, f := range found {
			res.Peers = append(res.Peers, &model.Peer{f.Address.String(), int32(f.PortStatus), f.Pbk, f.Sign, f.Value})
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
				go c.DoKnock(peer.Address.String(), remote.String(), m.Token)
				res = model.ReplyOk(remote, peer.Address.String())
			} else {
				res = model.ReplyError(remote, notFound)
			}
		}
		return reply(remote, res)
	})
}

func HandleDoKnock(c *client.Client, knocks *store.TSPendingOps) MessageQueryHandler {
	return QueryHandler(model.DoKnock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		//todo: protect from undesired usage.
		addrToKnock := m.Data
		knockToken := m.Token
		knock := knocks.Add(knockToken, func(remote string, m model.Message) bool {
			return m.Token == knockToken
		})
		go func() {
			defer knocks.Rm(knock)
			for i := 0; i < 5; i++ {
				_, err := knock.Run(func() error {
					_, err := c.Knock(addrToKnock, knock.Token)
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

func HandleKnock(knocks *store.TSPendingOps) MessageQueryHandler {
	return QueryHandler(model.Knock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		log.Println("knock q: ", remote.String(), m.Token)
		if knocks.Resolve(remote.String(), m.Token, m) {
			res = model.ReplyOk(remote, m.Token)
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

func HandleTestPort(registrations *store.TSRegistrations, client *client.Client) MessageQueryHandler {
	return QueryHandler(model.TestPort, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		//todo: protect from undesired usage.
		go func(remote string, token string) {
			for i := 0; i < 5; i++ {
				res, err := client.PortTest(remote, token)
				if err == nil {
					log.Println(res)
					registrations.SetPortStatus(remote, model.PortStatusOpen)
					return
				}
			}
		}(remote.String(), m.Token)
		return nil
	})
}

func HandlePortTest(portTests *store.TSPendingOps) MessageQueryHandler {
	return QueryHandler(model.PortTest, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
		var res *model.Message
		log.Println("porttest q: ", remote.String(), m.Token)
		if portTests.Resolve(remote.String(), m.Token, m) {
			res = model.ReplyOk(remote, m.Token)
			log.Println("porttest success")
		} else {
			log.Println("porttest fail")
		}
		return reply(remote, res)
	})
}

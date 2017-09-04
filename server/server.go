package server

import (
	"encoding/json"
	"net"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
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

// FromSocket ...
func FromSocket(s socket.Socket) *Server {
	return &Server{s, store.New(nil)}
}

// Server ...
type Server struct {
	s             socket.Socket
	registrations *store.TSRegistrations
}

//Close ...
func (s *Server) Close() error {
	return s.s.Close()
}

//Listen ...
func (s *Server) Listen() error {
	return s.s.Listen(s.handleQuery)
}

//HandleQuery ...
func (s *Server) handleQuery(remote net.Addr, data []byte, writer socket.ResponseWriter) error {

	var v model.Message
	err := json.Unmarshal(data, &v)
	if err != nil {
		return errors.WithMessage(err, "json unmarshal")
	}

	logger.Info(remote.String(), "<-", v)

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
				s.registrations.RemoveByAddr(addr)
				s.registrations.Add(addr, v.Pbk)
			}()
			res = model.ReplyOk(remote, "")
		}

	case model.Unregister:
		//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
		if len(v.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else {
			addr := remote.String() //is it a safe value ?
			go s.registrations.RemoveByAddr(addr)
			res = model.ReplyOk(remote, "")
		}

	case model.Find:
		if len(v.Pbk) == 0 {
			res = model.ReplyError(remote, missingPbk)

		} else if peer := s.registrations.GetByPbk(v.Pbk); peer != nil {
			res = model.ReplyOk(remote, peer.Address)

		} else {
			res = model.ReplyError(remote, notFound)
		}

	case model.Join:
		//todo: Join the swarm
	case model.Leave:
		//todo: leave the swarm
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

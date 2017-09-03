package server

import (
	"encoding/json"
	"net"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
	logging "github.com/op/go-logging"
)

var (
	okCode = 200

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
func FromSocket(s socket.Socket) Server {
	return Server{s, store.New(nil)}
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
		return err
	}

	logger.Info(remote.String(), "<-", v)

	var res *model.Message

	switch v.Query {

	case model.Ping:
		res = replyOk(remote, "")

	case model.Register:

		if len(v.Pbk) == 0 {
			res = replyError(remote, missingPbk)

		} else if len(v.Pbk) != 32 {
			res = replyError(remote, wrongPbk)

		} else if len(v.Sign) == 0 {
			res = replyError(remote, missingSign)

		} else if len(v.Value) > 100 {
			res = replyError(remote, invalidValue)

		} else if ed25519.Verify(v.Pbk, []byte(v.Value), v.Sign) == false {
			res = replyError(remote, invalidSign)

		} else {
			addr := remote.String() //is it a safe value ?
			go func() {
				s.registrations.RemoveByAddr(addr)
				s.registrations.Add(addr, v.Pbk)
			}()
			res = replyOk(remote, "")
		}

	case model.Unregister:
		if len(v.Pbk) == 0 {
			res = replyError(remote, missingPbk)

		} else {
			addr := remote.String() //is it a safe value ?
			go s.registrations.RemoveByAddr(addr)
			res = replyOk(remote, "")
		}

	case model.Find:
		if len(v.Pbk) == 0 {
			res = replyError(remote, missingPbk)

		} else if peer := s.registrations.GetByPbk(v.Pbk); peer != nil {
			res = replyOk(remote, peer.Address)

		} else {
			res = replyError(remote, notFound)
		}

	case model.Join:
		//todo: Join the swarm
	case model.Leave:
		//todo: leave the swarm
	}

	if res != nil {
		b, err := json.Marshal(*res)
		if err != nil {
			return err
		}
		return writer(b)
	}
	return nil
}

func reply(remote net.Addr) *model.Message {
	var m model.Message
	m.Address = remote.String()
	// m.Type = "r"
	return &m
}
func replyError(remote net.Addr, code int) *model.Message {
	m := reply(remote)
	m.Code = code
	return m
}
func replyOk(remote net.Addr, data string) *model.Message {
	m := reply(remote)
	m.Code = okCode
	m.Response = data
	return m
}

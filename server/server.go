package server

import (
	"encoding/json"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
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

// Handler of a rendez-vous server
func Handler(storage *store.TSStore) socket.Handler {
	return func(data []byte, remote socket.Socket) error {

		var v model.Message
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		switch v.Query {

		case model.Ping:
			return sendError(remote, okCode)

		case model.Register:
			if len(v.Pbk) == 0 {
				return sendError(remote, missingPbk)
			}
			if len(v.Pbk) != 32 {
				return sendError(remote, wrongPbk)
			}
			if len(v.Sign) == 0 {
				return sendError(remote, missingSign)
			}
			if len(v.Value) > 100 {
				return sendError(remote, invalidValue)
			}
			if ed25519.Verify(v.Pbk, []byte(v.Value), v.Sign) == false {
				return sendError(remote, invalidSign)
			}
			addr := remote.Addr() //is it a safe value ?
			storage.RemoveByAddr(addr)
			storage.Add(addr, v.Pbk)

			return sendError(remote, okCode)

		case model.Unregister:
			if len(v.Pbk) == 0 {
				return sendError(remote, missingPbk)
			}

			addr := remote.Addr() //is it a safe value ?
			storage.RemoveByAddr(addr)

			return sendError(remote, okCode)

		case model.Find:
			if len(v.Pbk) == 0 {
				return sendError(remote, missingPbk)
			}
			peer := storage.GetByPbk(v.Pbk)
			if peer == nil {
				return sendError(remote, notFound)
			}

			return sendOk(remote, peer.Address)

		default:
			return sendError(remote, wrongQuery)
		}
	}
}

func reply(remote socket.Socket, m model.Message) error {
	m.Address = remote.Addr()
	m.Type = "r"
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return remote.Write(b)
}

func sendError(remote socket.Socket, code int) error {
	return reply(remote, model.Message{Type: "r", Code: code})
}

func sendOk(remote socket.Socket, v string) error {
	return reply(remote, model.Message{Type: "r", Code: okCode, Response: v})
}

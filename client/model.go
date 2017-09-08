package client

import (
	"encoding/json"
	"net"

	bencode "github.com/anacrolix/torrent/bencode"
	"github.com/golang/protobuf/proto"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/pkg/errors"
)

// Protoc

// ProtoClient implements proto query / responses
type ProtoClient struct {
	socket.Socket
}

// Query using proto format
func (j ProtoClient) Query(q model.Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := proto.Marshal(&q)
	if err != nil {
		return errors.WithMessage(err, "query proto marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res model.Message
		if !timedout {
			if replyErr = proto.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response proto unmarshal")
			}
		}
		return h(res, timedout)
	})
}

// Reply using json format
// func (j ProtoClient) Reply(data model.Message, remote net.Addr, txID uint16) error {
// 	b, err := proto.Marshal(&data)
// 	if err != nil {
// 		return errors.WithMessage(err, "json marshal")
// 	}
// 	return j.Socket.Reply(b, remote, txID)
// }

// Bencode

// BencodeClient implements bencode query / responses
type BencodeClient struct {
	socket.Socket
}

// Query using bencode format
func (j BencodeClient) Query(q model.Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := bencode.Marshal(q)
	if err != nil {
		return errors.WithMessage(err, "query bencode marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res model.Message
		if !timedout {
			if replyErr = bencode.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response bencode unmarshal")
			}
		}
		return h(res, timedout)
	})
}

// Reply using json format
// func (j BencodeClient) Reply(data model.Message, remote net.Addr, txID uint16) error {
// 	b, err := bencode.Marshal(data)
// 	if err != nil {
// 		return errors.WithMessage(err, "json marshal")
// 	}
// 	return j.Socket.Reply(b, remote, txID)
// }

// JSON

// JSONClient implements json query / responses
type JSONClient struct {
	socket.Socket
}

// Query using json format
func (j JSONClient) Query(q model.Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := json.Marshal(q)
	if err != nil {
		return errors.WithMessage(err, "query json marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res model.Message
		if !timedout {
			if replyErr = json.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response json unmarshal")
			}
		}
		return h(res, timedout)
	})
}

// Reply using json format
// func (j JSONClient) Reply(data model.Message, remote net.Addr, txID uint16) error {
// 	b, err := json.Marshal(data)
// 	if err != nil {
// 		return errors.WithMessage(err, "json marshal")
// 	}
// 	return j.Socket.Reply(b, remote, txID)
// }

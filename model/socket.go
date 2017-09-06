package model

import (
	"encoding/json"
	"net"

	bencode "github.com/anacrolix/torrent/bencode"
	"github.com/golang/protobuf/proto"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/pkg/errors"
)

// MessageResponseWriter writes response
type MessageResponseWriter func(remote net.Addr, m Message) error

// MessageResponseHandler handles  query's response
type MessageResponseHandler func(data Message, timedout bool) error

// MessageQuerier handles bencode query/response
type MessageQuerier interface {
	Query(q Message, remote net.Addr, h MessageResponseHandler) error
}

// Additions to the socket packet.

// MessageQueryHandler handles json requests
type MessageQueryHandler func(remote net.Addr, m Message, reply MessageResponseWriter) error

// JSONHandler is a json query/response reader/writer
func JSONHandler(h MessageQueryHandler) socket.TxHandler {
	return func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v Message
		err := json.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "json unmarshal")
		}
		w := func(remote net.Addr, res Message) error {
			b, err := json.Marshal(res)
			if err != nil {
				return errors.WithMessage(err, "json marshal")
			}
			return writer(b)
		}
		if h != nil {
			return h(remote, v, w)
		}
		return nil
	}
}

// JSONClient implements json query / responses
type JSONClient struct {
	socket.Socket
}

// Query using json format
func (j JSONClient) Query(q Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := json.Marshal(q)
	if err != nil {
		return errors.WithMessage(err, "query json marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res Message
		if !timedout {
			if replyErr = json.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response json unmarshal")
			}
		}
		return h(res, timedout)
	})
}

// BencodeQueryHandler handles bencode requests
type BencodeQueryHandler func(remote net.Addr, m Message, reply MessageResponseWriter) error

// BencodeHandler is a bencode query/response reader/writer
func BencodeHandler(h MessageQueryHandler) socket.TxHandler {
	return func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v Message
		err := bencode.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "bencode unmarshal")
		}
		w := func(remote net.Addr, res Message) error {
			b, err := bencode.Marshal(res)
			if err != nil {
				return errors.WithMessage(err, "bencode marshal")
			}
			return writer(b)
		}
		if h != nil {
			return h(remote, v, w)
		}
		return nil
	}
}

// BencodeClient implements bencode query / responses
type BencodeClient struct {
	socket.Socket
}

// Query using bencode format
func (j BencodeClient) Query(q Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := bencode.Marshal(q)
	if err != nil {
		return errors.WithMessage(err, "query bencode marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res Message
		if !timedout {
			if replyErr = bencode.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response bencode unmarshal")
			}
		}
		return h(res, timedout)
	})
}

// ProtoQueryHandler handles proto requests
type ProtoQueryHandler func(remote net.Addr, m Message, reply MessageResponseWriter) error

// ProtoHandler is a proto query/response reader/writer
func ProtoHandler(h MessageQueryHandler) socket.TxHandler {
	return func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v Message
		err := proto.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "proto unmarshal")
		}
		w := func(remote net.Addr, res Message) error {
			b, err := proto.Marshal(&res)
			if err != nil {
				return errors.WithMessage(err, "proto marshal")
			}
			return writer(b)
		}
		if h != nil {
			return h(remote, v, w)
		}
		return nil
	}
}

// ProtoClient implements proto query / responses
type ProtoClient struct {
	socket.Socket
}

// Query using proto format
func (j ProtoClient) Query(q Message, remote net.Addr, h MessageResponseHandler) error {
	data, err := proto.Marshal(&q)
	if err != nil {
		return errors.WithMessage(err, "query proto marshal")
	}
	return j.Socket.Query(data, remote, func(data []byte, timedout bool) error {
		var replyErr error
		var res Message
		if !timedout {
			if replyErr = proto.Unmarshal(data, &res); replyErr != nil {
				return errors.WithMessage(replyErr, "response proto unmarshal")
			}
		}
		return h(res, timedout)
	})
}

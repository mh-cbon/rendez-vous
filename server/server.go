package server

import (
	"encoding/json"
	"net"

	bencode "github.com/anacrolix/torrent/bencode"
	"github.com/golang/protobuf/proto"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/pkg/errors"
)

// JSON server from a socket
func JSON(s socket.Socket, handler MessageQueryHandler) *JSONServer {
	return &JSONServer{
		s:       s,
		handler: handler,
	}
}

type JSONServer struct {
	s       socket.Socket
	handler MessageQueryHandler
}

// ListenAndServe the queries/responses
func (s *JSONServer) ListenAndServe() error {
	return s.s.Listen(func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v model.Message
		err := json.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "json unmarshal")
		}
		w := func(remote net.Addr, res *model.Message) error {
			b, err := json.Marshal(res)
			if err != nil {
				return errors.WithMessage(err, "json marshal")
			}
			return writer(b)
		}
		if s.handler != nil {
			return s.handler(remote, v, w)
		}
		return nil
	})
}

// Close the server and the underlying socket.
func (s *JSONServer) Close() error { return s.s.Close() }

// Bencoded server from a socket
func Bencoded(s socket.Socket, handler MessageQueryHandler) *BencodeServer {
	return &BencodeServer{
		s:       s,
		handler: handler,
	}
}

type BencodeServer struct {
	s       socket.Socket
	handler MessageQueryHandler
}

// ListenAndServe the queries/responses
func (s *BencodeServer) ListenAndServe() error {
	return s.s.Listen(func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v model.Message
		err := bencode.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "bencode unmarshal")
		}
		w := func(remote net.Addr, res *model.Message) error {
			b, err := bencode.Marshal(res)
			if err != nil {
				return errors.WithMessage(err, "bencode marshal")
			}
			return writer(b)
		}
		if s.handler != nil {
			return s.handler(remote, v, w)
		}
		return nil
	})
}

// Close the server and the underlying socket.
func (s *BencodeServer) Close() error { return s.s.Close() }

// Protobuf server from a socket
func Protobuf(s socket.Socket, handler MessageQueryHandler) *ProtoServer {
	return &ProtoServer{
		s:       s,
		handler: handler,
	}
}

type ProtoServer struct {
	s       socket.Socket
	handler MessageQueryHandler
}

// ListenAndServe the queries/responses
func (s *ProtoServer) ListenAndServe() error {
	return s.s.Listen(func(remote net.Addr, data []byte, writer socket.ResponseWriter) error {
		var v model.Message
		err := proto.Unmarshal(data, &v)
		if err != nil {
			return errors.WithMessage(err, "proto unmarshal")
		}
		w := func(remote net.Addr, res *model.Message) error {
			b, err := proto.Marshal(res)
			if err != nil {
				return errors.WithMessage(err, "proto marshal")
			}
			return writer(b)
		}
		if s.handler != nil {
			return s.handler(remote, v, w)
		}
		return nil
	})
}

// Close the server and the underlying socket.
func (s *ProtoServer) Close() error { return s.s.Close() }

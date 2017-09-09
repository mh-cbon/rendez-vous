package model

import (
	"net"
)

//go:generate protoc --go_out=. *.proto

// Message of rendez-vous servers and clients.
type Message struct {
	// Query value of the message
	Query string `json:"q,omitempty" bencode:"q,omitempty" protobuf:"bytes,1,opt,name=query"`
	// Code of the response
	Code int32 `json:"c,omitempty" bencode:"c,omitempty" protobuf:"varint,2,opt,name=code"`
	// ed25519 Public key (raw string)
	Pbk []byte `json:"p,omitempty" bencode:"p,omitempty" protobuf:"bytes,3,opt,name=pbk"`
	// Sign of Pbk(v) (raw string)
	Sign []byte `json:"s,omitempty" bencode:"s,omitempty" protobuf:"bytes,4,opt,name=sign"`
	// Value to sign
	Value string `json:"v,omitempty" bencode:"v,omitempty" protobuf:"bytes,5,opt,name=value"`
	// Address of the querier sent in a response
	Address string `json:"a,omitempty" bencode:"a,omitempty" protobuf:"bytes,6,opt,name=address"`
	// Data for query / response
	Data string `json:"d,omitempty" bencode:"d,omitempty" protobuf:"bytes,7,opt,name=data"`
	// Token for query / response
	Token string `json:"t,omitempty" bencode:"t,omitempty" protobuf:"bytes,8,opt,name=token"`
	// PortStatus 1=open 2=close
	PortStatus int32 `json:"u,omitempty" bencode:"u,omitempty" protobuf:"bytes,9,opt,name=portStatus"`
	// A list of peers
	Peers []*Peer `json:"z,omitempty" bencode:"z,omitempty" protobuf:"bytes,10,opt,name=peers"`
	// Start of the listing
	Start int32 `json:"e,omitempty" bencode:"e,omitempty" protobuf:"bytes,11,opt,name=start"`
	Limit int32 `json:"l,omitempty" bencode:"l,omitempty" protobuf:"bytes,12,opt,name=limit"`
}
type Peer struct {
	Address    string `json:"a,omitempty" bencode:"a,omitempty" protobuf:"bytes,1,opt,name=address"`
	PortStatus int32  `json:"u,omitempty" bencode:"u,omitempty" protobuf:"bytes,2,opt,name=portStatus"`
	Pbk        []byte `json:"p,omitempty" bencode:"u,omitempty" protobuf:"bytes,3,opt,name=pbk"`
	Sign       []byte `json:"s,omitempty" bencode:"s,omitempty" protobuf:"bytes,4,opt,name=sign"`
	Value      string `json:"v,omitempty" bencode:"v,omitempty" protobuf:"bytes,5,opt,name=value"`
}

// defines default response codes
var (
	OkCode = 200
)

// defines default verb
var (
	// central rendez-vous
	Ping       = "ping"
	Register   = "reg"
	Unregister = "unreg"
	Find       = "find"
	List       = "list"
	Join       = "join"
	Leave      = "leave"
	ReqKnock   = "reqknock"
	DoKnock    = "doknock"
	Knock      = "knock"
	TestPort   = "testport"
	PortTest   = "porttest"
	// leaf node
	Services = "svc"
	//
	Verbs = []string{
		Ping, Register, Unregister, Find, Join, Leave, Services,
	}
)

// OkVerb returns true for a correct verb.
func OkVerb(v string) bool {
	return v == Ping ||
		v == Register ||
		v == Unregister ||
		v == Find ||
		v == ReqKnock ||
		v == DoKnock ||
		v == Knock ||
		v == TestPort ||
		v == PortTest ||
		v == Join ||
		v == Join ||
		v == Leave
}

// Reply builds a reply message
func Reply(remote net.Addr) *Message {
	var m Message
	m.Address = remote.String()
	return &m
}

// ReplyError builds an error reply message
func ReplyError(remote net.Addr, code int) *Message {
	m := Reply(remote)
	m.Code = int32(code)
	return m
}

// ReplyOk builds an ok reply message
func ReplyOk(remote net.Addr, data string) *Message {
	m := ReplyError(remote, OkCode)
	m.Data = data
	return m
}

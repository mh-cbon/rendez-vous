package model

// Message of rendez-vous servers and clients.
type Message struct {
	// Type     string `json:"t,omitempty"` // type of the Message q or r (query or response)
	Query    string `json:"q,omitempty"` // Query value of the message
	Code     int    `json:"c,omitempty"` // Code of the response
	Pbk      []byte `json:"p,omitempty"` // ed25519 Public key (raw string)
	Sign     []byte `json:"s,omitempty"` // Sign of Pbk(v) (raw string)
	Value    string `json:"v,omitempty"` // Value to sign
	Address  string `json:"a,omitempty"` // Address of the querier sent in a response
	Response string `json:"r,omitempty"` // Value content of a response
}

// defines default verb
var (
	Ping       = "ping"
	Register   = "reg"
	Unregister = "unreg"
	Find       = "find"
	Join       = "join"
	Leave      = "leave"
)

// OkVerb returns true for a correct verb.
func OkVerb(v string) bool {
	return v == Ping ||
		v == Register ||
		v == Unregister ||
		v == Find ||
		v == Join ||
		v == Leave
}

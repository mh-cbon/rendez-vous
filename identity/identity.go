package identity

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/mh-cbon/dht/ed25519"
)

// PublicIdentity is a signed public identity
type PublicIdentity struct {
	Pbk, Value string
}

// Identity is a signed identity
type Identity struct {
	Pvk, Pbk, Sign, Value string
}

// FromPvk returns an Identity for registration
func FromPvk(pvk, value string) (*Identity, error) {
	if pvk == "" {
		pvkRaw, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		pvk = hex.EncodeToString(pvkRaw)
	}
	pvkRaw, pbkRaw, err := ed25519.PvkFromHex(pvk)
	if err != nil {
		return nil, err
	}
	signRaw := ed25519.Sign(pvkRaw, pbkRaw, []byte(value))
	pbkHex := hex.EncodeToString(pbkRaw)
	signHex := hex.EncodeToString(signRaw)

	return &Identity{
		Pvk:   pvk,
		Pbk:   pbkHex,
		Sign:  signHex,
		Value: value,
	}, nil
}

// FromPbk returns an Identity for lookup
func FromPbk(pbk, value string) (*PublicIdentity, error) {
	return &PublicIdentity{
		Pbk:   pbk,
		Value: value,
	}, nil
}

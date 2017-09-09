package identity

import (
	"encoding/hex"
	"log"

	"github.com/mh-cbon/dht/ed25519"
)

// PublicIdentity is a signed public identity
type PublicIdentity struct {
	Pbk, Value string
	BPbk       []byte
}

// Identity is a signed identity
type Identity struct {
	Pvk, Pbk, Sign, Value string
	BPvk, BPbk, BSign     []byte
}

// Derive given value
func (i Identity) Derive(value string) (*Identity, error) {
	return FromPvk(i.Pvk, value)
}

// FromPvk returns an Identity for registration
func FromPvk(pvk, value string) (*Identity, error) {
	if pvk == "" {
		pvkRaw, _, err := ed25519.GenerateKey(nil)
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
	log.Println(ed25519.Verify(pbkRaw, []byte(value), signRaw))
	log.Println(ed25519.Verify(pbkRaw, []byte(value), signRaw))

	return &Identity{
		Pvk:   pvk,
		Pbk:   pbkHex,
		Sign:  signHex,
		Value: value,
		BPvk:  pvkRaw,
		BPbk:  pbkRaw,
		BSign: signRaw,
	}, nil
}

// FromPbk returns an Identity for lookup
func FromPbk(pbk, value string) (*PublicIdentity, error) {
	bPbk, err := hex.DecodeString(pbk)
	if err != nil {
		return nil, err
	}
	return &PublicIdentity{
		Pbk:   pbk,
		Value: value,
		BPbk:  bPbk,
	}, nil
}

// Sign given value
func Sign(pvk string, value string) ([]byte, error) {
	pvkRaw, pbkRaw, err := ed25519.PvkFromHex(pvk)
	if err != nil {
		return nil, err
	}
	signRaw := ed25519.Sign(pvkRaw, pbkRaw, []byte(value))
	return signRaw, nil
}

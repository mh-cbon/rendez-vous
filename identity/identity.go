package identity

import (
	"encoding/hex"

	src "golang.org/x/crypto/ed25519"
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
		_, pvkRaw, err := src.GenerateKey(nil)
		if err != nil {
			return nil, err
		}
		pvk = hex.EncodeToString(pvkRaw)
	}
	pvkRaw, err := hex.DecodeString(pvk)
	if err != nil {
		return nil, err
	}
	pbkRaw := src.PrivateKey(pvkRaw).Public().(src.PublicKey)
	signRaw := src.Sign(pvkRaw, []byte(value))
	pbkHex := hex.EncodeToString(pbkRaw)
	signHex := hex.EncodeToString(signRaw)

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
	pvkRaw, err := hex.DecodeString(pvk)
	if err != nil {
		return nil, err
	}
	// pbkRaw := src.PrivateKey(pvkRaw).Public()
	signRaw := src.Sign(pvkRaw, []byte(value))
	return signRaw, nil
}

// Verify given value/sign
func Verify(pbkRaw, signRaw []byte, value string) bool {
	return src.Verify(pbkRaw, []byte(value), signRaw)
}

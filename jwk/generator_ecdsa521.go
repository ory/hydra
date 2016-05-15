package jwk

import (
	"crypto/rand"
	"github.com/go-errors/errors"
	"github.com/square/go-jose"
	"crypto/ecdsa"
	"crypto/elliptic"
)

type ECDSA521Generator struct{}

func (g *ECDSA521Generator) Generate(id string) (*jose.JsonWebKeySet, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{
			jose.JsonWebKey{
				Key:   key,
				KeyID: ider("private", id),
			},
			jose.JsonWebKey{
				Key:   &key.PublicKey,
				KeyID: ider("public", id),
			},
		},
	}, nil
}

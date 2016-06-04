package jwk

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/square/go-jose"
)

type RS256Generator struct{
	KeyLength int
}

func (g *RS256Generator) Generate(id string) (*jose.JsonWebKeySet, error) {
	if g.KeyLength < 4096 {
		g.KeyLength = 4096
	}

	key, err := rsa.GenerateKey(rand.Reader, g.KeyLength)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	// jose does not support this...
	key.Precomputed = rsa.PrecomputedValues{}
	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{
			{
				Key:   key,
				KeyID: ider("private", id),
			},
			{
				Key:   &key.PublicKey,
				KeyID: ider("public", id),
			},
		},
	}, nil
}

func ider(typ, id string) string {
	if id != "" {
		return fmt.Sprintf("%s:%s", typ, id)
	}
	return typ
}

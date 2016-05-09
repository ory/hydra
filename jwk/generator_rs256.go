package jwk

import (
	"crypto/rsa"
	"github.com/square/go-jose"
	"crypto/rand"
	"github.com/go-errors/errors"
	"fmt"
)

type RS256Generator struct {}

func (g *RS256Generator) Generate(id string) (set *jose.JsonWebKeySet, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	set = jose.JsonWebKeySet{
		&jose.JsonWebKey{
			Key: key,
			KeyID: fmt.Sprintf("private.%s", id),
		},
		&jose.JsonWebKey{
			Key: key.PublicKey,
			KeyID: fmt.Sprintf("public.%s", id),
		},
	}
	return set, nil
}
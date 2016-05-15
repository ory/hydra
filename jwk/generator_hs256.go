package jwk

import (
	"github.com/go-errors/errors"
	"github.com/square/go-jose"
	"github.com/ory-am/common/rand/sequence"
)

type HS256Generator struct {
	Length int
}

func (g *HS256Generator) Generate(id string) (*jose.JsonWebKeySet, error) {
	key, err := sequence.RuneSequence(g.Length, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789,.-;:_#+*!§$%&/()=?}][{<>"))
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{
			jose.JsonWebKey{
				Key:   []byte(string(key)),
				KeyID: id,
			},
		},
	}, nil
}

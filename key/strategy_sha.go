package key

import (
	"crypto/sha512"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/rand"
)

type SHAStrategy struct{}

func (s *SHAStrategy) SymmetricKey(id string) (*SymmetricKey, error) {
	key, err := rand.RandomBytes(32)
	if err != nil {
		return nil, errors.New(err)
	}

	hash := sha512.New()
	return &SymmetricKey{
		ID:  id,
		Key: hash.Sum(key),
	}, nil
}

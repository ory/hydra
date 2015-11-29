package hash

import (
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	WorkFactor int
}

func (b *BCrypt) Hash(data string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(data), b.WorkFactor)
	if err != nil {
		return "", errors.New(err)
	}
	return string(s), nil
}

func (b *BCrypt) Compare(hash string, data string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)); err != nil {
		return errors.New(err)
	}
	return nil
}

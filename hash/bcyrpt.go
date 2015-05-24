package hash

import "golang.org/x/crypto/bcrypt"

type BCrypt struct {
	WorkFactor int
}

func (b *BCrypt) Hash(data string) (string, error) {
	if b.WorkFactor == 0 {
		b.WorkFactor = 12
	}
	s, err := bcrypt.GenerateFromPassword([]byte(data), b.WorkFactor)
	return string(s), err
}

func (b *BCrypt) Compare(hash string, data string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
}

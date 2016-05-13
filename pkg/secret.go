package pkg

import "github.com/ory-am/common/rand/sequence"

var secretCharSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.,!$%&/()=?><")

func GenerateSecret(length int) ([]byte, error) {
	secret, err := sequence.RuneSequence(length, secretCharSet)
	if err != nil {
		return []byte{}, err
	}
	return []byte(string(secret)), nil
}

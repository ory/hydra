package jwk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

type AEAD struct {
	Key []byte
}

func (c *AEAD) Encrypt(plaintext []byte) (string, error) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	if len(c.Key) < 32 {
		return "", errors.Errorf("Key must be longer 32 bytes, got %d bytes", len(c.Key))
	}

	block, err := aes.NewCipher(c.Key[:32])
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return base64.URLEncoding.EncodeToString(append(ciphertext, nonce...)), nil
}

func (c *AEAD) Decrypt(ciphertext string) ([]byte, error) {
	if len(c.Key) < 32 {
		return []byte{}, errors.Errorf("Key must be longer 32 bytes, got %d bytes", len(c.Key))
	}

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return []byte{}, errors.Wrap(err, "")
	}

	n := len(raw)
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return []byte{}, errors.Wrap(err, "")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "")
	}

	plaintext, err := aesgcm.Open(nil, raw[n-12:n], raw[:n-12], nil)
	if err != nil {
		return []byte{}, errors.Wrap(err, "")
	}

	return plaintext, nil
}

// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"

	"github.com/ory/x/errorsx"
)

type AESGCM struct {
	c Dependencies
}

func NewAESGCM(c Dependencies) *AESGCM {
	return &AESGCM{c: c}
}

func aeadKey(key []byte) *[32]byte {
	var result [32]byte
	copy(result[:], key[:32])
	return &result
}

func (c *AESGCM) Encrypt(ctx context.Context, plaintext, additionalData []byte) (string, error) {
	global, err := c.c.GetGlobalSecret(ctx)
	if err != nil {
		return "", err
	}

	rotated, err := c.c.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return "", err
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return "", errors.Errorf("at least one encryption key must be defined but none were")
	}

	if len(keys[0]) < 32 {
		return "", errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(keys[0]))
	}

	ciphertext, err := aesGCMEncrypt(plaintext, aeadKey(keys[0]), additionalData)
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	return encode(ciphertext, additionalData), nil
}

func (c *AESGCM) Decrypt(ctx context.Context, s string) (plaintext, aad []byte, err error) {
	ciphertext, aad, err := decode(s)
	if err != nil {
		return nil, nil, err
	}

	global, err := c.c.GetGlobalSecret(ctx)
	if err != nil {
		return nil, nil, err
	}

	rotated, err := c.c.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return nil, nil, err
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return nil, nil, errors.Errorf("at least one decryption key must be defined but none were")
	}

	for _, key := range keys {
		if plaintext, err = c.decrypt(ciphertext, key, aad); err == nil {
			return plaintext, aad, nil
		}
	}

	return nil, nil, err
}

func (c *AESGCM) decrypt(ciphertext []byte, key, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	plaintext, err := aesGCMDecrypt(ciphertext, aeadKey(key), additionalData)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return plaintext, nil
}

// aesGCMEncrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func aesGCMEncrypt(plaintext []byte, key *[32]byte, additionalData []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, additionalData), nil
}

// aesGCMDecrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func aesGCMDecrypt(ciphertext []byte, key *[32]byte, additionalData []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		additionalData,
	)
}

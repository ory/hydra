// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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
	key, err := encryptionKey(ctx, c.c, 32)
	if err != nil {
		return "", err
	}

	ciphertext, err := aesGCMEncrypt(plaintext, aeadKey(key), additionalData)
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *AESGCM) Decrypt(ctx context.Context, ciphertext string, aad []byte) (plaintext []byte, err error) {
	msg, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	keys, err := allKeys(ctx, c.c)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	for _, key := range keys {
		if plaintext, err = c.decrypt(msg, key, aad); err == nil {
			return plaintext, nil
		}
	}

	return nil, err
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

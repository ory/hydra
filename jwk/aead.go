// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/errorsx"
)

type AEAD struct {
	c *config.DefaultProvider
}

func NewAEAD(c *config.DefaultProvider) *AEAD {
	return &AEAD{c: c}
}

func aeadKey(key []byte) *[32]byte {
	var result [32]byte
	copy(result[:], key[:32])
	return &result
}

func (c *AEAD) Encrypt(ctx context.Context, plaintext []byte) (string, error) {
	return c.EncryptWithAdditionalData(ctx, plaintext, nil)
}

func (c *AEAD) EncryptWithAdditionalData(ctx context.Context, plaintext, additionalData []byte) (string, error) {
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

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *AEAD) Decrypt(ctx context.Context, ciphertext string) (p []byte, err error) {
	return c.DecryptWithAdditionalData(ctx, ciphertext, nil)
}

func (c *AEAD) DecryptWithAdditionalData(ctx context.Context, ciphertext string, additionalData []byte) (p []byte, err error) {
	global, err := c.c.GetGlobalSecret(ctx)
	if err != nil {
		return nil, err
	}

	rotated, err := c.c.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return nil, err
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return nil, errors.Errorf("at least one decryption key must be defined but none were")
	}

	for _, key := range keys {
		if p, err = c.decrypt(ciphertext, key, additionalData); err == nil {
			return p, nil
		}
	}

	return nil, err
}

func (c *AEAD) decrypt(ciphertext string, key, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	plaintext, err := aesGCMDecrypt(raw, aeadKey(key), additionalData)
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

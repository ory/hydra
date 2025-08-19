// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

var _ Cipher = (*XChaCha20Poly1305)(nil)

type (
	XChaCha20Poly1305 struct {
		d Dependencies
	}
)

func NewXChaCha20Poly1305(d Dependencies) *XChaCha20Poly1305 {
	return &XChaCha20Poly1305{d}
}

func (x *XChaCha20Poly1305) Encrypt(ctx context.Context, plaintext, additionalData []byte) (string, error) {
	key, err := encryptionKey(ctx, x.d, chacha20poly1305.KeySize)
	if err != nil {
		return "", err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// Make sure the size calculation does not overflow.
	if len(plaintext) > math.MaxInt-aead.NonceSize()-aead.Overhead() {
		return "", errors.WithStack(fmt.Errorf("plaintext too large"))
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plaintext)+aead.Overhead())
	_, err = cryptorand.Read(nonce)
	if err != nil {
		return "", errors.WithStack(err)
	}

	ciphertext := aead.Seal(nonce, nonce, plaintext, additionalData)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (x *XChaCha20Poly1305) Decrypt(ctx context.Context, ciphertext string, aad []byte) (plaintext []byte, err error) {
	msg, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(msg) < chacha20poly1305.NonceSizeX {
		return nil, errors.WithStack(fmt.Errorf("malformed ciphertext: too short"))
	}
	nonce, ciphered := msg[:chacha20poly1305.NonceSizeX], msg[chacha20poly1305.NonceSizeX:]

	keys, err := allKeys(ctx, x.d)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var aead cipher.AEAD
	for _, key := range keys {
		aead, err = chacha20poly1305.NewX(key)
		if err != nil {
			continue
		}
		plaintext, err = aead.Open(nil, nonce, ciphered, aad)
		if err == nil {
			return plaintext, nil
		}
	}

	return nil, errors.WithStack(err)
}

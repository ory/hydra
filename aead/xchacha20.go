// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/pkg/errors"

	"github.com/ory/x/errorsx"
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
	global, err := x.d.GetGlobalSecret(ctx)
	if err != nil {
		return "", err
	}

	rotated, err := x.d.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return "", err
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return "", fmt.Errorf("at least one encryption key must be defined but none were")
	}

	if len(keys[0]) != chacha20poly1305.KeySize {
		return "", fmt.Errorf("key must be exactly %d bytes long, got %d bytes", chacha20poly1305.KeySize, len(keys[0]))
	}

	cipher, err := chacha20poly1305.NewX(keys[0])
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	nonce := make([]byte, cipher.NonceSize(), cipher.NonceSize()+len(plaintext)+cipher.Overhead())
	_, err = cryptorand.Read(nonce)
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	ciphertext := cipher.Seal(nonce, nonce, plaintext, additionalData)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (x *XChaCha20Poly1305) Decrypt(ctx context.Context, ciphertext string, aad []byte) (plaintext []byte, err error) {
	msg, err := base64.RawURLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	if len(msg) < chacha20poly1305.NonceSizeX {
		return nil, errorsx.WithStack(fmt.Errorf("malformed ciphertext: too short"))
	}
	nonce, ciphered := msg[:chacha20poly1305.NonceSizeX], msg[chacha20poly1305.NonceSizeX:]

	global, err := x.d.GetGlobalSecret(ctx)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	rotated, err := x.d.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return nil, errorsx.WithStack(errors.Errorf("at least one decryption key must be defined but none were"))
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

	return nil, errorsx.WithStack(err)
}

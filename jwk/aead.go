// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/base64"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/v2/driver/config"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"
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

	ciphertext, err := cryptopasta.Encrypt(plaintext, aeadKey(keys[0]))
	if err != nil {
		return "", errorsx.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *AEAD) Decrypt(ctx context.Context, ciphertext string) (p []byte, err error) {
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
		if p, err = c.decrypt(ciphertext, key); err == nil {
			return p, nil
		}
	}

	return nil, err
}

func (c *AEAD) decrypt(ciphertext string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	plaintext, err := cryptopasta.Decrypt(raw, aeadKey(key))
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return plaintext, nil
}

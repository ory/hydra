// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !commercial

// Package aead provides the authenticated cipher (AEAD) that Ory uses to
// seal opaque payloads such as pagination tokens.
package aead

import (
	"crypto/cipher"

	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

// New returns the XChaCha20-Poly1305 AEAD that seals opaque payloads under
// the given key. Its 192-bit random nonce imposes no practical limit on the
// number of payloads sealed per key.
func New(key [32]byte) (cipher.AEAD, error) {
	a, err := chacha20poly1305.NewX(key[:])
	return a, errors.Wrap(err, "cannot create XChaCha20-Poly1305 AEAD")
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"

	"github.com/ory/fosite"
)

// Cipher provides AEAD (authenticated encryption with associated data). The
// ciphertext is returned base64url-encoded.
type Cipher interface {
	// Encrypt encrypts and encodes the given plaintext, optionally using
	// additiona data.
	Encrypt(ctx context.Context, plaintext, additionalData []byte) (ciphertext string, err error)

	// Decrypt decodes, decrypts, and verifies the plaintext and additional data
	// from the ciphertext. The ciphertext must be given in the form as returned
	// by Encrypt.
	Decrypt(ctx context.Context, ciphertext string, additionalData []byte) (plaintext []byte, err error)
}

type Dependencies interface {
	fosite.GlobalSecretProvider
	fosite.RotatedGlobalSecretsProvider
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ory/fosite"
	"github.com/ory/x/errorsx"
)

// Cipher provides AEAD (authenticated encryption with associated data). The
// ciphertext is returned base64url-encoded. If the additional data is not
// empty, it is also base64url-encoded and appended to the ciphertext, separated
// by a dot (.)
type Cipher interface {
	// Encrypt encrypts and encodes the given plaintext and additiona data.
	//
	//	Encrypt(ctx, []byte("secret message"), []byte("non-confidential data"))
	//		// might return "dmuEZUYjJnrgOCm4D0edAzAFnQLqRxwa_Bug6IfGXVAoSCEmngaRxsH7"
	//	Encrypt(ctx, []byte("secret message"), []byte("non-confidential data"))
	//		// might return "3z-hDUqTwOe9rNX0ki-KXaBpDP5wGqEJQfyq15v3TfM1ndIpmR_c5UrH.bm9uLWNvbmZpZGVudGlhbCBkYXRh"
	Encrypt(ctx context.Context, plaintext, additionalData []byte) (ciphertext string, err error)

	// Decrypt decodes, decrypts, and verifies the plaintext and additional data
	// from the ciphertext. The ciphertext must be given in the form as returned
	// by Encrypt.
	Decrypt(ctx context.Context, ciphertext string) (plaintext, additionalData []byte, err error)
}

type Dependencies interface {
	fosite.GlobalSecretProvider
	fosite.RotatedGlobalSecretsProvider
}

func encode(ciphertext, additionalData []byte) string {
	var aad string
	if len(additionalData) > 0 {
		aad = "." + base64.RawURLEncoding.EncodeToString(additionalData)
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext) + aad
}

func decode(s string) (ciphertext, aad []byte, err error) {
	split := strings.SplitN(s, ".", 2)
	if len(split) == 2 {
		aad, err = base64.RawURLEncoding.DecodeString(split[1])
		if err != nil {
			return nil, nil, errorsx.WithStack(fmt.Errorf("malformed ciphertext: %w", err))
		}
	}
	split[0] = strings.TrimRight(split[0], "=") // We previously used padding in the base64-encoding; this is for compatibility.
	ciphertext, err = base64.RawURLEncoding.DecodeString(split[0])
	if err != nil {
		return nil, nil, errorsx.WithStack(fmt.Errorf("malformed ciphertext: %w", err))
	}
	return
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/aead"
)

// Decode decodes the given string to a value.
func Decode[T any](ctx context.Context, cipher *aead.XChaCha20Poly1305, encoded string) (*T, error) {
	plaintext, err := cipher.Decrypt(ctx, encoded, nil)
	if err != nil {
		return nil, err
	}

	rawBytes, err := gzip.NewReader(bytes.NewReader(plaintext))
	if err != nil {
		return nil, err
	}
	defer func() { _ = rawBytes.Close() }()

	var val T
	if err = json.NewDecoder(rawBytes).Decode(&val); err != nil {
		return nil, err
	}

	return &val, nil
}

// Encode encodes the given value to a string.
func Encode(ctx context.Context, cipher aead.Cipher, val any) (s string, err error) {
	// Steps:
	// 1. Encode to JSON
	// 2. GZIP
	// 3. Encrypt with AEAD (AES-GCM) + Base64 URL-encode
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)

	if err = json.NewEncoder(gz).Encode(val); err != nil {
		return "", err
	}
	if err = gz.Close(); err != nil {
		return "", err
	}

	return cipher.Encrypt(ctx, b.Bytes(), nil)
}

// SetCookie encrypts the given value and sets it in a cookie.
func SetCookie(ctx context.Context, w http.ResponseWriter, cipher aead.Cipher, cookieName string, value any) error {
	cookie, err := Encode(ctx, cipher, value)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// FromCookie looks up the value stored in the cookie and decodes it.
func FromCookie[T any](ctx context.Context, r *http.Request, cipher *jwk.AEAD, cookieName string) (*T, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return Decode[T](ctx, cipher, cookie.Value)
}

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
	plaintext, _, err := cipher.Decrypt(ctx, encoded)
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

// EncodeFromContext encodes the value stored in the context under the given cookie name.
func EncodeFromContext(ctx context.Context, cipher aead.Cipher, cookieName string) (s string, err error) {
	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return "", errors.WithStack(ErrNoValueInCtx)
	}

	return Encode(ctx, cipher, v.Ptr)
}

// SetCookie looks up the value stored in the context under the given cookie name and sets it as a cookie on the
// response writer.
func SetCookie(ctx context.Context, w http.ResponseWriter, cipher aead.Cipher, cookieName string) error {
	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return errors.WithStack(ErrNoValueInCtx)
	}

	cookie, err := Encode(ctx, cipher, v.Ptr)
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

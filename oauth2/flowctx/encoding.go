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

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
)

type (
	data struct {
		Purpose purpose `json:"p,omitempty"`
	}
	purpose     int
	CodecOption func(ad *data)
)

const (
	loginChallenge purpose = iota
	loginVerifier
	consentChallenge
	consentVerifier
)

func withPurpose(purpose purpose) CodecOption { return func(ad *data) { ad.Purpose = purpose } }

var (
	AsLoginChallenge   = withPurpose(loginChallenge)
	AsLoginVerifier    = withPurpose(loginVerifier)
	AsConsentChallenge = withPurpose(consentChallenge)
	AsConsentVerifier  = withPurpose(consentVerifier)
)

func additionalDataFromOpts(opts ...CodecOption) []byte {
	if len(opts) == 0 {
		return nil
	}
	ad := &data{}
	for _, o := range opts {
		o(ad)
	}
	b, err := json.Marshal(ad)
	if err != nil {
		// Panic is OK here because the struct and the parameters are all known.
		panic("failed to marshal additional data: " + errors.WithStack(err).Error())
	}

	return b
}

// Decode decodes the given string to a value.
func Decode[T any](ctx context.Context, cipher aead.Cipher, encoded string, opts ...CodecOption) (*T, error) {
	plaintext, err := cipher.Decrypt(ctx, encoded, additionalDataFromOpts(opts...))
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
func Encode(ctx context.Context, cipher aead.Cipher, val any, opts ...CodecOption) (s string, err error) {
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

	return cipher.Encrypt(ctx, b.Bytes(), additionalDataFromOpts(opts...))
}

// SetCookie encrypts the given value and sets it in a cookie.
func SetCookie(ctx context.Context, w http.ResponseWriter, reg interface {
	FlowCipher() *aead.XChaCha20Poly1305
	config.Provider
}, cookieName string, value any, opts ...CodecOption) error {
	cipher := reg.FlowCipher()
	cookie, err := Encode(ctx, cipher, value, opts...)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cookie,
		HttpOnly: true,
		Domain:   reg.Config().CookieDomain(ctx),
		Secure:   reg.Config().CookieSecure(ctx),
		SameSite: reg.Config().CookieSameSiteMode(ctx),
	})

	return nil
}

// FromCookie looks up the value stored in the cookie and decodes it.
func FromCookie[T any](ctx context.Context, r *http.Request, cipher aead.Cipher, cookieName string, opts ...CodecOption) (*T, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidClient.WithHint("No cookie found for this request. Please initiate a new flow and retry."))
	}

	return Decode[T](ctx, cipher, cookie.Value, opts...)
}

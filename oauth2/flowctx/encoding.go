// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/aead"
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
	// 3. Encrypt with AEAD (XChaCha20-Poly1305) + Base64 URL-encode
	var b bytes.Buffer

	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return "", err
	}

	if err = json.NewEncoder(gz).Encode(val); err != nil {
		return "", err
	}

	if err = gz.Close(); err != nil {
		return "", err
	}

	return cipher.Encrypt(ctx, b.Bytes(), additionalDataFromOpts(opts...))
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
)

// Set the aadAccessTokenPrefix to something unique to avoid ciphertext confusion with other usages of the AEAD cipher.
var aadAccessTokenPrefix = "vc-nonce-at:" // nolint:gosec

func (p *Persister) NewNonce(ctx context.Context, accessToken string, expiresIn time.Time) (res string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.NewNonce")
	defer otelx.End(span, &err)

	plaintext := x.IntToBytes(expiresIn.Unix())
	aad := []byte(aadAccessTokenPrefix + accessToken)

	return p.r.FlowCipher().Encrypt(ctx, plaintext, aad)
}

func (p *Persister) IsNonceValid(ctx context.Context, accessToken, nonce string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.IsNonceValid")
	defer otelx.End(span, &err)

	aad := []byte(aadAccessTokenPrefix + accessToken)
	plaintext, err := p.r.FlowCipher().Decrypt(ctx, nonce, aad)
	if err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("The nonce is invalid."))
	}

	exp, err := x.BytesToInt(plaintext)
	if err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("The nonce is invalid.")) // should never happen
	}

	if exp < time.Now().Unix() {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("The nonce has expired."))
	}

	return nil
}

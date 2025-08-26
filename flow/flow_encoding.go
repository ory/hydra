// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"
)

type decodeDependencies interface {
	CipherProvider
	x.NetworkProvider
	config.Provider
	x.TracingProvider
}

func decodeChallenge(ctx context.Context, d decodeDependencies, challenge string, p purpose) (_ *Flow, err error) {
	f, err := Decode[Flow](ctx, d.FlowCipher(), challenge, withPurpose(p))
	if err != nil {
		return nil, errors.WithStack(x.ErrNotFound.WithWrap(err))
	}

	if f.NID != d.Networker().NetworkID(ctx) {
		return nil, errors.WithStack(x.ErrNotFound.WithDescription("Network IDs are not matching."))
	}

	if f.RequestedAt.Add(d.Config().ConsentRequestMaxAge(ctx)).Before(time.Now()) {
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized.WithHint("The login request has expired, please try again."))
	}

	return f, nil
}

func DecodeFromLoginChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromLoginChallenge")
	defer otelx.End(span, &err)

	return decodeChallenge(ctx, d, challenge, loginChallenge)
}

func DecodeFromConsentChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromConsentChallenge")
	defer otelx.End(span, &err)

	return decodeChallenge(ctx, d, challenge, consentChallenge)
}

func DecodeFromDeviceChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromDeviceChallenge")
	defer otelx.End(span, &err)

	return decodeChallenge(ctx, d, challenge, deviceChallenge)
}

func DecodeAndInvalidateLoginVerifier(ctx context.Context, d decodeDependencies, verifier string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeAndInvalidateLoginVerifier")
	defer otelx.End(span, &err)

	f, err := Decode[Flow](ctx, d.FlowCipher(), verifier, AsLoginVerifier)
	if err != nil {
		return nil, errors.WithStack(sqlcon.ErrNoRows)
	}

	if f.NID != d.Networker().NetworkID(ctx) {
		return nil, errors.WithStack(sqlcon.ErrNoRows)
	}

	if err := f.InvalidateLoginRequest(); err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	return f, nil
}

func DecodeAndInvalidateConsentVerifier(ctx context.Context, d decodeDependencies, verifier string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeAndInvalidateLoginVerifier")
	defer otelx.End(span, &err)

	f, err := Decode[Flow](ctx, d.FlowCipher(), verifier, AsConsentVerifier)
	if err != nil {
		return nil, errors.WithStack(fosite.ErrAccessDenied.WithHint("The consent verifier has already been used, has not been granted, or is invalid."))
	}

	if f.NID != d.Networker().NetworkID(ctx) {
		return nil, errors.WithStack(sqlcon.ErrNoRows)
	}

	if err = f.InvalidateConsentRequest(); err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	return f, nil
}

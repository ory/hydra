// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
)

type decodeDependencies interface {
	CipherProvider
	x.NetworkProvider
	config.Provider
	x.TracingProvider
}

func decodeFlow(ctx context.Context, d decodeDependencies, enc string, p purpose) (_ *Flow, err error) {
	f, err := Decode[Flow](ctx, d.FlowCipher(), enc, withPurpose(p))
	if err != nil {
		return nil, errors.WithStack(x.ErrNotFound.WithWrap(err))
	}

	if f.NID != d.Networker().NetworkID(ctx) {
		return nil, errors.WithStack(x.ErrNotFound.WithDescription("Network IDs are not matching."))
	}

	if f.RequestedAt.Add(d.Config().ConsentRequestMaxAge(ctx)).Before(time.Now()) {
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized.WithHintf("The %s request has expired, please try again.", p.RequestType()))
	}

	return f, nil
}

func DecodeFromLoginChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromLoginChallenge")
	defer otelx.End(span, &err)

	return decodeFlow(ctx, d, challenge, loginChallenge)
}

func DecodeFromConsentChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromConsentChallenge")
	defer otelx.End(span, &err)

	return decodeFlow(ctx, d, challenge, consentChallenge)
}

func DecodeFromDeviceChallenge(ctx context.Context, d decodeDependencies, challenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromDeviceChallenge")
	defer otelx.End(span, &err)

	return decodeFlow(ctx, d, challenge, deviceChallenge)
}

func decodeVerifier(ctx context.Context, d decodeDependencies, verifier string, p purpose) (_ *Flow, err error) {
	f, err := decodeFlow(ctx, d, verifier, p)
	if err != nil {
		if errors.Is(err, x.ErrNotFound) {
			return nil, errors.WithStack(fosite.ErrAccessDenied.WithHintf("The %s verifier has already been used, has not been granted, or is invalid.", p.RequestType()))
		}
		return nil, err
	}

	return f, nil
}

func DecodeAndInvalidateLoginVerifier(ctx context.Context, d decodeDependencies, verifier string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeAndInvalidateLoginVerifier")
	defer otelx.End(span, &err)

	f, err := decodeVerifier(ctx, d, verifier, loginVerifier)
	if err != nil {
		return nil, err
	}

	if err := f.InvalidateLoginRequest(); err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	return f, nil
}

func DecodeAndInvalidateDeviceVerifier(ctx context.Context, d decodeDependencies, verifier string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeAndInvalidateDeviceVerifier")
	defer otelx.End(span, &err)

	f, err := decodeVerifier(ctx, d, verifier, deviceVerifier)
	if err != nil {
		return nil, err
	}

	if err = f.InvalidateDeviceRequest(); err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	return f, nil
}

func DecodeAndInvalidateConsentVerifier(ctx context.Context, d decodeDependencies, verifier string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeAndInvalidateLoginVerifier")
	defer otelx.End(span, &err)

	f, err := decodeVerifier(ctx, d, verifier, consentVerifier)
	if err != nil {
		return nil, err
	}

	if err = f.InvalidateConsentRequest(); err != nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	return f, nil
}

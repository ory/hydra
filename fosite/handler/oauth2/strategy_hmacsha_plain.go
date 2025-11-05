// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/v2/fosite"
	enigma "github.com/ory/hydra/v2/fosite/token/hmac"
)

var (
	_ AuthorizeCodeStrategy = (*HMACSHAStrategyUnPrefixed)(nil)
	_ AccessTokenStrategy   = (*HMACSHAStrategyUnPrefixed)(nil)
	_ RefreshTokenStrategy  = (*HMACSHAStrategyUnPrefixed)(nil)
)

type HMACSHAStrategyUnPrefixed struct {
	Enigma *enigma.HMACStrategy
	Config LifespanConfigProvider
}

func NewHMACSHAStrategyUnPrefixed(
	enigma *enigma.HMACStrategy,
	config LifespanConfigProvider,
) *HMACSHAStrategyUnPrefixed {
	return &HMACSHAStrategyUnPrefixed{
		Enigma: enigma,
		Config: config,
	}
}

func (h *HMACSHAStrategyUnPrefixed) AccessTokenSignature(ctx context.Context, token string) string {
	return h.Enigma.Signature(token)
}

func (h *HMACSHAStrategyUnPrefixed) RefreshTokenSignature(ctx context.Context, token string) string {
	return h.Enigma.Signature(token)
}

func (h *HMACSHAStrategyUnPrefixed) AuthorizeCodeSignature(ctx context.Context, token string) string {
	return h.Enigma.Signature(token)
}

func (h *HMACSHAStrategyUnPrefixed) GenerateAccessToken(ctx context.Context, _ fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.Enigma.Generate(ctx)
	if err != nil {
		return "", "", err
	}

	return token, sig, nil
}

func (h *HMACSHAStrategyUnPrefixed) ValidateAccessToken(ctx context.Context, r fosite.Requester, token string) (err error) {
	exp := r.GetSession().GetExpiresAt(fosite.AccessToken)
	if exp.IsZero() && r.GetRequestedAt().Add(h.Config.GetAccessTokenLifespan(ctx)).Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrTokenExpired.WithHintf("Access token expired at '%s'.", r.GetRequestedAt().Add(h.Config.GetAccessTokenLifespan(ctx))))
	}

	if !exp.IsZero() && exp.Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrTokenExpired.WithHintf("Access token expired at '%s'.", exp))
	}

	return h.Enigma.Validate(ctx, token)
}

func (h *HMACSHAStrategyUnPrefixed) GenerateRefreshToken(ctx context.Context, _ fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.Enigma.Generate(ctx)
	if err != nil {
		return "", "", err
	}

	return token, sig, nil
}

func (h *HMACSHAStrategyUnPrefixed) ValidateRefreshToken(ctx context.Context, r fosite.Requester, token string) (err error) {
	exp := r.GetSession().GetExpiresAt(fosite.RefreshToken)
	if exp.IsZero() {
		// Unlimited lifetime
		return h.Enigma.Validate(ctx, token)
	}

	if !exp.IsZero() && exp.Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrTokenExpired.WithHintf("Refresh token expired at '%s'.", exp))
	}

	return h.Enigma.Validate(ctx, token)
}

func (h *HMACSHAStrategyUnPrefixed) GenerateAuthorizeCode(ctx context.Context, _ fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.Enigma.Generate(ctx)
	if err != nil {
		return "", "", err
	}

	return token, sig, nil
}

func (h *HMACSHAStrategyUnPrefixed) ValidateAuthorizeCode(ctx context.Context, r fosite.Requester, token string) (err error) {
	exp := r.GetSession().GetExpiresAt(fosite.AuthorizeCode)
	if exp.IsZero() && r.GetRequestedAt().Add(h.Config.GetAuthorizeCodeLifespan(ctx)).Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrTokenExpired.WithHintf("Authorize code expired at '%s'.", r.GetRequestedAt().Add(h.Config.GetAuthorizeCodeLifespan(ctx))))
	}

	if !exp.IsZero() && exp.Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrTokenExpired.WithHintf("Authorize code expired at '%s'.", exp))
	}

	return h.Enigma.Validate(ctx, token)
}

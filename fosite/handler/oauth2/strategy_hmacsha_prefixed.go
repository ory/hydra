// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"strings"

	enigma "github.com/ory/hydra/v2/fosite/token/hmac"

	"github.com/ory/hydra/v2/fosite"
)

var (
	_ AuthorizeCodeStrategy = (*HMACSHAStrategy)(nil)
	_ AccessTokenStrategy   = (*HMACSHAStrategy)(nil)
	_ RefreshTokenStrategy  = (*HMACSHAStrategy)(nil)
)

type HMACSHAStrategy struct {
	*HMACSHAStrategyUnPrefixed
}

func NewHMACSHAStrategy(
	enigma *enigma.HMACStrategy,
	config LifespanConfigProvider,
) *HMACSHAStrategy {
	return &HMACSHAStrategy{
		HMACSHAStrategyUnPrefixed: NewHMACSHAStrategyUnPrefixed(enigma, config),
	}
}

func (h *HMACSHAStrategy) getPrefix(part string) string {
	return fmt.Sprintf("ory_%s_", part)
}

func (h *HMACSHAStrategy) trimPrefix(token, part string) string {
	return strings.TrimPrefix(token, h.getPrefix(part))
}

func (h *HMACSHAStrategy) setPrefix(token, part string) string {
	if token == "" {
		return ""
	}
	return h.getPrefix(part) + token
}

func (h *HMACSHAStrategy) GenerateAccessToken(ctx context.Context, r fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.HMACSHAStrategyUnPrefixed.GenerateAccessToken(ctx, r)
	return h.setPrefix(token, "at"), sig, err
}

func (h *HMACSHAStrategy) ValidateAccessToken(ctx context.Context, r fosite.Requester, token string) (err error) {
	return h.HMACSHAStrategyUnPrefixed.ValidateAccessToken(ctx, r, h.trimPrefix(token, "at"))
}

func (h *HMACSHAStrategy) GenerateRefreshToken(ctx context.Context, r fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.HMACSHAStrategyUnPrefixed.GenerateRefreshToken(ctx, r)
	return h.setPrefix(token, "rt"), sig, err
}

func (h *HMACSHAStrategy) ValidateRefreshToken(ctx context.Context, r fosite.Requester, token string) (err error) {
	return h.HMACSHAStrategyUnPrefixed.ValidateRefreshToken(ctx, r, h.trimPrefix(token, "rt"))
}

func (h *HMACSHAStrategy) GenerateAuthorizeCode(ctx context.Context, r fosite.Requester) (token string, signature string, err error) {
	token, sig, err := h.HMACSHAStrategyUnPrefixed.GenerateAuthorizeCode(ctx, r)
	return h.setPrefix(token, "ac"), sig, err
}

func (h *HMACSHAStrategy) ValidateAuthorizeCode(ctx context.Context, r fosite.Requester, token string) (err error) {
	return h.HMACSHAStrategyUnPrefixed.ValidateAuthorizeCode(ctx, r, h.trimPrefix(token, "ac"))
}

// Implements providers interface
func (h *HMACSHAStrategy) AuthorizeCodeStrategy() AuthorizeCodeStrategy {
	return h
}

func (h *HMACSHAStrategy) AccessTokenStrategy() AccessTokenStrategy {
	return h
}

func (h *HMACSHAStrategy) RefreshTokenStrategy() RefreshTokenStrategy {
	return h
}

// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628

import (
	"context"
	"strings"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/randx"

	"github.com/ory/hydra/v2/fosite"
	enigma "github.com/ory/hydra/v2/fosite/token/hmac"
)

var (
	_ DeviceRateLimitStrategy = (*DefaultDeviceStrategy)(nil)
	_ DeviceCodeStrategy      = (*DefaultDeviceStrategy)(nil)
	_ UserCodeStrategy        = (*DefaultDeviceStrategy)(nil)
)

// DefaultDeviceStrategy implements the default device strategy
type DefaultDeviceStrategy struct {
	Enigma *enigma.HMACStrategy
	Config interface {
		fosite.DeviceProvider
		fosite.DeviceAndUserCodeLifespanProvider
		fosite.UserCodeProvider
	}
}

// GenerateUserCode generates a user_code
func (h *DefaultDeviceStrategy) GenerateUserCode(ctx context.Context) (string, string, error) {
	seq, err := randx.RuneSequence(h.Config.GetUserCodeLength(ctx), h.Config.GetUserCodeSymbols(ctx))
	if err != nil {
		return "", "", err
	}
	userCode := string(seq)
	signUserCode, signErr := h.UserCodeSignature(ctx, userCode)
	if signErr != nil {
		return "", "", err
	}
	return userCode, signUserCode, nil
}

// UserCodeSignature generates a user_code signature
func (h *DefaultDeviceStrategy) UserCodeSignature(ctx context.Context, token string) (string, error) {
	return h.Enigma.GenerateHMACForString(ctx, token)
}

// ValidateUserCode validates a user_code
// This function only checks if the device request session is active as we cannot verify the authenticity of the token.
// Unlike other tokens, the user_code is of limited length, which means that we cannot include the HMAC signature in the token itself.
// The only way to check the validity of the user_code is to check if its signature is stored in storage.
func (h *DefaultDeviceStrategy) ValidateUserCode(ctx context.Context, r fosite.DeviceRequester, code string) error {
	exp := r.GetSession().GetExpiresAt(fosite.UserCode)
	if exp.IsZero() && r.GetRequestedAt().Add(h.Config.GetDeviceAndUserCodeLifespan(ctx)).Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrDeviceExpiredToken.WithHintf("User code expired at '%s'.", r.GetRequestedAt().Add(h.Config.GetDeviceAndUserCodeLifespan(ctx))))
	}
	if !exp.IsZero() && exp.Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrDeviceExpiredToken.WithHintf("User code expired at '%s'.", exp))
	}
	return nil
}

// GenerateDeviceCode generates a device_code
func (h *DefaultDeviceStrategy) GenerateDeviceCode(ctx context.Context) (string, string, error) {
	token, sig, err := h.Enigma.Generate(ctx)
	if err != nil {
		return "", "", err
	}

	return "ory_dc_" + token, sig, nil
}

// DeviceCodeSignature generates a device_code signature
func (h *DefaultDeviceStrategy) DeviceCodeSignature(ctx context.Context, token string) (string, error) {
	return h.Enigma.Signature(token), nil
}

// ValidateDeviceCode validates a device_code
func (h *DefaultDeviceStrategy) ValidateDeviceCode(ctx context.Context, r fosite.DeviceRequester, code string) error {
	exp := r.GetSession().GetExpiresAt(fosite.DeviceCode)
	if exp.IsZero() && r.GetRequestedAt().Add(h.Config.GetDeviceAndUserCodeLifespan(ctx)).Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrDeviceExpiredToken.WithHintf("Device code expired at '%s'.", r.GetRequestedAt().Add(h.Config.GetDeviceAndUserCodeLifespan(ctx))))
	}

	if !exp.IsZero() && exp.Before(time.Now().UTC()) {
		return errorsx.WithStack(fosite.ErrDeviceExpiredToken.WithHintf("Device code expired at '%s'.", exp))
	}

	return h.Enigma.Validate(ctx, strings.TrimPrefix(code, "ory_dc_"))
}

// ShouldRateLimit is used to decide whether a request should be rate-limited
func (h *DefaultDeviceStrategy) ShouldRateLimit(context context.Context, code string) (bool, error) {
	return false, nil
}

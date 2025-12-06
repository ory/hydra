// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628

import (
	"context"
	"fmt"
	"time"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
)

// MaxAttempts for retrying the generation of user codes.
const MaxAttempts = 3

// DeviceAuthHandler is a response handler for the Device Authorisation Grant as
// defined in https://tools.ietf.org/html/rfc8628#section-3.1
type DeviceAuthHandler struct {
	Storage interface {
		DeviceAuthStorageProvider
		oauth2.AccessTokenStorageProvider
		oauth2.RefreshTokenStorageProvider
	}
	Strategy interface {
		DeviceRateLimitStrategyProvider
		DeviceCodeStrategyProvider
		UserCodeStrategyProvider
	}
	Config interface {
		fosite.DeviceProvider
		fosite.DeviceAndUserCodeLifespanProvider
	}
}

// HandleDeviceEndpointRequest implements https://tools.ietf.org/html/rfc8628#section-3.1
func (d *DeviceAuthHandler) HandleDeviceEndpointRequest(ctx context.Context, dar fosite.DeviceRequester, resp fosite.DeviceResponder) error {
	var err error

	deviceCode, userCode, err := d.handleDeviceAuthSession(ctx, dar)
	if err != nil {
		return err
	}

	// Populate the response fields
	resp.SetDeviceCode(deviceCode)
	resp.SetUserCode(userCode)
	resp.SetVerificationURI(d.Config.GetDeviceVerificationURL(ctx))
	resp.SetVerificationURIComplete(d.Config.GetDeviceVerificationURL(ctx) + "?user_code=" + userCode)
	resp.SetExpiresIn(int64(time.Until(dar.GetSession().GetExpiresAt(fosite.UserCode)).Seconds()))
	resp.SetInterval(int(d.Config.GetDeviceAuthTokenPollingInterval(ctx).Seconds()))
	return nil
}

func (d *DeviceAuthHandler) handleDeviceAuthSession(ctx context.Context, dar fosite.DeviceRequester) (string, string, error) {
	var userCode, userCodeSignature string

	deviceCode, deviceCodeSignature, err := d.Strategy.DeviceCodeStrategy().GenerateDeviceCode(ctx)
	if err != nil {
		return "", "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	dar.GetSession().SetExpiresAt(fosite.UserCode, time.Now().UTC().Add(d.Config.GetDeviceAndUserCodeLifespan(ctx)).Round(time.Second))
	dar.GetSession().SetExpiresAt(fosite.DeviceCode, time.Now().UTC().Add(d.Config.GetDeviceAndUserCodeLifespan(ctx)).Round(time.Second))
	// Note: the retries are added here because we need to ensure uniqueness of user codes.
	// The chances of duplicates should however be diminishing, because they are the same
	// chance an attacker will be able to hit a valid code with few guesses. However, as
	// used codes will probably still be around for some time before they get cleaned,
	// the chances of hitting a duplicate here can be higher.
	// Three retries should be plenty, as otherwise the entropy is definitely off.
	for i := 0; i < MaxAttempts; i++ {
		userCode, userCodeSignature, err = d.Strategy.UserCodeStrategy().GenerateUserCode(ctx)
		if err != nil {
			return "", "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		}

		err = d.Storage.DeviceAuthStorage().CreateDeviceAuthSession(ctx, deviceCodeSignature, userCodeSignature, dar.Sanitize(nil).(fosite.DeviceRequester))
		if err == nil {
			break
		}
		if !errors.Is(err, fosite.ErrExistingUserCodeSignature) {
			return "", "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		}
	}

	if err != nil {
		errMsg := fmt.Sprintf("Exceeded user-code generation max attempts %v: %s", MaxAttempts, err.Error())
		return "", "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(errMsg))
	}
	return deviceCode, userCode, nil
}

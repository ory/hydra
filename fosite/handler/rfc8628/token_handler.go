// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/v2/fosite"
)

var _ fosite.TokenEndpointHandler = (*DeviceCodeTokenEndpointHandler)(nil)

// DeviceCodeTokenEndpointHandler is a token response handler for
// - the Authorize code grant using the explicit grant type as defined in https://tools.ietf.org/html/rfc6749#section-4.1
// - the Device Authorization Grant as defined in https://www.rfc-editor.org/rfc/rfc8628
type DeviceCodeTokenEndpointHandler struct {
	Storage interface {
		DeviceAuthStorageProvider
		oauth2.AccessTokenStorageProvider
		oauth2.RefreshTokenStorageProvider
		oauth2.TokenRevocationStorageProvider
	}
	Strategy interface {
		DeviceRateLimitStrategyProvider
		DeviceCodeStrategyProvider
		UserCodeStrategyProvider
		oauth2.AccessTokenStrategyProvider
		oauth2.RefreshTokenStrategyProvider
	}
	Config interface {
		fosite.AccessTokenLifespanProvider
		fosite.RefreshTokenLifespanProvider
		fosite.RefreshTokenScopesProvider
	}
}

func (c *DeviceCodeTokenEndpointHandler) CanSkipClientAuth(ctx context.Context, requester fosite.AccessRequester) bool {
	return false
}

func (c *DeviceCodeTokenEndpointHandler) CanHandleTokenEndpointRequest(ctx context.Context, requester fosite.AccessRequester) bool {
	return requester.GetGrantTypes().ExactOne(string(fosite.GrantTypeDeviceCode))
}

func (v DeviceCodeTokenEndpointHandler) CanHandleRequest(requester fosite.AccessRequester) bool {
	return requester.GetGrantTypes().ExactOne(string(fosite.GrantTypeDeviceCode))
}

func (c *DeviceCodeTokenEndpointHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !c.CanHandleTokenEndpointRequest(ctx, requester) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	var code, signature string
	var err error
	if code, signature, err = c.deviceCode(ctx, requester); err != nil {
		return err
	}

	var ar fosite.DeviceRequester
	if ar, err = c.session(ctx, requester, signature); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	if err = c.Strategy.DeviceCodeStrategy().ValidateDeviceCode(ctx, ar, code); err != nil {
		return errorsx.WithStack(err)
	}

	for _, scope := range ar.GetGrantedScopes() {
		requester.GrantScope(scope)
	}

	for _, audience := range ar.GetGrantedAudience() {
		requester.GrantAudience(audience)
	}

	var accessToken, accessTokenSignature string
	accessToken, accessTokenSignature, err = c.Strategy.AccessTokenStrategy().GenerateAccessToken(ctx, requester)
	if err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	var refreshToken, refreshTokenSignature string
	if c.canIssueRefreshToken(ctx, requester) {
		refreshToken, refreshTokenSignature, err = c.Strategy.RefreshTokenStrategy().GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		}
	}

	ctx, err = fosite.MaybeBeginTx(ctx, c.Storage)
	if err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}
	defer func() {
		if err != nil {
			if rollBackTxnErr := fosite.MaybeRollbackTx(ctx, c.Storage); rollBackTxnErr != nil {
				err = errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebugf("error: %s; rollback error: %s", err, rollBackTxnErr))
			}
		}
	}()

	if err = c.Storage.DeviceAuthStorage().InvalidateDeviceCodeSession(ctx, signature); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	if err = c.Storage.AccessTokenStorage().CreateAccessTokenSession(ctx, accessTokenSignature, requester.Sanitize([]string{})); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	if refreshTokenSignature != "" {
		if err = c.Storage.RefreshTokenStorage().CreateRefreshTokenSession(ctx, refreshTokenSignature, accessTokenSignature, requester.Sanitize([]string{})); err != nil {
			return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		}
	}

	lifeSpan := fosite.GetEffectiveLifespan(requester.GetClient(), c.getGrantType(requester), fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	responder.SetAccessToken(accessToken)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(getExpiresIn(requester, fosite.AccessToken, lifeSpan, time.Now().UTC()))
	responder.SetScopes(requester.GetGrantedScopes())
	if refreshToken != "" {
		responder.SetExtra("refresh_token", refreshToken)
	}

	if err = fosite.MaybeCommitTx(ctx, c.Storage); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	return nil
}

func (c *DeviceCodeTokenEndpointHandler) HandleTokenEndpointRequest(ctx context.Context, requester fosite.AccessRequester) error {
	if !c.CanHandleTokenEndpointRequest(ctx, requester) {
		return errorsx.WithStack(errorsx.WithStack(fosite.ErrUnknownRequest))
	}

	var err error
	if err = c.validateGrantTypes(requester); err != nil {
		return err
	}

	var code, signature string
	if code, signature, err = c.deviceCode(ctx, requester); err != nil {
		return err
	}

	if err = c.validateCode(ctx, requester, code); err != nil {
		return errorsx.WithStack(err)
	}

	var ar fosite.DeviceRequester
	if ar, err = c.session(ctx, requester, signature); err != nil {
		if ar != nil && (errors.Is(err, fosite.ErrInvalidatedAuthorizeCode) || errors.Is(err, fosite.ErrInvalidatedDeviceCode)) {
			return c.revokeTokens(ctx, requester.GetID())
		}

		return err
	}

	if err = c.Strategy.DeviceCodeStrategy().ValidateDeviceCode(ctx, ar, code); err != nil {
		return errorsx.WithStack(err)
	}

	// Override scopes
	requester.SetRequestedScopes(ar.GetRequestedScopes())

	// Override audiences
	requester.SetRequestedAudience(ar.GetRequestedAudience())

	// The authorization server MUST ensure that
	// the authorization code was issued to the authenticated confidential client,
	// or if the client is public, ensure that the code was issued to "client_id" in the request
	if ar.GetClient().GetID() != requester.GetClient().GetID() {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client ID from this request does not match the one from the authorize request."))
	}

	// Checking of POST client_id skipped, because
	// if the client type is confidential or the client was issued client credentials (or assigned other authentication requirements),
	// the client MUST authenticate with the authorization server as described in Section 3.2.1.
	requester.SetSession(ar.GetSession())
	requester.SetID(ar.GetID())

	atLifespan := fosite.GetEffectiveLifespan(requester.GetClient(), c.getGrantType(requester), fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	requester.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan).Round(time.Second))

	rtLifespan := fosite.GetEffectiveLifespan(requester.GetClient(), c.getGrantType(requester), fosite.RefreshToken, c.Config.GetRefreshTokenLifespan(ctx))
	if rtLifespan > -1 {
		requester.GetSession().SetExpiresAt(fosite.RefreshToken, time.Now().UTC().Add(rtLifespan).Round(time.Second))
	}

	return nil
}

func (c *DeviceCodeTokenEndpointHandler) canIssueRefreshToken(ctx context.Context, requester fosite.Requester) bool {
	scopes := c.Config.GetRefreshTokenScopes(ctx)

	// Require one of the refresh token scopes, if set.
	if len(scopes) > 0 && !requester.GetGrantedScopes().HasOneOf(scopes...) {
		return false
	}

	// Do not issue a refresh token to clients that cannot use the refresh token grant type.
	if !requester.GetClient().GetGrantTypes().Has("refresh_token") {
		return false
	}

	return true
}

func (c *DeviceCodeTokenEndpointHandler) revokeTokens(ctx context.Context, reqId string) error {
	hint := "The authorization code has already been used."
	var debug strings.Builder

	revokeAndAppendErr := func(tokenType string, revokeFunc func(context.Context, string) error) {
		if err := revokeFunc(ctx, reqId); err != nil {
			hint += fmt.Sprintf(" Additionally, an error occurred during processing the %s token revocation.", tokenType)
			debug.WriteString(fmt.Sprintf("Revocation of %s token lead to error %s.", tokenType, err.Error()))
		}
	}

	revokeAndAppendErr("access", c.Storage.TokenRevocationStorage().RevokeAccessToken)
	revokeAndAppendErr("refresh", c.Storage.TokenRevocationStorage().RevokeRefreshToken)

	return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint(hint).WithDebug(debug.String()))
}

func (c DeviceCodeTokenEndpointHandler) deviceCode(ctx context.Context, requester fosite.AccessRequester) (code string, signature string, err error) {
	code = requester.GetRequestForm().Get("device_code")

	signature, err = c.Strategy.DeviceCodeStrategy().DeviceCodeSignature(ctx, code)
	if err != nil {
		return "", "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	return
}

func (c DeviceCodeTokenEndpointHandler) validateCode(ctx context.Context, requester fosite.Requester, code string) error {
	shouldRateLimit, err := c.Strategy.DeviceRateLimitStrategy().ShouldRateLimit(ctx, code)
	if err != nil {
		return err
	}
	if shouldRateLimit {
		return errorsx.WithStack(fosite.ErrSlowDown)
	}
	return nil
}

func (s DeviceCodeTokenEndpointHandler) session(ctx context.Context, requester fosite.AccessRequester, codeSignature string) (fosite.DeviceRequester, error) {
	req, err := s.Storage.DeviceAuthStorage().GetDeviceCodeSession(ctx, codeSignature, requester.GetSession())

	if err != nil && errors.Is(err, fosite.ErrInvalidatedDeviceCode) {
		if req != nil {
			return req, err
		}

		return req, fosite.ErrServerError.
			WithHint("Misconfigured code lead to an error that prohibited the OAuth 2.0 Framework from processing this request.").
			WithDebug("\"GetDeviceCodeSession\" must return a value for \"fosite.Requester\" when returning \"ErrInvalidatedDeviceCode\".")
	}

	if err != nil && errors.Is(err, fosite.ErrNotFound) {
		return nil, errorsx.WithStack(fosite.ErrInvalidGrant.WithWrap(err).WithDebug(err.Error()))
	}

	if err != nil {
		return nil, errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	state := req.GetUserCodeState()

	if state == fosite.UserCodeUnused {
		return nil, fosite.ErrAuthorizationPending
	}
	if state == fosite.UserCodeRejected {
		return nil, fosite.ErrAccessDenied
	}

	return req, err
}

func (v DeviceCodeTokenEndpointHandler) validateGrantTypes(requester fosite.AccessRequester) error {
	if !requester.GetClient().GetGrantTypes().Has(string(fosite.GrantTypeDeviceCode)) {
		return errorsx.WithStack(fosite.ErrUnauthorizedClient.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant \"urn:ietf:params:oauth:grant-type:device_code\"."))
	}

	return nil
}

func (v DeviceCodeTokenEndpointHandler) getGrantType(requester fosite.AccessRequester) fosite.GrantType {
	return fosite.GrantTypeDeviceCode
}

func getExpiresIn(r fosite.Requester, key fosite.TokenType, defaultLifespan time.Duration, now time.Time) time.Duration {
	if r.GetSession().GetExpiresAt(key).IsZero() {
		return defaultLifespan
	}
	return time.Duration(r.GetSession().GetExpiresAt(key).UnixNano() - now.UnixNano())
}

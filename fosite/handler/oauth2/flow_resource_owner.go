// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
)

var _ fosite.TokenEndpointHandler = (*ResourceOwnerPasswordCredentialsGrantHandler)(nil)

// Deprecated: This handler is deprecated as a means to communicate that the ROPC grant type is widely discouraged and
// is at the time of this writing going to be omitted in the OAuth 2.1 spec. For more information on why this grant type
// is discouraged see: https://www.scottbrady91.com/oauth/why-the-resource-owner-password-credentials-grant-type-is-not-authentication-nor-suitable-for-modern-applications
type ResourceOwnerPasswordCredentialsGrantHandler struct {
	Storage interface {
		ResourceOwnerPasswordCredentialsGrantStorageProvider
		AccessTokenStorageProvider
		RefreshTokenStorageProvider
	}
	Strategy interface {
		AccessTokenStrategyProvider
		RefreshTokenStrategyProvider
	}
	Config interface {
		fosite.ScopeStrategyProvider
		fosite.AudienceStrategyProvider
		fosite.RefreshTokenScopesProvider
		fosite.RefreshTokenLifespanProvider
		fosite.AccessTokenLifespanProvider
	}
}

type Session interface {
	// SetSubject sets the session's subject.
	SetSubject(subject string)
}

// HandleTokenEndpointRequest implements https://tools.ietf.org/html/rfc6749#section-4.3.2
func (c *ResourceOwnerPasswordCredentialsGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !c.CanHandleTokenEndpointRequest(ctx, request) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("password") {
		return errorsx.WithStack(fosite.ErrUnauthorizedClient.WithHint("The client is not allowed to use authorization grant 'password'."))
	}

	client := request.GetClient()
	for _, scope := range request.GetRequestedScopes() {
		if !c.Config.GetScopeStrategy(ctx)(client.GetScopes(), scope) {
			return errorsx.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope '%s'.", scope))
		}
	}

	if err := c.Config.GetAudienceStrategy(ctx)(client.GetAudience(), request.GetRequestedAudience()); err != nil {
		return err
	}

	username := request.GetRequestForm().Get("username")
	password := request.GetRequestForm().Get("password")
	if username == "" || password == "" {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Username or password are missing from the POST body."))
	} else if sub, err := c.Storage.ResourceOwnerPasswordCredentialsGrantStorage().Authenticate(ctx, username, password); errors.Is(err, fosite.ErrNotFound) {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("Unable to authenticate the provided username and password credentials.").WithWrap(err).WithDebug(err.Error()))
	} else if err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	} else {
		if sess, ok := request.GetSession().(Session); ok {
			sess.SetSubject(sub)
		}
	}

	// Credentials must not be passed around, potentially leaking to the database!
	delete(request.GetRequestForm(), "password")

	atLifespan := fosite.GetEffectiveLifespan(request.GetClient(), fosite.GrantTypePassword, fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan).Round(time.Second))

	rtLifespan := fosite.GetEffectiveLifespan(request.GetClient(), fosite.GrantTypePassword, fosite.RefreshToken, c.Config.GetRefreshTokenLifespan(ctx))
	if rtLifespan > -1 {
		request.GetSession().SetExpiresAt(fosite.RefreshToken, time.Now().UTC().Add(rtLifespan).Round(time.Second))
	}

	return nil
}

// PopulateTokenEndpointResponse implements https://tools.ietf.org/html/rfc6749#section-4.3.3
func (c *ResourceOwnerPasswordCredentialsGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !c.CanHandleTokenEndpointRequest(ctx, requester) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	atLifespan := fosite.GetEffectiveLifespan(requester.GetClient(), fosite.GrantTypePassword, fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	accessTokenSignature, err := c.IssueAccessToken(ctx, atLifespan, requester, responder)
	if err != nil {
		return err
	}

	var refresh, refreshSignature string
	if len(c.Config.GetRefreshTokenScopes(ctx)) == 0 || requester.GetGrantedScopes().HasOneOf(c.Config.GetRefreshTokenScopes(ctx)...) {
		var err error
		refresh, refreshSignature, err = c.Strategy.RefreshTokenStrategy().GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		} else if err := c.Storage.RefreshTokenStorage().CreateRefreshTokenSession(ctx, refreshSignature, accessTokenSignature, requester.Sanitize([]string{})); err != nil {
			return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
		}
	}

	if refresh != "" {
		responder.SetExtra("refresh_token", refresh)
	}

	return nil
}

func (c *ResourceOwnerPasswordCredentialsGrantHandler) IssueAccessToken(ctx context.Context, atLifespan time.Duration, requester fosite.AccessRequester, responder fosite.AccessResponder) (signature string, err error) {
	token, signature, err := c.Strategy.AccessTokenStrategy().GenerateAccessToken(ctx, requester)
	if err != nil {
		return "", err
	} else if err := c.Storage.AccessTokenStorage().CreateAccessTokenSession(ctx, signature, requester.Sanitize([]string{})); err != nil {
		return "", err
	}

	if !requester.GetSession().GetExpiresAt(fosite.AccessToken).IsZero() {
		atLifespan = time.Duration(requester.GetSession().GetExpiresAt(fosite.AccessToken).UnixNano() - time.Now().UTC().UnixNano())
	}

	responder.SetAccessToken(token)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(atLifespan)
	responder.SetScopes(requester.GetGrantedScopes())

	return signature, nil
}

func (c *ResourceOwnerPasswordCredentialsGrantHandler) CanSkipClientAuth(ctx context.Context, _ fosite.AccessRequester) bool {
	return false
}

func (c *ResourceOwnerPasswordCredentialsGrantHandler) CanHandleTokenEndpointRequest(ctx context.Context, requester fosite.AccessRequester) bool {
	// grant_type REQUIRED.
	// Value MUST be set to "password".
	return requester.GetGrantTypes().ExactOne("password")
}

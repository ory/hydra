// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/v2/fosite"
)

var _ fosite.AuthorizeEndpointHandler = (*AuthorizeImplicitGrantHandler)(nil)

// AuthorizeImplicitGrantHandler is a response handler for the Authorize Code grant using the implicit grant type
// as defined in https://tools.ietf.org/html/rfc6749#section-4.2
type AuthorizeImplicitGrantHandler struct {
	Strategy AccessTokenStrategyProvider
	Storage  AccessTokenStorageProvider
	Config   interface {
		fosite.AccessTokenLifespanProvider
		fosite.ScopeStrategyProvider
		fosite.AudienceStrategyProvider
	}
}

func (c *AuthorizeImplicitGrantHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	// This let's us define multiple response types, for example open id connect's id_token
	if !ar.GetResponseTypes().ExactOne("token") {
		return nil
	}

	ar.SetDefaultResponseMode(fosite.ResponseModeFragment)

	// Disabled because this is already handled at the authorize_request_handler
	// if !ar.GetClient().GetResponseTypes().Has("token") {
	// 	 return errorsx.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type token"))
	// }

	if !ar.GetClient().GetGrantTypes().Has("implicit") {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant 'implicit'."))
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.Config.GetScopeStrategy(ctx)(client.GetScopes(), scope) {
			return errorsx.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope '%s'.", scope))
		}
	}

	if err := c.Config.GetAudienceStrategy(ctx)(client.GetAudience(), ar.GetRequestedAudience()); err != nil {
		return err
	}

	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2

	return c.IssueImplicitAccessToken(ctx, ar, resp)
}

func (c *AuthorizeImplicitGrantHandler) IssueImplicitAccessToken(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	// Only override expiry if none is set.
	atLifespan := fosite.GetEffectiveLifespan(ar.GetClient(), fosite.GrantTypeImplicit, fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	if ar.GetSession().GetExpiresAt(fosite.AccessToken).IsZero() {
		ar.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan).Round(time.Second))
	}

	// Generate the access token
	token, signature, err := c.Strategy.AccessTokenStrategy().GenerateAccessToken(ctx, ar)
	if err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	if err := c.Storage.AccessTokenStorage().CreateAccessTokenSession(ctx, signature, ar.Sanitize([]string{})); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebug(err.Error()))
	}

	resp.AddParameter("access_token", token)
	resp.AddParameter("expires_in", strconv.FormatInt(int64(getExpiresIn(ar, fosite.AccessToken, atLifespan, time.Now().UTC())/time.Second), 10))
	resp.AddParameter("token_type", "bearer")
	resp.AddParameter("state", ar.GetState())
	resp.AddParameter("scope", strings.Join(ar.GetGrantedScopes(), " "))

	ar.SetResponseTypeHandled("token")

	return nil
}

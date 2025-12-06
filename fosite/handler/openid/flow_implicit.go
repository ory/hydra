// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

type OpenIDConnectImplicitHandler struct {
	*IDTokenHandleHelper

	AuthorizeImplicitGrantTypeHandler *oauth2.AuthorizeImplicitGrantHandler
	OpenIDConnectRequestValidator     *OpenIDConnectRequestValidator
	RS256JWTStrategy                  *jwt.DefaultSigner

	Config interface {
		fosite.IDTokenLifespanProvider
		fosite.MinParameterEntropyProvider
		fosite.ScopeStrategyProvider
	}
}

func (c *OpenIDConnectImplicitHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if !(ar.GetGrantedScopes().Has("openid") && (ar.GetResponseTypes().Has("token", "id_token") || ar.GetResponseTypes().ExactOne("id_token"))) {
		return nil
	} else if ar.GetResponseTypes().Has("code") {
		// hybrid flow
		return nil
	}

	ar.SetDefaultResponseMode(fosite.ResponseModeFragment)

	if !ar.GetClient().GetGrantTypes().Has("implicit") {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant 'implicit'."))
	}

	// Disabled because this is already handled at the authorize_request_handler
	//if ar.GetResponseTypes().ExactOne("id_token") && !ar.GetClient().GetResponseTypes().Has("id_token") {
	//	return errorsx.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type id_token"))
	//} else if ar.GetResponseTypes().Matches("token", "id_token") && !ar.GetClient().GetResponseTypes().Has("token", "id_token") {
	//	return errorsx.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type token and id_token"))
	//}

	// This ensures that the 'redirect_uri' parameter is present for OpenID Connect 1.0 authorization requests as per:
	//
	// Authorization Code Flow - https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	// Implicit Flow - https://openid.net/specs/openid-connect-core-1_0.html#ImplicitAuthRequest
	// Hybrid Flow - https://openid.net/specs/openid-connect-core-1_0.html#HybridAuthRequest
	//
	// Note: as per the Hybrid Flow documentation the Hybrid Flow has the same requirements as the Authorization Code Flow.
	rawRedirectURI := ar.GetRequestForm().Get("redirect_uri")
	if len(rawRedirectURI) == 0 {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("The 'redirect_uri' parameter is required when using OpenID Connect 1.0."))
	}

	if nonce := ar.GetRequestForm().Get("nonce"); len(nonce) == 0 {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Parameter 'nonce' must be set when using the OpenID Connect Implicit Flow."))
	} else if len(nonce) < c.Config.GetMinParameterEntropy(ctx) {
		return errorsx.WithStack(fosite.ErrInsufficientEntropy.WithHintf("Parameter 'nonce' is set but does not satisfy the minimum entropy of %d characters.", c.Config.GetMinParameterEntropy(ctx)))
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.Config.GetScopeStrategy(ctx)(client.GetScopes(), scope) {
			return errorsx.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope '%s'.", scope))
		}
	}

	sess, ok := ar.GetSession().(Session)
	if !ok {
		return errorsx.WithStack(ErrInvalidSession)
	}

	if err := c.OpenIDConnectRequestValidator.ValidatePrompt(ctx, ar); err != nil {
		return err
	}

	claims := sess.IDTokenClaims()
	if ar.GetResponseTypes().Has("token") {
		if err := c.AuthorizeImplicitGrantTypeHandler.IssueImplicitAccessToken(ctx, ar, resp); err != nil {
			return errorsx.WithStack(err)
		}

		ar.SetResponseTypeHandled("token")
		hash, err := c.ComputeHash(ctx, sess, resp.GetParameters().Get("access_token"))
		if err != nil {
			return err
		}

		claims.AccessTokenHash = hash
	} else {
		resp.AddParameter("state", ar.GetState())
	}

	idTokenLifespan := fosite.GetEffectiveLifespan(ar.GetClient(), fosite.GrantTypeImplicit, fosite.IDToken, c.Config.GetIDTokenLifespan(ctx))
	if err := c.IssueImplicitIDToken(ctx, idTokenLifespan, ar, resp); err != nil {
		return errorsx.WithStack(err)
	}

	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2

	ar.SetResponseTypeHandled("id_token")
	return nil
}

// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
)

type OpenIDConnectRefreshHandler struct {
	*IDTokenHandleHelper

	Config interface {
		fosite.IDTokenLifespanProvider
	}
}

func (c *OpenIDConnectRefreshHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !c.CanHandleTokenEndpointRequest(ctx, request) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetGrantedScopes().Has("openid") {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("refresh_token") {
		return errorsx.WithStack(fosite.ErrUnauthorizedClient.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"refresh_token\"."))
	}

	// Refresh tokens can only be issued by an authorize_code which in turn disables the need to check if the id_token
	// response type is enabled by the client.
	//
	// if !request.GetClient().GetResponseTypes().Has("id_token") {
	// 	return errorsx.WithStack(fosite.ErrUnknownRequest.WithDebug("The client is not allowed to use response type id_token"))
	// }

	sess, ok := request.GetSession().(Session)
	if !ok {
		return errors.New("Failed to generate id token because session must be of type fosite/handler/openid.Session")
	}

	// We need to reset the expires at value as this would be the previous expiry.
	sess.IDTokenClaims().ExpiresAt = time.Time{}

	// These will be recomputed in PopulateTokenEndpointResponse
	sess.IDTokenClaims().JTI = ""
	sess.IDTokenClaims().AccessTokenHash = ""

	// We are not issuing a code so there is no need for this field.
	sess.IDTokenClaims().CodeHash = ""

	return nil
}

func (c *OpenIDConnectRefreshHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !c.CanHandleTokenEndpointRequest(ctx, requester) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	if !requester.GetGrantedScopes().Has("openid") {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	if !requester.GetClient().GetGrantTypes().Has("refresh_token") {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"refresh_token\"."))
	}

	// Disabled because this is already handled at the authorize_request_handler
	// if !requester.GetClient().GetResponseTypes().Has("id_token") {
	// 	 return errorsx.WithStack(fosite.ErrUnknownRequest.WithDebug("The client is not allowed to use response type id_token"))
	// }

	sess, ok := requester.GetSession().(Session)
	if !ok {
		return errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because session must be of type fosite/handler/openid.Session."))
	}

	claims := sess.IDTokenClaims()
	if claims.Subject == "" {
		return errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because subject is an empty string."))
	}

	claims.AccessTokenHash = c.GetAccessTokenHash(ctx, requester, responder)
	claims.JTI = uuid.New().String()
	claims.CodeHash = ""
	claims.IssuedAt = time.Now().Truncate(time.Second)

	idTokenLifespan := fosite.GetEffectiveLifespan(requester.GetClient(), fosite.GrantTypeRefreshToken, fosite.IDToken, c.Config.GetIDTokenLifespan(ctx))
	return c.IssueExplicitIDToken(ctx, idTokenLifespan, requester, responder)
}

func (c *OpenIDConnectRefreshHandler) CanSkipClientAuth(ctx context.Context, requester fosite.AccessRequester) bool {
	return false
}

func (c *OpenIDConnectRefreshHandler) CanHandleTokenEndpointRequest(ctx context.Context, requester fosite.AccessRequester) bool {
	// grant_type REQUIRED.
	// Value MUST be set to "refresh_token"
	return requester.GetGrantTypes().ExactOne("refresh_token")
}

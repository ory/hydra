/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package oauth2

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

// HandleTokenEndpointRequest implements
// * https://tools.ietf.org/html/rfc6749#section-4.1.3 (everything)
func (c *AuthorizeExplicitGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code".
	if !request.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(errors.WithStack(fosite.ErrUnknownRequest))
	}

	if !request.GetClient().GetGrantTypes().Has("authorization_code") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant \"authorization_code\"."))
	}

	code := request.GetRequestForm().Get("code")
	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	authorizeRequest, err := c.CoreStorage.GetAuthorizeCodeSession(ctx, signature, request.GetSession())
	if errors.Cause(err) == fosite.ErrInvalidatedAuthorizeCode {
		if authorizeRequest == nil {
			return fosite.ErrServerError.
				WithHint("Misconfigured code lead to an error that prohibited the OAuth 2.0 Framework from processing this request.").
				WithDebug("GetAuthorizeCodeSession must return a value for \"fosite.Requester\" when returning \"ErrInvalidatedAuthorizeCode\".")
		}

		//If an authorize code is used twice, we revoke all refresh and access tokens associated with this request.
		reqID := authorizeRequest.GetID()
		hint := "The authorization code has already been used."
		debug := ""
		if revErr := c.TokenRevocationStorage.RevokeAccessToken(ctx, reqID); revErr != nil {
			hint += " Additionally, an error occurred during processing the access token revocation."
			debug += "Revokation of access_token lead to error " + revErr.Error() + "."
		}
		if revErr := c.TokenRevocationStorage.RevokeRefreshToken(ctx, reqID); revErr != nil {
			hint += " Additionally, an error occurred during processing the refresh token revocation."
			debug += "Revokation of refresh_token lead to error " + revErr.Error() + "."
		}
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint(hint).WithDebug(debug))
	} else if err != nil && errors.Cause(err).Error() == fosite.ErrNotFound.Error() {
		return errors.WithStack(fosite.ErrInvalidGrant.WithDebug(err.Error()))
	} else if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	// The authorization server MUST verify that the authorization code is valid
	// This needs to happen after store retrieval for the session to be hydrated properly
	if err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, request, code); err != nil {
		return errors.WithStack(fosite.ErrInvalidGrant.WithDebug(err.Error()))
	}

	// Override scopes
	request.SetRequestedScopes(authorizeRequest.GetRequestedScopes())

	// The authorization server MUST ensure that the authorization code was issued to the authenticated
	// confidential client, or if the client is public, ensure that the
	// code was issued to "client_id" in the request,
	if authorizeRequest.GetClient().GetID() != request.GetClient().GetID() {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("The OAuth 2.0 Client ID from this request does not match the one from the authorize request."))
	}

	// ensure that the "redirect_uri" parameter is present if the
	// "redirect_uri" parameter was included in the initial authorization
	// request as described in Section 4.1.1, and if included ensure that
	// their values are identical.
	forcedRedirectURI := authorizeRequest.GetRequestForm().Get("redirect_uri")
	if forcedRedirectURI != "" && forcedRedirectURI != request.GetRequestForm().Get("redirect_uri") {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("The \"redirect_uri\" from this request does not match the one from the authorize request."))
	}

	// Checking of POST client_id skipped, because:
	// If the client type is confidential or the client was issued client
	// credentials (or assigned other authentication requirements), the
	// client MUST authenticate with the authorization server as described
	// in Section 3.2.1.
	request.SetSession(authorizeRequest.GetSession())
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(c.AccessTokenLifespan))
	request.SetID(authorizeRequest.GetID())
	return nil
}

func (c *AuthorizeExplicitGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code", as this is the explicit grant handler.
	if !requester.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	code := requester.GetRequestForm().Get("code")
	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	authorizeRequest, err := c.CoreStorage.GetAuthorizeCodeSession(ctx, signature, requester.GetSession())
	if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	} else if err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, requester, code); err != nil {
		// This needs to happen after store retrieval for the session to be hydrated properly
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	}

	for _, scope := range authorizeRequest.GetGrantedScopes() {
		requester.GrantScope(scope)
	}

	access, accessSignature, err := c.AccessTokenStrategy.GenerateAccessToken(ctx, requester)
	if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	var refresh, refreshSignature string
	if authorizeRequest.GetGrantedScopes().HasOneOf("offline", "offline_access") {
		refresh, refreshSignature, err = c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}
	}

	if err := c.CoreStorage.InvalidateAuthorizeCodeSession(ctx, signature); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	} else if err := c.CoreStorage.CreateAccessTokenSession(ctx, accessSignature, requester.Sanitize([]string{})); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	} else if refreshSignature != "" {
		if err := c.CoreStorage.CreateRefreshTokenSession(ctx, refreshSignature, requester.Sanitize([]string{})); err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}
	}

	responder.SetAccessToken(access)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(getExpiresIn(requester, fosite.AccessToken, c.AccessTokenLifespan, time.Now().UTC()))
	responder.SetScopes(requester.GetGrantedScopes())
	if refresh != "" {
		responder.SetExtra("refresh_token", refresh)
	}

	return nil
}

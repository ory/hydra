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

package openid

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
)

type OpenIDConnectImplicitHandler struct {
	AuthorizeImplicitGrantTypeHandler *oauth2.AuthorizeImplicitGrantTypeHandler
	*IDTokenHandleHelper
	ScopeStrategy                 fosite.ScopeStrategy
	OpenIDConnectRequestValidator *OpenIDConnectRequestValidator

	RS256JWTStrategy *jwt.RS256JWTStrategy
}

func (c *OpenIDConnectImplicitHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if !(ar.GetGrantedScopes().Has("openid") && (ar.GetResponseTypes().Has("token", "id_token") || ar.GetResponseTypes().Exact("id_token"))) {
		return nil
	} else if ar.GetResponseTypes().Has("code") {
		// hybrid flow
		return nil
	}

	if !ar.GetClient().GetGrantTypes().Has("implicit") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"implicit\"."))
	}

	// Disabled because this is already handled at the authorize_request_handler
	//if ar.GetResponseTypes().Exact("id_token") && !ar.GetClient().GetResponseTypes().Has("id_token") {
	//	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type id_token"))
	//} else if ar.GetResponseTypes().Matches("token", "id_token") && !ar.GetClient().GetResponseTypes().Has("token", "id_token") {
	//	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type token and id_token"))
	//}

	if nonce := ar.GetRequestForm().Get("nonce"); len(nonce) == 0 {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Parameter \"nonce\" must be set when using the OpenID Connect Hybrid Flow."))
	} else if len(nonce) < fosite.MinParameterEntropy {
		return errors.WithStack(fosite.ErrInsufficientEntropy.WithHintf("Parameter \"nonce\" is set but does not satisfy the minimum entropy of %d characters.", fosite.MinParameterEntropy))
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope \"%s\".", scope))
		}
	}

	sess, ok := ar.GetSession().(Session)
	if !ok {
		return errors.WithStack(ErrInvalidSession)
	}

	if err := c.OpenIDConnectRequestValidator.ValidatePrompt(ctx, ar); err != nil {
		return err
	}

	claims := sess.IDTokenClaims()
	if ar.GetResponseTypes().Has("token") {
		if err := c.AuthorizeImplicitGrantTypeHandler.IssueImplicitAccessToken(ctx, ar, resp); err != nil {
			return errors.WithStack(err)
		}

		ar.SetResponseTypeHandled("token")
		hash, err := c.RS256JWTStrategy.Hash(ctx, []byte(resp.GetFragment().Get("access_token")))
		if err != nil {
			return err
		}

		claims.AccessTokenHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.RS256JWTStrategy.GetSigningMethodLength()/2]))
	} else {
		resp.AddFragment("state", ar.GetState())
	}

	if err := c.IssueImplicitIDToken(ctx, ar, resp); err != nil {
		return errors.WithStack(err)
	}

	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2

	ar.SetResponseTypeHandled("id_token")
	return nil
}

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

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
	"github.com/pkg/errors"
)

type OpenIDConnectHybridHandler struct {
	AuthorizeImplicitGrantTypeHandler *oauth2.AuthorizeImplicitGrantTypeHandler
	AuthorizeExplicitGrantHandler     *oauth2.AuthorizeExplicitGrantHandler
	IDTokenHandleHelper               *IDTokenHandleHelper
	ScopeStrategy                     fosite.ScopeStrategy
	OpenIDConnectRequestValidator     *OpenIDConnectRequestValidator
	OpenIDConnectRequestStorage       OpenIDConnectRequestStorage

	Enigma *jwt.RS256JWTStrategy
}

func (c *OpenIDConnectHybridHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if len(ar.GetResponseTypes()) < 2 {
		return nil
	}

	if !(ar.GetResponseTypes().Matches("token", "id_token", "code") || ar.GetResponseTypes().Matches("token", "code") || ar.GetResponseTypes().Matches("id_token", "code")) {
		return nil
	}

	// Disabled because this is already handled at the authorize_request_handler
	//if ar.GetResponseTypes().Matches("token") && !ar.GetClient().GetResponseTypes().Has("token") {
	//	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use the token response type"))
	//} else if ar.GetResponseTypes().Matches("code") && !ar.GetClient().GetResponseTypes().Has("code") {
	//	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use the code response type"))
	//} else if ar.GetResponseTypes().Matches("id_token") && !ar.GetClient().GetResponseTypes().Has("id_token") {
	//	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use the id_token response type"))
	//}

	if nonce := ar.GetRequestForm().Get("nonce"); len(nonce) == 0 {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Parameter \"nonce\" must be set when using the OpenID Connect Hybrid Flow."))
	} else if len(nonce) < fosite.MinParameterEntropy {
		return errors.WithStack(fosite.ErrInsufficientEntropy.WithHintf("Parameter \"nonce\" is set but does not satisfy the minimum entropy of %d characters.", fosite.MinParameterEntropy))
	}

	sess, ok := ar.GetSession().(Session)
	if !ok {
		return errors.WithStack(ErrInvalidSession)
	}

	if err := c.OpenIDConnectRequestValidator.ValidatePrompt(ctx, ar); err != nil {
		return err
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope \"%s\".", scope))
		}
	}

	claims := sess.IDTokenClaims()
	if ar.GetResponseTypes().Has("code") {
		if !ar.GetClient().GetGrantTypes().Has("authorization_code") {
			return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant \"authorization_code\"."))
		}

		code, signature, err := c.AuthorizeExplicitGrantHandler.AuthorizeCodeStrategy.GenerateAuthorizeCode(ctx, ar)
		if err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		} else if err := c.AuthorizeExplicitGrantHandler.CoreStorage.CreateAuthorizeCodeSession(ctx, signature, ar.Sanitize(c.AuthorizeExplicitGrantHandler.GetSanitationWhiteList())); err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}

		resp.AddFragment("code", code)
		ar.SetResponseTypeHandled("code")

		hash, err := c.Enigma.Hash(ctx, []byte(resp.GetFragment().Get("code")))
		if err != nil {
			return err
		}
		claims.CodeHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.Enigma.GetSigningMethodLength()/2]))

		if ar.GetGrantedScopes().Has("openid") {
			if err := c.OpenIDConnectRequestStorage.CreateOpenIDConnectSession(ctx, resp.GetCode(), ar.Sanitize(oidcParameters)); err != nil {
				return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
			}
		}
	}

	if ar.GetResponseTypes().Has("token") {
		if !ar.GetClient().GetGrantTypes().Has("implicit") {
			return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"implicit\"."))
		} else if err := c.AuthorizeImplicitGrantTypeHandler.IssueImplicitAccessToken(ctx, ar, resp); err != nil {
			return errors.WithStack(err)
		}
		ar.SetResponseTypeHandled("token")

		hash, err := c.Enigma.Hash(ctx, []byte(resp.GetFragment().Get("access_token")))
		if err != nil {
			return err
		}
		claims.AccessTokenHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.Enigma.GetSigningMethodLength()/2]))
	}

	if resp.GetFragment().Get("state") == "" {
		resp.AddFragment("state", ar.GetState())
	}

	if !ar.GetGrantedScopes().Has("openid") || !ar.GetResponseTypes().Has("id_token") {
		ar.SetResponseTypeHandled("id_token")
		return nil
	}

	if err := c.IDTokenHandleHelper.IssueImplicitIDToken(ctx, ar, resp); err != nil {
		return errors.WithStack(err)
	}

	ar.SetResponseTypeHandled("id_token")
	return nil
	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2
}

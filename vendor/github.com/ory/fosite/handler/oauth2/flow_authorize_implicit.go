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
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

// AuthorizeImplicitGrantTypeHandler is a response handler for the Authorize Code grant using the implicit grant type
// as defined in https://tools.ietf.org/html/rfc6749#section-4.2
type AuthorizeImplicitGrantTypeHandler struct {
	AccessTokenStrategy AccessTokenStrategy

	// AccessTokenStorage is used to persist session data across requests.
	AccessTokenStorage AccessTokenStorage

	// AccessTokenLifespan defines the lifetime of an access token.
	AccessTokenLifespan time.Duration

	ScopeStrategy fosite.ScopeStrategy
}

func (c *AuthorizeImplicitGrantTypeHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	// This let's us define multiple response types, for example open id connect's id_token
	if !ar.GetResponseTypes().Exact("token") {
		return nil
	}

	// Disabled because this is already handled at the authorize_request_handler
	// if !ar.GetClient().GetResponseTypes().Has("token") {
	// 	 return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type token"))
	// }

	if !ar.GetClient().GetGrantTypes().Has("implicit") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"implicit\"."))
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope \"%s\".", scope))
		}
	}

	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2

	return c.IssueImplicitAccessToken(ctx, ar, resp)
}

func (c *AuthorizeImplicitGrantTypeHandler) IssueImplicitAccessToken(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	ar.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(c.AccessTokenLifespan))

	// Generate the code
	token, signature, err := c.AccessTokenStrategy.GenerateAccessToken(ctx, ar)
	if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	if err := c.AccessTokenStorage.CreateAccessTokenSession(ctx, signature, ar.Sanitize([]string{})); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	resp.AddFragment("access_token", token)
	resp.AddFragment("expires_in", strconv.FormatInt(int64(getExpiresIn(ar, fosite.AccessToken, c.AccessTokenLifespan, time.Now().UTC())/time.Second), 10))
	resp.AddFragment("token_type", "bearer")
	resp.AddFragment("state", ar.GetState())
	resp.AddFragment("scope", strings.Join(ar.GetGrantedScopes(), " "))
	ar.SetResponseTypeHandled("token")

	return nil
}

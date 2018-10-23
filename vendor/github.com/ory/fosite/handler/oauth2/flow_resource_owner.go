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

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

type ResourceOwnerPasswordCredentialsGrantHandler struct {
	// ResourceOwnerPasswordCredentialsGrantStorage is used to persist session data across requests.
	ResourceOwnerPasswordCredentialsGrantStorage ResourceOwnerPasswordCredentialsGrantStorage

	RefreshTokenStrategy RefreshTokenStrategy
	ScopeStrategy        fosite.ScopeStrategy

	*HandleHelper
}

// HandleTokenEndpointRequest implements https://tools.ietf.org/html/rfc6749#section-4.3.2
func (c *ResourceOwnerPasswordCredentialsGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "password".
	if !request.GetGrantTypes().Exact("password") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("password") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The client is not allowed to use authorization grant \"password\"."))
	}

	username := request.GetRequestForm().Get("username")
	password := request.GetRequestForm().Get("password")
	if username == "" || password == "" {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Username or password are missing from the POST body."))
	} else if err := c.ResourceOwnerPasswordCredentialsGrantStorage.Authenticate(ctx, username, password); errors.Cause(err) == fosite.ErrNotFound {
		return errors.WithStack(fosite.ErrRequestUnauthorized.WithHint("Unable to authenticate the provided username and password credentials.").WithDebug(err.Error()))
	} else if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	client := request.GetClient()
	for _, scope := range request.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope \"%s\".", scope))
		}
	}

	// Credentials must not be passed around, potentially leaking to the database!
	delete(request.GetRequestForm(), "password")

	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(c.AccessTokenLifespan))
	return nil
}

// PopulateTokenEndpointResponse implements https://tools.ietf.org/html/rfc6749#section-4.3.3
func (c *ResourceOwnerPasswordCredentialsGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("password") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	var refresh, refreshSignature string
	if requester.GetGrantedScopes().HasOneOf("offline", "offline_access") {
		var err error
		refresh, refreshSignature, err = c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		} else if err := c.ResourceOwnerPasswordCredentialsGrantStorage.CreateRefreshTokenSession(ctx, refreshSignature, requester.Sanitize([]string{})); err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}
	}

	if err := c.IssueAccessToken(ctx, requester, responder); err != nil {
		return err
	}

	if refresh != "" {
		responder.SetExtra("refresh_token", refresh)
	}

	return nil
}

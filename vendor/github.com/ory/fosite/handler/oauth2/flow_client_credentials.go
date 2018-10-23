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

type ClientCredentialsGrantHandler struct {
	*HandleHelper
	ScopeStrategy fosite.ScopeStrategy
}

// IntrospectTokenEndpointRequest implements https://tools.ietf.org/html/rfc6749#section-4.4.2
func (c *ClientCredentialsGrantHandler) HandleTokenEndpointRequest(_ context.Context, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "client_credentials".
	if !request.GetGrantTypes().Exact("client_credentials") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	client := request.GetClient()
	for _, scope := range request.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope \"%s\".", scope))
		}
	}

	// The client MUST authenticate with the authorization server as described in Section 3.2.1.
	// This requirement is already fulfilled because fosite requires all token requests to be authenticated as described
	// in https://tools.ietf.org/html/rfc6749#section-3.2.1
	if client.IsPublic() {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is marked as public and is thus not allowed to use authorization grant \"client_credentials\"."))
	}
	// if the client is not public, he has already been authenticated by the access request handler.

	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(c.AccessTokenLifespan))
	return nil
}

// PopulateTokenEndpointResponse implements https://tools.ietf.org/html/rfc6749#section-4.4.3
func (c *ClientCredentialsGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, request fosite.AccessRequester, response fosite.AccessResponder) error {
	if !request.GetGrantTypes().Exact("client_credentials") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("client_credentials") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant \"client_credentials\"."))
	}

	return c.IssueAccessToken(ctx, request, response)
}

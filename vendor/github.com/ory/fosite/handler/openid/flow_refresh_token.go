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
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

type OpenIDConnectRefreshHandler struct {
	*IDTokenHandleHelper
}

func (c *OpenIDConnectRefreshHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !request.GetGrantTypes().Exact("refresh_token") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetGrantedScopes().Has("openid") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("refresh_token") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"refresh_token\"."))
	}

	// Refresh tokens can only be issued by an authorize_code which in turn disables the need to check if the id_token
	// response type is enabled by the client.
	//
	// if !request.GetClient().GetResponseTypes().Has("id_token") {
	// 	return errors.WithStack(fosite.ErrUnknownRequest.WithDebug("The client is not allowed to use response type id_token"))
	// }

	sess, ok := request.GetSession().(Session)
	if !ok {
		return errors.New("Failed to generate id token because session must be of type fosite/handler/openid.Session")
	}

	// We need to reset the expires at value
	sess.IDTokenClaims().ExpiresAt = time.Time{}
	sess.IDTokenClaims().Nonce = ""
	return nil
}

func (c *OpenIDConnectRefreshHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("refresh_token") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !requester.GetGrantedScopes().Has("openid") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !requester.GetClient().GetGrantTypes().Has("refresh_token") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"refresh_token\"."))
	}

	// Disabled because this is already handled at the authorize_request_handler
	// if !requester.GetClient().GetResponseTypes().Has("id_token") {
	// 	 return errors.WithStack(fosite.ErrUnknownRequest.WithDebug("The client is not allowed to use response type id_token"))
	// }

	return c.IssueExplicitIDToken(ctx, requester, responder)
}

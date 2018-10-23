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

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

func (c *OpenIDConnectExplicitHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	return errors.WithStack(fosite.ErrUnknownRequest)
}

func (c *OpenIDConnectExplicitHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	authorize, err := c.OpenIDConnectRequestStorage.GetOpenIDConnectSession(ctx, requester.GetRequestForm().Get("code"), requester)
	if errors.Cause(err) == ErrNoSessionFound {
		return errors.WithStack(fosite.ErrUnknownRequest.WithDebug(err.Error()))
	} else if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	if !authorize.GetGrantedScopes().Has("openid") {
		return errors.WithStack(fosite.ErrMisconfiguration.WithDebug("An OpenID Connect session was found but the openid scope is missing, probably due to a broken code configuration."))
	}

	if !requester.GetClient().GetGrantTypes().Has("authorization_code") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client is not allowed to use the authorization grant \"authorization_code\"."))
	}

	// The response type `id_token` is only required when performing the implicit or hybrid flow, see:
	// https://openid.net/specs/openid-connect-registration-1_0.html
	//
	// if !requester.GetClient().GetResponseTypes().Has("id_token") {
	// 	return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use response type id_token"))
	// }

	return c.IssueExplicitIDToken(ctx, authorize, responder)
}

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

type OpenIDConnectExplicitHandler struct {
	// OpenIDConnectRequestStorage is the storage for open id connect sessions.
	OpenIDConnectRequestStorage   OpenIDConnectRequestStorage
	OpenIDConnectRequestValidator *OpenIDConnectRequestValidator

	*IDTokenHandleHelper
}

var oidcParameters = []string{"grant_type",
	"max_age",
	"prompt",
	"acr_values",
	"id_token_hint",
	"nonce",
}

func (c *OpenIDConnectExplicitHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if !(ar.GetGrantedScopes().Has("openid") && ar.GetResponseTypes().Exact("code")) {
		return nil
	}

	//if !ar.GetClient().GetResponseTypes().Has("id_token", "code") {
	//	return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("The client is not allowed to use response type id_token and code"))
	//}

	if len(resp.GetCode()) == 0 {
		return errors.WithStack(fosite.ErrMisconfiguration.WithDebug("The authorization code has not been issued yet, indicating a broken code configuration."))
	}

	if err := c.OpenIDConnectRequestValidator.ValidatePrompt(ctx, ar); err != nil {
		return err
	}

	if err := c.OpenIDConnectRequestStorage.CreateOpenIDConnectSession(ctx, resp.GetCode(), ar.Sanitize(oidcParameters)); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	// there is no need to check for https, because it has already been checked by core.explicit

	return nil
}

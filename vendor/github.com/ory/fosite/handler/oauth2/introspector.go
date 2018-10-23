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

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

type CoreValidator struct {
	CoreStrategy
	CoreStorage
	ScopeStrategy                 fosite.ScopeStrategy
	DisableRefreshTokenValidation bool
}

func (c *CoreValidator) IntrospectToken(ctx context.Context, token string, tokenType fosite.TokenType, accessRequest fosite.AccessRequester, scopes []string) (fosite.TokenType, error) {
	if c.DisableRefreshTokenValidation {
		if err := c.introspectAccessToken(ctx, token, accessRequest, scopes); err != nil {
			return "", err
		}
		return fosite.AccessToken, nil
	}

	var err error
	switch tokenType {
	case fosite.RefreshToken:
		if err = c.introspectRefreshToken(ctx, token, accessRequest, scopes); err == nil {
			return fosite.RefreshToken, nil
		} else if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
			return fosite.AccessToken, nil
		}
		return "", err
	}

	if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
		return fosite.AccessToken, nil
	} else if err := c.introspectRefreshToken(ctx, token, accessRequest, scopes); err == nil {
		return fosite.RefreshToken, nil
	}

	return "", err
}

func matchScopes(ss fosite.ScopeStrategy, granted, scopes []string) error {
	for _, scope := range scopes {
		if scope == "" {
			continue
		}

		if !ss(granted, scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithHintf("The request scope \"%s\" has not been granted or is not allowed to be requested.", scope))
		}
	}

	return nil
}

func (c *CoreValidator) introspectAccessToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.CoreStrategy.AccessTokenSignature(token)
	or, err := c.CoreStorage.GetAccessTokenSession(ctx, sig, accessRequest.GetSession())
	if err != nil {
		return errors.WithStack(fosite.ErrRequestUnauthorized.WithDebug(err.Error()))
	} else if err := c.CoreStrategy.ValidateAccessToken(ctx, or, token); err != nil {
		return err
	}

	if err := matchScopes(c.ScopeStrategy, or.GetGrantedScopes(), scopes); err != nil {
		return err
	}

	accessRequest.Merge(or)
	return nil
}

func (c *CoreValidator) introspectRefreshToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.CoreStrategy.RefreshTokenSignature(token)
	or, err := c.CoreStorage.GetRefreshTokenSession(ctx, sig, accessRequest.GetSession())

	if err != nil {
		return errors.WithStack(fosite.ErrRequestUnauthorized.WithDebug(err.Error()))
	} else if err := c.CoreStrategy.ValidateRefreshToken(ctx, or, token); err != nil {
		return err
	}

	if err := matchScopes(c.ScopeStrategy, or.GetGrantedScopes(), scopes); err != nil {
		return err
	}

	accessRequest.Merge(or)
	return nil
}

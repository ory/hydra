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

type JWTAccessTokenStrategy interface {
	AccessTokenStrategy
	JWTStrategy
}

type StatelessJWTValidator struct {
	JWTAccessTokenStrategy
	ScopeStrategy fosite.ScopeStrategy
}

func (v *StatelessJWTValidator) IntrospectToken(ctx context.Context, token string, tokenType fosite.TokenType, accessRequest fosite.AccessRequester, scopes []string) (fosite.TokenType, error) {
	or, err := v.JWTAccessTokenStrategy.ValidateJWT(ctx, fosite.AccessToken, token)
	if err != nil {
		return "", err
	}

	for _, scope := range scopes {
		if scope == "" {
			continue
		}

		if !v.ScopeStrategy(or.GetGrantedScopes(), scope) {
			return "", errors.WithStack(fosite.ErrInvalidScope)
		}
	}

	accessRequest.Merge(or)
	return "", nil
}

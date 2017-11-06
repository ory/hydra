// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warden

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/pkg/errors"
)

type TokenValidator struct {
	oauth2.CoreStrategy
	oauth2.CoreStorage
	ScopeStrategy fosite.ScopeStrategy
}

func (c *TokenValidator) IntrospectToken(ctx context.Context, token string, tokenType fosite.TokenType, accessRequest fosite.AccessRequester, scopes []string) (err error) {
	if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
		return nil
	}

	return err
}

func matchScopes(ss fosite.ScopeStrategy, granted, scopes []string) error {
	for _, scope := range scopes {
		if scope == "" {
			continue
		}

		if !ss(granted, scope) {
			return errors.Wrapf(fosite.ErrInvalidScope, "Scope %s was not granted", scope)
		}
	}

	return nil
}

func (c *TokenValidator) introspectAccessToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.CoreStrategy.AccessTokenSignature(token)
	or, err := c.CoreStorage.GetAccessTokenSession(ctx, sig, accessRequest.GetSession())
	if err != nil {
		return errors.Wrap(fosite.ErrRequestUnauthorized, err.Error())
	} else if err := c.CoreStrategy.ValidateAccessToken(ctx, or, token); err != nil {
		return err
	}

	if err := matchScopes(c.ScopeStrategy, or.GetGrantedScopes(), scopes); err != nil {
		return err
	}

	accessRequest.Merge(or)
	return nil
}

func OAuth2TokenIntrospectionFactory(config *compose.Config, storage interface{}, strategy interface{}) interface{} {
	return &TokenValidator{
		CoreStrategy:  strategy.(oauth2.CoreStrategy),
		CoreStorage:   storage.(oauth2.CoreStorage),
		ScopeStrategy: config.GetScopeStrategy(),
	}
}

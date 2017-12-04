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

package oauth2

import (
	"context"
	"testing"

	"net/url"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultRequest = fosite.Request{
	RequestedAt:   time.Now().Round(time.Second),
	Client:        &client.Client{ID: "foobar"},
	Scopes:        fosite.Arguments{"fa", "ba"},
	GrantedScopes: fosite.Arguments{"fa", "ba"},
	Form:          url.Values{"foo": []string{"bar", "baz"}},
	Session:       &fosite.DefaultSession{Subject: "bar"},
}

func TestHelperCreateGetDeleteOpenIDConnectSession(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		assert.NotNil(t, err)

		err = m.CreateOpenIDConnectSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{Session: &fosite.DefaultSession{}})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteOpenIDConnectSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		assert.NotNil(t, err)
	}
}

func TestHelperCreateGetDeleteRefreshTokenSession(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateRefreshTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteRefreshTokenSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}

}
func TestHelperRevokeRefreshToken(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New()
		_, err := m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1111", &fosite.Request{ID: id, Client: &client.Client{ID: "foobar"}, RequestedAt: time.Now().Round(time.Second), Session: &fosite.DefaultSession{}})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", &fosite.Request{ID: id, Client: &client.Client{ID: "foobar"}, RequestedAt: time.Now().Round(time.Second), Session: &fosite.DefaultSession{}})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, id)
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1122", &fosite.DefaultSession{})
		assert.NotNil(t, err)

	}

}
func TestHelperCreateGetDeleteAuthorizeCodes(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateAuthorizeCodeSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAuthorizeCodeSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}
}

func TestHelperCreateGetDeleteAccessTokenSession(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}
}

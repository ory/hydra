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
 */

package oauth2

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/sqlcon"
)

var defaultRequest = fosite.Request{
	RequestedAt:   time.Now().UTC().Round(time.Second),
	Client:        &client.Client{ClientID: "foobar"},
	Scopes:        fosite.Arguments{"fa", "ba"},
	GrantedScopes: fosite.Arguments{"fa", "ba"},
	Form:          url.Values{"foo": []string{"bar", "baz"}},
	Session:       &fosite.DefaultSession{Subject: "bar"},
}

func TestHelperUniqueConstraints(m pkg.FositeStorer, storageType string) func(t *testing.T) {
	return func(t *testing.T) {
		dbErrorIsConstraintError := func(dbErr error) {
			assert.Error(t, dbErr)
			switch err := errors.Cause(dbErr).(type) {
			case *herodot.DefaultError:
				assert.Equal(t, sqlcon.ErrUniqueViolation, err)
			default:
				t.Errorf("unexpected error type %s", err)
			}
		}

		requestId := uuid.New()
		signatureOne := uuid.New()
		signatureTwo := uuid.New()
		fositeRequest := &fosite.Request{
			ID:          requestId,
			Client:      &client.Client{ClientID: "foobar"},
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     &fosite.DefaultSession{},
		}

		err := m.CreateRefreshTokenSession(context.TODO(), signatureOne, fositeRequest)
		assert.NoError(t, err)
		err = m.CreateAccessTokenSession(context.TODO(), signatureOne, fositeRequest)
		assert.NoError(t, err)

		// attempting to insert new records with the SAME requestID should fail as there is a unique index
		// on the request_id column

		err = m.CreateRefreshTokenSession(context.TODO(), signatureTwo, fositeRequest)
		dbErrorIsConstraintError(err)
		err = m.CreateAccessTokenSession(context.TODO(), signatureTwo, fositeRequest)
		dbErrorIsConstraintError(err)
	}
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
		_, err := m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
		assert.Error(t, err)

		reqIdOne := uuid.New()
		reqIdTwo := uuid.New()

		err = m.CreateRefreshTokenSession(ctx, "1111", &fosite.Request{ID: reqIdOne, Client: &client.Client{ClientID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &fosite.DefaultSession{}})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", &fosite.Request{ID: reqIdTwo, Client: &client.Client{ClientID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &fosite.DefaultSession{}})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdOne)
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdTwo)
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
		res, err := m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		assert.Error(t, err)
		assert.Nil(t, res)

		err = m.CreateAuthorizeCodeSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.InvalidateAuthorizeCodeSession(ctx, "4321")
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		require.Error(t, err)
		assert.EqualError(t, err, fosite.ErrInvalidatedAuthorizeCode.Error())
		assert.NotNil(t, res)
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

func TestHelperCreateGetDeletePKCERequestSession(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetPKCERequestSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreatePKCERequestSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetPKCERequestSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeletePKCERequestSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetPKCERequestSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}
}

var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:            "flush-1",
		RequestedAt:   time.Now().Round(time.Second),
		Client:        &client.Client{ClientID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-2",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:        &client.Client{ClientID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-3",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:        &client.Client{ClientID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
}

func TestHelperFlushTokens(m pkg.FositeStorer, lifespan time.Duration) func(t *testing.T) {

	ds := &fosite.DefaultSession{}

	return func(t *testing.T) {
		ctx := context.Background()
		for _, r := range flushRequests {
			require.NoError(t, m.CreateAccessTokenSession(ctx, r.ID, r))
			_, err := m.GetAccessTokenSession(ctx, r.ID, ds)
			require.NoError(t, err)
		}

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now().Add(-time.Hour*24)))
		_, err := m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.NoError(t, err)

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now().Add(-(lifespan+time.Hour/2))))
		_, err = m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.Error(t, err)

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now()))
		_, err = m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.Error(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.Error(t, err)
	}
}

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

package oauth2_test

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
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/x/sqlcon"
)

var defaultRequest = fosite.Request{
	ID:                "blank",
	RequestedAt:       time.Now().UTC().Round(time.Second),
	Client:            &client.Client{ClientID: "foobar"},
	RequestedScope:    fosite.Arguments{"fa", "ba"},
	GrantedScope:      fosite.Arguments{"fa", "ba"},
	RequestedAudience: fosite.Arguments{"ad1", "ad2"},
	GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
	Form:              url.Values{"foo": []string{"bar", "baz"}},
	Session:           &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
}

func mockRequestForeignKey(t *testing.T, id string, x managerTestSetup, createClient bool) {
	cl := &client.Client{ClientID: "foobar"}
	cr := &consent.ConsentRequest{
		Client: cl, OpenIDConnectContext: new(consent.OpenIDConnectContext), LoginChallenge: id,
		Challenge: id, Verifier: id, AuthenticatedAt: time.Now(), RequestedAt: time.Now(),
	}

	if createClient {
		require.NoError(t, x.cl.CreateClient(context.Background(), cl))
	}

	require.NoError(t, x.co.CreateAuthenticationRequest(context.Background(), &consent.AuthenticationRequest{Client: cl, OpenIDConnectContext: new(consent.OpenIDConnectContext), Challenge: id, Verifier: id, AuthenticatedAt: time.Now(), RequestedAt: time.Now()}))
	require.NoError(t, x.co.CreateConsentRequest(context.Background(), cr))
	_, err := x.co.HandleConsentRequest(context.Background(), id, &consent.HandledConsentRequest{
		ConsentRequest: cr, Session: new(consent.ConsentRequestSessionData), AuthenticatedAt: time.Now(),
		Challenge:   id,
		RequestedAt: time.Now(),
	})
	require.NoError(t, err)
}

func testHelperUniqueConstraints(m managerTestSetup, storageType string) func(t *testing.T) {
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
		mockRequestForeignKey(t, requestId, m, true)
		cl := &client.Client{ClientID: "foobar"}

		signatureOne := uuid.New()
		signatureTwo := uuid.New()
		fositeRequest := &fosite.Request{
			ID:          requestId,
			Client:      cl,
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     &Session{},
		}

		err := m.f.CreateRefreshTokenSession(context.TODO(), signatureOne, fositeRequest)
		assert.NoError(t, err)
		err = m.f.CreateAccessTokenSession(context.TODO(), signatureOne, fositeRequest)
		assert.NoError(t, err)

		// attempting to insert new records with the SAME requestID should fail as there is a unique index
		// on the request_id column

		err = m.f.CreateRefreshTokenSession(context.TODO(), signatureTwo, fositeRequest)
		dbErrorIsConstraintError(err)
		err = m.f.CreateAccessTokenSession(context.TODO(), signatureTwo, fositeRequest)
		dbErrorIsConstraintError(err)
	}
}

func testHelperCreateGetDeleteOpenIDConnectSession(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		ctx := context.Background()
		_, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		assert.NotNil(t, err)

		err = m.CreateOpenIDConnectSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{Session: &Session{}})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteOpenIDConnectSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		assert.NotNil(t, err)
	}
}

func testHelperCreateGetDeleteRefreshTokenSession(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, "4321", &Session{})
		assert.NotNil(t, err)

		err = m.CreateRefreshTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetRefreshTokenSession(ctx, "4321", &Session{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteRefreshTokenSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "4321", &Session{})
		assert.NotNil(t, err)
	}
}

func testHelperRevokeRefreshToken(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, "1111", &Session{})
		assert.Error(t, err)

		reqIdOne := uuid.New()
		reqIdTwo := uuid.New()

		mockRequestForeignKey(t, reqIdOne, x, false)
		mockRequestForeignKey(t, reqIdTwo, x, false)

		err = m.CreateRefreshTokenSession(ctx, "1111", &fosite.Request{ID: reqIdOne, Client: &client.Client{ClientID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &Session{}})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", &fosite.Request{ID: reqIdTwo, Client: &client.Client{ClientID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &Session{}})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &Session{})
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdOne)
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdTwo)
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &Session{})
		assert.NotNil(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1122", &Session{})
		assert.NotNil(t, err)

	}
}

func testHelperCreateGetDeleteAuthorizeCodes(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		mockRequestForeignKey(t, "blank", x, false)

		ctx := context.Background()
		res, err := m.GetAuthorizeCodeSession(ctx, "4321", &Session{})
		assert.Error(t, err)
		assert.Nil(t, res)

		err = m.CreateAuthorizeCodeSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, "4321", &Session{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.InvalidateAuthorizeCodeSession(ctx, "4321")
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, "4321", &Session{})
		require.Error(t, err)
		assert.EqualError(t, err, fosite.ErrInvalidatedAuthorizeCode.Error())
		assert.NotNil(t, res)
	}
}

func testHelperNilAccessToken(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f
		c := &client.Client{ClientID: "nil-request-client-id-123"}
		require.NoError(t, x.cl.CreateClient(context.Background(), c))
		err := m.CreateAccessTokenSession(context.TODO(), "nil-request-id", &fosite.Request{
			ID:                "",
			RequestedAt:       time.Now().UTC().Round(time.Second),
			Client:            c,
			RequestedScope:    fosite.Arguments{"fa", "ba"},
			GrantedScope:      fosite.Arguments{"fa", "ba"},
			RequestedAudience: fosite.Arguments{"ad1", "ad2"},
			GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
			Form:              url.Values{"foo": []string{"bar", "baz"}},
			Session:           &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		})
		require.NoError(t, err)
	}
}

func testHelperCreateGetDeleteAccessTokenSession(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		ctx := context.Background()
		_, err := m.GetAccessTokenSession(ctx, "4321", &Session{})
		assert.Error(t, err)

		err = m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAccessTokenSession(ctx, "4321", &Session{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, "4321", &Session{})
		assert.Error(t, err)
	}
}

func testHelperCreateGetDeletePKCERequestSession(x managerTestSetup) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.f

		ctx := context.Background()
		_, err := m.GetPKCERequestSession(ctx, "4321", &Session{})
		assert.NotNil(t, err)

		err = m.CreatePKCERequestSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetPKCERequestSession(ctx, "4321", &Session{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeletePKCERequestSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetPKCERequestSession(ctx, "4321", &Session{})
		assert.NotNil(t, err)
	}
}

func testHelperFlushTokens(x managerTestSetup, lifespan time.Duration) func(t *testing.T) {
	m := x.f
	ds := &Session{}

	return func(t *testing.T) {
		ctx := context.Background()
		for _, r := range flushRequests {
			mockRequestForeignKey(t, r.ID, x, false)
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

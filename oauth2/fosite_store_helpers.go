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
	"crypto/sha256"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gobuffalo/pop/v6"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite/handler/rfc7523"

	"github.com/ory/hydra/oauth2/trust"

	"github.com/ory/hydra/x"

	"github.com/ory/fosite/storage"
	"github.com/ory/x/sqlxx"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/x/sqlcon"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
)

func signatureFromJTI(jti string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(jti)))
}

type BlacklistedJTI struct {
	JTI    string    `db:"-"`
	ID     string    `db:"signature"`
	Expiry time.Time `db:"expires_at"`
}

func (j *BlacklistedJTI) AfterFind(_ *pop.Connection) error {
	j.Expiry = j.Expiry.UTC()
	return nil
}

func (BlacklistedJTI) TableName() string {
	return "hydra_oauth2_jti_blacklist"
}

func NewBlacklistedJTI(jti string, exp time.Time) *BlacklistedJTI {
	return &BlacklistedJTI{
		JTI: jti,
		ID:  signatureFromJTI(jti),
		// because the database timestamp types are not as accurate as time.Time we truncate to seconds (which should always work)
		Expiry: exp.UTC().Truncate(time.Second),
	}
}

type AssertionJWTReader interface {
	x.FositeStorer

	GetClientAssertionJWT(ctx context.Context, jti string) (*BlacklistedJTI, error)

	SetClientAssertionJWTRaw(context.Context, *BlacklistedJTI) error
}

var defaultRequest = fosite.Request{
	ID:                "blank",
	RequestedAt:       time.Now().UTC().Round(time.Second),
	Client:            &client.Client{OutfacingID: "foobar"},
	RequestedScope:    fosite.Arguments{"fa", "ba"},
	GrantedScope:      fosite.Arguments{"fa", "ba"},
	RequestedAudience: fosite.Arguments{"ad1", "ad2"},
	GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
	Form:              url.Values{"foo": []string{"bar", "baz"}},
	Session:           &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
}

var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:             "flush-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-2",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-3",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
}

func mockRequestForeignKey(t *testing.T, id string, x InternalRegistry, createClient bool) {
	cl := &client.Client{OutfacingID: "foobar"}
	cr := &consent.ConsentRequest{
		Client: cl, OpenIDConnectContext: new(consent.OpenIDConnectContext), LoginChallenge: sqlxx.NullString(id),
		ID: id, Verifier: id, AuthenticatedAt: sqlxx.NullTime(time.Now()), RequestedAt: time.Now(),
	}

	if createClient {
		require.NoError(t, x.ClientManager().CreateClient(context.Background(), cl))
	}

	require.NoError(t, x.ConsentManager().CreateLoginRequest(context.Background(), &consent.LoginRequest{Client: cl, OpenIDConnectContext: new(consent.OpenIDConnectContext), ID: id, Verifier: id, AuthenticatedAt: sqlxx.NullTime(time.Now()), RequestedAt: time.Now()}))
	require.NoError(t, x.ConsentManager().CreateConsentRequest(context.Background(), cr))
	_, err := x.ConsentManager().HandleConsentRequest(context.Background(), id, &consent.HandledConsentRequest{
		ConsentRequest: cr, Session: new(consent.ConsentRequestSessionData), AuthenticatedAt: sqlxx.NullTime(time.Now()),
		ID:          id,
		RequestedAt: time.Now(),
		HandledAt:   sqlxx.NullTime(time.Now()),
	})
	require.NoError(t, err)
}

// TestHelperRunner is used to run the database suite of tests in this package.
// KEEP EXPORTED AND AVAILABLE FOR THIRD PARTIES TO TEST PLUGINS!
func TestHelperRunner(t *testing.T, store InternalRegistry, k string) {
	t.Helper()
	if k != "memory" {
		t.Run(fmt.Sprintf("case=testHelperUniqueConstraints/db=%s", k), testHelperRequestIDMultiples(store, k))
		t.Run("case=testFositeSqlStoreTransactionsCommitAccessToken", testFositeSqlStoreTransactionCommitAccessToken(store))
		t.Run("case=testFositeSqlStoreTransactionsRollbackAccessToken", testFositeSqlStoreTransactionRollbackAccessToken(store))
		t.Run("case=testFositeSqlStoreTransactionCommitRefreshToken", testFositeSqlStoreTransactionCommitRefreshToken(store))
		t.Run("case=testFositeSqlStoreTransactionRollbackRefreshToken", testFositeSqlStoreTransactionRollbackRefreshToken(store))
		t.Run("case=testFositeSqlStoreTransactionCommitAuthorizeCode", testFositeSqlStoreTransactionCommitAuthorizeCode(store))
		t.Run("case=testFositeSqlStoreTransactionRollbackAuthorizeCode", testFositeSqlStoreTransactionRollbackAuthorizeCode(store))
		t.Run("case=testFositeSqlStoreTransactionCommitPKCERequest", testFositeSqlStoreTransactionCommitPKCERequest(store))
		t.Run("case=testFositeSqlStoreTransactionRollbackPKCERequest", testFositeSqlStoreTransactionRollbackPKCERequest(store))
		t.Run("case=testFositeSqlStoreTransactionCommitOpenIdConnectSession", testFositeSqlStoreTransactionCommitOpenIdConnectSession(store))
		t.Run("case=testFositeSqlStoreTransactionRollbackOpenIdConnectSession", testFositeSqlStoreTransactionRollbackOpenIdConnectSession(store))

	}
	t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteAuthorizeCodes/db=%s", k), testHelperCreateGetDeleteAuthorizeCodes(store))
	t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteAccessTokenSession/db=%s", k), testHelperCreateGetDeleteAccessTokenSession(store))
	t.Run(fmt.Sprintf("case=testHelperNilAccessToken/db=%s", k), testHelperNilAccessToken(store))
	t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteOpenIDConnectSession/db=%s", k), testHelperCreateGetDeleteOpenIDConnectSession(store))
	t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteRefreshTokenSession/db=%s", k), testHelperCreateGetDeleteRefreshTokenSession(store))
	t.Run(fmt.Sprintf("case=testHelperRevokeRefreshToken/db=%s", k), testHelperRevokeRefreshToken(store))
	t.Run(fmt.Sprintf("case=testHelperCreateGetDeletePKCERequestSession/db=%s", k), testHelperCreateGetDeletePKCERequestSession(store))
	t.Run(fmt.Sprintf("case=testHelperFlushTokens/db=%s", k), testHelperFlushTokens(store, time.Hour))
	t.Run(fmt.Sprintf("case=testHelperFlushTokensWithLimitAndBatchSize/db=%s", k), testHelperFlushTokensWithLimitAndBatchSize(store, 3, 2))
	t.Run(fmt.Sprintf("case=testFositeStoreSetClientAssertionJWT/db=%s", k), testFositeStoreSetClientAssertionJWT(store))
	t.Run(fmt.Sprintf("case=testFositeStoreClientAssertionJWTValid/db=%s", k), testFositeStoreClientAssertionJWTValid(store))
	t.Run(fmt.Sprintf("case=testHelperDeleteAccessTokens/db=%s", k), testHelperDeleteAccessTokens(store))
	t.Run(fmt.Sprintf("case=testHelperRevokeAccessToken/db=%s", k), testHelperRevokeAccessToken(store))
	t.Run(fmt.Sprintf("case=testFositeJWTBearerGrantStorage/db=%s", k), testFositeJWTBearerGrantStorage(store))
}

func testHelperRequestIDMultiples(m InternalRegistry, _ string) func(t *testing.T) {
	return func(t *testing.T) {
		requestId := uuid.New()
		mockRequestForeignKey(t, requestId, m, true)
		cl := &client.Client{OutfacingID: "foobar"}

		fositeRequest := &fosite.Request{
			ID:          requestId,
			Client:      cl,
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     &Session{},
		}

		for i := 0; i < 4; i++ {
			signature := uuid.New()
			err := m.OAuth2Storage().CreateRefreshTokenSession(context.TODO(), signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateAccessTokenSession(context.TODO(), signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateOpenIDConnectSession(context.TODO(), signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreatePKCERequestSession(context.TODO(), signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateAuthorizeCodeSession(context.TODO(), signature, fositeRequest)
			assert.NoError(t, err)
		}
	}
}

func testHelperCreateGetDeleteOpenIDConnectSession(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

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

func testHelperCreateGetDeleteRefreshTokenSession(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

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

func testHelperRevokeRefreshToken(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, "1111", &Session{})
		assert.Error(t, err)

		reqIdOne := uuid.New()
		reqIdTwo := uuid.New()

		mockRequestForeignKey(t, reqIdOne, x, false)
		mockRequestForeignKey(t, reqIdTwo, x, false)

		err = m.CreateRefreshTokenSession(ctx, "1111", &fosite.Request{ID: reqIdOne, Client: &client.Client{OutfacingID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &Session{}})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", &fosite.Request{ID: reqIdTwo, Client: &client.Client{OutfacingID: "foobar"}, RequestedAt: time.Now().UTC().Round(time.Second), Session: &Session{}})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", &Session{})
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdOne)
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdTwo)
		require.NoError(t, err)

		req, err := m.GetRefreshTokenSession(ctx, "1111", &Session{})
		assert.NotNil(t, req)
		assert.EqualError(t, err, fosite.ErrInactiveToken.Error())

		req, err = m.GetRefreshTokenSession(ctx, "1122", &Session{})
		assert.NotNil(t, req)
		assert.EqualError(t, err, fosite.ErrInactiveToken.Error())

	}
}

func testHelperCreateGetDeleteAuthorizeCodes(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

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

func testHelperNilAccessToken(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		c := &client.Client{OutfacingID: "nil-request-client-id-123"}
		require.NoError(t, x.ClientManager().CreateClient(context.Background(), c))
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

func testHelperCreateGetDeleteAccessTokenSession(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

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

func testHelperDeleteAccessTokens(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := context.Background()

		err := m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, "4321", &Session{})
		require.NoError(t, err)

		err = m.DeleteAccessTokens(ctx, defaultRequest.Client.GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, "4321", &Session{})
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperRevokeAccessToken(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := context.Background()

		err := m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, "4321", &Session{})
		require.NoError(t, err)

		err = m.RevokeAccessToken(ctx, defaultRequest.GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, "4321", &Session{})
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperCreateGetDeletePKCERequestSession(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

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

func testHelperFlushTokens(x InternalRegistry, lifespan time.Duration) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &Session{}

	return func(t *testing.T) {
		ctx := context.Background()
		for _, r := range flushRequests {
			mockRequestForeignKey(t, r.ID, x, false)
			require.NoError(t, m.CreateAccessTokenSession(ctx, r.ID, r))
			_, err := m.GetAccessTokenSession(ctx, r.ID, ds)
			require.NoError(t, err)
		}

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now().Add(-time.Hour*24), 100, 10))
		_, err := m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.NoError(t, err)

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now().Add(-(lifespan+time.Hour/2)), 100, 10))
		_, err = m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.Error(t, err)

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now(), 100, 10))
		_, err = m.GetAccessTokenSession(ctx, "flush-1", ds)
		require.NoError(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-2", ds)
		require.Error(t, err)
		_, err = m.GetAccessTokenSession(ctx, "flush-3", ds)
		require.Error(t, err)
	}
}

func testHelperFlushTokensWithLimitAndBatchSize(x InternalRegistry, limit int, batchSize int) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &Session{}

	return func(t *testing.T) {
		ctx := context.Background()
		var requests []*fosite.Request

		// create five expired requests
		id := uuid.New()
		for i := 0; i < 5; i++ {
			r := createTestRequest(fmt.Sprintf("%s-%d", id, i+1))
			r.RequestedAt = time.Now().Add(-2 * time.Hour)
			mockRequestForeignKey(t, r.ID, x, false)
			require.NoError(t, m.CreateAccessTokenSession(ctx, r.ID, r))
			_, err := m.GetAccessTokenSession(ctx, r.ID, ds)
			require.NoError(t, err)
			requests = append(requests, r)
		}

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now(), limit, batchSize))
		for i := range requests {
			_, err := m.GetAccessTokenSession(ctx, requests[i].ID, ds)
			if i >= limit {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
	}
}

func testFositeSqlStoreTransactionCommitAccessToken(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionRollbackAccessToken(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionCommitRefreshToken(m InternalRegistry) func(t *testing.T) {

	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestCommit(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionRollbackRefreshToken(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestRollback(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionCommitAuthorizeCode(m InternalRegistry) func(t *testing.T) {

	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionRollbackAuthorizeCode(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionCommitPKCERequest(m InternalRegistry) func(t *testing.T) {

	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

func testFositeSqlStoreTransactionRollbackPKCERequest(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

// OpenIdConnect tests can't use the helper functions, due to the signature of GetOpenIdConnectSession being
// different from the other getter methods
func testFositeSqlStoreTransactionCommitOpenIdConnectSession(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		txnStore, ok := m.OAuth2Storage().(storage.Transactional)
		require.True(t, ok)
		ctx := context.Background()
		ctx, err := txnStore.BeginTX(ctx)
		require.NoError(t, err)
		signature := uuid.New()
		testRequest := createTestRequest(signature)
		err = m.OAuth2Storage().CreateOpenIDConnectSession(ctx, signature, testRequest)
		require.NoError(t, err)
		err = txnStore.Commit(ctx)
		require.NoError(t, err)

		// Require a new context, since the old one contains the transaction.
		res, err := m.OAuth2Storage().GetOpenIDConnectSession(context.Background(), signature, testRequest)
		// session should have been created successfully because Commit did not return an error
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		// test delete within a transaction
		ctx, err = txnStore.BeginTX(context.Background())
		require.NoError(t, err)
		err = m.OAuth2Storage().DeleteOpenIDConnectSession(ctx, signature)
		require.NoError(t, err)
		err = txnStore.Commit(ctx)
		require.NoError(t, err)

		// Require a new context, since the old one contains the transaction.
		_, err = m.OAuth2Storage().GetOpenIDConnectSession(context.Background(), signature, testRequest)
		// Since commit worked for delete, we should get an error here.
		require.Error(t, err)
	}
}

func testFositeSqlStoreTransactionRollbackOpenIdConnectSession(m InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		txnStore, ok := m.OAuth2Storage().(storage.Transactional)
		require.True(t, ok)
		ctx := context.Background()
		ctx, err := txnStore.BeginTX(ctx)
		require.NoError(t, err)

		signature := uuid.New()
		testRequest := createTestRequest(signature)
		err = m.OAuth2Storage().CreateOpenIDConnectSession(ctx, signature, testRequest)
		require.NoError(t, err)
		err = txnStore.Rollback(ctx)
		require.NoError(t, err)

		// Require a new context, since the old one contains the transaction.
		ctx = context.Background()
		_, err = m.OAuth2Storage().GetOpenIDConnectSession(ctx, signature, testRequest)
		// Since we rolled back above, the session should not exist and getting it should result in an error
		require.Error(t, err)

		// create a new session, delete it, then rollback the delete. We should be able to then get it.
		signature2 := uuid.New()
		testRequest2 := createTestRequest(signature2)
		err = m.OAuth2Storage().CreateOpenIDConnectSession(ctx, signature2, testRequest2)
		require.NoError(t, err)
		_, err = m.OAuth2Storage().GetOpenIDConnectSession(ctx, signature2, testRequest2)
		require.NoError(t, err)

		ctx, err = txnStore.BeginTX(context.Background())
		require.NoError(t, err)
		err = m.OAuth2Storage().DeleteOpenIDConnectSession(ctx, signature2)
		require.NoError(t, err)
		err = txnStore.Rollback(ctx)

		require.NoError(t, err)
		_, err = m.OAuth2Storage().GetOpenIDConnectSession(context.Background(), signature2, testRequest2)
		require.NoError(t, err)
	}
}

func testFositeStoreSetClientAssertionJWT(m InternalRegistry) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=basic setting works", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			jti := NewBlacklistedJTI("basic jti", time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))

			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			require.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})

		t.Run("case=errors when the JTI is blacklisted", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			jti := NewBlacklistedJTI("already set jti", time.Now().Add(time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.ErrorIs(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry), fosite.ErrJTIKnown)
		})

		t.Run("case=deletes expired JTIs", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			expiredJTI := NewBlacklistedJTI("expired jti", time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), expiredJTI))
			newJTI := NewBlacklistedJTI("some new jti", time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), newJTI.JTI, newJTI.Expiry))

			_, err := store.GetClientAssertionJWT(context.Background(), expiredJTI.JTI)
			assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
			cmp, err := store.GetClientAssertionJWT(context.Background(), newJTI.JTI)
			require.NoError(t, err)
			assert.Equal(t, newJTI, cmp)
		})

		t.Run("case=inserts same JTI if expired", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			jti := NewBlacklistedJTI("going to be reused jti", time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			jti.Expiry = jti.Expiry.Add(2 * time.Minute)
			assert.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))
			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			assert.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})
	}
}

func testFositeStoreClientAssertionJWTValid(m InternalRegistry) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=returns valid on unknown JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), "unknown jti"))
		})

		t.Run("case=returns invalid on known JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			jti := NewBlacklistedJTI("known jti", time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.True(t, errors.Is(store.ClientAssertionJWTValid(context.Background(), jti.JTI), fosite.ErrJTIKnown))
		})

		t.Run("case=returns valid on expired JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(AssertionJWTReader)
			require.True(t, ok)
			jti := NewBlacklistedJTI("expired jti 2", time.Now().Add(-time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), jti.JTI))
		})
	}
}

func testFositeJWTBearerGrantStorage(x InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		grantManager := x.GrantManager()
		keyManager := x.KeyManager()
		keyGenerators := x.KeyGenerators()
		keyGenerator, ok := keyGenerators[string(jose.RS256)]
		require.True(t, ok)
		grantStorage := x.OAuth2Storage().(rfc7523.RFC7523KeyStorage)

		t.Run("case=associated key added with grant", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("token-service-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "token-service"
			subject := "bob@example.com"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			storedKeySet, err := grantStorage.GetPublicKeys(context.TODO(), issuer, subject)
			require.NoError(t, err)
			require.Len(t, storedKeySet.Keys, 0)

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			storedKeySet, err = grantStorage.GetPublicKeys(context.TODO(), issuer, subject)
			require.NoError(t, err)
			assert.Len(t, storedKeySet.Keys, 1)

			storedKey, err := grantStorage.GetPublicKey(context.TODO(), issuer, subject, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, publicKey.KeyID, storedKey.KeyID)
			assert.Equal(t, publicKey.Use, storedKey.Use)
			assert.Equal(t, publicKey.Key, storedKey.Key)

			storedScopes, err := grantStorage.GetPublicKeyScopes(context.TODO(), issuer, subject, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, grant.Scope, storedScopes)

			storedKeySet, err = keyManager.GetKey(context.TODO(), issuer, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, publicKey.KeyID, storedKeySet.Keys[0].KeyID)
			assert.Equal(t, publicKey.Use, storedKeySet.Keys[0].Use)
			assert.Equal(t, publicKey.Key, storedKeySet.Keys[0].Key)
		})

		t.Run("case=only associated key returns", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("some-key", "sig")
			require.NoError(t, err)

			err = keyManager.AddKeySet(context.TODO(), "some-set", keySet)
			require.NoError(t, err)

			keySet, err = keyGenerator.Generate("maria-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "maria"
			subject := "maria@example.com"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			storedKeySet, err := grantStorage.GetPublicKeys(context.TODO(), issuer, subject)
			require.NoError(t, err)
			assert.Len(t, storedKeySet.Keys, 1)
			assert.Equal(t, publicKey.KeyID, storedKeySet.Keys[0].KeyID)
			assert.Equal(t, publicKey.Use, storedKeySet.Keys[0].Use)
			assert.Equal(t, publicKey.Key, storedKeySet.Keys[0].Key)

			storedKeySet, err = grantStorage.GetPublicKeys(context.TODO(), issuer, "non-existing-subject")
			require.NoError(t, err)
			assert.Len(t, storedKeySet.Keys, 0)

			_, err = grantStorage.GetPublicKeyScopes(context.TODO(), issuer, "non-existing-subject", publicKey.KeyID)
			require.Error(t, err)
		})

		t.Run("case=associated key is deleted, when granted is deleted", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("hackerman-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "aeneas"
			subject := "aeneas@example.com"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, subject, grant.PublicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(context.TODO(), issuer, publicKey.KeyID)
			require.NoError(t, err)

			err = grantManager.DeleteGrant(context.TODO(), grant.ID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, subject, publicKey.KeyID)
			assert.Error(t, err)

			_, err = keyManager.GetKey(context.TODO(), issuer, publicKey.KeyID)
			assert.Error(t, err)
		})

		t.Run("case=associated grant is deleted, when key is deleted", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("vladimir-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "vladimir"
			subject := "vladimir@example.com"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(context.TODO(), issuer, publicKey.KeyID)
			require.NoError(t, err)

			err = keyManager.DeleteKey(context.TODO(), issuer, publicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(context.TODO(), issuer, publicKey.KeyID)
			assert.Error(t, err)

			_, err = grantManager.GetConcreteGrant(context.TODO(), grant.ID)
			assert.Error(t, err)
		})

		t.Run("case=only returns the key when subject matches", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("issuer-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "limited-issuer"
			subject := "jagoba"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			// All three get methods should only return the public key when using the valid subject
			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, "any-subject-1", publicKey.KeyID)
			require.Error(t, err)
			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKeyScopes(context.TODO(), issuer, "any-subject-2", publicKey.KeyID)
			require.Error(t, err)
			_, err = grantStorage.GetPublicKeyScopes(context.TODO(), issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			jwks, err := grantStorage.GetPublicKeys(context.TODO(), issuer, "any-subject-3")
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.Empty(t, jwks.Keys)
			jwks, err = grantStorage.GetPublicKeys(context.TODO(), issuer, subject)
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.NotEmpty(t, jwks.Keys)
		})

		t.Run("case=returns the key when any subject is allowed", func(t *testing.T) {
			keySet, err := keyGenerator.Generate("issuer-key", "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[1]
			issuer := "unlimited-issuer"
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         "",
				AllowAnySubject: true,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			err = grantManager.CreateGrant(context.TODO(), grant, publicKey)
			require.NoError(t, err)

			// All three get methods should always return the public key
			_, err = grantStorage.GetPublicKey(context.TODO(), issuer, "any-subject-1", publicKey.KeyID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKeyScopes(context.TODO(), issuer, "any-subject-2", publicKey.KeyID)
			require.NoError(t, err)

			jwks, err := grantStorage.GetPublicKeys(context.TODO(), issuer, "any-subject-3")
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.NotEmpty(t, jwks.Keys)
		})
	}
}

func doTestCommit(m InternalRegistry, t *testing.T,
	createFn func(context.Context, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {

	txnStore, ok := m.OAuth2Storage().(storage.Transactional)
	require.True(t, ok)
	ctx := context.Background()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.New()
	err = createFn(ctx, signature, createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	res, err := getFn(context.Background(), signature, &Session{})
	// token should have been created successfully because Commit did not return an error
	require.NoError(t, err)
	AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

	// testrevoke within a transaction
	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature)
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	_, err = getFn(context.Background(), signature, &Session{})
	// Since commit worked for revoke, we should get an error here.
	require.Error(t, err)
}

func doTestRollback(m InternalRegistry, t *testing.T,
	createFn func(context.Context, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(storage.Transactional)
	require.True(t, ok)

	ctx := context.Background()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.New()
	err = createFn(ctx, signature, createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	ctx = context.Background()
	_, err = getFn(ctx, signature, &Session{})
	// Since we rolled back above, the token should not exist and getting it should result in an error
	require.Error(t, err)

	// create a new token, revoke it, then rollback the revoke. We should be able to then get it successfully.
	signature2 := uuid.New()
	err = createFn(ctx, signature2, createTestRequest(signature2))
	require.NoError(t, err)
	_, err = getFn(ctx, signature2, &Session{})
	require.NoError(t, err)

	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature2)
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	_, err = getFn(context.Background(), signature2, &Session{})
	require.NoError(t, err)
}

func createTestRequest(id string) *fosite.Request {
	return &fosite.Request{
		ID:                id,
		RequestedAt:       time.Now().UTC().Round(time.Second),
		Client:            &client.Client{OutfacingID: "foobar"},
		RequestedScope:    fosite.Arguments{"fa", "ba"},
		GrantedScope:      fosite.Arguments{"fa", "ba"},
		RequestedAudience: fosite.Arguments{"ad1", "ad2"},
		GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
		Form:              url.Values{"foo": []string{"bar", "baz"}},
		Session:           &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	}
}

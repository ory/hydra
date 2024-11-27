// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"testing"
	"time"

	"github.com/ory/hydra/v2/persistence/sql"

	"github.com/go-jose/go-jose/v3"
	gofrsuuid "github.com/gofrs/uuid"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/rfc7523"
	"github.com/ory/fosite/storage"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/assertx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

var defaultIgnoreKeys = []string{
	"id",
	"session",
	"requested_scope",
	"granted_scope",
	"form",
	"created_at",
	"updated_at",
	"client.created_at",
	"client.updated_at",
	"requestedAt",
	"client.client_secret",
}

func newDefaultRequest(id string) fosite.Request {
	return fosite.Request{
		ID:          id,
		RequestedAt: time.Now().UTC().Round(time.Second),
		Client: &client.Client{
			ID:                 "foobar",
			Contacts:           []string{},
			RedirectURIs:       []string{},
			Audience:           []string{},
			AllowedCORSOrigins: []string{},
			ResponseTypes:      []string{},
			GrantTypes:         []string{},
			JSONWebKeys:        &x.JoseJSONWebKeySet{},
			Metadata:           sqlxx.JSONRawMessage("{}"),
		},
		RequestedScope:    fosite.Arguments{"fa", "ba"},
		GrantedScope:      fosite.Arguments{"fa", "ba"},
		RequestedAudience: fosite.Arguments{"ad1", "ad2"},
		GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
		Form:              url.Values{"foo": []string{"bar", "baz"}},
		Session:           oauth2.NewSession("bar"),
	}
}

var defaultRequest = newDefaultRequest("blank")

// var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:             "flush-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{ID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-2",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:         &client.Client{ID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-3",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:         &client.Client{ID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
}

func mockRequestForeignKey(t *testing.T, id string, x oauth2.InternalRegistry) {
	cl := &client.Client{ID: "foobar"}
	cr := &flow.OAuth2ConsentRequest{
		Client:               cl,
		OpenIDConnectContext: new(flow.OAuth2ConsentRequestOpenIDConnectContext),
		LoginChallenge:       sqlxx.NullString(id),
		ID:                   id,
		Verifier:             id,
		CSRF:                 id,
		AuthenticatedAt:      sqlxx.NullTime(time.Now()),
		RequestedAt:          time.Now(),
	}

	ctx := context.Background()
	if _, err := x.ClientManager().GetClient(ctx, cl.ID); errors.Is(err, sqlcon.ErrNoRows) {
		require.NoError(t, x.ClientManager().CreateClient(ctx, cl))
	}

	f, err := x.ConsentManager().CreateLoginRequest(
		ctx, &flow.LoginRequest{
			Client:               cl,
			OpenIDConnectContext: new(flow.OAuth2ConsentRequestOpenIDConnectContext),
			ID:                   id,
			Verifier:             id,
			AuthenticatedAt:      sqlxx.NullTime(time.Now()),
			RequestedAt:          time.Now(),
		})
	require.NoError(t, err)
	err = x.ConsentManager().CreateConsentRequest(ctx, f, cr)
	require.NoError(t, err)

	encodedFlow, err := f.ToConsentVerifier(ctx, x)
	require.NoError(t, err)

	_, err = x.ConsentManager().HandleConsentRequest(ctx, f, &flow.AcceptOAuth2ConsentRequest{
		ConsentRequest:  cr,
		Session:         new(flow.AcceptOAuth2ConsentRequestSession),
		AuthenticatedAt: sqlxx.NullTime(time.Now()),
		ID:              encodedFlow,
		RequestedAt:     time.Now(),
		HandledAt:       sqlxx.NullTime(time.Now()),
	})

	require.NoError(t, err)
}

func TestHelperRunner(t *testing.T) {
}

func testHelperRequestIDMultiples(m oauth2.InternalRegistry, _ string) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		requestID := uuid.New()
		mockRequestForeignKey(t, requestID, m)
		cl := &client.Client{ID: "foobar"}

		fositeRequest := &fosite.Request{
			ID:          requestID,
			Client:      cl,
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewSession("bar"),
		}

		for i := 0; i < 4; i++ {
			signature := uuid.New()
			accessSignature := uuid.New()
			err := m.OAuth2Storage().CreateRefreshTokenSession(ctx, signature, accessSignature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateAccessTokenSession(ctx, signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateOpenIDConnectSession(ctx, signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreatePKCERequestSession(ctx, signature, fositeRequest)
			assert.NoError(t, err)
			err = m.OAuth2Storage().CreateAuthorizeCodeSession(ctx, signature, fositeRequest)
			assert.NoError(t, err)
		}
	}
}

func testHelperCreateGetDeleteOpenIDConnectSession(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.New()
		ctx := context.Background()
		_, err := m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewSession("bar")})
		assert.NotNil(t, err)

		err = m.CreateOpenIDConnectSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewSession("bar")})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteOpenIDConnectSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewSession("bar")})
		assert.NotNil(t, err)
	}
}

func testHelperCreateGetDeleteRefreshTokenSession(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.New()
		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.NotNil(t, err)

		err = m.CreateRefreshTokenSession(ctx, code, "", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetRefreshTokenSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteRefreshTokenSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.NotNil(t, err)
	}
}

func testHelperRevokeRefreshToken(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		ctx := context.Background()
		_, err := m.GetRefreshTokenSession(ctx, "1111", oauth2.NewSession("bar"))
		assert.Error(t, err)

		reqIdOne := uuid.New()
		reqIdTwo := uuid.New()

		mockRequestForeignKey(t, reqIdOne, x)
		mockRequestForeignKey(t, reqIdTwo, x)

		err = m.CreateRefreshTokenSession(ctx, "1111", "", &fosite.Request{
			ID:          reqIdOne,
			Client:      &client.Client{ID: "foobar"},
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewSession("user"),
		})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", "", &fosite.Request{
			ID:          reqIdTwo,
			Client:      &client.Client{ID: "foobar"},
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewSession("user"),
		})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", oauth2.NewSession("bar"))
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdOne)
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdTwo)
		require.NoError(t, err)

		req, err := m.GetRefreshTokenSession(ctx, "1111", oauth2.NewSession("bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())

		req, err = m.GetRefreshTokenSession(ctx, "1122", oauth2.NewSession("bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperCreateGetDeleteAuthorizeCodes(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		mockRequestForeignKey(t, "blank", x)

		code := uuid.New()

		ctx := context.Background()
		res, err := m.GetAuthorizeCodeSession(ctx, code, oauth2.NewSession("bar"))
		assert.Error(t, err)
		assert.Nil(t, res)

		err = m.CreateAuthorizeCodeSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.InvalidateAuthorizeCodeSession(ctx, code)
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, code, oauth2.NewSession("bar"))
		require.Error(t, err)
		assert.EqualError(t, err, fosite.ErrInvalidatedAuthorizeCode.Error())
		assert.NotNil(t, res)
	}
}

type testHelperExpiryFieldsResult struct {
	ExpiresAt time.Time `db:"expires_at"`
	name      string
}

func (r testHelperExpiryFieldsResult) TableName() string {
	return "hydra_oauth2_" + r.name
}

func testHelperExpiryFields(reg oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := reg.OAuth2Storage()
		t.Parallel()

		mockRequestForeignKey(t, "blank", reg)

		ctx := context.Background()

		s := oauth2.NewSession("bar")
		s.SetExpiresAt(fosite.AccessToken, time.Now().Add(time.Hour).Round(time.Minute))
		s.SetExpiresAt(fosite.RefreshToken, time.Now().Add(time.Hour*2).Round(time.Minute))
		s.SetExpiresAt(fosite.AuthorizeCode, time.Now().Add(time.Hour*3).Round(time.Minute))
		request := fosite.Request{
			ID:          uuid.New(),
			RequestedAt: time.Now().UTC().Round(time.Second),
			Client: &client.Client{
				ID:       "foobar",
				Metadata: sqlxx.JSONRawMessage("{}"),
			},
			RequestedScope:    fosite.Arguments{"fa", "ba"},
			GrantedScope:      fosite.Arguments{"fa", "ba"},
			RequestedAudience: fosite.Arguments{"ad1", "ad2"},
			GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
			Form:              url.Values{"foo": []string{"bar", "baz"}},
			Session:           s,
		}

		t.Run("case=CreateAccessTokenSession", func(t *testing.T) {
			id := uuid.New()
			err := m.CreateAccessTokenSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "access"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", x.SignatureHash(id)).First(&r))

			assert.EqualValues(t, s.GetExpiresAt(fosite.AccessToken).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateRefreshTokenSession", func(t *testing.T) {
			id := uuid.New()
			err := m.CreateRefreshTokenSession(ctx, id, "", &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "refresh"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.RefreshToken).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateAuthorizeCodeSession", func(t *testing.T) {
			id := uuid.New()
			err := m.CreateAuthorizeCodeSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "code"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreatePKCERequestSession", func(t *testing.T) {
			id := uuid.New()
			err := m.CreatePKCERequestSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "pkce"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateOpenIDConnectSession", func(t *testing.T) {
			id := uuid.New()
			err := m.CreateOpenIDConnectSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "oidc"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})
	}
}

func testHelperNilAccessToken(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		c := &client.Client{ID: uuid.New()}
		require.NoError(t, x.ClientManager().CreateClient(context.Background(), c))
		err := m.CreateAccessTokenSession(context.Background(), uuid.New(), &fosite.Request{
			ID:                "",
			RequestedAt:       time.Now().UTC().Round(time.Second),
			Client:            c,
			RequestedScope:    fosite.Arguments{"fa", "ba"},
			GrantedScope:      fosite.Arguments{"fa", "ba"},
			RequestedAudience: fosite.Arguments{"ad1", "ad2"},
			GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
			Form:              url.Values{"foo": []string{"bar", "baz"}},
			Session:           oauth2.NewSession("bar"),
		})
		require.NoError(t, err)
	}
}

func testHelperCreateGetDeleteAccessTokenSession(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.New()
		ctx := context.Background()
		_, err := m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.Error(t, err)

		err = m.CreateAccessTokenSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.Error(t, err)
	}
}

func testHelperDeleteAccessTokens(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := context.Background()

		code := uuid.New()
		err := m.CreateAccessTokenSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)

		err = m.DeleteAccessTokens(ctx, defaultRequest.Client.GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperRevokeAccessToken(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := context.Background()

		code := uuid.New()
		err := m.CreateAccessTokenSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)

		err = m.RevokeAccessToken(ctx, defaultRequest.GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, code, oauth2.NewSession("bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperRotateRefreshToken(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		createTokens := func(t *testing.T, r *fosite.Request) (refreshTokenSession string, accessTokenSession string) {
			refreshTokenSession = fmt.Sprintf("refresh_token_%s", uuid.New())
			accessTokenSession = fmt.Sprintf("access_token_%s", uuid.New())
			err := x.OAuth2Storage().CreateAccessTokenSession(ctx, accessTokenSession, r)
			require.NoError(t, err)

			err = x.OAuth2Storage().CreateRefreshTokenSession(ctx, refreshTokenSession, accessTokenSession, r)
			require.NoError(t, err)

			// Sanity check
			req, err := x.OAuth2Storage().GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			require.NoError(t, err)
			require.EqualValues(t, r.GetID(), req.GetID())

			req, err = x.OAuth2Storage().GetAccessTokenSession(ctx, accessTokenSession, nil)
			require.NoError(t, err)
			require.EqualValues(t, r.GetID(), req.GetID())
			return
		}

		t.Run("Revokes refresh token when grace period not configured", func(t *testing.T) {
			m := x.OAuth2Storage()
			r := newDefaultRequest(uuid.New())
			refreshTokenSession, accessTokenSession := createTokens(t, &r)

			err := m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession)
			require.NoError(t, err)

			_, err = m.GetAccessTokenSession(ctx, accessTokenSession, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound, "Token is no longer active because it was refreshed")

			_, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.ErrorIs(t, err, fosite.ErrInactiveToken, "Token is no longer active because it was refreshed")
		})

		t.Run("refresh token is valid until the grace period has ended", func(t *testing.T) {
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")

			// By setting this to one hour we ensure that using the refresh token triggers the start of the grace period.
			x.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1h")
			t.Cleanup(func() {
				x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
			})

			m := x.OAuth2Storage()
			r := newDefaultRequest(uuid.New())
			refreshTokenSession, accessTokenSession1 := createTokens(t, &r)
			accessTokenSession2 := fmt.Sprintf("access_token_%s", uuid.New())
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession2, &r))

			// Create a second access token
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			req, err := m.GetAccessTokenSession(ctx, accessTokenSession1, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			req, err = m.GetAccessTokenSession(ctx, accessTokenSession2, nil)
			assert.NoError(t, err, "The second access token is still valid.")

			req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.NoError(t, err)
			assert.Equal(t, r.GetID(), req.GetID())

			// We only wait a second, meaning that the token is theoretically still within TTL, but since the
			// grace period was issued, the token is still valid.
			time.Sleep(time.Second * 2)
			req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.Error(t, err)
		})

		t.Run("the used at time does not change", func(t *testing.T) {
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")

			// By setting this to one hour we ensure that using the refresh token triggers the start of the grace period.
			x.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1h")
			t.Cleanup(func() {
				x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
			})

			m := x.OAuth2Storage()
			r := newDefaultRequest(uuid.New())

			refreshTokenSession, _ := createTokens(t, &r)
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			var expected sql.OAuth2RefreshTable
			require.NoError(t, x.Persister().Connection(ctx).Where("signature=?", refreshTokenSession).First(&expected))
			assert.False(t, expected.FirstUsedAt.Time.IsZero())
			assert.True(t, expected.FirstUsedAt.Valid)

			// Refresh does not change the time
			time.Sleep(time.Second * 2)
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			var actual sql.OAuth2RefreshTable
			require.NoError(t, x.Persister().Connection(ctx).Where("signature=?", refreshTokenSession).First(&actual))
			assert.Equal(t, expected.FirstUsedAt.Time, actual.FirstUsedAt.Time)
		})

		t.Run("refresh token revokes all access tokens from the request if the access token signature is not found", func(t *testing.T) {
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")
			t.Cleanup(func() {
				x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
			})

			m := x.OAuth2Storage()
			r := newDefaultRequest(uuid.New())

			refreshTokenSession := fmt.Sprintf("refresh_token_%s", uuid.New())
			accessTokenSession1 := fmt.Sprintf("access_token_%s", uuid.New())
			accessTokenSession2 := fmt.Sprintf("access_token_%s", uuid.New())
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession1, &r))
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession2, &r))

			require.NoError(t, m.CreateRefreshTokenSession(ctx, refreshTokenSession, "", &r),
				"precondition failed: could not create refresh token session")

			// ACT
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			req, err := m.GetAccessTokenSession(ctx, accessTokenSession1, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			req, err = m.GetAccessTokenSession(ctx, accessTokenSession2, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.NoError(t, err)
			assert.Equal(t, r.GetID(), req.GetID())

			time.Sleep(time.Second * 2)

			req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.Error(t, err)
		})
	}
}

func testHelperCreateGetDeletePKCERequestSession(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.New()
		ctx := context.Background()
		_, err := m.GetPKCERequestSession(ctx, code, oauth2.NewSession("bar"))
		assert.NotNil(t, err)

		err = m.CreatePKCERequestSession(ctx, code, &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetPKCERequestSession(ctx, code, oauth2.NewSession("bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeletePKCERequestSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetPKCERequestSession(ctx, code, oauth2.NewSession("bar"))
		assert.NotNil(t, err)
	}
}

func testHelperFlushTokens(x oauth2.InternalRegistry, lifespan time.Duration) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &oauth2.Session{}

	return func(t *testing.T) {
		ctx := context.Background()
		for _, r := range flushRequests {
			mockRequestForeignKey(t, r.ID, x)
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
		require.NoError(t, m.DeleteAccessTokens(ctx, "foobar"))
	}
}

func testHelperFlushTokensWithLimitAndBatchSize(x oauth2.InternalRegistry, limit int, batchSize int) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &oauth2.Session{}

	return func(t *testing.T) {
		ctx := context.Background()
		var requests []*fosite.Request

		// create five expired requests
		id := uuid.New()
		totalCount := 5
		for i := 0; i < totalCount; i++ {
			r := createTestRequest(fmt.Sprintf("%s-%d", id, i+1))
			r.RequestedAt = time.Now().Add(-2 * time.Hour)
			mockRequestForeignKey(t, r.ID, x)
			require.NoError(t, m.CreateAccessTokenSession(ctx, r.ID, r))
			_, err := m.GetAccessTokenSession(ctx, r.ID, ds)
			require.NoError(t, err)
			requests = append(requests, r)
		}

		require.NoError(t, m.FlushInactiveAccessTokens(ctx, time.Now(), limit, batchSize))
		var notFoundCount, foundCount int
		for i := range requests {
			if _, err := m.GetAccessTokenSession(ctx, requests[i].ID, ds); err == nil {
				foundCount++
			} else {
				require.ErrorIs(t, err, fosite.ErrNotFound)
				notFoundCount++
			}
		}
		assert.Equal(t, limit, notFoundCount, "should have deleted %d tokens", limit)
		assert.Equal(t, totalCount-limit, foundCount, "should have found %d tokens", totalCount-limit)
	}
}

func testFositeSqlStoreTransactionCommitAccessToken(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionRollbackAccessToken(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionCommitRefreshToken(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestCommitRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestCommitRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionRollbackRefreshToken(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollbackRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestRollbackRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionCommitAuthorizeCode(m oauth2.InternalRegistry) func(t *testing.T) {

	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionRollbackAuthorizeCode(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionCommitPKCERequest(m oauth2.InternalRegistry) func(t *testing.T) {

	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

func testFositeSqlStoreTransactionRollbackPKCERequest(m oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

// OpenIdConnect tests can't use the helper functions, due to the signature of GetOpenIdConnectSession being
// different from the other getter methods
func testFositeSqlStoreTransactionCommitOpenIdConnectSession(m oauth2.InternalRegistry) func(t *testing.T) {
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
		assertx.EqualAsJSONExcept(t, &defaultRequest, res, defaultIgnoreKeys)

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

func testFositeSqlStoreTransactionRollbackOpenIdConnectSession(m oauth2.InternalRegistry) func(t *testing.T) {
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

func testFositeStoreSetClientAssertionJWT(m oauth2.InternalRegistry) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=basic setting works", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))

			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			require.NotEqual(t, cmp.NID, gofrsuuid.Nil)
			cmp.NID = gofrsuuid.Nil
			require.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})

		t.Run("case=errors when the JTI is blacklisted", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.ErrorIs(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry), fosite.ErrJTIKnown)
		})

		t.Run("case=deletes expired JTIs", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			expiredJTI := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), expiredJTI))
			newJTI := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), newJTI.JTI, newJTI.Expiry))

			_, err := store.GetClientAssertionJWT(context.Background(), expiredJTI.JTI)
			assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
			cmp, err := store.GetClientAssertionJWT(context.Background(), newJTI.JTI)
			require.NoError(t, err)
			require.NotEqual(t, cmp.NID, gofrsuuid.Nil)
			cmp.NID = gofrsuuid.Nil
			assert.Equal(t, newJTI, cmp)
		})

		t.Run("case=inserts same JTI if expired", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			jti.Expiry = jti.Expiry.Add(2 * time.Minute)
			assert.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))
			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			assert.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})
	}
}

func testFositeStoreClientAssertionJWTValid(m oauth2.InternalRegistry) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=returns valid on unknown JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), uuid.New()))
		})

		t.Run("case=returns invalid on known JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.True(t, errors.Is(store.ClientAssertionJWTValid(context.Background(), jti.JTI), fosite.ErrJTIKnown))
		})

		t.Run("case=returns valid on expired JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.New(), time.Now().Add(-time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), jti.JTI))
		})
	}
}

func testFositeJWTBearerGrantStorage(x oauth2.InternalRegistry) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		grantManager := x.GrantManager()
		keyManager := x.KeyManager()
		grantStorage := x.OAuth2Storage().(rfc7523.RFC7523KeyStorage)

		t.Run("case=associated key added with grant", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
			subject := "bob+" + uuid.New() + "@example.com"
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

			storedKeySet, err := grantStorage.GetPublicKeys(ctx, issuer, subject)
			require.NoError(t, err)
			require.Len(t, storedKeySet.Keys, 0)

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			storedKeySet, err = grantStorage.GetPublicKeys(ctx, issuer, subject)
			require.NoError(t, err)
			assert.Len(t, storedKeySet.Keys, 1)

			storedKey, err := grantStorage.GetPublicKey(ctx, issuer, subject, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, publicKey.KeyID, storedKey.KeyID)
			assert.Equal(t, publicKey.Use, storedKey.Use)
			assert.Equal(t, publicKey.Key, storedKey.Key)

			storedScopes, err := grantStorage.GetPublicKeyScopes(ctx, issuer, subject, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, grant.Scope, storedScopes)

			storedKeySet, err = keyManager.GetKey(ctx, issuer, publicKey.KeyID)
			require.NoError(t, err)
			assert.Equal(t, publicKey.KeyID, storedKeySet.Keys[0].KeyID)
			assert.Equal(t, publicKey.Use, storedKeySet.Keys[0].Use)
			assert.Equal(t, publicKey.Key, storedKeySet.Keys[0].Key)
		})

		t.Run("case=only associated key returns", func(t *testing.T) {
			keySetToNotReturn, err := jwk.GenerateJWK(context.Background(), jose.ES256, uuid.New(), "sig")
			require.NoError(t, err)
			require.NoError(t, keyManager.AddKeySet(context.Background(), uuid.New(), keySetToNotReturn), "adding a random key should not fail")

			issuer := uuid.New()
			subject := "maria+" + uuid.New() + "@example.com"

			keySet1ToReturn, err := jwk.GenerateJWK(context.Background(), jose.ES256, uuid.New(), "sig")
			require.NoError(t, err)
			require.NoError(t, grantManager.CreateGrant(context.Background(), trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: keySet1ToReturn.Keys[0].Public().KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}, keySet1ToReturn.Keys[0].Public()))

			keySet2ToReturn, err := jwk.GenerateJWK(context.Background(), jose.ES256, uuid.New(), "sig")
			require.NoError(t, err)
			require.NoError(t, grantManager.CreateGrant(ctx, trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: keySet2ToReturn.Keys[0].Public().KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}, keySet2ToReturn.Keys[0].Public()))

			storedKeySet, err := grantStorage.GetPublicKeys(context.Background(), issuer, subject)
			require.NoError(t, err)
			require.Len(t, storedKeySet.Keys, 2)

			// Cannot rely on sort order because the created_at timestamps may alias.
			idx1 := slices.IndexFunc(storedKeySet.Keys, func(k jose.JSONWebKey) bool {
				return k.KeyID == keySet1ToReturn.Keys[0].Public().KeyID
			})
			require.GreaterOrEqual(t, idx1, 0)
			idx2 := slices.IndexFunc(storedKeySet.Keys, func(k jose.JSONWebKey) bool {
				return k.KeyID == keySet2ToReturn.Keys[0].Public().KeyID
			})
			require.GreaterOrEqual(t, idx2, 0)

			assert.Equal(t, keySet1ToReturn.Keys[0].Public().KeyID, storedKeySet.Keys[idx1].KeyID)
			assert.Equal(t, keySet1ToReturn.Keys[0].Public().Use, storedKeySet.Keys[idx1].Use)
			assert.Equal(t, keySet1ToReturn.Keys[0].Public().Key, storedKeySet.Keys[idx1].Key)
			assert.Equal(t, keySet2ToReturn.Keys[0].Public().KeyID, storedKeySet.Keys[idx2].KeyID)
			assert.Equal(t, keySet2ToReturn.Keys[0].Public().Use, storedKeySet.Keys[idx2].Use)
			assert.Equal(t, keySet2ToReturn.Keys[0].Public().Key, storedKeySet.Keys[idx2].Key)

			storedKeySet, err = grantStorage.GetPublicKeys(context.Background(), issuer, "non-existing-subject")
			require.NoError(t, err)
			assert.Len(t, storedKeySet.Keys, 0)

			_, err = grantStorage.GetPublicKeyScopes(context.Background(), issuer, "non-existing-subject", keySet2ToReturn.Keys[0].Public().KeyID)
			require.Error(t, err)
		})

		t.Run("case=associated key is deleted, when granted is deleted", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
			subject := "aeneas+" + uuid.New() + "@example.com"
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

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(ctx, issuer, subject, grant.PublicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(ctx, issuer, publicKey.KeyID)
			require.NoError(t, err)

			err = grantManager.DeleteGrant(ctx, grant.ID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(ctx, issuer, subject, publicKey.KeyID)
			assert.Error(t, err)

			_, err = keyManager.GetKey(ctx, issuer, publicKey.KeyID)
			assert.Error(t, err)
		})

		t.Run("case=associated grant is deleted, when key is deleted", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
			subject := "vladimir+" + uuid.New() + "@example.com"
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

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKey(ctx, issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(ctx, issuer, publicKey.KeyID)
			require.NoError(t, err)

			err = keyManager.DeleteKey(ctx, issuer, publicKey.KeyID)
			require.NoError(t, err)

			_, err = keyManager.GetKey(ctx, issuer, publicKey.KeyID)
			assert.Error(t, err)

			_, err = grantManager.GetConcreteGrant(ctx, grant.ID)
			assert.Error(t, err)
		})

		t.Run("case=only returns the key when subject matches", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
			subject := "jagoba+" + uuid.New() + "@example.com"
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

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			// All three get methods should only return the public key when using the valid subject
			_, err = grantStorage.GetPublicKey(ctx, issuer, "any-subject-1", publicKey.KeyID)
			require.Error(t, err)
			_, err = grantStorage.GetPublicKey(ctx, issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKeyScopes(ctx, issuer, "any-subject-2", publicKey.KeyID)
			require.Error(t, err)
			_, err = grantStorage.GetPublicKeyScopes(ctx, issuer, subject, publicKey.KeyID)
			require.NoError(t, err)

			jwks, err := grantStorage.GetPublicKeys(ctx, issuer, "any-subject-3")
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.Empty(t, jwks.Keys)
			jwks, err = grantStorage.GetPublicKeys(ctx, issuer, subject)
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.NotEmpty(t, jwks.Keys)
		})

		t.Run("case=returns the key when any subject is allowed", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
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

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			// All three get methods should always return the public key
			_, err = grantStorage.GetPublicKey(ctx, issuer, "any-subject-1", publicKey.KeyID)
			require.NoError(t, err)

			_, err = grantStorage.GetPublicKeyScopes(ctx, issuer, "any-subject-2", publicKey.KeyID)
			require.NoError(t, err)

			jwks, err := grantStorage.GetPublicKeys(ctx, issuer, "any-subject-3")
			require.NoError(t, err)
			require.NotNil(t, jwks)
			require.NotEmpty(t, jwks.Keys)
		})

		t.Run("case=does not return expired values", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(context.Background(), jose.RS256, uuid.New(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.New()
			grant := trust.Grant{
				ID:              uuid.New(),
				Issuer:          issuer,
				Subject:         "",
				AllowAnySubject: true,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(-1, 0, 0),
			}

			err = grantManager.CreateGrant(ctx, grant, publicKey)
			require.NoError(t, err)

			keys, err := grantStorage.GetPublicKeys(ctx, issuer, "any-subject-3")
			require.NoError(t, err)
			assert.Len(t, keys.Keys, 0)
		})
	}
}

func doTestCommit(m oauth2.InternalRegistry, t *testing.T,
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
	res, err := getFn(context.Background(), signature, oauth2.NewSession("bar"))
	// token should have been created successfully because Commit did not return an error
	require.NoError(t, err)
	assertx.EqualAsJSONExcept(t, &defaultRequest, res, defaultIgnoreKeys)
	// AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

	// testrevoke within a transaction
	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature)
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	_, err = getFn(context.Background(), signature, oauth2.NewSession("bar"))
	// Since commit worked for revoke, we should get an error here.
	require.Error(t, err)
}

func doTestCommitRefresh(m oauth2.InternalRegistry, t *testing.T,
	createFn func(context.Context, string, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(storage.Transactional)
	require.True(t, ok)
	ctx := context.Background()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.New()
	err = createFn(ctx, signature, "", createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	res, err := getFn(context.Background(), signature, oauth2.NewSession("bar"))
	// token should have been created successfully because Commit did not return an error
	require.NoError(t, err)
	assertx.EqualAsJSONExcept(t, &defaultRequest, res, defaultIgnoreKeys)
	// AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

	// testrevoke within a transaction
	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature)
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	_, err = getFn(context.Background(), signature, oauth2.NewSession("bar"))
	// Since commit worked for revoke, we should get an error here.
	require.Error(t, err)
}

func doTestRollback(m oauth2.InternalRegistry, t *testing.T,
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
	_, err = getFn(ctx, signature, oauth2.NewSession("bar"))
	// Since we rolled back above, the token should not exist and getting it should result in an error
	require.Error(t, err)

	// create a new token, revoke it, then rollback the revoke. We should be able to then get it successfully.
	signature2 := uuid.New()
	err = createFn(ctx, signature2, createTestRequest(signature2))
	require.NoError(t, err)
	_, err = getFn(ctx, signature2, oauth2.NewSession("bar"))
	require.NoError(t, err)

	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature2)
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	_, err = getFn(context.Background(), signature2, oauth2.NewSession("bar"))
	require.NoError(t, err)
}

func doTestRollbackRefresh(m oauth2.InternalRegistry, t *testing.T,
	createFn func(context.Context, string, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(storage.Transactional)
	require.True(t, ok)

	ctx := context.Background()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.New()
	err = createFn(ctx, signature, "", createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	ctx = context.Background()
	_, err = getFn(ctx, signature, oauth2.NewSession("bar"))
	// Since we rolled back above, the token should not exist and getting it should result in an error
	require.Error(t, err)

	// create a new token, revoke it, then rollback the revoke. We should be able to then get it successfully.
	signature2 := uuid.New()
	err = createFn(ctx, signature2, "", createTestRequest(signature2))
	require.NoError(t, err)
	_, err = getFn(ctx, signature2, oauth2.NewSession("bar"))
	require.NoError(t, err)

	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature2)
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	_, err = getFn(context.Background(), signature2, oauth2.NewSession("bar"))
	require.NoError(t, err)
}

func createTestRequest(id string) *fosite.Request {
	return &fosite.Request{
		ID:                id,
		RequestedAt:       time.Now().UTC().Round(time.Second),
		Client:            &client.Client{ID: "foobar"},
		RequestedScope:    fosite.Arguments{"fa", "ba"},
		GrantedScope:      fosite.Arguments{"fa", "ba"},
		RequestedAudience: fosite.Arguments{"ad1", "ad2"},
		GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
		Form:              url.Values{"foo": []string{"bar", "baz"}},
		Session:           &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	}
}

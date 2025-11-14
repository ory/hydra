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

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/rfc7523"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence/sql"
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

func newDefaultRequest(t testing.TB, id string) *fosite.Request {
	return &fosite.Request{
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
		Session:           oauth2.NewTestSession(t, "bar"),
	}
}

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

func mockRequestForeignKey(t *testing.T, _ string, x *driver.RegistrySQL) {
	cl := &client.Client{ID: "foobar"}
	if _, err := x.ClientManager().GetClient(t.Context(), cl.ID); errors.Is(err, sqlcon.ErrNoRows) {
		require.NoError(t, x.ClientManager().CreateClient(t.Context(), cl))
	}
}

func testHelperRequestIDMultiples(m *driver.RegistrySQL, _ string) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()
		requestID := uuid.Must(uuid.NewV4()).String()
		mockRequestForeignKey(t, requestID, m)
		cl := &client.Client{ID: "foobar"}

		fositeRequest := &fosite.Request{
			ID:          requestID,
			Client:      cl,
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewTestSession(t, "bar"),
		}

		for range 4 {
			signature := uuid.Must(uuid.NewV4()).String()
			accessSignature := uuid.Must(uuid.NewV4()).String()
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

func testHelperCreateGetDeleteOpenIDConnectSession(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.Must(uuid.NewV4()).String()
		ctx := t.Context()
		_, err := m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewTestSession(t, "bar")})
		assert.NotNil(t, err)

		err = m.CreateOpenIDConnectSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		res, err := m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewTestSession(t, "bar")})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, newDefaultRequest(t, "blank"), res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteOpenIDConnectSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetOpenIDConnectSession(ctx, code, &fosite.Request{Session: oauth2.NewTestSession(t, "bar")})
		assert.NotNil(t, err)
	}
}

func testHelperCreateGetDeleteRefreshTokenSession(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.Must(uuid.NewV4()).String()
		ctx := t.Context()
		_, err := m.GetRefreshTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.NotNil(t, err)

		err = m.CreateRefreshTokenSession(ctx, code, "", newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		res, err := m.GetRefreshTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, newDefaultRequest(t, "blank"), res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteRefreshTokenSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.NotNil(t, err)
	}
}

func testHelperRevokeRefreshToken(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		ctx := t.Context()
		_, err := m.GetRefreshTokenSession(ctx, "1111", oauth2.NewTestSession(t, "bar"))
		assert.Error(t, err)

		reqIdOne := uuid.Must(uuid.NewV4()).String()
		reqIdTwo := uuid.Must(uuid.NewV4()).String()

		mockRequestForeignKey(t, reqIdOne, x)
		mockRequestForeignKey(t, reqIdTwo, x)

		err = m.CreateRefreshTokenSession(ctx, "1111", "", &fosite.Request{
			ID:          reqIdOne,
			Client:      &client.Client{ID: "foobar"},
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewTestSession(t, "user"),
		})
		require.NoError(t, err)

		err = m.CreateRefreshTokenSession(ctx, "1122", "", &fosite.Request{
			ID:          reqIdTwo,
			Client:      &client.Client{ID: "foobar"},
			RequestedAt: time.Now().UTC().Round(time.Second),
			Session:     oauth2.NewTestSession(t, "user"),
		})
		require.NoError(t, err)

		_, err = m.GetRefreshTokenSession(ctx, "1111", oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdOne)
		require.NoError(t, err)

		err = m.RevokeRefreshToken(ctx, reqIdTwo)
		require.NoError(t, err)

		req, err := m.GetRefreshTokenSession(ctx, "1111", oauth2.NewTestSession(t, "bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())

		req, err = m.GetRefreshTokenSession(ctx, "1122", oauth2.NewTestSession(t, "bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperCreateGetDeleteAuthorizeCodes(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		mockRequestForeignKey(t, "blank", x)

		code := uuid.Must(uuid.NewV4()).String()

		ctx := t.Context()
		res, err := m.GetAuthorizeCodeSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.Error(t, err)
		assert.Nil(t, res)

		err = m.CreateAuthorizeCodeSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, newDefaultRequest(t, "blank"), res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.InvalidateAuthorizeCodeSession(ctx, code)
		require.NoError(t, err)

		res, err = m.GetAuthorizeCodeSession(ctx, code, oauth2.NewTestSession(t, "bar"))
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

func testHelperExpiryFields(reg *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := reg.OAuth2Storage()
		t.Parallel()

		mockRequestForeignKey(t, "blank", reg)

		ctx := t.Context()

		s := oauth2.NewTestSession(t, "bar")
		s.SetExpiresAt(fosite.AccessToken, time.Now().Add(time.Hour).Round(time.Minute))
		s.SetExpiresAt(fosite.RefreshToken, time.Now().Add(time.Hour*2).Round(time.Minute))
		s.SetExpiresAt(fosite.AuthorizeCode, time.Now().Add(time.Hour*3).Round(time.Minute))
		request := fosite.Request{
			ID:          uuid.Must(uuid.NewV4()).String(),
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
			id := uuid.Must(uuid.NewV4()).String()
			err := m.CreateAccessTokenSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "access"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", x.SignatureHash(id)).First(&r))

			assert.EqualValues(t, s.GetExpiresAt(fosite.AccessToken).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateRefreshTokenSession", func(t *testing.T) {
			id := uuid.Must(uuid.NewV4()).String()
			err := m.CreateRefreshTokenSession(ctx, id, "", &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "refresh"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.RefreshToken).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateAuthorizeCodeSession", func(t *testing.T) {
			id := uuid.Must(uuid.NewV4()).String()
			err := m.CreateAuthorizeCodeSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "code"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreatePKCERequestSession", func(t *testing.T) {
			id := uuid.Must(uuid.NewV4()).String()
			err := m.CreatePKCERequestSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "pkce"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})

		t.Run("case=CreateOpenIDConnectSession", func(t *testing.T) {
			id := uuid.Must(uuid.NewV4()).String()
			err := m.CreateOpenIDConnectSession(ctx, id, &request)
			require.NoError(t, err)

			r := testHelperExpiryFieldsResult{name: "oidc"}
			require.NoError(t, reg.Persister().Connection(ctx).Select("expires_at").Where("signature = ?", id).First(&r))
			assert.EqualValues(t, s.GetExpiresAt(fosite.AuthorizeCode).UTC(), r.ExpiresAt.UTC())
		})
	}
}

func testHelperNilAccessToken(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		c := &client.Client{ID: uuid.Must(uuid.NewV4()).String()}
		require.NoError(t, x.ClientManager().CreateClient(context.Background(), c))
		err := m.CreateAccessTokenSession(context.Background(), uuid.Must(uuid.NewV4()).String(), &fosite.Request{
			ID:                "",
			RequestedAt:       time.Now().UTC().Round(time.Second),
			Client:            c,
			RequestedScope:    fosite.Arguments{"fa", "ba"},
			GrantedScope:      fosite.Arguments{"fa", "ba"},
			RequestedAudience: fosite.Arguments{"ad1", "ad2"},
			GrantedAudience:   fosite.Arguments{"ad1", "ad2"},
			Form:              url.Values{"foo": []string{"bar", "baz"}},
			Session:           oauth2.NewTestSession(t, "bar"),
		})
		require.NoError(t, err)
	}
}

func testHelperCreateGetDeleteAccessTokenSession(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.Must(uuid.NewV4()).String()
		ctx := t.Context()
		_, err := m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.Error(t, err)

		err = m.CreateAccessTokenSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		res, err := m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, newDefaultRequest(t, "blank"), res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.Error(t, err)
	}
}

func testHelperDeleteAccessTokens(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := t.Context()

		code := uuid.Must(uuid.NewV4()).String()
		err := m.CreateAccessTokenSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)

		err = m.DeleteAccessTokens(ctx, newDefaultRequest(t, "blank").Client.GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperRevokeAccessToken(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()
		ctx := t.Context()

		code := uuid.Must(uuid.NewV4()).String()
		err := m.CreateAccessTokenSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)

		err = m.RevokeAccessToken(ctx, newDefaultRequest(t, "blank").GetID())
		require.NoError(t, err)

		req, err := m.GetAccessTokenSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.Nil(t, req)
		assert.EqualError(t, err, fosite.ErrNotFound.Error())
	}
}

func testHelperRotateRefreshToken(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()

		createTokens := func(t *testing.T, r *fosite.Request) (refreshTokenSession string, accessTokenSession string) {
			refreshTokenSession = fmt.Sprintf("refresh_token_%s", uuid.Must(uuid.NewV4()).String())
			accessTokenSession = fmt.Sprintf("access_token_%s", uuid.Must(uuid.NewV4()).String())
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
			r := newDefaultRequest(t, uuid.Must(uuid.NewV4()).String())
			refreshTokenSession, accessTokenSession := createTokens(t, r)

			err := m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession)
			require.NoError(t, err)

			_, err = m.GetAccessTokenSession(ctx, accessTokenSession, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound, "Token is no longer active because it was refreshed")

			_, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.ErrorIs(t, err, fosite.ErrInactiveToken, "Token is no longer active because it was refreshed")
		})

		t.Run("Rotation works when access token is already pruned", func(t *testing.T) {
			// Test both with and without grace period
			testCases := []struct {
				name              string
				configureGrace    bool
				expectTokenActive bool
			}{
				{
					name:              "with grace period",
					configureGrace:    true,
					expectTokenActive: true,
				},
				{
					name:              "without grace period",
					configureGrace:    false,
					expectTokenActive: false,
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					if tc.configureGrace {
						x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")
					} else {
						x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
					}
					t.Cleanup(func() {
						x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
					})

					m := x.OAuth2Storage()
					r := newDefaultRequest(t, uuid.Must(uuid.NewV4()).String())

					// Create tokens
					refreshTokenSession := fmt.Sprintf("refresh_token_%s", uuid.Must(uuid.NewV4()).String())
					accessTokenSession := fmt.Sprintf("access_token_%s", uuid.Must(uuid.NewV4()).String())

					// Create access token
					err := m.CreateAccessTokenSession(ctx, accessTokenSession, r)
					require.NoError(t, err)

					// Create refresh token linked to the access token
					err = m.CreateRefreshTokenSession(ctx, refreshTokenSession, accessTokenSession, r)
					require.NoError(t, err)

					// Verify tokens were created successfully
					req, err := m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
					require.NoError(t, err)
					require.Equal(t, r.GetID(), req.GetID())

					req, err = m.GetAccessTokenSession(ctx, accessTokenSession, nil)
					require.NoError(t, err)
					require.Equal(t, r.GetID(), req.GetID())

					// Delete the access token (simulating it being pruned)
					err = m.DeleteAccessTokenSession(ctx, accessTokenSession)
					require.NoError(t, err)

					// Verify access token is gone
					_, err = m.GetAccessTokenSession(ctx, accessTokenSession, nil)
					assert.Error(t, err)

					// Rotation should still work even though the access token is gone
					err = m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession)
					require.NoError(t, err)

					// Check refresh token state based on grace period configuration
					req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
					if tc.expectTokenActive {
						assert.NoError(t, err)
						assert.Equal(t, r.GetID(), req.GetID())
					} else {
						assert.ErrorIs(t, err, fosite.ErrInactiveToken, "Token should be inactive when no grace period is configured")
					}
				})
			}
		})

		t.Run("refresh token is valid until the grace period has ended", func(t *testing.T) {
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")

			// By setting this to one hour we ensure that using the refresh token triggers the start of the grace period.
			x.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1h")
			t.Cleanup(func() {
				x.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
			})

			m := x.OAuth2Storage()
			r := newDefaultRequest(t, uuid.Must(uuid.NewV4()).String())
			refreshTokenSession, accessTokenSession1 := createTokens(t, r)
			accessTokenSession2 := fmt.Sprintf("access_token_%s", uuid.Must(uuid.NewV4()).String())
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession2, r))

			// Create a second access token
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			_, err := m.GetAccessTokenSession(ctx, accessTokenSession1, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			_, err = m.GetAccessTokenSession(ctx, accessTokenSession2, nil)
			assert.NoError(t, err, "The second access token is still valid.")

			req, err := m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.NoError(t, err)
			assert.Equal(t, r.GetID(), req.GetID())

			// We only wait a second, meaning that the token is theoretically still within TTL, but since the
			// grace period was issued, the token is still valid.
			time.Sleep(time.Second * 2)
			_, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
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
			r := newDefaultRequest(t, uuid.Must(uuid.NewV4()).String())

			refreshTokenSession, _ := createTokens(t, r)
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
			r := newDefaultRequest(t, uuid.Must(uuid.NewV4()).String())

			refreshTokenSession := fmt.Sprintf("refresh_token_%s", uuid.Must(uuid.NewV4()).String())
			accessTokenSession1 := fmt.Sprintf("access_token_%s", uuid.Must(uuid.NewV4()).String())
			accessTokenSession2 := fmt.Sprintf("access_token_%s", uuid.Must(uuid.NewV4()).String())
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession1, r))
			require.NoError(t, m.CreateAccessTokenSession(ctx, accessTokenSession2, r))

			require.NoError(t, m.CreateRefreshTokenSession(ctx, refreshTokenSession, "", r),
				"precondition failed: could not create refresh token session")

			// ACT
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))
			require.NoError(t, m.RotateRefreshToken(ctx, r.GetID(), refreshTokenSession))

			_, err := m.GetAccessTokenSession(ctx, accessTokenSession1, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			_, err = m.GetAccessTokenSession(ctx, accessTokenSession2, nil)
			assert.ErrorIs(t, err, fosite.ErrNotFound)

			req, err := m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.NoError(t, err)
			assert.Equal(t, r.GetID(), req.GetID())

			time.Sleep(time.Second * 2)

			req, err = m.GetRefreshTokenSession(ctx, refreshTokenSession, nil)
			assert.Error(t, err)
		})
	}
}

func testHelperCreateGetDeletePKCERequestSession(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		m := x.OAuth2Storage()

		code := uuid.Must(uuid.NewV4()).String()
		ctx := t.Context()
		_, err := m.GetPKCERequestSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.NotNil(t, err)

		err = m.CreatePKCERequestSession(ctx, code, newDefaultRequest(t, "blank"))
		require.NoError(t, err)

		res, err := m.GetPKCERequestSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		require.NoError(t, err)
		AssertObjectKeysEqual(t, newDefaultRequest(t, "blank"), res, "RequestedScope", "GrantedScope", "Form", "Session")

		err = m.DeletePKCERequestSession(ctx, code)
		require.NoError(t, err)

		_, err = m.GetPKCERequestSession(ctx, code, oauth2.NewTestSession(t, "bar"))
		assert.NotNil(t, err)
	}
}

func testHelperFlushTokens(x *driver.RegistrySQL, lifespan time.Duration) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &oauth2.Session{}

	return func(t *testing.T) {
		ctx := t.Context()
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

func testHelperFlushTokensWithLimitAndBatchSize(x *driver.RegistrySQL, limit int, batchSize int) func(t *testing.T) {
	m := x.OAuth2Storage()
	ds := &oauth2.Session{}

	return func(t *testing.T) {
		ctx := t.Context()
		var requests []*fosite.Request

		// create five expired requests
		id := uuid.Must(uuid.NewV4()).String()
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

func testFositeSqlStoreTransactionCommitAccessToken(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestCommit(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionRollbackAccessToken(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		{
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().RevokeAccessToken)
			doTestRollback(m, t, m.OAuth2Storage().CreateAccessTokenSession, m.OAuth2Storage().GetAccessTokenSession, m.OAuth2Storage().DeleteAccessTokenSession)
		}
	}
}

func testFositeSqlStoreTransactionCommitRefreshToken(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestCommitRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestCommitRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionRollbackRefreshToken(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollbackRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().RevokeRefreshToken)
		doTestRollbackRefresh(m, t, m.OAuth2Storage().CreateRefreshTokenSession, m.OAuth2Storage().GetRefreshTokenSession, m.OAuth2Storage().DeleteRefreshTokenSession)
	}
}

func testFositeSqlStoreTransactionCommitAuthorizeCode(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionRollbackAuthorizeCode(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreateAuthorizeCodeSession, m.OAuth2Storage().GetAuthorizeCodeSession, m.OAuth2Storage().InvalidateAuthorizeCodeSession)
	}
}

func testFositeSqlStoreTransactionCommitPKCERequest(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestCommit(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

func testFositeSqlStoreTransactionRollbackPKCERequest(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		doTestRollback(m, t, m.OAuth2Storage().CreatePKCERequestSession, m.OAuth2Storage().GetPKCERequestSession, m.OAuth2Storage().DeletePKCERequestSession)
	}
}

// OpenIdConnect tests can't use the helper functions, due to the signature of GetOpenIdConnectSession being
// different from the other getter methods
func testFositeSqlStoreTransactionCommitOpenIdConnectSession(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
		require.True(t, ok)
		ctx := t.Context()
		ctx, err := txnStore.BeginTX(ctx)
		require.NoError(t, err)
		signature := uuid.Must(uuid.NewV4()).String()
		testRequest := createTestRequest(signature)
		err = m.OAuth2Storage().CreateOpenIDConnectSession(ctx, signature, testRequest)
		require.NoError(t, err)
		err = txnStore.Commit(ctx)
		require.NoError(t, err)

		// Require a new context, since the old one contains the transaction.
		res, err := m.OAuth2Storage().GetOpenIDConnectSession(context.Background(), signature, testRequest)
		// session should have been created successfully because Commit did not return an error
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, newDefaultRequest(t, "blank"), res, defaultIgnoreKeys)

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

func testFositeSqlStoreTransactionRollbackOpenIdConnectSession(m *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
		require.True(t, ok)
		ctx := t.Context()
		ctx, err := txnStore.BeginTX(ctx)
		require.NoError(t, err)

		signature := uuid.Must(uuid.NewV4()).String()
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
		signature2 := uuid.Must(uuid.NewV4()).String()
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

func testFositeStoreSetClientAssertionJWT(m *driver.RegistrySQL) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=basic setting works", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))

			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			require.NotEqual(t, cmp.NID, uuid.Nil)
			cmp.NID = uuid.Nil
			require.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})

		t.Run("case=errors when the JTI is blacklisted", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.ErrorIs(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry), fosite.ErrJTIKnown)
		})

		t.Run("case=deletes expired JTIs", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			expiredJTI := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), expiredJTI))
			newJTI := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWT(context.Background(), newJTI.JTI, newJTI.Expiry))

			_, err := store.GetClientAssertionJWT(context.Background(), expiredJTI.JTI)
			assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
			cmp, err := store.GetClientAssertionJWT(context.Background(), newJTI.JTI)
			require.NoError(t, err)
			require.NotEqual(t, cmp.NID, uuid.Nil)
			cmp.NID = uuid.Nil
			assert.Equal(t, newJTI, cmp)
		})

		t.Run("case=inserts same JTI if expired", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(-time.Minute))
			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			jti.Expiry = jti.Expiry.Add(2 * time.Minute)
			assert.NoError(t, store.SetClientAssertionJWT(context.Background(), jti.JTI, jti.Expiry))
			cmp, err := store.GetClientAssertionJWT(context.Background(), jti.JTI)
			assert.NoError(t, err)
			assert.Equal(t, jti, cmp)
		})
	}
}

func testFositeStoreClientAssertionJWTValid(m *driver.RegistrySQL) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("case=returns valid on unknown JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), uuid.Must(uuid.NewV4()).String()))
		})

		t.Run("case=returns invalid on known JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.True(t, errors.Is(store.ClientAssertionJWTValid(context.Background(), jti.JTI), fosite.ErrJTIKnown))
		})

		t.Run("case=returns valid on expired JTI", func(t *testing.T) {
			store, ok := m.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(-time.Minute))

			require.NoError(t, store.SetClientAssertionJWTRaw(context.Background(), jti))

			assert.NoError(t, store.ClientAssertionJWTValid(context.Background(), jti.JTI))
		})
	}
}

func testFositeJWTBearerGrantStorage(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()
		grantManager := x.GrantManager()
		keyManager := x.KeyManager()
		grantStorage := x.OAuth2Storage().(rfc7523.RFC7523KeyStorage)

		t.Run("case=associated key added with grant", func(t *testing.T) {
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			subject := "bob+" + uuid.Must(uuid.NewV4()).String() + "@example.com"
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
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

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

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
			keySetToNotReturn, err := jwk.GenerateJWK(jose.ES256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)
			require.NoError(t, keyManager.AddKeySet(context.Background(), uuid.Must(uuid.NewV4()).String(), keySetToNotReturn), "adding a random key should not fail")

			issuer := uuid.Must(uuid.NewV4()).String()
			subject := "maria+" + uuid.Must(uuid.NewV4()).String() + "@example.com"

			keySet1ToReturn, err := jwk.GenerateJWK(jose.ES256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)
			require.NoError(t, grantManager.CreateGrant(t.Context(), trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: keySet1ToReturn.Keys[0].Public().KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}, keySet1ToReturn.Keys[0].Public()))

			keySet2ToReturn, err := jwk.GenerateJWK(jose.ES256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)
			require.NoError(t, grantManager.CreateGrant(ctx, trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
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
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			subject := "aeneas+" + uuid.Must(uuid.NewV4()).String() + "@example.com"
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

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
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			subject := "vladimir+" + uuid.Must(uuid.NewV4()).String() + "@example.com"
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

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
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			subject := "jagoba+" + uuid.Must(uuid.NewV4()).String() + "@example.com"
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         subject,
				AllowAnySubject: false,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

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
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         "",
				AllowAnySubject: true,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(1, 0, 0),
			}

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

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
			keySet, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
			require.NoError(t, err)

			publicKey := keySet.Keys[0].Public()
			issuer := uuid.Must(uuid.NewV4()).String()
			grant := trust.Grant{
				ID:              uuid.Must(uuid.NewV4()),
				Issuer:          issuer,
				Subject:         "",
				AllowAnySubject: true,
				Scope:           []string{"openid", "offline"},
				PublicKey:       trust.PublicKey{Set: issuer, KeyID: publicKey.KeyID},
				CreatedAt:       time.Now().UTC().Round(time.Second),
				ExpiresAt:       time.Now().UTC().Round(time.Second).AddDate(-1, 0, 0),
			}

			require.NoError(t, grantManager.CreateGrant(ctx, grant, publicKey))

			keys, err := grantStorage.GetPublicKeys(ctx, issuer, "any-subject-3")
			require.NoError(t, err)
			assert.Len(t, keys.Keys, 0)
		})
	}
}

func doTestCommit(m *driver.RegistrySQL, t *testing.T,
	createFn func(context.Context, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
	require.True(t, ok)
	ctx := t.Context()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature, createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	res, err := getFn(context.Background(), signature, oauth2.NewTestSession(t, "bar"))
	// token should have been created successfully because Commit did not return an error
	require.NoError(t, err)
	assertx.EqualAsJSONExcept(t, newDefaultRequest(t, "blank"), res, defaultIgnoreKeys)
	// AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

	// testrevoke within a transaction
	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature)
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	_, err = getFn(context.Background(), signature, oauth2.NewTestSession(t, "bar"))
	// Since commit worked for revoke, we should get an error here.
	require.Error(t, err)
}

func doTestCommitRefresh(m *driver.RegistrySQL, t *testing.T,
	createFn func(context.Context, string, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
	require.True(t, ok)
	ctx := t.Context()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature, "", createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	res, err := getFn(context.Background(), signature, oauth2.NewTestSession(t, "bar"))
	// token should have been created successfully because Commit did not return an error
	require.NoError(t, err)
	assertx.EqualAsJSONExcept(t, newDefaultRequest(t, "blank"), res, defaultIgnoreKeys)
	// AssertObjectKeysEqual(t, &defaultRequest, res, "RequestedScope", "GrantedScope", "Form", "Session")

	// testrevoke within a transaction
	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature)
	require.NoError(t, err)
	err = txnStore.Commit(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	_, err = getFn(context.Background(), signature, oauth2.NewTestSession(t, "bar"))
	// Since commit worked for revoke, we should get an error here.
	require.Error(t, err)
}

func doTestRollback(m *driver.RegistrySQL, t *testing.T,
	createFn func(context.Context, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
	require.True(t, ok)

	ctx := t.Context()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature, createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	ctx = context.Background()
	_, err = getFn(ctx, signature, oauth2.NewTestSession(t, "bar"))
	// Since we rolled back above, the token should not exist and getting it should result in an error
	require.Error(t, err)

	// create a new token, revoke it, then rollback the revoke. We should be able to then get it successfully.
	signature2 := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature2, createTestRequest(signature2))
	require.NoError(t, err)
	_, err = getFn(ctx, signature2, oauth2.NewTestSession(t, "bar"))
	require.NoError(t, err)

	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature2)
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	_, err = getFn(context.Background(), signature2, oauth2.NewTestSession(t, "bar"))
	require.NoError(t, err)
}

func doTestRollbackRefresh(m *driver.RegistrySQL, t *testing.T,
	createFn func(context.Context, string, string, fosite.Requester) error,
	getFn func(context.Context, string, fosite.Session) (fosite.Requester, error),
	revokeFn func(context.Context, string) error,
) {
	txnStore, ok := m.OAuth2Storage().(fosite.Transactional)
	require.True(t, ok)

	ctx := t.Context()
	ctx, err := txnStore.BeginTX(ctx)
	require.NoError(t, err)
	signature := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature, "", createTestRequest(signature))
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	// Require a new context, since the old one contains the transaction.
	ctx = context.Background()
	_, err = getFn(ctx, signature, oauth2.NewTestSession(t, "bar"))
	// Since we rolled back above, the token should not exist and getting it should result in an error
	require.Error(t, err)

	// create a new token, revoke it, then rollback the revoke. We should be able to then get it successfully.
	signature2 := uuid.Must(uuid.NewV4()).String()
	err = createFn(ctx, signature2, "", createTestRequest(signature2))
	require.NoError(t, err)
	_, err = getFn(ctx, signature2, oauth2.NewTestSession(t, "bar"))
	require.NoError(t, err)

	ctx, err = txnStore.BeginTX(context.Background())
	require.NoError(t, err)
	err = revokeFn(ctx, signature2)
	require.NoError(t, err)
	err = txnStore.Rollback(ctx)
	require.NoError(t, err)

	_, err = getFn(context.Background(), signature2, oauth2.NewTestSession(t, "bar"))
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

func testHelperRefreshTokenExpiryUpdate(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()

		// Create client
		cl := &client.Client{ID: "refresh-expiry-client"}
		require.NoError(t, x.ClientManager().CreateClient(ctx, cl))

		// Create a request with a long expiry
		initialRequest := fosite.Request{
			ID:          uuid.Must(uuid.NewV4()).String(),
			RequestedAt: time.Now().UTC().Round(time.Second),
			Client:      cl,
			Session:     oauth2.NewTestSession(t, "sub"),
		}

		// Set a long expiry time (24 hours)
		initialExpiry := time.Now().Add(24 * time.Hour)
		initialRequest.Session.SetExpiresAt(fosite.RefreshToken, initialExpiry)

		t.Run("regular rotation", func(t *testing.T) {
			// Create original refresh token
			regularSignature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, x.OAuth2Storage().CreateRefreshTokenSession(ctx, regularSignature, "", &initialRequest))

			// Verify initial expiry is set correctly
			originalToken, err := x.OAuth2Storage().GetRefreshTokenSession(ctx, regularSignature, oauth2.NewTestSession(t, "sub"))
			require.NoError(t, err)
			require.Equal(t, initialExpiry.Unix(), originalToken.GetSession().GetExpiresAt(fosite.RefreshToken).Unix())

			// Set up a connection to directly query the database
			var actualExpiresAt time.Time
			require.NoError(t, x.Persister().Connection(ctx).RawQuery("SELECT expires_at FROM hydra_oauth2_refresh WHERE signature=?", regularSignature).First(&actualExpiresAt))
			require.Equal(t, initialExpiry.UTC().Round(time.Second), actualExpiresAt.UTC().Round(time.Second))

			// Rotate the token
			err = x.OAuth2Storage().RotateRefreshToken(ctx, initialRequest.ID, regularSignature)
			require.NoError(t, err)

			// Check that the original token's expiry was updated to be closer to now
			var revokedData struct {
				ExpiresAt time.Time `db:"expires_at"`
				Active    bool      `db:"active"`
			}
			require.NoError(t, x.Persister().Connection(ctx).RawQuery("SELECT expires_at, active FROM hydra_oauth2_refresh WHERE signature=?", regularSignature).First(&revokedData))

			// Verify the token is now inactive
			require.False(t, revokedData.Active)

			// Verify the expiry is updated to be closer to now than the original expiry
			require.True(t, revokedData.ExpiresAt.Before(initialExpiry), "Expiry should be updated to be sooner than original")
			require.True(t, revokedData.ExpiresAt.After(time.Now()), "Expiry should still be in the future")
			require.True(t, time.Until(revokedData.ExpiresAt) < time.Until(initialExpiry), "New expiry should be closer to now than original expiry")

			t.Logf("Original expiry: %v, Updated expiry: %v, Now: %v", initialExpiry, revokedData.ExpiresAt, time.Now())
		})

		t.Run("graceful rotation", func(t *testing.T) {
			// Create refresh token for graceful rotation
			gracefulSignature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, x.OAuth2Storage().CreateRefreshTokenSession(ctx, gracefulSignature, "", &initialRequest))

			// Set config to graceful rotation
			oldPeriod := x.Config().GracefulRefreshTokenRotation(ctx).Period
			t.Cleanup(func() {
				x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, oldPeriod)
				x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGraceReuseCount, 0)
			})
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, time.Minute*30)
			x.Config().MustSet(ctx, config.KeyRefreshTokenRotationGraceReuseCount, 3)

			// Record time before rotation
			beforeRotation := time.Now().UTC().Add(-time.Second) // Ensure we have a different timestamp for first_used_at

			// Rotate the token
			err := x.OAuth2Storage().RotateRefreshToken(ctx, initialRequest.ID, gracefulSignature)
			require.NoError(t, err)

			// Check the token's expiry and status
			var rotatedData struct {
				ExpiresAt   time.Time       `db:"expires_at"`
				Active      bool            `db:"active"`
				FirstUsedAt sqlxx.NullTime  `db:"first_used_at"`
				UsedTimes   sqlxx.NullInt64 `db:"used_times"`
			}
			require.NoError(t, x.Persister().Connection(ctx).RawQuery("SELECT expires_at, active, first_used_at, used_times FROM hydra_oauth2_refresh WHERE signature=?", gracefulSignature).First(&rotatedData))

			// Token is used
			require.False(t, rotatedData.Active)

			// Verify first_used_at is set and reasonable
			assert.True(t, time.Time(rotatedData.FirstUsedAt).After(beforeRotation) || time.Time(rotatedData.FirstUsedAt).Equal(beforeRotation), "%s should be after or equal to %s", time.Time(rotatedData.FirstUsedAt), beforeRotation)

			now := time.Now().UTC().Add(time.Second)
			assert.True(t, time.Time(rotatedData.FirstUsedAt).Before(now) || time.Time(rotatedData.FirstUsedAt).Equal(now), "%s should be before or equal to %s", time.Time(rotatedData.FirstUsedAt), now)

			// Verify used_times was incremented
			assert.True(t, rotatedData.UsedTimes.Valid)
			assert.Equal(t, int64(1), rotatedData.UsedTimes.Int)

			// Verify the expiry is updated and is in the future
			assert.True(t, rotatedData.ExpiresAt.Before(initialExpiry), "Expiry should be updated to be sooner than original")
			assert.True(t, rotatedData.ExpiresAt.After(time.Now().UTC()), "Expiry should still be in the future")
			assert.True(t, time.Until(rotatedData.ExpiresAt) < time.Until(initialExpiry), "New expiry should be closer to now than original expiry")

			t.Logf("Original expiry: %v, Updated expiry: %v, Now: %v", initialExpiry, rotatedData.ExpiresAt, time.Now())
		})
	}
}

func testHelperAuthorizeCodeInvalidation(x *driver.RegistrySQL) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()

		// Create client
		cl := &client.Client{ID: "auth-code-client"}
		require.NoError(t, x.ClientManager().CreateClient(ctx, cl))

		// Create a request with a long expiry
		initialRequest := fosite.Request{
			ID:          uuid.Must(uuid.NewV4()).String(),
			RequestedAt: time.Now().UTC().Round(time.Second),
			Client:      cl,
			Session:     oauth2.NewTestSession(t, "sub"),
		}

		// Set a long expiry time (1 hour)
		initialExpiry := time.Now().Add(1 * time.Hour)
		initialRequest.Session.SetExpiresAt(fosite.AuthorizeCode, initialExpiry)

		// Create authorize code session
		authCodeSignature := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, x.OAuth2Storage().CreateAuthorizeCodeSession(ctx, authCodeSignature, &initialRequest))

		// Verify initial state
		originalCode, err := x.OAuth2Storage().GetAuthorizeCodeSession(ctx, authCodeSignature, oauth2.NewTestSession(t, "sub"))
		require.NoError(t, err)
		require.Equal(t, initialExpiry.Unix(), originalCode.GetSession().GetExpiresAt(fosite.AuthorizeCode).Unix())

		// Check database directly
		var codeData struct {
			ExpiresAt time.Time `db:"expires_at"`
			Active    bool      `db:"active"`
		}
		require.NoError(t, x.Persister().Connection(ctx).RawQuery(
			"SELECT expires_at, active FROM hydra_oauth2_code WHERE signature=?",
			authCodeSignature).First(&codeData))
		require.Equal(t, initialExpiry.UTC().Round(time.Second), codeData.ExpiresAt.UTC().Round(time.Second))
		require.True(t, codeData.Active)

		// Invalidate the code
		err = x.OAuth2Storage().InvalidateAuthorizeCodeSession(ctx, authCodeSignature)
		require.NoError(t, err)

		// Check that the code was invalidated but is still retrievable
		invalidatedCode, err := x.OAuth2Storage().GetAuthorizeCodeSession(ctx, authCodeSignature, oauth2.NewTestSession(t, "sub"))
		require.Error(t, err)
		require.ErrorIs(t, err, fosite.ErrInvalidatedAuthorizeCode)
		require.NotNil(t, invalidatedCode) // Should still be retrievable

		// Verify database state after invalidation
		var invalidatedData struct {
			ExpiresAt time.Time `db:"expires_at"`
			Active    bool      `db:"active"`
		}
		require.NoError(t, x.Persister().Connection(ctx).RawQuery(
			"SELECT expires_at, active FROM hydra_oauth2_code WHERE signature=?",
			authCodeSignature).First(&invalidatedData))

		// Verify the code is now inactive
		require.False(t, invalidatedData.Active)

		// Verify the expiry is updated to be closer to now than the original expiry
		require.True(t, invalidatedData.ExpiresAt.Before(initialExpiry),
			"Expiry should be updated to be sooner than original")
		require.True(t, invalidatedData.ExpiresAt.After(time.Now()),
			"Expiry should still be in the future")
		require.True(t, time.Until(invalidatedData.ExpiresAt) < time.Until(initialExpiry),
			"New expiry should be closer to now than original expiry")

		t.Logf("Original expiry: %v, Updated expiry: %v, Now: %v",
			initialExpiry, invalidatedData.ExpiresAt, time.Now())
	}
}

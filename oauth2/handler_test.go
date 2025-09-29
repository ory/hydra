// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/driver"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/snapshotx"
)

var lifespan = time.Hour

func TestHandlerDeleteHandler(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(config.KeyIssuerURL, "http://hydra.localhost")))

	cm := reg.ClientManager()
	store := reg.OAuth2Storage()

	h := oauth2.NewHandler(reg)

	deleteRequest := &fosite.Request{
		ID:             "del-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{ID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	}
	require.NoError(t, cm.CreateClient(ctx, deleteRequest.Client.(*client.Client)))
	require.NoError(t, store.CreateAccessTokenSession(ctx, deleteRequest.ID, deleteRequest))

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	r := httprouterx.NewRouterAdminWithPrefix(metrics)
	h.SetPublicRoutes(r.ToPublic(), func(h http.Handler) http.Handler { return h })
	h.SetAdminRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	c := hydra.NewAPIClient(hydra.NewConfiguration())
	c.GetConfig().Servers = hydra.ServerConfigurations{{URL: ts.URL}}

	_, err := c.
		OAuth2API.DeleteOAuth2Token(ctx).
		ClientId("foobar").Execute()
	require.NoError(t, err)

	ds := new(oauth2.Session)
	_, err = store.GetAccessTokenSession(ctx, "del-1", ds)
	require.Error(t, err, "not_found")
}

func TestUserinfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	op := NewMockOAuth2Provider(ctrl)
	t.Cleanup(ctrl.Finish)

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyScopeStrategy:    "",
		config.KeyAuthCodeLifespan: lifespan,
		config.KeyIssuerURL:        "http://hydra.localhost",
	})), driver.WithOAuth2Provider(op))
	testhelpers.MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)

	h := oauth2.NewHandler(reg)

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	router := httprouterx.NewRouterAdminWithPrefix(metrics)
	h.SetPublicRoutes(router.ToPublic(), func(h http.Handler) http.Handler { return h })
	h.SetAdminRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, tc := range []struct {
		setup                func(t *testing.T)
		checkForSuccess      func(t *testing.T, body []byte)
		checkForUnauthorized func(t *testing.T, body []byte, header http.Header)
		expectStatusCode     int
	}{
		{
			setup: func(t *testing.T) {
				op.EXPECT().IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).Return(fosite.AccessToken, nil, errors.New("asdf"))
			},
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					Return(fosite.RefreshToken, nil, nil)
			},
			checkForUnauthorized: func(t *testing.T, body []byte, headers http.Header) {
				assert.True(t, headers.Get("WWW-Authenticate") == `Bearer error="invalid_token",error_description="Only access tokens are allowed in the authorization header."`, "%s", headers)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					Return(fosite.AccessToken, nil, fosite.ErrRequestUnauthorized)
			},
			checkForUnauthorized: func(t *testing.T, body []byte, headers http.Header) {
				assert.True(t, headers.Get("WWW-Authenticate") == `Bearer error="request_unauthorized",error_description="The request could not be authorized. Check that you provided valid credentials in the right format."`, "%s", headers)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, _ fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session := &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									ID: "foobar",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				bodyString := string(body)
				assert.True(t, strings.Contains(bodyString, `"sub":"alice"`), "%s", body)
				assert.True(t, strings.Contains(bodyString, `"aud":["foobar"]`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, _ fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session := &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject:  "another-alice",
									Audience: []string{"something-else"},
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									ID: "foobar",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				bodyString := string(body)
				assert.False(t, strings.Contains(bodyString, `"sub":"alice"`), "%s", body)
				assert.True(t, strings.Contains(bodyString, `"sub":"another-alice"`), "%s", body)
				assert.True(t, strings.Contains(bodyString, `"aud":["something-else","foobar"]`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, _ fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session := &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject:  "alice",
									Audience: []string{"foobar"},
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									ID:                        "foobar",
									UserinfoSignedResponseAlg: "none",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				bodyString := string(body)
				assert.True(t, strings.Contains(bodyString, `"sub":"alice"`), "%s", body)
				assert.True(t, strings.Contains(bodyString, `"aud":["foobar"]`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, _ fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session := &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									UserinfoSignedResponseAlg: "asdfasdf",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, _ fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session := &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									ID:                        "foobar-client",
									UserinfoSignedResponseAlg: "RS256",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				claims, err := jwt.Parse(string(body), func(token *jwt.Token) (interface{}, error) {
					keys, err := reg.KeyManager().GetKeySet(t.Context(), x.OpenIDConnectKeyName)
					require.NoError(t, err)
					t.Logf("%+v", keys)
					key, _ := jwk.FindPublicKey(keys)
					return key.Key, nil
				})
				require.NoError(t, err)
				assert.EqualValues(t, "alice", claims.Claims["sub"])
				assert.EqualValues(t, []interface{}{"foobar-client"}, claims.Claims["aud"], "%#v", claims.Claims)
				assert.NotEmpty(t, claims.Claims["jti"])
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			tc.setup(t)

			req, err := http.NewRequest("GET", ts.URL+"/userinfo", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer access-token")
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck
			require.EqualValues(t, tc.expectStatusCode, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tc.expectStatusCode == http.StatusOK {
				tc.checkForSuccess(t, body)
			} else if tc.expectStatusCode == http.StatusUnauthorized {
				tc.checkForUnauthorized(t, body, resp.Header)
			}
		})
	}
}

func TestHandlerWellKnown(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyScopeStrategy:                 "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY",
		config.KeyIssuerURL:                     "http://hydra.localhost",
		config.KeySubjectTypesSupported:         []string{"pairwise", "public"},
		config.KeyOIDCDiscoverySupportedClaims:  []string{"sub"},
		config.KeyOAuth2ClientRegistrationURL:   "http://client-register/registration",
		config.KeyOIDCDiscoveryUserinfoEndpoint: "/userinfo",
	})))
	t.Run(fmt.Sprintf("hsm_enabled=%v", reg.Config().HSMEnabled()), func(t *testing.T) {
		testhelpers.MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)

		h := oauth2.NewHandler(reg)

		metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
		r := httprouterx.NewRouterAdminWithPrefix(metrics)
		h.SetPublicRoutes(r.ToPublic(), func(h http.Handler) http.Handler { return h })
		h.SetAdminRoutes(r)
		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/.well-known/openid-configuration")
		require.NoError(t, err)
		defer res.Body.Close() //nolint:errcheck

		var wellKnownResp hydra.OidcConfiguration
		err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
		require.NoError(t, err, "problem decoding wellknown json response: %+v", err)

		snapshotOpts := []snapshotx.Opt{}
		if reg.Config().HSMEnabled() {
			// The signing algorithm is not stable in the HSM tests, because the key is kept
			// in the HSM and persists across test runs.
			snapshotOpts = append(snapshotOpts, snapshotx.ExceptPaths(
				"id_token_signed_response_alg",
				"id_token_signing_alg_values_supported",
				"userinfo_signed_response_alg",
				"userinfo_signing_alg_values_supported",
			))
		}
		snapshotx.SnapshotT(t, wellKnownResp, snapshotOpts...)
	})
}

func TestHandlerOauthAuthorizationServer(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyScopeStrategy:                 "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY",
		config.KeyIssuerURL:                     "http://hydra.localhost",
		config.KeySubjectTypesSupported:         []string{"pairwise", "public"},
		config.KeyOIDCDiscoverySupportedClaims:  []string{"sub"},
		config.KeyOAuth2ClientRegistrationURL:   "http://client-register/registration",
		config.KeyOIDCDiscoveryUserinfoEndpoint: "/userinfo",
	})))
	t.Run(fmt.Sprintf("hsm_enabled=%v", reg.Config().HSMEnabled()), func(t *testing.T) {
		testhelpers.MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)

		h := oauth2.NewHandler(reg)

		metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
		r := httprouterx.NewRouterAdminWithPrefix(metrics)
		h.SetPublicRoutes(r.ToPublic(), func(h http.Handler) http.Handler { return h })
		h.SetAdminRoutes(r)
		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/.well-known/oauth-authorization-server")
		require.NoError(t, err)
		defer res.Body.Close() //nolint:errcheck

		var wellKnownResp hydra.OidcConfiguration
		err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
		require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
		snapshotOpts := []snapshotx.Opt{}
		if reg.Config().HSMEnabled() {
			// The signing algorithm is not stable in the HSM tests, because the key is kept
			// in the HSM and persists across test runs.
			snapshotOpts = append(snapshotOpts, snapshotx.ExceptPaths(
				"id_token_signed_response_alg",
				"id_token_signing_alg_values_supported",
				"userinfo_signed_response_alg",
				"userinfo_signing_alg_values_supported",
			))
		}
		snapshotx.SnapshotT(t, wellKnownResp, snapshotOpts...)
	})
}

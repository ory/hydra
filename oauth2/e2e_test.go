// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/uuidx"
)

func TestAuthCodeFlowE2E(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyAccessTokenStrategy:  "opaque",
		config.KeyRefreshTokenHook:     "",
		config.KeyLoginURL:             x.LoginURL,
		config.KeyConsentURL:           x.ConsentURL,
		config.KeyAccessTokenLifespan:  10 * time.Minute, // allow to debug
		config.KeyRefreshTokenLifespan: 20 * time.Minute, // allow to debug
		config.KeyScopeStrategy:        "exact",
		config.KeyIssuerURL:            "https://hydra.ory",
	})))

	jwk.EnsureAsymmetricKeypairExists(t, reg, string(jose.ES256), x.OpenIDConnectKeyName)
	jwk.EnsureAsymmetricKeypairExists(t, reg, string(jose.ES256), x.OAuth2JWTKeyName)

	publicTS, adminTS := testhelpers.NewConfigurableOAuth2Server(t.Context(), t, reg)
	publicClient := hydra.NewAPIClient(hydra.NewConfiguration())
	publicClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: publicTS.URL}}
	adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
	adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	t.Run("auth code flow", func(t *testing.T) {
		t.Run("rejects invalid audience", func(t *testing.T) {
			cl := x.NewEmptyJarClient(t)
			cl.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
			_, conf := newOAuth2Client(t, reg, x.ClientCallbackURL)
			loc := x.GetExpectRedirect(t, cl, conf.AuthCodeURL(uuidx.NewV4().String(), oauth2.SetAuthURLParam("audience", "invalid-audience")))
			require.Equal(t, x.ClientCallbackURL, fmt.Sprintf("%s://%s%s", loc.Scheme, loc.Host, loc.Path))
			assert.Equal(t, "invalid_request", loc.Query().Get("error"))
			assert.Contains(t, loc.Query().Get("error_description"), "Requested audience 'invalid-audience' has not been whitelisted by the OAuth 2.0 Client.")
		})

		for _, accessTokenStrategy := range []string{"opaque", "jwt"} {
			t.Run("strategy="+accessTokenStrategy, func(t *testing.T) {
				cl, conf := newOAuth2Client(t, reg, x.ClientCallbackURL, func(c *client.Client) {
					c.AccessTokenStrategy = accessTokenStrategy
					c.Audience = []string{"audience-1", "audience-2"}
					c.ID = "64f78bf1-f388-4eeb-9fee-e7207226c6be-" + accessTokenStrategy
				})
				sub := "c6a8ee1c-e0c4-404c-bba7-6a5b8702a2e9"

				t.Run("access and id tokens with extra claims", func(t *testing.T) {
					token := x.PerformAuthCodeFlow(t.Context(), t, nil, conf, adminClient, func(t *testing.T, req *hydra.OAuth2LoginRequest) hydra.AcceptOAuth2LoginRequest {
						snapshotx.SnapshotT(t, req,
							snapshotx.ExceptPaths("challenge", "client.created_at", "client.updated_at", "session_id", "request_url"),
							snapshotx.WithName("login_request"))
						return hydra.AcceptOAuth2LoginRequest{
							Amr:     []string{"amr1", "amr2"},
							Acr:     pointerx.Ptr("acr-value"),
							Subject: sub,
						}
					}, func(t *testing.T, req *hydra.OAuth2ConsentRequest) hydra.AcceptOAuth2ConsentRequest {
						snapshotx.SnapshotT(t, req,
							snapshotx.ExceptPaths("challenge", "client.created_at", "client.updated_at", "consent_request_id", "login_challenge", "login_session_id", "request_url"),
							snapshotx.WithName("consent_request"))
						return hydra.AcceptOAuth2ConsentRequest{
							GrantScope: []string{"openid"},
							Session: &hydra.AcceptOAuth2ConsentRequestSession{
								AccessToken: map[string]any{"key_access": "extra access token value"},
								IdToken:     map[string]any{"key_id": "extra id token value"},
							},
						}
					})

					// check access token
					introspected := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
					require.True(t, introspected.Get("active").Bool())
					testhelpers.AssertAccessToken(t, introspected, sub, cl.ID)
					assert.Equal(t, "extra access token value", introspected.Get("ext.key_access").Str)

					if accessTokenStrategy == "jwt" {
						dec := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))
						testhelpers.AssertAccessToken(t, dec, sub, cl.ID)
						assert.Equal(t, "extra access token value", dec.Get("ext.key_access").Str)
					} else {
						assert.Len(t, strings.Split(token.AccessToken, "."), 2)
					}

					idToken := testhelpers.DecodeIDToken(t, token)
					testhelpers.AssertIDToken(t, idToken, sub, cl.ID)
					assert.Equal(t, "extra id token value", idToken.Get("key_id").Str)
					assert.JSONEq(t, `["amr1", "amr2"]`, idToken.Get("amr").Raw)
					assert.Equal(t, "acr-value", idToken.Get("acr").Str)
				})

				t.Run("refreshed access and id tokens with extra claims", func(t *testing.T) {
					token := x.PerformAuthCodeFlow(t.Context(), t, nil, conf, adminClient, func(*testing.T, *hydra.OAuth2LoginRequest) hydra.AcceptOAuth2LoginRequest {
						return hydra.AcceptOAuth2LoginRequest{Subject: sub}
					}, func(*testing.T, *hydra.OAuth2ConsentRequest) hydra.AcceptOAuth2ConsentRequest {
						return hydra.AcceptOAuth2ConsentRequest{
							GrantScope: []string{"openid", "offline"},
							Session: &hydra.AcceptOAuth2ConsentRequestSession{
								AccessToken: map[string]any{"key_access": "extra access token value"},
								IdToken:     map[string]any{"key_id": "extra id token value"},
							},
						}
					})

					token.Expiry = time.Now().Add(-time.Hour)
					refreshed, err := conf.TokenSource(t.Context(), token).Token()
					require.NoError(t, err)
					require.NotEqual(t, token.AccessToken, refreshed.AccessToken)
					require.NotEqual(t, token.RefreshToken, refreshed.RefreshToken)
					require.NotEqual(t, token.Extra("id_token"), refreshed.Extra("id_token"))

					// check access token
					introspected := testhelpers.IntrospectToken(t, refreshed.AccessToken, adminTS)
					require.True(t, introspected.Get("active").Bool())
					testhelpers.AssertAccessToken(t, introspected, sub, cl.ID)
					assert.Equal(t, "extra access token value", introspected.Get("ext.key_access").Str)

					if accessTokenStrategy == "jwt" {
						dec := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, refreshed.AccessToken))
						testhelpers.AssertAccessToken(t, dec, sub, cl.ID)
						assert.Equal(t, "extra access token value", dec.Get("ext.key_access").Str)
					} else {
						assert.Len(t, strings.Split(refreshed.AccessToken, "."), 2)
					}

					// check id token
					idToken := testhelpers.DecodeIDToken(t, refreshed)
					testhelpers.AssertIDToken(t, idToken, sub, cl.ID)
					assert.Equal(t, "extra id token value", idToken.Get("key_id").Str)

					t.Run("original tokens are invalidated", func(t *testing.T) {
						introspected := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
						assert.False(t, introspected.Get("active").Bool(), introspected.Raw)
						introspected = testhelpers.IntrospectToken(t, token.RefreshToken, adminTS)
						assert.False(t, introspected.Get("active").Bool(), introspected.Raw)
					})
				})

				t.Run("audience is forwarded to access token", func(t *testing.T) {
					token := x.PerformAuthCodeFlow(t.Context(), t, nil, conf, adminClient, func(t *testing.T, req *hydra.OAuth2LoginRequest) hydra.AcceptOAuth2LoginRequest {
						assert.EqualValues(t, cl.Audience, req.RequestedAccessTokenAudience)
						return hydra.AcceptOAuth2LoginRequest{Subject: sub}
					}, func(t *testing.T, req *hydra.OAuth2ConsentRequest) hydra.AcceptOAuth2ConsentRequest {
						assert.EqualValues(t, cl.Audience, req.RequestedAccessTokenAudience)
						return hydra.AcceptOAuth2ConsentRequest{
							GrantScope:               []string{"openid"},
							GrantAccessTokenAudience: req.RequestedAccessTokenAudience,
						}
					}, oauth2.SetAuthURLParam("audience", strings.Join(cl.Audience, " ")))

					expectedAud, err := json.Marshal(cl.Audience)
					require.NoError(t, err)

					introspected := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
					require.True(t, introspected.Get("active").Bool())
					testhelpers.AssertAccessToken(t, introspected, sub, cl.ID)
					assert.JSONEq(t, string(expectedAud), introspected.Get("aud").Raw)

					if accessTokenStrategy == "jwt" {
						decoded := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))
						testhelpers.AssertAccessToken(t, decoded, sub, cl.ID)
						assert.JSONEq(t, string(expectedAud), decoded.Get("aud").Raw)
					} else {
						assert.Len(t, strings.Split(token.AccessToken, "."), 2)
					}

					idToken := testhelpers.DecodeIDToken(t, token)
					testhelpers.AssertIDToken(t, idToken, sub, cl.ID)
					require.Len(t, idToken.Get("aud").Array(), 1)
					assert.Equal(t, cl.ID, idToken.Get("aud.0").Str)
				})
			})
		}
	})
}

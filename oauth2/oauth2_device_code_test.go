// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/pointerx"
)

func TestDeviceAuthRequest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t)
	testhelpers.NewOAuth2Server(ctx, t, reg)

	secret := uuid.New()
	c := &client.Client{
		ID:                      "device-client",
		Secret:                  secret,
		GrantTypes:              []string{"urn:ietf:params:oauth:grant-type:device_code"},
		Scope:                   "hydra offline openid",
		Audience:                []string{"https://api.ory.sh/"},
		TokenEndpointAuthMethod: "client_secret_post",
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID:     c.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:     oauth2.AuthStyleInParams,
		},
		Scopes: strings.Split(c.Scope, " "),
	}

	testCases := []struct {
		description string
		setUp       func()
		check       func(t *testing.T, resp *oauth2.DeviceAuthResponse, err error)
		cleanUp     func()
	}{
		{
			description: "should pass",
			check: func(t *testing.T, resp *oauth2.DeviceAuthResponse, _ error) {
				assert.NotEmpty(t, resp.DeviceCode)
				assert.NotEmpty(t, resp.UserCode)
				assert.NotEmpty(t, resp.Interval)
				assert.NotEmpty(t, resp.VerificationURI)
				assert.NotEmpty(t, resp.VerificationURIComplete)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("case="+testCase.description, func(t *testing.T) {
			if testCase.setUp != nil {
				testCase.setUp()
			}

			resp, err := oauthClient.DeviceAuth(context.Background(), []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("client_secret", secret)}...)

			if testCase.check != nil {
				testCase.check(t, resp, err)
			}

			if testCase.cleanUp != nil {
				testCase.cleanUp()
			}
		})
	}
}

func TestDeviceTokenRequest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t)
	testhelpers.NewOAuth2Server(ctx, t, reg)

	secret := uuid.New()
	c := &client.Client{
		ID:     "device-client",
		Secret: secret,
		GrantTypes: []string{
			string(fosite.GrantTypeDeviceCode),
			string(fosite.GrantTypeRefreshToken),
		},
		Scope:    "hydra offline openid",
		Audience: []string{"https://api.ory.sh/"},
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID:     c.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:     oauth2.AuthStyleInHeader,
		},
		Scopes: strings.Split(c.Scope, " "),
	}

	testCases := []struct {
		description string
		setUp       func(signature, userCodeSignature string)
		check       func(t *testing.T, token *oauth2.Token, err error)
		cleanUp     func()
	}{
		{
			description: "should pass with refresh token",
			setUp: func(signature, userCodeSignature string) {
				authreq := &fosite.DeviceRequest{
					UserCodeState: fosite.UserCodeAccepted,
					Request: fosite.Request{
						Client: &fosite.DefaultClient{
							ID:         c.GetID(),
							GrantTypes: []string{string(fosite.GrantTypeDeviceCode)},
						},
						RequestedScope: []string{"hydra", "offline"},
						GrantedScope:   []string{"hydra", "offline"},
						Session: &hydraoauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "hydra",
								},
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(time.Hour).UTC(),
								},
							},
						},
						RequestedAt: time.Now(),
					},
				}

				require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
			},
			check: func(t *testing.T, token *oauth2.Token, err error) {
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
			},
		},
		{
			description: "should pass with ID token",
			setUp: func(signature, userCodeSignature string) {
				authreq := &fosite.DeviceRequest{
					UserCodeState: fosite.UserCodeAccepted,
					Request: fosite.Request{
						Client: &fosite.DefaultClient{
							ID:         c.GetID(),
							GrantTypes: []string{string(fosite.GrantTypeDeviceCode)},
						},
						RequestedScope: []string{"hydra", "offline", "openid"},
						GrantedScope:   []string{"hydra", "offline", "openid"},
						Session: &hydraoauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "hydra",
								},
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(time.Hour).UTC(),
								},
							},
						},
						RequestedAt: time.Now(),
					},
				}

				require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
				require.NoError(t, reg.OAuth2Storage().CreateOpenIDConnectSession(context.TODO(), signature, authreq))
			},
			check: func(t *testing.T, token *oauth2.Token, err error) {
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
				assert.NotEmpty(t, token.Extra("id_token"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("case="+testCase.description, func(t *testing.T) {
			code, signature, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(context.TODO())
			require.NoError(t, err)
			_, userCodeSignature, err := reg.RFC8628HMACStrategy().GenerateUserCode(context.TODO())
			require.NoError(t, err)

			if testCase.setUp != nil {
				testCase.setUp(signature, userCodeSignature)
			}

			var token *oauth2.Token
			token, err = oauthClient.DeviceAccessToken(context.Background(), &oauth2.DeviceAuthResponse{DeviceCode: code})

			if testCase.check != nil {
				testCase.check(t, token, err)
			}

			if testCase.cleanUp != nil {
				testCase.cleanUp()
			}
		})
	}
}

func TestDeviceCodeWithDefaultStrategy(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyAccessTokenStrategy: "opaque",
		config.KeyRefreshTokenHook:    "",
	})))
	publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)

	publicClient := hydra.NewAPIClient(hydra.NewConfiguration())
	publicClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: publicTS.URL}}
	adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
	adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	getDeviceCode := func(t *testing.T, conf *oauth2.Config, c *http.Client, params ...oauth2.AuthCodeOption) (*oauth2.DeviceAuthResponse, error) {
		return conf.DeviceAuth(ctx, params...)
	}

	acceptUserCode := func(t *testing.T, conf *oauth2.Config, c *http.Client, devResp *oauth2.DeviceAuthResponse) *http.Response {
		if c == nil {
			c = testhelpers.NewEmptyJarClient(t)
		}

		resp, err := c.Get(devResp.VerificationURIComplete)
		require.NoError(t, err)
		require.Contains(t, reg.Config().DeviceDoneURL(ctx).String(), resp.Request.URL.Path, "did not end up in post device URL")
		require.Equal(t, resp.Request.URL.Query().Get("client_id"), conf.ClientID)

		return resp
	}

	acceptDeviceHandler := func(t *testing.T, c *client.Client) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userCode := r.URL.Query().Get("user_code")
			payload := hydra.AcceptDeviceUserCodeRequest{
				UserCode: &userCode,
			}

			v, _, err := adminClient.OAuth2API.AcceptUserCodeRequest(context.Background()).
				DeviceChallenge(r.URL.Query().Get("device_challenge")).
				AcceptDeviceUserCodeRequest(payload).
				Execute()
			require.NoError(t, err)
			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}

	acceptLoginHandler := func(t *testing.T, c *client.Client, subject string, scopes []string, checkRequestPayload func(request *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, _, err := adminClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
			require.NoError(t, err)

			assert.EqualValues(t, c.GetID(), pointerx.Deref(rr.Client.ClientId))
			assert.Empty(t, pointerx.Deref(rr.Client.ClientSecret))
			assert.EqualValues(t, c.GrantTypes, rr.Client.GrantTypes)
			assert.EqualValues(t, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
			assert.EqualValues(t, r.URL.Query().Get("login_challenge"), rr.Challenge)
			assert.EqualValues(t, scopes, rr.RequestedScope)
			assert.Contains(t, rr.RequestUrl, hydraoauth2.DeviceVerificationPath)

			acceptBody := hydra.AcceptOAuth2LoginRequest{
				Subject:  subject,
				Remember: pointerx.Ptr(!rr.Skip),
				Acr:      pointerx.Ptr("1"),
				Amr:      []string{"pwd"},
				Context:  map[string]interface{}{"context": "bar"},
			}
			if checkRequestPayload != nil {
				if b := checkRequestPayload(rr); b != nil {
					acceptBody = *b
				}
			}

			v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
				LoginChallenge(r.URL.Query().Get("login_challenge")).
				AcceptOAuth2LoginRequest(acceptBody).
				Execute()
			require.NoError(t, err)
			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}

	acceptConsentHandler := func(t *testing.T, c *client.Client, subject string, scopes []string, checkRequestPayload func(*hydra.OAuth2ConsentRequest)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
			require.NoError(t, err)

			assert.EqualValues(t, c.GetID(), pointerx.Deref(rr.Client.ClientId))
			assert.Empty(t, pointerx.Deref(rr.Client.ClientSecret))
			assert.EqualValues(t, c.GrantTypes, rr.Client.GrantTypes)
			assert.EqualValues(t, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
			assert.EqualValues(t, subject, pointerx.Deref(rr.Subject))
			assert.EqualValues(t, scopes, rr.RequestedScope)
			assert.Contains(t, *rr.RequestUrl, hydraoauth2.DeviceVerificationPath)
			if checkRequestPayload != nil {
				checkRequestPayload(rr)
			}

			assert.Equal(t, map[string]interface{}{"context": "bar"}, rr.Context)
			v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
				ConsentChallenge(r.URL.Query().Get("consent_challenge")).
				AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
					GrantScope: scopes, Remember: pointerx.Ptr(true), RememberFor: pointerx.Ptr[int64](0),
					GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}).
				Execute()
			require.NoError(t, err)
			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}

	assertRefreshToken := func(t *testing.T, token *oauth2.Token, c *oauth2.Config, expectedExp time.Time) {
		actualExp := testhelpers.IntrospectToken(t, token.RefreshToken, adminTS).Get("exp").Int()
		assert.WithinDuration(t, expectedExp, time.Unix(actualExp, 0), time.Second)
	}

	assertIDToken := func(t *testing.T, token *oauth2.Token, c *oauth2.Config, expectedSubject, expectedNonce string, expectedExp time.Time) gjson.Result {
		idt, ok := token.Extra("id_token").(string)
		require.True(t, ok)
		assert.NotEmpty(t, idt)

		claims := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, idt))
		assert.True(t, time.Now().After(time.Unix(claims.Get("iat").Int(), 0)), "%s", claims)
		assert.True(t, time.Now().After(time.Unix(claims.Get("nbf").Int(), 0)), "%s", claims)
		assert.True(t, time.Now().Before(time.Unix(claims.Get("exp").Int(), 0)), "%s", claims)
		assert.WithinDuration(t, expectedExp, time.Unix(claims.Get("exp").Int(), 0), 2*time.Second)
		assert.NotEmpty(t, claims.Get("jti").String(), "%s", claims)
		assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), claims.Get("iss").String(), "%s", claims)
		assert.NotEmpty(t, claims.Get("sid").String(), "%s", claims)
		assert.Equal(t, "1", claims.Get("acr").String(), "%s", claims)
		require.Len(t, claims.Get("amr").Array(), 1, "%s", claims)
		assert.EqualValues(t, "pwd", claims.Get("amr").Array()[0].String(), "%s", claims)

		require.Len(t, claims.Get("aud").Array(), 1, "%s", claims)
		assert.EqualValues(t, c.ClientID, claims.Get("aud").Array()[0].String(), "%s", claims)
		assert.EqualValues(t, expectedSubject, claims.Get("sub").String(), "%s", claims)
		assert.EqualValues(t, `baz`, claims.Get("bar").String(), "%s", claims)

		return claims
	}

	introspectAccessToken := func(t *testing.T, conf *oauth2.Config, token *oauth2.Token, expectedSubject string) gjson.Result {
		require.NotEmpty(t, token.AccessToken)
		i := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
		assert.True(t, i.Get("active").Bool(), "%s", i)
		assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
		assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
		assert.EqualValues(t, `bar`, i.Get("ext.foo").String(), "%s", i)
		return i
	}

	assertJWTAccessToken := func(t *testing.T, strat string, conf *oauth2.Config, token *oauth2.Token, expectedSubject string, expectedExp time.Time, scopes string) gjson.Result {
		require.NotEmpty(t, token.AccessToken)
		parts := strings.Split(token.AccessToken, ".")
		if strat != "jwt" {
			require.Len(t, parts, 2)
			return gjson.Parse("null")
		}
		require.Len(t, parts, 3)

		i := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))
		assert.NotEmpty(t, i.Get("jti").String())
		assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
		assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
		assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), i.Get("iss").String(), "%s", i)
		assert.True(t, time.Now().After(time.Unix(i.Get("iat").Int(), 0)), "%s", i)
		assert.True(t, time.Now().After(time.Unix(i.Get("nbf").Int(), 0)), "%s", i)
		assert.True(t, time.Now().Before(time.Unix(i.Get("exp").Int(), 0)), "%s", i)
		assert.WithinDuration(t, expectedExp, time.Unix(i.Get("exp").Int(), 0), time.Second)
		assert.EqualValues(t, `bar`, i.Get("ext.foo").String(), "%s", i)
		assert.EqualValues(t, scopes, i.Get("scp").Raw, "%s", i)
		return i
	}

	waitForRefreshTokenExpiry := func() {
		time.Sleep(reg.Config().GetRefreshTokenLifespan(ctx) + time.Second)
	}

	t.Run("case=checks if request fails when audience does not match", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		_, conf := newDeviceClient(t, reg)
		resp, err := conf.DeviceAuth(ctx, oauth2.SetAuthURLParam("audience", "https://not-ory-api/"))
		require.Error(t, err)
		var devErr *oauth2.RetrieveError
		require.ErrorAs(t, err, &devErr)
		require.Nil(t, resp)
		require.Equal(t, devErr.Response.StatusCode, http.StatusBadRequest)
	})

	subject := "aeneas-rekkas"
	nonce := uuid.New()
	t.Run("case=perform device flow without ID and refresh tokens", func(t *testing.T) {

		c, conf := newDeviceClient(t, reg)
		conf.Scopes = []string{"hydra"}
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			acceptDeviceHandler(t, c),
			acceptLoginHandler(t, c, subject, conf.Scopes, nil),
			acceptConsentHandler(t, c, subject, conf.Scopes, nil),
		)

		resp, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)
		loginFlowResp := acceptUserCode(t, conf, nil, resp)
		require.NotNil(t, loginFlowResp)
		token, err := conf.DeviceAccessToken(context.Background(), resp)
		require.NoError(t, err)

		assert.Empty(t, token.Extra("c_nonce_draft_00"), "should not be set if not requested")
		assert.Empty(t, token.Extra("c_nonce_expires_in_draft_00"), "should not be set if not requested")
		introspectAccessToken(t, conf, token, subject)
		assert.Empty(t, token.Extra("id_token"))
		assert.Empty(t, token.RefreshToken)
	})
	t.Run("case=perform device flow with ID token", func(t *testing.T) {

		c, conf := newDeviceClient(t, reg)
		conf.Scopes = []string{"openid", "hydra"}
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			acceptDeviceHandler(t, c),
			acceptLoginHandler(t, c, subject, conf.Scopes, nil),
			acceptConsentHandler(t, c, subject, conf.Scopes, nil),
		)

		resp, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)
		loginFlowResp := acceptUserCode(t, conf, nil, resp)
		require.NotNil(t, loginFlowResp)
		token, err := conf.DeviceAccessToken(context.Background(), resp)
		iat := time.Now()
		require.NoError(t, err)

		assert.Empty(t, token.Extra("c_nonce_draft_00"), "should not be set if not requested")
		assert.Empty(t, token.Extra("c_nonce_expires_in_draft_00"), "should not be set if not requested")
		introspectAccessToken(t, conf, token, subject)
		assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
		assert.Empty(t, token.RefreshToken)
	})
	t.Run("case=perform device flow with refresh token", func(t *testing.T) {

		c, conf := newDeviceClient(t, reg)
		conf.Scopes = []string{"hydra", "offline"}
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			acceptDeviceHandler(t, c),
			acceptLoginHandler(t, c, subject, conf.Scopes, nil),
			acceptConsentHandler(t, c, subject, conf.Scopes, nil),
		)

		resp, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)
		loginFlowResp := acceptUserCode(t, conf, nil, resp)
		require.NotNil(t, loginFlowResp)
		token, err := conf.DeviceAccessToken(context.Background(), resp)
		iat := time.Now()
		require.NoError(t, err)

		assert.Empty(t, token.Extra("c_nonce_draft_00"), "should not be set if not requested")
		assert.Empty(t, token.Extra("c_nonce_expires_in_draft_00"), "should not be set if not requested")
		introspectAccessToken(t, conf, token, subject)
		assert.Empty(t, token.Extra("id_token"))
		assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))
	})
	t.Run("case=perform device flow with ID token and refresh tokens", func(t *testing.T) {
		run := func(t *testing.T, strategy string) {
			c, conf := newDeviceClient(t, reg)
			testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
				acceptDeviceHandler(t, c),
				acceptLoginHandler(t, c, subject, conf.Scopes, nil),
				acceptConsentHandler(t, c, subject, conf.Scopes, nil),
			)

			resp, err := getDeviceCode(t, conf, nil)
			require.NoError(t, err)
			require.NotEmpty(t, resp.DeviceCode)
			require.NotEmpty(t, resp.UserCode)
			loginFlowResp := acceptUserCode(t, conf, nil, resp)
			require.NotNil(t, loginFlowResp)
			token, err := conf.DeviceAccessToken(context.Background(), resp)
			iat := time.Now()
			require.NoError(t, err)

			assert.Empty(t, token.Extra("c_nonce_draft_00"), "should not be set if not requested")
			assert.Empty(t, token.Extra("c_nonce_expires_in_draft_00"), "should not be set if not requested")
			introspectAccessToken(t, conf, token, subject)
			assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
			assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
			assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))

			t.Run("followup=successfully perform refresh token flow", func(t *testing.T) {
				require.NotEmpty(t, token.RefreshToken)
				token.Expiry = token.Expiry.Add(-time.Hour * 24)
				iat = time.Now()
				refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
				require.NoError(t, err)

				require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
				require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
				require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))
				introspectAccessToken(t, conf, refreshedToken, subject)

				t.Run("followup=refreshed tokens contain valid tokens", func(t *testing.T) {
					assertJWTAccessToken(t, strategy, conf, refreshedToken, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
					assertIDToken(t, refreshedToken, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
					assertRefreshToken(t, refreshedToken, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))
				})

				t.Run("followup=original access token is no longer valid", func(t *testing.T) {
					i := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
					assert.False(t, i.Get("active").Bool(), "%s", i)
				})

				t.Run("followup=original refresh token is no longer valid", func(t *testing.T) {
					_, err := conf.TokenSource(context.Background(), token).Token()
					assert.Error(t, err)
				})

				t.Run("followup=but fail subsequent refresh because expiry was reached", func(t *testing.T) {
					waitForRefreshTokenExpiry()

					// Force golang to refresh token
					refreshedToken.Expiry = refreshedToken.Expiry.Add(-time.Hour * 24)
					_, err := conf.TokenSource(context.Background(), refreshedToken).Token()
					require.Error(t, err)
				})
			})
		}

		t.Run("strategy=jwt", func(t *testing.T) {
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
			run(t, "jwt")
		})

		t.Run("strategy=opaque", func(t *testing.T) {
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
			run(t, "opaque")
		})
	})
	t.Run("case=perform flow with audience", func(t *testing.T) {
		expectAud := "https://api.ory.sh/"
		c, conf := newDeviceClient(t, reg)
		testhelpers.NewDeviceLoginConsentUI(
			t,
			reg.Config(),
			acceptDeviceHandler(t, c),
			acceptLoginHandler(t, c, subject, conf.Scopes, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
				assert.False(t, r.Skip)
				assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
				return nil
			}),
			acceptConsentHandler(t, c, subject, conf.Scopes, func(r *hydra.OAuth2ConsentRequest) {
				assert.False(t, *r.Skip)
				assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
			}),
		)

		resp, err := conf.DeviceAuth(ctx, oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"))
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)
		loginFlowResp := acceptUserCode(t, conf, nil, resp)
		require.NotNil(t, loginFlowResp)

		token, err := conf.DeviceAccessToken(context.Background(), resp)
		require.NoError(t, err)

		claims := introspectAccessToken(t, conf, token, subject)
		aud := claims.Get("aud").Array()
		require.Len(t, aud, 1)
		assert.EqualValues(t, aud[0].String(), expectAud)

		assertIDToken(t, token, conf, subject, nonce, time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))
	})

	t.Run("case=respects client token lifespan configuration", func(t *testing.T) {
		run := func(t *testing.T, strategy string, c *client.Client, conf *oauth2.Config, expectedLifespans client.Lifespans) {
			testhelpers.NewDeviceLoginConsentUI(
				t,
				reg.Config(),
				acceptDeviceHandler(t, c),
				acceptLoginHandler(t, c, subject, conf.Scopes, nil),
				acceptConsentHandler(t, c, subject, conf.Scopes, nil),
			)

			resp, err := getDeviceCode(t, conf, nil)
			require.NoError(t, err)
			require.NotEmpty(t, resp.DeviceCode)
			require.NotEmpty(t, resp.UserCode)
			loginFlowResp := acceptUserCode(t, conf, nil, resp)
			require.NotNil(t, loginFlowResp)

			token, err := conf.DeviceAccessToken(context.Background(), resp)
			iat := time.Now()
			require.NoError(t, err)

			body := introspectAccessToken(t, conf, token, subject)
			assert.WithinDuration(t, iat.Add(expectedLifespans.DeviceAuthorizationGrantAccessTokenLifespan.Duration), time.Unix(body.Get("exp").Int(), 0), time.Second)

			assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(expectedLifespans.DeviceAuthorizationGrantAccessTokenLifespan.Duration), `["hydra","offline","openid"]`)
			assertIDToken(t, token, conf, subject, nonce, iat.Add(expectedLifespans.DeviceAuthorizationGrantIDTokenLifespan.Duration))
			assertRefreshToken(t, token, conf, iat.Add(expectedLifespans.DeviceAuthorizationGrantRefreshTokenLifespan.Duration))

			t.Run("followup=successfully perform refresh token flow", func(t *testing.T) {
				require.NotEmpty(t, token.RefreshToken)
				token.Expiry = token.Expiry.Add(-time.Hour * 24)
				refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
				iat = time.Now()
				require.NoError(t, err)
				assertRefreshToken(t, refreshedToken, conf, iat.Add(expectedLifespans.RefreshTokenGrantRefreshTokenLifespan.Duration))
				assertJWTAccessToken(t, strategy, conf, refreshedToken, subject, iat.Add(expectedLifespans.RefreshTokenGrantAccessTokenLifespan.Duration), `["hydra","offline","openid"]`)
				assertIDToken(t, refreshedToken, conf, subject, nonce, iat.Add(expectedLifespans.RefreshTokenGrantIDTokenLifespan.Duration))

				require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
				require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
				require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))

				body := introspectAccessToken(t, conf, refreshedToken, subject)
				assert.WithinDuration(t, iat.Add(expectedLifespans.RefreshTokenGrantAccessTokenLifespan.Duration), time.Unix(body.Get("exp").Int(), 0), time.Second)

				t.Run("followup=original access token is no longer valid", func(t *testing.T) {
					i := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
					assert.False(t, i.Get("active").Bool(), "%s", i)
				})

				t.Run("followup=original refresh token is no longer valid", func(t *testing.T) {
					_, err := conf.TokenSource(context.Background(), token).Token()
					assert.Error(t, err)
				})
			})
		}

		t.Run("case=custom-lifespans-active-jwt", func(t *testing.T) {
			c, conf := newDeviceClient(t, reg)
			ls := testhelpers.TestLifespans
			ls.DeviceAuthorizationGrantAccessTokenLifespan = x.NullDuration{Valid: true, Duration: 6 * time.Second}
			testhelpers.UpdateClientTokenLifespans(
				t,
				&oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret},
				c.GetID(),
				ls, adminTS,
			)
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
			run(t, "jwt", c, conf, ls)
		})

		t.Run("case=custom-lifespans-active-opaque", func(t *testing.T) {
			c, conf := newDeviceClient(t, reg)
			ls := testhelpers.TestLifespans
			ls.DeviceAuthorizationGrantAccessTokenLifespan = x.NullDuration{Valid: true, Duration: 6 * time.Second}
			testhelpers.UpdateClientTokenLifespans(
				t,
				&oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret},
				c.GetID(),
				ls, adminTS,
			)
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
			run(t, "opaque", c, conf, ls)
		})

		t.Run("case=custom-lifespans-unset", func(t *testing.T) {
			c, conf := newDeviceClient(t, reg)
			testhelpers.UpdateClientTokenLifespans(t, &oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret}, c.GetID(), testhelpers.TestLifespans, adminTS)
			testhelpers.UpdateClientTokenLifespans(t, &oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret}, c.GetID(), client.Lifespans{}, adminTS)
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")

			//goland:noinspection GoDeprecation
			expectedLifespans := client.Lifespans{
				AuthorizationCodeGrantAccessTokenLifespan:    x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				AuthorizationCodeGrantIDTokenLifespan:        x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
				AuthorizationCodeGrantRefreshTokenLifespan:   x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
				ClientCredentialsGrantAccessTokenLifespan:    x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				ImplicitGrantAccessTokenLifespan:             x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				ImplicitGrantIDTokenLifespan:                 x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
				JwtBearerGrantAccessTokenLifespan:            x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				PasswordGrantAccessTokenLifespan:             x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				PasswordGrantRefreshTokenLifespan:            x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
				RefreshTokenGrantIDTokenLifespan:             x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
				RefreshTokenGrantAccessTokenLifespan:         x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				RefreshTokenGrantRefreshTokenLifespan:        x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
				DeviceAuthorizationGrantIDTokenLifespan:      x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
				DeviceAuthorizationGrantAccessTokenLifespan:  x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
				DeviceAuthorizationGrantRefreshTokenLifespan: x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
			}
			run(t, "opaque", c, conf, expectedLifespans)
		})
	})
	t.Run("case=cannot reuse user_code", func(t *testing.T) {
		c, conf := newDeviceClient(t, reg)
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				userCode := r.URL.Query().Get("user_code")
				payload := hydra.AcceptDeviceUserCodeRequest{
					UserCode: &userCode,
				}

				v, _, err := adminClient.OAuth2API.AcceptUserCodeRequest(context.Background()).
					DeviceChallenge(r.URL.Query().Get("device_challenge")).
					AcceptDeviceUserCodeRequest(payload).
					Execute()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				require.NotEmpty(t, v.RedirectTo)
				http.Redirect(w, r, v.RedirectTo, http.StatusFound)
			},
			acceptLoginHandler(t, c, subject, conf.Scopes, nil),
			acceptConsentHandler(t, c, subject, conf.Scopes, nil),
		)

		resp, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)
		loginFlowResp := acceptUserCode(t, conf, nil, resp)
		require.NotNil(t, loginFlowResp)
		token, err := conf.DeviceAccessToken(context.Background(), resp)
		iat := time.Now()
		require.NoError(t, err)

		introspectAccessToken(t, conf, token, subject)
		assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
		assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))

		hc := testhelpers.NewEmptyJarClient(t)

		loginFlowResp2, err := hc.Get(resp.VerificationURIComplete)
		require.NoError(t, err)
		require.Equal(t, loginFlowResp2.StatusCode, http.StatusBadRequest)
	})
	t.Run("case=cannot reuse device_challenge", func(t *testing.T) {
		var deviceChallenge string
		c, conf := newDeviceClient(t, reg)
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				userCode := r.URL.Query().Get("user_code")
				payload := hydra.AcceptDeviceUserCodeRequest{
					UserCode: &userCode,
				}

				if deviceChallenge == "" {
					deviceChallenge = r.URL.Query().Get("device_challenge")
				}
				v, _, err := adminClient.OAuth2API.AcceptUserCodeRequest(context.Background()).
					DeviceChallenge(deviceChallenge).
					AcceptDeviceUserCodeRequest(payload).
					Execute()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				require.NoError(t, err)
				require.NotEmpty(t, v.RedirectTo)
				http.Redirect(w, r, v.RedirectTo, http.StatusFound)
			},
			acceptLoginHandler(t, c, subject, conf.Scopes, nil),
			acceptConsentHandler(t, c, subject, conf.Scopes, nil),
		)

		resp, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.DeviceCode)
		require.NotEmpty(t, resp.UserCode)

		hc := testhelpers.NewEmptyJarClient(t)
		loginFlowResp := acceptUserCode(t, conf, hc, resp)
		require.NoError(t, err)
		require.Contains(t, reg.Config().DeviceDoneURL(ctx).String(), loginFlowResp.Request.URL.Path, "did not end up in post device URL")
		require.Equal(t, loginFlowResp.Request.URL.Query().Get("client_id"), conf.ClientID)

		require.NotNil(t, loginFlowResp)
		token, err := conf.DeviceAccessToken(context.Background(), resp)
		iat := time.Now()
		require.NoError(t, err)

		introspectAccessToken(t, conf, token, subject)
		assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
		assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))

		resp2, err := getDeviceCode(t, conf, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp2.DeviceCode)
		require.NotEmpty(t, resp2.UserCode)

		payload := hydra.AcceptDeviceUserCodeRequest{
			UserCode: &resp2.UserCode,
		}

		acceptResp, _, err := adminClient.OAuth2API.AcceptUserCodeRequest(context.Background()).
			DeviceChallenge(deviceChallenge).
			AcceptDeviceUserCodeRequest(payload).
			Execute()
		require.NoError(t, err)

		loginFlowResp2, err := hc.Get(acceptResp.RedirectTo)
		require.NoError(t, err)
		require.Equalf(t, http.StatusForbidden, loginFlowResp2.StatusCode, "requested %q", acceptResp.RedirectTo)
	})
}

func newDeviceClient(
	t *testing.T,
	reg interface {
		config.Provider
		client.Registry
	},
	opts ...func(*client.Client),
) (*client.Client, *oauth2.Config) {
	ctx := context.Background()
	c := &client.Client{
		GrantTypes: []string{
			"refresh_token",
			"urn:ietf:params:oauth:grant-type:device_code",
		},
		Scope:                   "hydra offline openid",
		Audience:                []string{"https://api.ory.sh/"},
		TokenEndpointAuthMethod: "none",
	}

	// apply options
	for _, o := range opts {
		o(c)
	}

	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))
	return c, &oauth2.Config{
		ClientID: c.GetID(),
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:     oauth2.AuthStyleInHeader,
		},
		Scopes: strings.Split(c.Scope, " "),
	}
}

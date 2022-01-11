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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/internal/testhelpers"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/fosite"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	hydraoauth2 "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"
)

func noopHandler(t *testing.T) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type clientCreator interface {
	CreateClient(cxt context.Context, client *hc.Client) error
}

// TestAuthCodeWithDefaultStrategy runs proper integration tests against in-memory and database connectors, specifically
// we test:
//
// - [x] If the flow - in general - works
// - [x] If `authenticatedAt` is properly managed across the lifecycle
//   - [x] The value `authenticatedAt` should be an old time if no user interaction wrt login was required
//   - [x] The value `authenticatedAt` should be a recent time if user interaction wrt login was required
// - [x] If `requestedAt` is properly managed across the lifecycle
//   - [x] The value of `requestedAt` must be the initial request time, not some other time (e.g. when accepting login)
// - [x] If `id_token_hint` is handled properly
//   - [x] What happens if `id_token_hint` does not match the value from the handled authentication request ("accept login")
func TestAuthCodeWithDefaultStrategy(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	reg.Config().MustSet(config.KeyAccessTokenStrategy, "opaque")
	publicTS, adminTS := testhelpers.NewOAuth2Server(t, reg)

	newOAuth2Client := func(t *testing.T, cb string) (*hc.Client, *oauth2.Config) {
		secret := uuid.New()
		c := &hc.Client{
			OutfacingID:   uuid.New(),
			Secret:        secret,
			RedirectURIs:  []string{cb},
			ResponseTypes: []string{"id_token", "code", "token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scope:         "hydra offline openid",
			Audience:      []string{"https://api.ory.sh/"},
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.TODO(), c))
		return c, &oauth2.Config{
			ClientID:     c.OutfacingID,
			ClientSecret: secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   reg.Config().OAuth2AuthURL().String(),
				TokenURL:  reg.Config().OAuth2TokenURL().String(),
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			Scopes: strings.Split(c.Scope, " "),
		}
	}

	adminClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(adminTS.URL).Host})

	getAuthorizeCode := func(t *testing.T, conf *oauth2.Config, c *http.Client, params ...oauth2.AuthCodeOption) (string, *http.Response) {
		if c == nil {
			c = testhelpers.NewEmptyJarClient(t)
		}

		state := uuid.New()
		resp, err := c.Get(conf.AuthCodeURL(state, params...))
		require.NoError(t, err)
		defer resp.Body.Close()

		q := resp.Request.URL.Query()
		require.EqualValues(t, state, q.Get("state"))
		return q.Get("code"), resp
	}

	acceptLoginHandler := func(t *testing.T, c *client.Client, subject string, checkRequestPayload func(*admin.GetLoginRequestOK) *models.AcceptLoginRequest) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, err := adminClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
			require.NoError(t, err)

			assert.EqualValues(t, c.GetID(), rr.Payload.Client.ClientID)
			assert.Empty(t, rr.Payload.Client.ClientSecret)
			assert.EqualValues(t, c.GrantTypes, rr.Payload.Client.GrantTypes)
			assert.EqualValues(t, c.LogoURI, rr.Payload.Client.LogoURI)
			assert.EqualValues(t, c.RedirectURIs, rr.Payload.Client.RedirectUris)
			assert.EqualValues(t, r.URL.Query().Get("login_challenge"), *rr.Payload.Challenge)
			assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.Payload.RequestedScope)
			assert.Contains(t, *rr.Payload.RequestURL, reg.Config().OAuth2AuthURL().String())

			acceptBody := &models.AcceptLoginRequest{
				Subject:  &subject,
				Remember: !*rr.Payload.Skip,
				Acr:      "1",
				Amr:      models.StringSlicePipeDelimiter{"pwd"},
				Context:  map[string]interface{}{"context": "bar"},
			}
			if checkRequestPayload != nil {
				if b := checkRequestPayload(rr); b != nil {
					acceptBody = b
				}
			}

			v, err := adminClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(acceptBody))
			require.NoError(t, err)
			require.NotEmpty(t, *v.Payload.RedirectTo)
			http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
		}
	}

	acceptConsentHandler := func(t *testing.T, c *client.Client, subject string, checkRequestPayload func(*admin.GetConsentRequestOK)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, err := adminClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
			require.NoError(t, err)

			assert.EqualValues(t, c.GetID(), rr.Payload.Client.ClientID)
			assert.Empty(t, rr.Payload.Client.ClientSecret)
			assert.EqualValues(t, c.GrantTypes, rr.Payload.Client.GrantTypes)
			assert.EqualValues(t, c.LogoURI, rr.Payload.Client.LogoURI)
			assert.EqualValues(t, c.RedirectURIs, rr.Payload.Client.RedirectUris)
			assert.EqualValues(t, subject, rr.Payload.Subject)
			assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.Payload.RequestedScope)
			assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), *rr.Payload.Challenge)
			assert.Contains(t, rr.Payload.RequestURL, reg.Config().OAuth2AuthURL().String())
			if checkRequestPayload != nil {
				checkRequestPayload(rr)
			}

			assert.Equal(t, map[string]interface{}{"context": "bar"}, rr.Payload.Context)
			v, err := adminClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
				WithBody(&models.AcceptConsentRequest{
					GrantScope: []string{"hydra", "offline", "openid"}, Remember: true, RememberFor: 0,
					GrantAccessTokenAudience: rr.Payload.RequestedAccessTokenAudience,
					Session: &models.ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IDToken:     map[string]interface{}{"bar": "baz"},
					},
				}))
			require.NoError(t, err)
			require.NotEmpty(t, *v.Payload.RedirectTo)
			http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
		}
	}

	assertIDToken := func(t *testing.T, token *oauth2.Token, c *oauth2.Config, expectedSubject, expectedNonce string) gjson.Result {
		idt, ok := token.Extra("id_token").(string)
		require.True(t, ok)
		assert.NotEmpty(t, idt)

		body, err := x.DecodeSegment(strings.Split(idt, ".")[1])
		require.NoError(t, err)

		claims := gjson.ParseBytes(body)
		assert.True(t, time.Now().After(time.Unix(claims.Get("iat").Int(), 0)), "%s", claims)
		assert.True(t, time.Now().After(time.Unix(claims.Get("nbf").Int(), 0)), "%s", claims)
		assert.True(t, time.Now().Before(time.Unix(claims.Get("exp").Int(), 0)), "%s", claims)
		assert.NotEmpty(t, claims.Get("jti").String(), "%s", claims)
		assert.EqualValues(t, reg.Config().IssuerURL().String(), claims.Get("iss").String(), "%s", claims)
		assert.NotEmpty(t, claims.Get("sid").String(), "%s", claims)
		assert.Equal(t, "1", claims.Get("acr").String(), "%s", claims)
		require.Len(t, claims.Get("amr").Array(), 1, "%s", claims)
		assert.EqualValues(t, "pwd", claims.Get("amr").Array()[0].String(), "%s", claims)

		require.Len(t, claims.Get("aud").Array(), 1, "%s", claims)
		assert.EqualValues(t, c.ClientID, claims.Get("aud").Array()[0].String(), "%s", claims)
		assert.EqualValues(t, expectedSubject, claims.Get("sub").String(), "%s", claims)
		assert.EqualValues(t, expectedNonce, claims.Get("nonce").String(), "%s", claims)
		assert.EqualValues(t, `baz`, claims.Get("bar").String(), "%s", claims)

		return claims
	}

	introspectAccessToken := func(t *testing.T, conf *oauth2.Config, token *oauth2.Token, expectedSubject string) gjson.Result {
		require.NotEmpty(t, token.AccessToken)
		i := testhelpers.IntrospectToken(t, conf, token, adminTS)
		assert.True(t, i.Get("active").Bool(), "%s", i)
		assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
		assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
		assert.EqualValues(t, `{"foo":"bar"}`, i.Get("ext").Raw, "%s", i)
		return i
	}

	assertJWTAccessToken := func(t *testing.T, strat string, conf *oauth2.Config, token *oauth2.Token, expectedSubject string) gjson.Result {
		require.NotEmpty(t, token.AccessToken)
		parts := strings.Split(token.AccessToken, ".")
		if strat != "jwt" {
			require.Len(t, parts, 2)
			return gjson.Parse("null")
		}
		require.Len(t, parts, 3)

		body, err := x.DecodeSegment(parts[1])
		require.NoError(t, err)

		i := gjson.ParseBytes(body)
		assert.NotEmpty(t, i.Get("jti").String())
		assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
		assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
		assert.EqualValues(t, reg.Config().IssuerURL().String(), i.Get("iss").String(), "%s", i)
		assert.True(t, time.Now().After(time.Unix(i.Get("iat").Int(), 0)), "%s", i)
		assert.True(t, time.Now().After(time.Unix(i.Get("nbf").Int(), 0)), "%s", i)
		assert.True(t, time.Now().Before(time.Unix(i.Get("exp").Int(), 0)), "%s", i)
		assert.EqualValues(t, `{"foo":"bar"}`, i.Get("ext").Raw, "%s", i)
		assert.EqualValues(t, `["hydra","offline","openid"]`, i.Get("scp").Raw, "%s", i)
		return i
	}

	waitForRefreshTokenExpiry := func() {
		time.Sleep(reg.Config().RefreshTokenLifespan() + time.Second)
	}

	t.Run("case=checks if request fails when audience does not match", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		_, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("audience", "https://not-ory-api/"))
		require.Empty(t, code)
	})

	subject := "aeneas-rekkas"
	nonce := uuid.New()
	t.Run("case=perform authorize code flow with ID token and refresh tokens", func(t *testing.T) {
		run := func(t *testing.T, strategy string) {
			c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				acceptLoginHandler(t, c, subject, nil),
				acceptConsentHandler(t, c, subject, nil),
			)

			code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
			require.NotEmpty(t, code)
			token, err := conf.Exchange(context.Background(), code)
			require.NoError(t, err)

			introspectAccessToken(t, conf, token, subject)
			assertJWTAccessToken(t, strategy, conf, token, subject)
			assertIDToken(t, token, conf, subject, nonce)

			t.Run("followup=successfully perform refresh token flow", func(t *testing.T) {
				require.NotEmpty(t, token.RefreshToken)
				token.Expiry = token.Expiry.Add(-time.Hour * 24)
				refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
				require.NoError(t, err)

				require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
				require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
				require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))
				introspectAccessToken(t, conf, refreshedToken, subject)

				t.Run("followup=refreshed tokens contain ID token", func(t *testing.T) {
					assertIDToken(t, refreshedToken, conf, subject, nonce)
				})

				t.Run("followup=original access token is no longer valid", func(t *testing.T) {
					i := testhelpers.IntrospectToken(t, conf, token, adminTS)
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
			reg.Config().MustSet(config.KeyAccessTokenStrategy, "jwt")
			run(t, "jwt")
		})

		t.Run("strategy=opaque", func(t *testing.T) {
			reg.Config().MustSet(config.KeyAccessTokenStrategy, "opaque")
			run(t, "opaque")
		})
	})

	t.Run("case=checks if request fails when subject is empty", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), func(w http.ResponseWriter, r *http.Request) {
			_, err := adminClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(&models.AcceptLoginRequest{Subject: pointerx.String(""), Remember: true}))
			require.Error(t, err) // expects 400
			assert.Contains(t, err.(*admin.AcceptLoginRequestBadRequest).Payload.ErrorDescription, "Field 'subject' must not be empty", "%+v", *err.(*admin.AcceptLoginRequestBadRequest).Payload)
		}, testhelpers.HTTPServerNoExpectedCallHandler(t))
		_, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

		_, err := testhelpers.NewEmptyJarClient(t).Get(conf.AuthCodeURL(uuid.New()))
		require.NoError(t, err)
	})

	t.Run("case=perform flow with audience", func(t *testing.T) {
		expectAud := "https://api.ory.sh/"
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, func(r *admin.GetLoginRequestOK) *models.AcceptLoginRequest {
				assert.False(t, *r.Payload.Skip)
				assert.EqualValues(t, []string{expectAud}, r.Payload.RequestedAccessTokenAudience)
				return nil
			}),
			acceptConsentHandler(t, c, subject, func(r *admin.GetConsentRequestOK) {
				assert.False(t, r.Payload.Skip)
				assert.EqualValues(t, []string{expectAud}, r.Payload.RequestedAccessTokenAudience)
			}))

		code, _ := getAuthorizeCode(t, conf, nil,
			oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
			oauth2.SetAuthURLParam("nonce", nonce))
		require.NotEmpty(t, code)

		token, err := conf.Exchange(context.Background(), code)
		require.NoError(t, err)

		claims := introspectAccessToken(t, conf, token, subject)
		aud := claims.Get("aud").Array()
		require.Len(t, aud, 1)
		assert.EqualValues(t, aud[0].String(), expectAud)

		assertIDToken(t, token, conf, subject, nonce)
	})

	t.Run("case=use remember feature and prompt=none", func(t *testing.T) {
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, nil),
			acceptConsentHandler(t, c, subject, nil),
		)

		oc := testhelpers.NewEmptyJarClient(t)
		code, _ := getAuthorizeCode(t, conf, oc,
			oauth2.SetAuthURLParam("nonce", nonce),
			oauth2.SetAuthURLParam("prompt", "login consent"),
			oauth2.SetAuthURLParam("max_age", "1"),
		)
		require.NotEmpty(t, code)
		token, err := conf.Exchange(context.Background(), code)
		require.NoError(t, err)
		introspectAccessToken(t, conf, token, subject)

		// Reset UI to check for skip values
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, func(r *admin.GetLoginRequestOK) *models.AcceptLoginRequest {
				require.True(t, *r.Payload.Skip)
				require.EqualValues(t, subject, *r.Payload.Subject)
				return nil
			}),
			acceptConsentHandler(t, c, subject, func(r *admin.GetConsentRequestOK) {
				require.True(t, r.Payload.Skip)
				require.EqualValues(t, subject, r.Payload.Subject)
			}),
		)

		t.Run("followup=checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=none works", func(t *testing.T) {
			// In order to check if authenticatedAt/requestedAt works, we'll sleep first in order to ensure that authenticatedAt is in the past
			// if handled correctly.
			time.Sleep(time.Second + time.Nanosecond)

			code, _ := getAuthorizeCode(t, conf, oc,
				oauth2.SetAuthURLParam("nonce", nonce),
				oauth2.SetAuthURLParam("prompt", "none"),
				oauth2.SetAuthURLParam("max_age", "60"),
			)
			require.NotEmpty(t, code)
			token, err := conf.Exchange(context.Background(), code)
			require.NoError(t, err)
			original := introspectAccessToken(t, conf, token, subject)

			t.Run("followup=run the flow three more times", func(t *testing.T) {
				for i := 0; i < 3; i++ {
					t.Run(fmt.Sprintf("run=%d", i), func(t *testing.T) {
						code, _ := getAuthorizeCode(t, conf, oc,
							oauth2.SetAuthURLParam("nonce", nonce),
							oauth2.SetAuthURLParam("prompt", "none"),
							oauth2.SetAuthURLParam("max_age", "60"),
						)
						require.NotEmpty(t, code)
						token, err := conf.Exchange(context.Background(), code)
						require.NoError(t, err)
						followup := introspectAccessToken(t, conf, token, subject)
						assert.Equal(t, original.Get("auth_time").Int(), followup.Get("auth_time").Int())
					})
				}
			})

			t.Run("followup=fails when max age is reached and prompt is none", func(t *testing.T) {
				code, _ := getAuthorizeCode(t, conf, oc,
					oauth2.SetAuthURLParam("nonce", nonce),
					oauth2.SetAuthURLParam("prompt", "none"),
					oauth2.SetAuthURLParam("max_age", "1"),
				)
				require.Empty(t, code)
			})

			t.Run("followup=passes and resets skip when prompt=login", func(t *testing.T) {
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, subject, func(r *admin.GetLoginRequestOK) *models.AcceptLoginRequest {
						require.False(t, *r.Payload.Skip)
						require.Empty(t, *r.Payload.Subject)
						return nil
					}),
					acceptConsentHandler(t, c, subject, func(r *admin.GetConsentRequestOK) {
						require.True(t, r.Payload.Skip)
						require.EqualValues(t, subject, r.Payload.Subject)
					}),
				)
				code, _ := getAuthorizeCode(t, conf, oc,
					oauth2.SetAuthURLParam("nonce", nonce),
					oauth2.SetAuthURLParam("prompt", "login"),
					oauth2.SetAuthURLParam("max_age", "1"),
				)
				require.NotEmpty(t, code)
				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)
				introspectAccessToken(t, conf, token, subject)
				assertIDToken(t, token, conf, subject, nonce)
			})
		})
	})

	t.Run("case=should fail if prompt=none but no auth session given", func(t *testing.T) {
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, nil),
			acceptConsentHandler(t, c, subject, nil),
		)

		oc := testhelpers.NewEmptyJarClient(t)
		code, _ := getAuthorizeCode(t, conf, oc,
			oauth2.SetAuthURLParam("prompt", "none"),
		)
		require.Empty(t, code)
	})

	t.Run("case=requires re-authentication when id_token_hint is set to a user 'patrik-neu' but the session is 'aeneas-rekkas' and then fails because the user id from the log in endpoint is 'aeneas-rekkas'", func(t *testing.T) {
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, func(r *admin.GetLoginRequestOK) *models.AcceptLoginRequest {
				require.False(t, *r.Payload.Skip)
				require.Empty(t, *r.Payload.Subject)
				return nil
			}),
			acceptConsentHandler(t, c, subject, nil),
		)

		oc := testhelpers.NewEmptyJarClient(t)

		// Create login session for aeneas-rekkas
		code, _ := getAuthorizeCode(t, conf, oc)
		require.NotEmpty(t, code)

		// Perform authentication for aeneas-rekkas which fails because id_token_hint is patrik-neu
		code, _ = getAuthorizeCode(t, conf, oc,
			oauth2.SetAuthURLParam("id_token_hint", testhelpers.NewIDToken(t, reg, "patrik-neu")),
		)
		require.Empty(t, code)
	})

	t.Run("case=should not cause issues if max_age is very low and consent takes a long time", func(t *testing.T) {
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, func(r *admin.GetLoginRequestOK) *models.AcceptLoginRequest {
				time.Sleep(time.Second * 2)
				return nil
			}),
			acceptConsentHandler(t, c, subject, nil),
		)

		code, _ := getAuthorizeCode(t, conf, nil)
		require.NotEmpty(t, code)
	})

	t.Run("case=ensure consistent claims returned for userinfo", func(t *testing.T) {
		c, conf := newOAuth2Client(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, c, subject, nil),
			acceptConsentHandler(t, c, subject, nil),
		)

		code, _ := getAuthorizeCode(t, conf, nil)
		require.NotEmpty(t, code)

		token, err := conf.Exchange(context.Background(), code)
		require.NoError(t, err)

		idClaims := assertIDToken(t, token, conf, subject, "")

		time.Sleep(time.Second)
		uiClaims := testhelpers.Userinfo(t, token, publicTS)

		for _, f := range []string{
			"sub",
			"iss",
			"aud",
			"bar",
			"auth_time",
		} {
			assert.NotEmpty(t, uiClaims.Get(f).Raw, "%s: %s", f, uiClaims)
			assert.EqualValues(t, idClaims.Get(f).Raw, uiClaims.Get(f).Raw, "%s\nuserinfo: %s\nidtoken: %s", f, uiClaims, idClaims)
		}

		for _, f := range []string{
			"at_hash",
			"c_hash",
			"nonce",
			"sid",
			"jti",
		} {
			assert.Empty(t, uiClaims.Get(f).Raw, "%s: %s", f, uiClaims)
		}
	})
}

// TestAuthCodeWithMockStrategy runs the authorization_code flow against various ConsentStrategy scenarios.
// For that purpose, the consent strategy is mocked so all scenarios can be applied properly. This test suite checks:
//
// - [x] should pass request if strategy passes
// - [x] should fail because prompt=none and max_age > auth_time
// - [x] should pass because prompt=none and max_age < auth_time
// - [x] should fail because prompt=none but auth_time suggests recent authentication
// - [x] should fail because consent strategy fails
// - [x] should pass with prompt=login when authentication time is recent
// - [x] should fail with prompt=login when authentication time is in the past
func TestAuthCodeWithMockStrategy(t *testing.T) {
	for _, strat := range []struct{ d string }{{d: "opaque"}, {d: "jwt"}} {
		t.Run("strategy="+strat.d, func(t *testing.T) {
			conf := internal.NewConfigurationWithDefaults()
			conf.MustSet(config.KeyAccessTokenLifespan, time.Second*2)
			conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
			conf.MustSet(config.KeyAccessTokenStrategy, strat.d)
			reg := internal.NewRegistryMemory(t, conf)
			internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
			internal.MustEnsureRegistryKeys(reg, x.OAuth2JWTKeyName)

			consentStrategy := &consentMock{}
			router := x.NewRouterPublic()
			ts := httptest.NewServer(router)
			defer ts.Close()

			reg.WithConsentStrategy(consentStrategy)
			handler := reg.OAuth2Handler()
			handler.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
				return h
			})

			var callbackHandler *httprouter.Handle
			router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				(*callbackHandler)(w, r, ps)
			})
			var mutex sync.Mutex

			require.NoError(t, reg.ClientManager().CreateClient(context.TODO(), &hc.Client{
				OutfacingID:   "app-client",
				Secret:        "secret",
				RedirectURIs:  []string{ts.URL + "/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scope:         "hydra.* offline openid",
			}))

			oauthConfig := &oauth2.Config{
				ClientID:     "app-client",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:  ts.URL + "/oauth2/auth",
					TokenURL: ts.URL + "/oauth2/token",
				},
				RedirectURL: ts.URL + "/callback",
				Scopes:      []string{"hydra.*", "offline", "openid"},
			}

			var code string
			for k, tc := range []struct {
				cj                        http.CookieJar
				d                         string
				cb                        func(t *testing.T) httprouter.Handle
				authURL                   string
				shouldPassConsentStrategy bool
				expectOAuthAuthError      bool
				expectOAuthTokenError     bool
				checkExpiry               bool
				authTime                  time.Time
				requestTime               time.Time
				assertAccessToken         func(*testing.T, string)
			}{
				{
					d:                         "should pass request if strategy passes",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
					shouldPassConsentStrategy: true,
					checkExpiry:               true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
					assertAccessToken: func(t *testing.T, token string) {
						if strat.d != "jwt" {
							return
						}

						body, err := x.DecodeSegment(strings.Split(token, ".")[1])
						require.NoError(t, err)

						data := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &data))

						assert.EqualValues(t, "app-client", data["client_id"])
						assert.EqualValues(t, "foo", data["sub"])
						assert.NotEmpty(t, data["iss"])
						assert.NotEmpty(t, data["jti"])
						assert.NotEmpty(t, data["exp"])
						assert.NotEmpty(t, data["iat"])
						assert.NotEmpty(t, data["nbf"])
						assert.EqualValues(t, data["nbf"], data["iat"])
						assert.EqualValues(t, []interface{}{"offline", "openid", "hydra.*"}, data["scp"])
					},
				},
				{
					d:                         "should fail because prompt=none and max_age > auth_time",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none&max_age=1",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							err := r.URL.Query().Get("error")
							require.Empty(t, code)
							require.EqualValues(t, fosite.ErrLoginRequired.Error(), err)
						}
					},
					expectOAuthAuthError: true,
				},
				{
					d:                         "should pass because prompt=none and max_age < auth_time",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none&max_age=3600",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
				},
				{
					d:                         "should fail because prompt=none but auth_time suggests recent authentication",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC().Add(-time.Hour),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							err := r.URL.Query().Get("error")
							require.Empty(t, code)
							require.EqualValues(t, fosite.ErrLoginRequired.Error(), err)
						}
					},
					expectOAuthAuthError: true,
				},
				{
					d:                         "should fail because consent strategy fails",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
					expectOAuthAuthError:      true,
					shouldPassConsentStrategy: false,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							require.Empty(t, r.URL.Query().Get("code"))
							assert.Equal(t, fosite.ErrRequestForbidden.Error(), r.URL.Query().Get("error"))
						}
					},
				},
				{
					d:                         "should pass with prompt=login when authentication time is recent",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=login",
					authTime:                  time.Now().UTC().Add(-time.Second),
					requestTime:               time.Now().UTC().Add(-time.Minute),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
				},
				{
					d:                         "should fail with prompt=login when authentication time is in the past",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=login",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					expectOAuthAuthError:      true,
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.Empty(t, code)
							assert.Equal(t, fosite.ErrLoginRequired.Error(), r.URL.Query().Get("error"))
						}
					},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
					mutex.Lock()
					defer mutex.Unlock()
					if tc.cb == nil {
						tc.cb = noopHandler
					}

					consentStrategy.deny = !tc.shouldPassConsentStrategy
					consentStrategy.authTime = tc.authTime
					consentStrategy.requestTime = tc.requestTime

					cb := tc.cb(t)
					callbackHandler = &cb

					req, err := http.NewRequest("GET", tc.authURL, nil)
					require.NoError(t, err)

					if tc.cj == nil {
						tc.cj = testhelpers.NewEmptyCookieJar(t)
					}

					resp, err := (&http.Client{Jar: tc.cj}).Do(req)
					require.NoError(t, err, tc.authURL, ts.URL)
					defer resp.Body.Close()

					if tc.expectOAuthAuthError {
						require.Empty(t, code)
						return
					}

					require.NotEmpty(t, code)

					token, err := oauthConfig.Exchange(oauth2.NoContext, code)

					if tc.expectOAuthTokenError {
						require.Error(t, err)
						return
					}

					if tc.assertAccessToken != nil {
						tc.assertAccessToken(t, token.AccessToken)
					}

					require.NoError(t, err, code)

					t.Run("case=userinfo", func(t *testing.T) {
						var makeRequest = func(req *http.Request) *http.Response {
							resp, err = http.DefaultClient.Do(req)
							require.NoError(t, err)
							return resp
						}

						var testSuccess = func(response *http.Response) {
							defer resp.Body.Close()

							require.Equal(t, http.StatusOK, resp.StatusCode)

							var claims map[string]interface{}
							require.NoError(t, json.NewDecoder(resp.Body).Decode(&claims))
							assert.Equal(t, "foo", claims["sub"])
						}

						req, err = http.NewRequest("GET", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("POST", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("POST", ts.URL+"/userinfo", bytes.NewBuffer([]byte("access_token="+token.AccessToken)))
						req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("GET", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer asdfg")
						resp := makeRequest(req)
						require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
					})
					t.Logf("Got token: %s", token.AccessToken)

					// time.Sleep(time.Millisecond * 1200) // Makes sure exp/iat/nbf time is different later on

					res, err := testRefresh(t, token, ts.URL, tc.checkExpiry)
					require.NoError(t, err)
					assert.Equal(t, http.StatusOK, res.StatusCode)

					body, err := ioutil.ReadAll(res.Body)
					require.NoError(t, err)

					var refreshedToken oauth2.Token
					require.NoError(t, json.Unmarshal(body, &refreshedToken))

					t.Logf("Got refresh token: %s", refreshedToken.AccessToken)

					if tc.assertAccessToken != nil {
						tc.assertAccessToken(t, refreshedToken.AccessToken)
					}

					t.Run("the tokens should be different", func(t *testing.T) {
						if strat.d != "jwt" {
							t.Skip()
						}

						body, err := x.DecodeSegment(strings.Split(token.AccessToken, ".")[1])
						require.NoError(t, err)

						origPayload := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &origPayload))

						body, err = x.DecodeSegment(strings.Split(refreshedToken.AccessToken, ".")[1])
						require.NoError(t, err)

						refreshedPayload := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &refreshedPayload))

						if tc.checkExpiry {
							assert.NotEqual(t, refreshedPayload["exp"], origPayload["exp"])
							assert.NotEqual(t, refreshedPayload["iat"], origPayload["iat"])
							assert.NotEqual(t, refreshedPayload["nbf"], origPayload["nbf"])
						}
						assert.NotEqual(t, refreshedPayload["jti"], origPayload["jti"])
						assert.Equal(t, refreshedPayload["client_id"], origPayload["client_id"])
					})

					require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)

					t.Run("old token should no longer be usable", func(t *testing.T) {
						req, err := http.NewRequest("GET", ts.URL+"/userinfo", nil)
						require.NoError(t, err)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						res, err := http.DefaultClient.Do(req)
						require.NoError(t, err)
						assert.EqualValues(t, http.StatusUnauthorized, res.StatusCode)
					})

					t.Run("refreshing new refresh token should work", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusOK, res.StatusCode)

						body, err := ioutil.ReadAll(res.Body)
						require.NoError(t, err)
						require.NoError(t, json.Unmarshal(body, &refreshedToken))
					})

					t.Run("should call refresh token hook if configured", func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

							var hookReq hydraoauth2.RefreshTokenHookRequest
							require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
							require.Equal(t, hookReq.Subject, "foo")
							require.ElementsMatch(t, hookReq.GrantedScopes, []string{"openid", "offline", "hydra.*"})
							require.ElementsMatch(t, hookReq.GrantedAudience, []string{})
							require.Equal(t, hookReq.ClientID, oauthConfig.ClientID)

							claims := map[string]interface{}{
								"hooked": true,
							}

							hookResp := hydraoauth2.RefreshTokenHookResponse{
								Session: consent.ConsentRequestSessionData{
									AccessToken: claims,
									IDToken:     claims,
								},
							}

							w.WriteHeader(http.StatusOK)
							require.NoError(t, json.NewEncoder(w).Encode(&hookResp))
						}))
						defer hs.Close()

						conf.MustSet(config.KeyRefreshTokenHookURL, hs.URL)
						defer conf.MustSet(config.KeyRefreshTokenHookURL, nil)

						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusOK, res.StatusCode)

						body, err := ioutil.ReadAll(res.Body)
						require.NoError(t, err)
						require.NoError(t, json.Unmarshal(body, &refreshedToken))

						accessTokenClaims := testhelpers.IntrospectToken(t, oauthConfig, &refreshedToken, ts)
						require.True(t, accessTokenClaims.Get("ext.hooked").Bool())

						idTokenBody, err := x.DecodeSegment(
							strings.Split(
								gjson.GetBytes(body, "id_token").String(),
								".",
							)[1],
						)
						require.NoError(t, err)

						require.True(t, gjson.GetBytes(idTokenBody, "hooked").Bool())
					})

					t.Run("should fail token refresh with `server_error` if hook fails", func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusInternalServerError)
						}))
						defer hs.Close()

						conf.MustSet(config.KeyRefreshTokenHookURL, hs.URL)
						defer conf.MustSet(config.KeyRefreshTokenHookURL, nil)

						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

						var errBody fosite.RFC6749ErrorJson
						require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
						require.Equal(t, fosite.ErrServerError.Error(), errBody.Name)
						require.Equal(t, fosite.ErrServerError.GetDescription(), errBody.Description)
					})

					t.Run("should fail token refresh with `access_denied` if hook denied the request", func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusForbidden)
						}))
						defer hs.Close()

						conf.MustSet(config.KeyRefreshTokenHookURL, hs.URL)
						defer conf.MustSet(config.KeyRefreshTokenHookURL, nil)

						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusForbidden, res.StatusCode)

						var errBody fosite.RFC6749ErrorJson
						require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
						require.Equal(t, fosite.ErrAccessDenied.Error(), errBody.Name)
						require.Equal(t, fosite.ErrAccessDenied.GetDescription(), errBody.Description)
					})

					t.Run("should fail token refresh with `server_error` if hook response is malformed", func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
						}))
						defer hs.Close()

						conf.MustSet(config.KeyRefreshTokenHookURL, hs.URL)
						defer conf.MustSet(config.KeyRefreshTokenHookURL, nil)

						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

						var errBody fosite.RFC6749ErrorJson
						require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
						require.Equal(t, fosite.ErrServerError.Error(), errBody.Name)
						require.Equal(t, fosite.ErrServerError.GetDescription(), errBody.Description)
					})

					t.Run("refreshing old token should no longer work", func(t *testing.T) {
						res, err := testRefresh(t, token, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
					})

					t.Run("attempt to refresh old token should revoke new token", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
					})

					t.Run("duplicate code exchange fails", func(t *testing.T) {
						token, err := oauthConfig.Exchange(oauth2.NoContext, code)
						require.Error(t, err)
						require.Nil(t, token)
					})

					code = ""
				})
			}
		})
	}
}

func testRefresh(t *testing.T, token *oauth2.Token, u string, sleep bool) (*http.Response, error) {
	if sleep {
		time.Sleep(time.Millisecond * 1001)
	}

	oauthClientConfig := &clientcredentials.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		TokenURL:     u + "/oauth2/token",
		Scopes:       []string{"foobar"},
	}

	req, err := http.NewRequest("POST", oauthClientConfig.TokenURL, strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{token.RefreshToken},
	}.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthClientConfig.ClientID, oauthClientConfig.ClientSecret)

	return http.DefaultClient.Do(req)
}

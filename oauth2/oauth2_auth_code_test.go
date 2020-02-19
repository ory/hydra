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
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	djwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/x"
	"github.com/ory/viper"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/ory/x/urlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func newCookieJar() http.CookieJar {
	c, _ := cookiejar.New(nil)
	return c
}

func noopHandler(t *testing.T) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
func noopHandlerDefaultStrategy(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func mockProvider(h *func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(*h)(w, r)
	}))
}

type clientCreator interface {
	CreateClient(cxt context.Context, client *hc.Client) error
}

func acceptLogin(apiClient *hydra.OryHydra, subject string, expectSkip bool, expectSubject string) func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			rrr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
			require.NoError(t, err)

			rr := rrr.Payload
			assert.Equal(t, expectSkip, rr.Skip)
			assert.EqualValues(t, expectSubject, rr.Subject)

			vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(&models.AcceptLoginRequest{Subject: pointerx.String(subject)}))
			require.NoError(t, err)

			v := vr.Payload
			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}
}

func acceptConsent(apiClient *hydra.OryHydra, scope []string, expectSkip bool, expectSubject string) func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
			require.NoError(t, err)

			rr := rrr.Payload
			assert.Equal(t, expectSkip, rr.Skip)
			assert.EqualValues(t, expectSubject, rr.Subject)

			vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
				WithBody(&models.AcceptConsentRequest{
					GrantScope: scope,
					Session: &models.ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IDToken:     map[string]interface{}{"bar": "baz"},
					},
				}))
			require.NoError(t, err)
			v := vr.Payload

			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}
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
	conf := internal.NewConfigurationWithDefaults()
	regs := map[string]driver.Registry{
		"memory": internal.NewRegistry(conf),
	}

	if !testing.Short() {
		var p, m, c *sqlx.DB
		dockertest.Parallel([]func(){
			func() {
				p = connectToPG(t)
			},
			func() {
				m = connectToMySQL(t)
			},
			func() {
				c = connectToCRDB(t)
			},
		})
		pr := internal.NewRegistrySQL(conf, p)
		_, err := pr.CreateSchemas("postgres")
		require.NoError(t, err)
		regs["postgres"] = pr

		mr := internal.NewRegistrySQL(conf, m)
		_, err = mr.CreateSchemas("mysql")
		require.NoError(t, err)
		regs["mysql"] = mr

		cr := internal.NewRegistrySQL(conf, c)
		_, err = cr.CreateSchemas("cockroach")
		require.NoError(t, err)
		regs["cockroach"] = cr
	}

	for km, reg := range regs {
		t.Run("manager="+km, func(t *testing.T) {
			for _, strat := range []struct{ d string }{{d: "opaque"}, {d: "jwt"}} {
				t.Run("strategy="+strat.d, func(t *testing.T) {
					conf := internal.NewConfigurationWithDefaults()
					viper.Set(configuration.ViperKeyAccessTokenStrategy, strat.d)
					viper.Set(configuration.ViperKeyScopeStrategy, "exact")
					viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
					viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Second+time.Millisecond*500)
					viper.Set(configuration.ViperKeyRefreshTokenLifespan, time.Second+time.Millisecond*800)
					// SendDebugMessagesToClients: true,
					internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)

					reg.WithConfig(conf)

					var m sync.Mutex
					l := logrus.New()
					l.Level = logrus.DebugLevel
					var lph, cph func(w http.ResponseWriter, r *http.Request)
					lp := mockProvider(&lph)
					defer lp.Close()
					cp := mockProvider(&cph)
					defer lp.Close()

					fooUserIDToken, _, err := reg.OpenIDJWTStrategy().Generate(context.TODO(), jwt.IDTokenClaims{
						Subject:   "foouser",
						ExpiresAt: time.Now().Add(time.Hour),
						IssuedAt:  time.Now(),
					}.ToMapClaims(), jwt.NewHeaders())
					require.NoError(t, err)

					router := x.NewRouterPublic()
					var callbackHandler *httprouter.Handle
					router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
						(*callbackHandler)(w, r, ps)
					})
					reg.OAuth2Handler().SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
						return h
					})

					ts := httptest.NewServer(router)
					defer ts.Close()

					apiRouter := x.NewRouterAdmin()
					reg.ConsentHandler().SetRoutes(apiRouter)
					api := httptest.NewServer(apiRouter)
					defer api.Close()

					viper.Set(configuration.ViperKeyLoginURL, lp.URL)
					viper.Set(configuration.ViperKeyConsentURL, cp.URL)
					viper.Set(configuration.ViperKeyIssuerURL, ts.URL)
					viper.Set(configuration.ViperKeyConsentRequestMaxAge, time.Hour)

					client := hc.Client{
						ClientID: "e2e-app-client" + km + strat.d, Secret: "secret", RedirectURIs: []string{ts.URL + "/callback"},
						ResponseTypes: []string{"id_token", "code", "token"},
						GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
						Scope:         "hydra offline openid",
						Audience:      []string{"https://api.ory.sh/"},
					}
					oauthConfig := &oauth2.Config{
						ClientID: client.GetID(), ClientSecret: client.Secret,
						Endpoint:    oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"},
						RedirectURL: client.RedirectURIs[0], Scopes: []string{"hydra", "offline", "openid"},
					}

					require.NoError(t, reg.OAuth2Storage().(clientCreator).CreateClient(context.TODO(), &client))
					apiClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(api.URL).Host})

					persistentCJ := newCookieJar()
					var code string
					for k, tc := range []struct {
						authURL               string
						cj                    http.CookieJar
						d                     string
						cb                    func(t *testing.T) httprouter.Handle
						expectOAuthAuthError  bool
						expectOAuthTokenError bool
						lph, cph              func(t *testing.T) func(w http.ResponseWriter, r *http.Request)
						setup                 func()
						expectRefreshToken    bool
						expectIDToken         bool

						assertAccessToken, assertIDToken func(*testing.T, string)
						assertRefreshToken               func(*testing.T, *oauth2.Token)
					}{
						{
							d:       "Checks if request fails when audience doesn't match",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("audience", "https://not-ory-api/")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) { t.Fatal("This should not have been called") }
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) { t.Fatal("This should not have been called") }
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.Empty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  true,
							expectOAuthTokenError: false,
							expectIDToken:         false,
							expectRefreshToken:    false,
							assertAccessToken: func(t *testing.T, token string) {
								require.Empty(t, token)
							},
						},
						{
							d:       "Perform OAuth2 flow with openid connect id token and verify the id token",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("nonce", "what-a-cool-nonce")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)
									assert.False(t, rr.Payload.Skip)
									assert.Empty(t, rr.Payload.Subject)
									v, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: true, RememberFor: 0, Acr: "1",
											Context: map[string]interface{}{"context": "bar"},
										}))
									require.NoError(t, err)
									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)
									assert.False(t, rr.Payload.Skip)
									assert.Equal(t, map[string]interface{}{"context": "bar"}, rr.Payload.Context)
									v, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
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
									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									assert.NotEmpty(t, code, "%s", r.URL.String())
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    true,
							assertIDToken: func(t *testing.T, token string) {
								body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
								require.NoError(t, err)

								data := map[string]interface{}{}
								require.NoError(t, json.Unmarshal(body, &data))
								assert.EqualValues(t, "user-a", data["sub"])
								assert.EqualValues(t, "1", data["acr"])
								assert.EqualValues(t, fmt.Sprintf("%v", []string{client.GetID()}), fmt.Sprintf("%v", data["aud"]))
								assert.NotEmpty(t, client.GetID(), data["exp"])
								assert.NotEmpty(t, client.GetID(), data["iat"])
								assert.NotEmpty(t, client.GetID(), data["jti"])
								assert.NotEmpty(t, "what-a-cool-nonce", data["nonce"])
							},
						},
						{
							d:       "Perform OAuth2 flow with refreshing which fails due to expiry",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("nonce", "what-a-cool-nonce")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)
									v, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: true, RememberFor: 0, Acr: "1",
										}))
									require.NoError(t, err)
									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "offline", "openid"}, Remember: false, RememberFor: 0,
											GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
											Session:                  &models.ConsentRequestSession{},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    true,
							assertRefreshToken: func(t *testing.T, token *oauth2.Token) {
								time.Sleep(viper.GetDuration(configuration.ViperKeyRefreshTokenLifespan) + time.Second)
								token.Expiry = token.Expiry.Add(-time.Hour * 24)
								_, err := oauthConfig.TokenSource(oauth2.NoContext, token).Token()
								require.Error(t, err)
							},
						}, {
							d:       "Checks if request fails when subject is empty",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("audience", "https://not-ory-api/")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									challenge, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)
									assert.False(t, challenge.Payload.Skip)

									v, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String(""), Remember: false, RememberFor: 0, Acr: "1",
										}))
									require.Error(t, err)
									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) { t.Fatal("This should not have been called") }
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.Empty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  true,
							expectOAuthTokenError: false,
							expectIDToken:         false,
							expectRefreshToken:    false,
							assertAccessToken: func(t *testing.T, token string) {
								require.Empty(t, token)
							},
						},
						{
							d:       "Perform OAuth2 flow with refreshing which works just fine",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("nonce", "what-a-cool-nonce")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)

									v, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: false, RememberFor: 0, Acr: "1",
										}))
									require.NoError(t, err)

									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "offline", "openid"}, Remember: false, RememberFor: 0,
											GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
											Session:                  &models.ConsentRequestSession{},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    true,
							assertRefreshToken: func(t *testing.T, token *oauth2.Token) {
								token.Expiry = token.Expiry.Add(-time.Hour * 24)
								n, err := oauthConfig.TokenSource(oauth2.NoContext, token).Token()
								require.NoError(t, err)
								require.NotEqual(t, token.AccessToken, n.AccessToken)
								require.NotEqual(t, token.RefreshToken, n.RefreshToken)
							},
						},
						{
							d:       "Perform OAuth2 flow with audience",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("audience", "https://api.ory.sh/")),
							cj:      newCookieJar(),
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload

									require.NoError(t, err)
									assert.False(t, rr.Skip)
									assert.Empty(t, rr.Subject)
									assert.EqualValues(t, []string{"https://api.ory.sh/"}, rr.RequestedAccessTokenAudience)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientID)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoURI)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, r.URL.Query().Get("login_challenge"), rr.Challenge)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("audience", "https://api.ory.sh/")), rr.RequestURL)

									v, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: true, RememberFor: 0,
										}))

									require.NoError(t, err)
									require.NotEmpty(t, v.Payload.RedirectTo)
									http.Redirect(w, r, v.Payload.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)

									rr := rrr.Payload
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientID)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoURI)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, []string{"https://api.ory.sh/"}, rr.RequestedAccessTokenAudience)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), rr.Challenge)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state", oauth2.SetAuthURLParam("audience", "https://api.ory.sh/")), rr.RequestURL)

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "offline", "openid"}, Remember: false, RememberFor: 0,
											GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
											Session: &models.ConsentRequestSession{
												AccessToken: map[string]interface{}{"foo": "bar"},
												IDToken:     map[string]interface{}{"bar": "baz"},
											},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    true,
							assertAccessToken: func(t *testing.T, token string) {
								if strat.d != "jwt" {
									res, err := ts.Client().PostForm(ts.URL+"/oauth2/introspect", url.Values{"token": {token}})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)

									body, err := ioutil.ReadAll(res.Body)
									require.NoError(t, err)

									var r models.OAuth2TokenIntrospection
									require.NoError(t, json.Unmarshal(body, &r))
									assert.EqualValues(t, "e2e-app-client"+km+strat.d, r.ClientID)
									assert.EqualValues(t, "user-a", r.Sub)
									assert.EqualValues(t, []string{"https://api.ory.sh/"}, r.Aud)
									assert.EqualValues(t, "hydra offline openid", r.Scope)
									assert.EqualValues(t, "map[foo:bar]", fmt.Sprintf("%s", r.Ext))
									return
								}

								body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
								require.NoError(t, err)

								data := map[string]interface{}{}
								require.NoError(t, json.Unmarshal(body, &data))

								assert.EqualValues(t, "e2e-app-client"+km+strat.d, data["client_id"])
								assert.EqualValues(t, "user-a", data["sub"])
								assert.NotEmpty(t, data["iss"])
								assert.NotEmpty(t, data["jti"])
								assert.NotEmpty(t, data["exp"])
								assert.NotEmpty(t, data["iat"])
								assert.NotEmpty(t, data["nbf"])
								assert.EqualValues(t, data["nbf"], data["iat"])
								assert.EqualValues(t, []interface{}{"https://api.ory.sh/"}, data["aud"])
								assert.EqualValues(t, []interface{}{"hydra", "offline", "openid"}, data["scp"])
								assert.EqualValues(t, "map[foo:bar]", fmt.Sprintf("%s", data["ext"]))
							},
						},
						{
							// First we need to create a persistent session in order to check if the other things work
							// as expected
							d:       "Creates a persisting session for the next test cases",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=login+consent&max_age=1",
							cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload

									assert.False(t, rr.Skip)
									assert.Empty(t, rr.Subject)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientID)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoURI)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, r.URL.Query().Get("login_challenge"), rr.Challenge)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state")+"&prompt=login+consent&max_age=1", rr.RequestURL)

									vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: true, RememberFor: 0,
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)

									rr := rrr.Payload
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientID)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoURI)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), rr.Challenge)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state")+"&prompt=login+consent&max_age=1", rr.RequestURL)

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "offline", "openid"}, Remember: true, RememberFor: 0,
											Session: &models.ConsentRequestSession{
												AccessToken: map[string]interface{}{"foo": "bar"},
												IDToken:     map[string]interface{}{"bar": "baz"},
											},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    true,
							assertAccessToken: func(t *testing.T, token string) {
								if strat.d != "jwt" {
									return
								}

								body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
								require.NoError(t, err)

								data := map[string]interface{}{}
								require.NoError(t, json.Unmarshal(body, &data))

								assert.EqualValues(t, "e2e-app-client"+km+strat.d, data["client_id"])
								assert.EqualValues(t, "user-a", data["sub"])
								assert.NotEmpty(t, data["iss"])
								assert.NotEmpty(t, data["jti"])
								assert.NotEmpty(t, data["exp"])
								assert.NotEmpty(t, data["iat"])
								assert.NotEmpty(t, data["nbf"])
								assert.EqualValues(t, data["nbf"], data["iat"])
								assert.EqualValues(t, []interface{}{"hydra", "offline", "openid"}, data["scp"])
								assert.EqualValues(t, "map[foo:bar]", fmt.Sprintf("%s", data["ext"]))
							},
						},
						{
							d: "checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=none works",
							setup: func() {
								// In order to check if authenticatedAt/requestedAt works, we'll sleep first in order to ensure that authenticatedAt is in the past
								// if handled correctly.
								time.Sleep(time.Second + time.Millisecond)
							},
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=60",
							cj:      persistentCJ,
							lph:     acceptLogin(apiClient, "user-a", true, "user-a"),
							cph:     acceptConsent(apiClient, []string{"hydra", "offline"}, true, "user-a"),
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         false,
							expectRefreshToken:    true,
							assertAccessToken: func(t *testing.T, token string) {
								if strat.d != "jwt" {
									res, err := ts.Client().PostForm(ts.URL+"/oauth2/introspect", url.Values{"token": {token}})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)

									body, err := ioutil.ReadAll(res.Body)
									require.NoError(t, err)

									var r models.OAuth2TokenIntrospection
									require.NoError(t, json.Unmarshal(body, &r))
									assert.EqualValues(t, "e2e-app-client"+km+strat.d, r.ClientID)
									assert.EqualValues(t, "user-a", r.Sub)
									assert.Empty(t, r.Aud)
									assert.EqualValues(t, "hydra offline", r.Scope)
									assert.EqualValues(t, "map[foo:bar]", fmt.Sprintf("%s", r.Ext))
									return
								}

								body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
								require.NoError(t, err)

								data := map[string]interface{}{}
								require.NoError(t, json.Unmarshal(body, &data))

								assert.EqualValues(t, "e2e-app-client"+km+strat.d, data["client_id"])
								assert.EqualValues(t, "user-a", data["sub"])
								assert.NotEmpty(t, data["iss"])
								assert.NotEmpty(t, data["jti"])
								assert.NotEmpty(t, data["exp"])
								assert.NotEmpty(t, data["iat"])
								assert.NotEmpty(t, data["nbf"])
								assert.Empty(t, data["aud"])
								assert.EqualValues(t, data["nbf"], data["iat"])
								assert.EqualValues(t, []interface{}{"hydra", "offline"}, data["scp"])
								assert.EqualValues(t, "map[foo:bar]", fmt.Sprintf("%s", data["ext"]))
							},
						},
						{
							d:                     "checks if prompt=none fails when no session is set",
							authURL:               oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none",
							cj:                    newCookieJar(),
							expectOAuthAuthError:  true,
							expectOAuthTokenError: false,
						},
						{
							d:                     "checks if consecutive authentications cause any issues (1/3)",
							authURL:               oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=60",
							cj:                    persistentCJ,
							lph:                   acceptLogin(apiClient, "user-a", true, "user-a"),
							cph:                   acceptConsent(apiClient, []string{"hydra", "offline"}, true, "user-a"),
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectRefreshToken:    true,
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
						},
						{
							d:                     "checks if consecutive authentications cause any issues (2/3)",
							authURL:               oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=60",
							cj:                    persistentCJ,
							lph:                   acceptLogin(apiClient, "user-a", true, "user-a"),
							cph:                   acceptConsent(apiClient, []string{"hydra", "offline"}, true, "user-a"),
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectRefreshToken:    true,
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
						},
						{
							d:                     "checks if consecutive authentications cause any issues (3/3)",
							authURL:               oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=60",
							cj:                    persistentCJ,
							lph:                   acceptLogin(apiClient, "user-a", true, "user-a"),
							cph:                   acceptConsent(apiClient, []string{"hydra", "offline"}, true, "user-a"),
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectRefreshToken:    true,
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
						},
						{
							d:                     "checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=none fails when maxAge is reached",
							authURL:               oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=1",
							cj:                    persistentCJ,
							expectOAuthAuthError:  true,
							expectOAuthTokenError: false,
						},
						{
							d:       "checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=login works",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=login",
							cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)

									rr := rrr.Payload
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "", rr.Subject)

									vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"),
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload

									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "openid"},
											Session: &models.ConsentRequestSession{
												AccessToken: map[string]interface{}{"foo": "bar"},
												IDToken:     map[string]interface{}{"bar": "baz"},
											},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    false,
						},
						{
							d:       "requires re-authentication when id_token_hint is set to a user (\"foouser\") but the session is \"user-a\" and it also fails because the user id from the log in endpoint is not foouser",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&id_token_hint=" + fooUserIDToken,
							cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)

									rr := rrr.Payload
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "", rr.Subject)
									assert.Empty(t, rr.Client.ClientSecret)

									vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"),
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)
									rr := rrr.Payload
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)
									assert.Empty(t, rr.Client.ClientSecret)

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "openid"},
											Session: &models.ConsentRequestSession{
												AccessToken: map[string]interface{}{"foo": "bar"},
												IDToken:     map[string]interface{}{"bar": "baz"},
											},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.Empty(t, code)
									require.EqualValues(t, "login_required", r.URL.Query().Get("error"), r.URL.Query().Get("error_debug"))
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  true,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    false,
							assertAccessToken: func(t *testing.T, token string) {
								if strat.d != "jwt" {
									return
								}

								body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
								require.NoError(t, err)

								data := map[string]interface{}{}
								require.NoError(t, json.Unmarshal(body, &data))

								assert.EqualValues(t, "e2e-app-client"+km+strat.d, data["client_id"])
								assert.EqualValues(t, "user-a", data["sub"])
								assert.NotEmpty(t, data["iss"])
								assert.NotEmpty(t, data["jti"])
								assert.NotEmpty(t, data["exp"])
								assert.NotEmpty(t, data["iat"])
								assert.NotEmpty(t, data["nbf"])
								assert.EqualValues(t, data["nbf"], data["iat"])
								assert.EqualValues(t, []interface{}{"hydra", "openid"}, data["scp"])
								assert.EqualValues(t, "bar", data["foo"])
							},
						},
						{
							d:       "should not cause issues if max_age is very low and consent takes a long time",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&max_age=1",
							// cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
									require.NoError(t, err)

									time.Sleep(time.Second * 2)

									vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
										WithLoginChallenge(r.URL.Query().Get("login_challenge")).
										WithBody(&models.AcceptLoginRequest{
											Subject: pointerx.String("user-a"), Remember: true, RememberFor: 0, Acr: "1",
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
									require.NoError(t, err)

									vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
										WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
										WithBody(&models.AcceptConsentRequest{
											GrantScope: []string{"hydra", "openid"},
											Session: &models.ConsentRequestSession{
												AccessToken: map[string]interface{}{"foo": "bar"},
												IDToken:     map[string]interface{}{"bar": "baz"},
											},
										}))
									require.NoError(t, err)
									v := vr.Payload

									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cb: func(t *testing.T) httprouter.Handle {
								return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
									code = r.URL.Query().Get("code")
									require.NotEmpty(t, code)
									w.WriteHeader(http.StatusOK)
								}
							},
							expectOAuthAuthError:  false,
							expectOAuthTokenError: false,
							expectIDToken:         true,
							expectRefreshToken:    false,
						},
					} {
						t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
							m.Lock()
							defer m.Unlock()

							code = ""

							if tc.setup != nil {
								tc.setup()
							}

							if tc.lph != nil {
								lph = tc.lph(t)
							} else {
								lph = noopHandlerDefaultStrategy(t)
							}

							if tc.cph != nil {
								cph = tc.cph(t)
							} else {
								cph = noopHandlerDefaultStrategy(t)
							}

							if tc.cb == nil {
								tc.cb = noopHandler
							}

							cb := tc.cb(t)
							callbackHandler = &cb

							req, err := http.NewRequest("GET", tc.authURL, nil)
							require.NoError(t, err)

							if tc.cj == nil {
								tc.cj = newCookieJar()
							}

							resp, err := (&http.Client{Jar: tc.cj}).Do(req)
							require.NoError(t, err)
							defer resp.Body.Close()

							t.Logf("Cookies: %+v", tc.cj)

							time.Sleep(time.Millisecond * 5)

							if tc.expectOAuthAuthError {
								require.Empty(t, code)
								return
							}

							var body []byte
							if code == "" {
								body, _ = ioutil.ReadAll(resp.Body)
							}
							require.NotEmpty(t, code, "body: %s\nreq: %s\nts: %s", body, req.URL.String(), ts.URL)

							token, err := oauthConfig.Exchange(oauth2.NoContext, code)

							if tc.expectOAuthTokenError {
								require.Error(t, err)
								return
							}

							require.NoError(t, err, code)
							assert.NotEmpty(t, token.AccessToken)
							if tc.expectRefreshToken {
								require.NotEmpty(t, token.RefreshToken)
							} else {
								require.Empty(t, token.RefreshToken)
							}
							if tc.expectIDToken {
								require.NotEmpty(t, token.Extra("id_token"))
							} else {
								require.Empty(t, token.Extra("id_token"))
							}
							if tc.assertAccessToken != nil {
								tc.assertAccessToken(t, token.AccessToken)
							}
							if tc.assertIDToken != nil {
								tc.assertIDToken(t, token.Extra("id_token").(string))
							}
							if tc.assertRefreshToken != nil {
								tc.assertRefreshToken(t, token)
							}
						})
					}
				})
			}
		})
	}
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
			viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Second*2)
			viper.Set(configuration.ViperKeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
			viper.Set(configuration.ViperKeyAccessTokenStrategy, strat.d)
			reg := internal.NewRegistry(conf)
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
				ClientID:      "app-client",
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

						body, err := djwt.DecodeSegment(strings.Split(token, ".")[1])
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
						tc.cj = newCookieJar()
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

						body, err := djwt.DecodeSegment(strings.Split(token.AccessToken, ".")[1])
						require.NoError(t, err)

						origPayload := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &origPayload))

						body, err = djwt.DecodeSegment(strings.Split(refreshedToken.AccessToken, ".")[1])
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

					t.Run("refreshing old token should no longer work", func(t *testing.T) {
						res, err := testRefresh(t, token, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, res.StatusCode)
					})

					t.Run("refreshing new refresh token should work", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusOK, res.StatusCode)
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

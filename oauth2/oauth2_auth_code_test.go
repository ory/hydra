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
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/herodot"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/square/go-jose.v2"
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
	CreateClient(client *hc.Client) error
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
	for _, strat := range []struct {
		d string
		s foauth2.CoreStrategy
	}{
		{
			d: "opaque",
			s: oauth2OpqaueStrategy,
		},
		{
			d: "jwt",
			s: oauth2JWTStrategy,
		},
	} {
		t.Run("strategy="+strat.d, func(t *testing.T) {
			var m sync.Mutex
			l := logrus.New()
			l.Level = logrus.DebugLevel
			var lph, cph func(w http.ResponseWriter, r *http.Request)
			lp := mockProvider(&lph)
			cp := mockProvider(&cph)
			jwts := &jwt.RS256JWTStrategy{
				PrivateKey: pkg.MustINSECURELOWENTROPYRSAKEYFORTEST(),
			}
			hasher := &fosite.BCrypt{
				WorkFactor: 4,
			}

			fooUserIDToken, _, err := jwts.Generate((jwt.IDTokenClaims{
				Subject:   "foouser",
				ExpiresAt: time.Now().Add(time.Hour),
				IssuedAt:  time.Now(),
			}).ToMapClaims(), jwt.NewHeaders())
			require.NoError(t, err)

			// we create a new fositeStore here because the old one

			for km, fs := range fositeStores {
				t.Run("manager="+km, func(t *testing.T) {
					var cm consent.Manager
					switch km {
					case "memory":
						cm = consent.NewMemoryManager(fs)
						fs.(*FositeMemoryStore).Manager = hc.NewMemoryManager(hasher)
					case "mysql":
						fallthrough
					case "postgres":
						scm := consent.NewSQLManager(databases[km], fs.(*FositeSQLStore).Manager, fs)
						_, err := scm.CreateSchemas()
						require.NoError(t, err)

						_, err = (fs.(*FositeSQLStore)).CreateSchemas()
						require.NoError(t, err)

						cm = scm
					}

					router := httprouter.New()
					ts := httptest.NewServer(router)
					cookieStore := sessions.NewCookieStore([]byte("foo-secret"))

					consentStrategy := consent.NewStrategy(
						lp.URL, cp.URL, ts.URL, "/oauth2/auth", cm,
						cookieStore,
						fosite.ExactScopeStrategy, false, time.Hour, jwts,
						openid.NewOpenIDConnectRequestValidator(nil, jwts),
					)

					jm := &jwk.MemoryManager{Keys: map[string]*jose.JSONWebKeySet{}}
					keys, err := (&jwk.RS256Generator{}).Generate("", "sig")
					require.NoError(t, err)
					require.NoError(t, jm.AddKeySet(OpenIDConnectKeyName, keys))
					jwtStrategy, err := jwk.NewRS256JWTStrategy(jm, OpenIDConnectKeyName)

					handler := &Handler{
						OAuth2: compose.Compose(
							fc, fs, strat.s, hasher,
							compose.OAuth2AuthorizeExplicitFactory,
							compose.OAuth2AuthorizeImplicitFactory,
							compose.OAuth2ClientCredentialsGrantFactory,
							compose.OAuth2RefreshTokenGrantFactory,
							compose.OpenIDConnectExplicitFactory,
							compose.OpenIDConnectHybridFactory,
							compose.OpenIDConnectImplicitFactory,
							compose.OAuth2TokenRevocationFactory,
							compose.OAuth2TokenIntrospectionFactory,
						),
						Consent:         consentStrategy,
						CookieStore:     cookieStore,
						H:               herodot.NewJSONWriter(l),
						ScopeStrategy:   fosite.ExactScopeStrategy,
						IDTokenLifespan: time.Minute, IssuerURL: ts.URL, ForcedHTTP: true, L: l,
						OpenIDJWTStrategy: jwtStrategy,
					}
					handler.SetRoutes(router,router)

					apiHandler := consent.NewHandler(herodot.NewJSONWriter(l), cm)
					apiRouter := httprouter.New()
					apiHandler.SetRoutes(apiRouter)
					api := httptest.NewServer(apiRouter)

					client := hc.Client{
						ClientID: "e2e-app-client" + km + strat.d, Secret: "secret", RedirectURIs: []string{ts.URL + "/callback"},
						ResponseTypes: []string{"id_token", "code", "token"},
						GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
						Scope:         "hydra offline openid",
					}
					oauthConfig := &oauth2.Config{
						ClientID: client.GetID(), ClientSecret: client.Secret,
						Endpoint:    oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"},
						RedirectURL: client.RedirectURIs[0], Scopes: []string{"hydra", "offline", "openid"},
					}

					require.NoError(t, fs.(clientCreator).CreateClient(&client))
					apiClient := swagger.NewOAuth2ApiWithBasePath(api.URL)

					var callbackHandler *httprouter.Handle
					router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
						(*callbackHandler)(w, r, ps)
					})

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

						assertAccessToken func(*testing.T, string)
					}{
						{
							// First we need to create a persistent session in order to check if the other things work
							// as expected
							d:       "Creates a persisting session for the next test cases",
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=login+consent&max_age=1",
							cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetLoginRequest(r.URL.Query().Get("login_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.False(t, rr.Skip)
									assert.Empty(t, rr.Subject)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientId)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoUri)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, r.URL.Query().Get("login_challenge"), rr.Challenge)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state")+"&prompt=login+consent&max_age=1", rr.RequestUrl)

									v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{
										Subject: "user-a", Remember: true, RememberFor: 0,
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetConsentRequest(r.URL.Query().Get("consent_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)
									assert.EqualValues(t, client.GetID(), rr.Client.ClientId)
									assert.EqualValues(t, client.GrantTypes, rr.Client.GrantTypes)
									assert.EqualValues(t, client.LogoURI, rr.Client.LogoUri)
									assert.EqualValues(t, client.RedirectURIs, rr.Client.RedirectUris)
									assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
									assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), rr.Challenge)
									assert.EqualValues(t, oauthConfig.AuthCodeURL("some-hardcoded-state")+"&prompt=login+consent&max_age=1", rr.RequestUrl)

									v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
										GrantScope: []string{"hydra", "offline", "openid"}, Remember: true, RememberFor: 0,
										Session: swagger.ConsentRequestSession{
											AccessToken: map[string]interface{}{"foo": "bar"},
											IdToken:     map[string]interface{}{"bar": "baz"},
										},
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
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
								assert.EqualValues(t, "bar", data["foo"])
							},
						},
						{
							d: "checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=none works",
							setup: func() {
								// In order to check if authenticatedAt/requestedAt works, we'll sleep first in order to ensure that authenticatedAt is in the past
								// if handled correctly.
								time.Sleep(time.Second * 2)
							},
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&prompt=none&max_age=60",
							cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetLoginRequest(r.URL.Query().Get("login_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)

									v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{Subject: "user-a"})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetConsentRequest(r.URL.Query().Get("consent_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)

									v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
										GrantScope: []string{"hydra", "offline"},
										Session: swagger.ConsentRequestSession{
											AccessToken: map[string]interface{}{"foo": "bar"},
											IdToken:     map[string]interface{}{"bar": "baz"},
										},
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
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
							expectIDToken:         false,
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
								assert.EqualValues(t, []interface{}{"hydra", "offline"}, data["scp"])
								assert.EqualValues(t, "bar", data["foo"])
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
									rr, res, err := apiClient.GetLoginRequest(r.URL.Query().Get("login_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "", rr.Subject)

									v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{Subject: "user-a"})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetConsentRequest(r.URL.Query().Get("consent_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)

									v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
										GrantScope: []string{"hydra", "openid"},
										Session: swagger.ConsentRequestSession{
											AccessToken: map[string]interface{}{"foo": "bar"},
											IdToken:     map[string]interface{}{"bar": "baz"},
										},
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
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
									rr, res, err := apiClient.GetLoginRequest(r.URL.Query().Get("login_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.False(t, rr.Skip)
									assert.EqualValues(t, "", rr.Subject)
									assert.Empty(t, rr.Client.ClientSecret)

									v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{Subject: "user-a"})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									rr, res, err := apiClient.GetConsentRequest(r.URL.Query().Get("consent_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									assert.True(t, rr.Skip)
									assert.EqualValues(t, "user-a", rr.Subject)
									assert.Empty(t, rr.Client.ClientSecret)

									v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
										GrantScope: []string{"hydra", "openid"},
										Session: swagger.ConsentRequestSession{
											AccessToken: map[string]interface{}{"foo": "bar"},
											IdToken:     map[string]interface{}{"bar": "baz"},
										},
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
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
							authURL: oauthConfig.AuthCodeURL("some-hardcoded-state") + "&max_age=3",
							//cj:      persistentCJ,
							lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, res, err := apiClient.GetLoginRequest(r.URL.Query().Get("login_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)

									time.Sleep(time.Second * 5)

									v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{Subject: "user-a"})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
									require.NotEmpty(t, v.RedirectTo)
									http.Redirect(w, r, v.RedirectTo, http.StatusFound)
								}
							},
							cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
								return func(w http.ResponseWriter, r *http.Request) {
									_, res, err := apiClient.GetConsentRequest(r.URL.Query().Get("consent_challenge"))
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)

									time.Sleep(time.Second * 5)

									v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
										GrantScope: []string{"hydra", "openid"},
										Session: swagger.ConsentRequestSession{
											AccessToken: map[string]interface{}{"foo": "bar"},
											IdToken:     map[string]interface{}{"bar": "baz"},
										},
									})
									require.NoError(t, err)
									require.EqualValues(t, http.StatusOK, res.StatusCode)
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
						})
					}
				})
			}
		})
	}
}

// TestAuthCodeWithMockStrategy runs the authorize_code flow against various ConsentStrategy scenarios.
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
	for _, strat := range []struct {
		d string
		s foauth2.CoreStrategy
	}{
		{
			d: "opaque",
			s: oauth2OpqaueStrategy,
		},
		{
			d: "jwt",
			s: oauth2JWTStrategy,
		},
	} {
		t.Run("strategy="+strat.d, func(t *testing.T) {
			consentStrategy := &consentMock{}
			router := httprouter.New()
			ts := httptest.NewServer(router)
			store := NewFositeMemoryStore(hc.NewMemoryManager(hasher), time.Second)

			l := logrus.New()
			l.Level = logrus.DebugLevel

			jm := &jwk.MemoryManager{Keys: map[string]*jose.JSONWebKeySet{}}
			keys, err := (&jwk.RS256Generator{}).Generate("", "sig")
			require.NoError(t, err)
			require.NoError(t, jm.AddKeySet(OpenIDConnectKeyName, keys))
			jwtStrategy, err := jwk.NewRS256JWTStrategy(jm, OpenIDConnectKeyName)

			handler := &Handler{
				OAuth2: compose.Compose(
					fc,
					store,
					strat.s,
					nil,
					compose.OAuth2AuthorizeExplicitFactory,
					compose.OAuth2AuthorizeImplicitFactory,
					compose.OAuth2ClientCredentialsGrantFactory,
					compose.OAuth2RefreshTokenGrantFactory,
					compose.OpenIDConnectExplicitFactory,
					compose.OpenIDConnectHybridFactory,
					compose.OpenIDConnectImplicitFactory,
					compose.OAuth2TokenRevocationFactory,
					compose.OAuth2TokenIntrospectionFactory,
				),
				Consent:           consentStrategy,
				CookieStore:       sessions.NewCookieStore([]byte("foo-secret")),
				ForcedHTTP:        true,
				L:                 l,
				H:                 herodot.NewJSONWriter(l),
				ScopeStrategy:     fosite.HierarchicScopeStrategy,
				IDTokenLifespan:   time.Minute,
				IssuerURL:         ts.URL,
				OpenIDJWTStrategy: jwtStrategy,
			}
			handler.SetRoutes(router,router)

			var callbackHandler *httprouter.Handle
			router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				(*callbackHandler)(w, r, ps)
			})
			m := sync.Mutex{}

			store.CreateClient(&hc.Client{
				ClientID:      "app-client",
				Secret:        "secret",
				RedirectURIs:  []string{ts.URL + "/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scope:         "hydra.* offline openid",
			})

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
				authTime                  time.Time
				requestTime               time.Time
				assertAccessToken         func(*testing.T, string)
			}{
				{
					d:                         "should pass request if strategy passes",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
					shouldPassConsentStrategy: true,
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
					m.Lock()
					defer m.Unlock()
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
					require.NoError(t, err)
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

					time.Sleep(time.Second * 2)

					res, err := testRefresh(t, token, ts.URL)
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

						assert.NotEqual(t, refreshedPayload["exp"], origPayload["exp"])
						assert.NotEqual(t, refreshedPayload["iat"], origPayload["iat"])
						assert.NotEqual(t, refreshedPayload["nbf"], origPayload["nbf"])
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
						res, err := testRefresh(t, token, ts.URL)
						require.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, res.StatusCode)
					})

					t.Run("refreshing new refresh token should work", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL)
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

func testRefresh(t *testing.T, token *oauth2.Token, u string) (*http.Response, error) {
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

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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent_test

import (
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

	"github.com/tidwall/gjson"

	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/ioutilx"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	"github.com/ory/x/sqlxx"

	"github.com/ory/viper"

	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"

	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/x"
)

func mustParseURL(t *testing.T, u string) *url.URL {
	uu, err := url.Parse(u)
	require.NoError(t, err)
	return uu
}

func mockProvider(h *func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(*h)(w, r)
	}))
}

func noopHandler(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func newCookieJar(t *testing.T) *cookiejar.Jar {
	c, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	return c
}

func acceptRequest(apiClient *hydra.OryHydra, consent *models.AcceptConsentRequest) func(t *testing.T) func(http.ResponseWriter, *http.Request) {
	if consent == nil {
		consent = &models.AcceptConsentRequest{
			GrantScope: []string{"scope-a"},
		}
	}
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
				WithBody(consent))
			require.NoError(t, err)
			v := vr.Payload
			http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
		}
	}
}

func newAuthCookieJar(t *testing.T, reg driver.Registry, u, sessionID string) http.CookieJar {
	cj, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	secrets := viper.GetStringSlice(configuration.ViperKeyGetCookieSecrets)
	bs := make([][]byte, len(secrets))
	for k, s := range secrets {
		bs[k] = []byte(s)
	}

	hr := &http.Request{Header: map[string][]string{}, URL: urlx.ParseOrPanic(u), RequestURI: u}
	cookie, _ := reg.CookieStore().Get(hr, CookieName(reg.Config().ServesHTTPS(), CookieAuthenticationName))

	cookie.Values[CookieAuthenticationSIDName] = sessionID
	cookie.Options.HttpOnly = true

	rw := httptest.NewRecorder()
	require.NoError(t, cookie.Save(hr, rw))

	cj.SetCookies(urlx.ParseOrPanic(u), rw.Result().Cookies())
	return cj
}

func newValidAuthCookieJar(t *testing.T, reg driver.Registry, u, sessionID, subject string) http.CookieJar {
	cj := newAuthCookieJar(t, reg, u, sessionID)
	require.NoError(t, reg.ConsentManager().CreateLoginSession(context.TODO(), &LoginSession{
		ID:              sessionID,
		Subject:         subject,
		AuthenticatedAt: sqlxx.NullTime(time.Now()),
		Remember:        true,
	}))
	return cj
}

func acceptLogoutChallenge(api *hydra.OryHydra, key string) func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c := r.URL.Query().Get("logout_challenge")
			assert.NotEmpty(t, c)
			logout, err := api.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().WithLogoutChallenge(c))
			require.NoError(t, err)
			assert.EqualValues(t, "logout-subject-"+key, logout.Payload.Subject)
			assert.EqualValues(t, "logout-session-"+key, logout.Payload.Sid)

			redir, err := api.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().WithLogoutChallenge(c))
			require.NoError(t, err)

			assert.Contains(t, *redir.Payload.RedirectTo, "?logout_verifier")
			http.Redirect(w, r, *redir.Payload.RedirectTo, http.StatusFound)
		}
	}
}

func genIDToken(t *testing.T, reg driver.Registry, c jwtgo.Claims) string {
	r, _, err := reg.OpenIDJWTStrategy().Generate(context.TODO(), c, jwt.NewHeaders())
	require.NoError(t, err)
	return r
}

func logoutHandler(strategy Strategy, writer herodot.Writer, w http.ResponseWriter, r *http.Request) {
	res, err := strategy.HandleOpenIDConnectLogout(w, r)
	if errors.Is(err, ErrAbortOAuth2Request) {
		// Do nothing
		return
	} else if err != nil {
		writer.WriteError(w, r, err)
		return
	}

	http.Redirect(w, r,
		urlx.CopyWithQuery(
			urlx.ParseOrPanic(res.RedirectTo),
			url.Values{},
		).String(),
		http.StatusFound,
	)
}

func runLogout(t *testing.T, method string) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
	// jwts := reg.OpenIDJWTStrategy()

	var lph func(w http.ResponseWriter, r *http.Request)
	logoutProviderServer := mockProvider(&lph)

	writer := reg.Writer()
	handler := reg.ConsentHandler()
	router := x.NewRouterAdmin()
	handler.SetRoutes(router)
	n := negroni.Classic()
	n.UseHandler(router)
	logoutApi := httptest.NewServer(n)
	defer logoutApi.Close()

	defaultRedirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		state := r.Form.Get("state")
		_, _ = w.Write([]byte(fmt.Sprintf("redirected to default server%s%s", string(state), strings.TrimLeft(r.URL.Path, "/"))))
	}))
	defer defaultRedirServer.Close()

	strategy := reg.ConsentStrategy()
	logoutRouter := x.NewRouterPublic()
	logoutRouter.GET("/oauth2/sessions/logout", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logoutHandler(strategy, writer, w, r)
	})
	logoutRouter.POST("/oauth2/sessions/logout", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logoutHandler(strategy, writer, w, r)
	})

	logoutServer := httptest.NewServer(logoutRouter)
	defer logoutServer.Close()

	viper.Set(configuration.ViperKeyIssuerURL, logoutServer.URL)
	viper.Set(configuration.ViperKeyLogoutURL, logoutProviderServer.URL)
	viper.Set(configuration.ViperKeyLogoutRedirectURL, defaultRedirServer.URL)

	defaultClient := &client.Client{OutfacingID: uuid.New(), PostLogoutRedirectURIs: []string{defaultRedirServer.URL + "/custom"}}
	require.NoError(t, reg.ClientManager().CreateClient(context.TODO(), defaultClient))

	jar1 := newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-1", "logout-subject-1")
	jar3 := newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-3", "logout-subject-3")
	apiClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(logoutApi.URL).Host})

	for k, tc := range []struct {
		d                string
		params           url.Values
		subject          string
		sessionID        string
		lph              func(t *testing.T) func(w http.ResponseWriter, r *http.Request)
		expectBody       string
		expectRequestURI string
		backChannels     []func(t *testing.T) func(w http.ResponseWriter, r *http.Request)
		expectStatusCode int
		jar              http.CookieJar
	}{
		{
			d:                "should ignore / redirect non-rp initiated logout if no session exists",
			lph:              noopHandler,
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
		},
		{
			d:                "should fail if non-rp initiated logout is initiated with state (indicating rp-flow)",
			params:           url.Values{"state": {"foobar"}},
			lph:              noopHandler,
			expectBody:       "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Logout failed because query parameter state is set but id_token_hint is missing.\"}",
			expectStatusCode: http.StatusBadRequest,
		},
		{
			d:                "should fail if non-rp initiated logout is initiated with post_logout_redirect_uri (indicating rp-flow)",
			params:           url.Values{"post_logout_redirect_uri": {"foobar"}},
			lph:              noopHandler,
			expectBody:       "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Logout failed because query parameter post_logout_redirect_uri is set but id_token_hint is missing.\"}",
			expectStatusCode: http.StatusBadRequest,
		},
		{
			d:                "should ignore / redirect non-rp initiated logout if a session cookie exists but the session itself is no longer active",
			lph:              noopHandler,
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
			jar:              newAuthCookieJar(t, reg, logoutServer.URL, "i-do-not-exist"),
		},
		{
			d:                "should redirect to logout provider if session exists and it's not rp-flow",
			lph:              acceptLogoutChallenge(apiClient, "1"),
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
			jar:              jar1,
			subject:          "logout-subject-1",
		},
		{
			d:                "should redirect to logout provider because the session has been removed previously",
			lph:              noopHandler,
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
			jar:              jar1,
		},
		{
			d:                "should execute backchannel logout if issued without rp-involvement",
			lph:              acceptLogoutChallenge(apiClient, "2"),
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
			backChannels: []func(t *testing.T) func(w http.ResponseWriter, r *http.Request){
				func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
					return func(w http.ResponseWriter, r *http.Request) {
						require.NoError(t, r.ParseForm())
						lt := r.PostFormValue("logout_token")
						assert.NotEmpty(t, lt)
						token, err := reg.OpenIDJWTStrategy().Decode(context.TODO(), lt)
						require.NoError(t, err)

						claims := token.Claims.(jwtgo.MapClaims)
						assert.EqualValues(t, "logout-session-2", claims["sid"])
						assert.Empty(t, claims["sub"]) // The sub claim should be empty because it doesn't work with forced obfuscation and thus we can't easily recover it.
						assert.Empty(t, claims["nonce"])

						w.WriteHeader(http.StatusOK)
					}
				},
			},
			jar:       newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-2", "logout-subject-2"),
			subject:   "logout-subject-2",
			sessionID: "logout-session-2",
		},
		{
			d:                "should error when rp-flow without valid id token",
			lph:              acceptLogoutChallenge(apiClient, "3"),
			params:           url.Values{"state": {"1234"}, "post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"}, "id_token_hint": {"i am not valid"}},
			expectBody:       "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. token contains an invalid number of segments\"}",
			expectStatusCode: http.StatusBadRequest,
		},
		{
			d:   "should fail rp-inititated flow because id token hint is missing issuer",
			lph: acceptLogoutChallenge(apiClient, "temp1"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"sub": "logout-subject-temp1",
					"sid": "logout-session-temp1",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusBadRequest,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp1", "logout-subject-temp1"),
			expectBody:       fmt.Sprintf("{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Logout failed because issuer claim value '' from query parameter id_token_hint does not match with issuer value from configuration '%s'.\"}", conf.IssuerURL().String()),
		},
		{
			d:   "should fail rp-inititated flow because id token hint is using wrong issuer",
			lph: acceptLogoutChallenge(apiClient, "temp2"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"iss": "some-issuer",
					"sub": "logout-subject-temp2",
					"sid": "logout-session-temp2",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusBadRequest,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp2", "logout-subject-temp2"),
			expectBody:       fmt.Sprintf("{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Logout failed because issuer claim value 'some-issuer' from query parameter id_token_hint does not match with issuer value from configuration '%s'.\"}", conf.IssuerURL().String()),
		},
		{
			d:   "should fail rp-inititated flow because iat is in the future",
			lph: acceptLogoutChallenge(apiClient, "temp3"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-temp3",
					"sid": "logout-session-temp3",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(+time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusBadRequest,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp3", "logout-subject-temp3"),
			expectBody:       "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Token used before issued\"}",
		},
		{
			d:   "should fail because post-logout url is not registered",
			lph: acceptLogoutChallenge(apiClient, "temp4"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {"https://this-is-not-a-valid-redirect-url/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-temp4",
					"sid": "logout-session-temp4",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusBadRequest,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp4", "logout-subject-temp4"),
			expectBody:       "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Logout failed because query parameter post_logout_redirect_uri is not a whitelisted as a post_logout_redirect_uri for the client.\"}",
		},
		{
			d:   "should pass rp-inititated even when expiry is in the past",
			lph: acceptLogoutChallenge(apiClient, "temp5"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-temp5",
					"sid": "logout-session-temp5",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusOK,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp5", "logout-subject-temp5"),
			expectBody:       "redirected to default server1234custom",
			expectRequestURI: "/custom?state=1234",
		},
		{
			d:   "should pass rp-inititated flow",
			lph: acceptLogoutChallenge(apiClient, "3"),
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": []string{defaultClient.OutfacingID}, // make sure this works with string slices too
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-3",
					"sid": "logout-session-3",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusOK,
			jar:              jar3,
			subject:          "logout-subject-2",
			expectBody:       "redirected to default server1234custom",
		},
		{
			d:                "should redirect to logout provider because the session has been removed previously",
			lph:              noopHandler,
			expectBody:       "redirected to default server",
			expectStatusCode: http.StatusOK,
			jar:              jar3,
		},
		{
			d: "should pass rp-inititated flow without any action because SID is unknown",
			params: url.Values{
				"state":                    {"1234"},
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": []string{defaultClient.OutfacingID}, // make sure this works with string slices too
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-3",
					"sid": "i-do-not-exist",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			lph:              noopHandler,
			expectStatusCode: http.StatusOK,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp6", "logout-subject-temp6"),
			expectBody:       "redirected to default server1234custom",
		},
		{
			d:   "should not append a state param if no state was passed to logout server",
			lph: acceptLogoutChallenge(apiClient, "temp7"),
			params: url.Values{
				"post_logout_redirect_uri": {defaultRedirServer.URL + "/custom"},
				"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
					"aud": defaultClient.OutfacingID,
					"iss": conf.IssuerURL().String(),
					"sub": "logout-subject-temp7",
					"sid": "logout-session-temp7",
					"exp": time.Now().Add(-time.Hour).Unix(),
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				})},
			},
			expectStatusCode: http.StatusOK,
			jar:              newValidAuthCookieJar(t, reg, logoutServer.URL, "logout-session-temp7", "logout-subject-temp7"),
			subject:          "logout-subject-7",
			expectBody:       "redirected to default servercustom",
			expectRequestURI: "/custom",
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			if tc.lph != nil {
				lph = tc.lph(t)
			} else {
				lph = noopHandler(t)
			}

			var bcWg sync.WaitGroup
			servers := make([]*httptest.Server, len(tc.backChannels))
			for k, bc := range tc.backChannels {
				bcWg.Add(1)
				n := negroni.Classic()
				n.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer bcWg.Done()
					bc(t)(w, r)
				}))
				servers[k] = httptest.NewServer(n)
				c, hc := MockConsentRequest(uuid.New(), true, 100, false, false, true)
				c.LoginSessionID = sqlxx.NullString(tc.sessionID)
				c.Client.BackChannelLogoutURI = servers[k].URL
				c.Subject = tc.subject
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c.Client))
				require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), &LoginRequest{ID: c.LoginChallenge.String(), Client: c.Client, Verifier: c.ID}))
				require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), c))
				_, err := reg.ConsentManager().HandleConsentRequest(context.Background(), c.ID, hc)
				require.NoError(t, err)
			}

			cl := &http.Client{
				Jar: tc.jar,
			}

			var err error
			var resp *http.Response

			if method == http.MethodGet {
				resp, err = cl.Get(
					logoutServer.URL + "/oauth2/sessions/logout?" + tc.params.Encode(),
				)
			} else if method == http.MethodPost {
				resp, err = cl.PostForm(
					logoutServer.URL+"/oauth2/sessions/logout",
					tc.params,
				)
			}
			require.NoError(t, err)
			out, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			bcWg.Wait()

			for _, s := range servers {
				s.Close()
			}

			assert.EqualValues(t, tc.expectStatusCode, resp.StatusCode, "%s\n%s", resp.Request.URL.String(), out)
			assert.EqualValues(t, tc.expectBody, strings.Trim(string(out), "\n"), "%s\n%s", resp.Request.URL.String(), out)

			if tc.expectRequestURI != "" {
				assert.EqualValues(t, tc.expectRequestURI, resp.Request.URL.RequestURI(), "%s\n%s", resp.Request.URL.String(), out)
			}
		})
	}
}

func TestStrategyLogout(t *testing.T) {
	runLogout(t, http.MethodGet)
}

func TestStrategyLogoutPost(t *testing.T) {
	runLogout(t, http.MethodPost)
}

func TestStrategyLoginConsentNext(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	viper.Set(configuration.ViperKeyAccessTokenStrategy, "opaque")
	viper.Set(configuration.ViperKeyConsentRequestMaxAge, time.Hour)
	publicTS, adminTS := testhelpers.NewOAuth2Server(t, reg)
	t.Logf("Public URL: %s", publicTS.URL)
	t.Logf("Admin URL: %s", adminTS.URL)

	adminClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(adminTS.URL).Host})

	//nonexistentCJ := newAuthCookieJar(t, reg, publicTS.URL, "i-do-not-exist")

	acceptLoginHandler := func(t *testing.T, subject string, payload *models.AcceptLoginRequest) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// login payloads are checked in ../oauth/auth_code tests
			if payload == nil {
				payload = &models.AcceptLoginRequest{Subject: &subject}
			}
			payload.Subject = &subject

			v, err := adminClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(payload))
			require.NoError(t, err)
			require.NotEmpty(t, *v.Payload.RedirectTo)
			http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
		}
	}

	acceptConsentHandler := func(t *testing.T, payload *models.AcceptConsentRequest) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// consent payloads are checked in ../oauth/auth_code tests
			if payload == nil {
				payload = &models.AcceptConsentRequest{}
			}

			v, err := adminClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
				WithBody(payload))
			require.NoError(t, err)
			require.NotEmpty(t, *v.Payload.RedirectTo)
			http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
		}
	}

	newClient := func(t *testing.T, redir string) *client.Client {
		secret := uuid.New()
		c := &client.Client{
			OutfacingID:  uuid.New(),
			Secret:       secret,
			RedirectURIs: []string{redir},
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))
		c.Secret = secret
		return c
	}

	makeRequest := func(t *testing.T, hc *http.Client, oc *client.Client, values url.Values) (gjson.Result, *http.Response) {
		if hc == nil {
			hc = testhelpers.NewEmptyJarClient(t)
		}

		values.Add("response_type", "code")
		values.Add("state", uuid.New())
		values.Add("client_id", oc.OutfacingID)
		res, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(), values).String())
		require.NoError(t, err)
		defer res.Body.Close()

		return gjson.ParseBytes(ioutilx.MustReadAll(res.Body)), res
	}

	t.Run("case=should fail because a login verifier was given that doesn't exist in the store", func(t *testing.T) {
		testhelpers.NewUI(t, testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

		_, res := makeRequest(t, nil, c, url.Values{"login_verifier": {"does-not-exist"}})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "The login verifier has already been used, has not been granted, or is invalid.")
	})

	t.Run("case=should fail because a non-existing consent verifier was given", func(t *testing.T) {
		testhelpers.NewUI(t, testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

		_, res := makeRequest(t, nil, c, url.Values{"consent_verifier": {"does-not-exist"}})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "The consent verifier has already been used, has not been granted, or is invalid.")
	})

	t.Run("case=should fail because the request was redirected but the login endpoint doesn't do anything (like redirecting back)", func(t *testing.T) {
		testhelpers.NewUI(t, testhelpers.HTTPServerNotImplementedHandler, testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNoExpectedCallHandler(t)))

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("login_challenge"), "%s", res.Request.URL)
	})

	t.Run("case=should fail because the request was redirected but consent endpoint doesn't do anything (like redirecting back)", func(t *testing.T) {
		testhelpers.NewUI(t, acceptLoginHandler(t, "aeneas-rekkas", nil), testhelpers.HTTPServerNotImplementedHandler)
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNoExpectedCallHandler(t)))

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("consent_challenge"), "%s", res.Request.URL)
	})

	t.Run("case=should fail because the request was redirected but the login endpoint rejected the request", func(t *testing.T) {
		testhelpers.NewUI(t, func(w http.ResponseWriter, r *http.Request) {
			vr, err := adminClient.Admin.RejectLoginRequest(admin.NewRejectLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(&models.RejectRequest{
					Error:            fosite.ErrInteractionRequired.ErrorField,
					ErrorDescription: "expect-reject-login",
					StatusCode:       int64(fosite.ErrInteractionRequired.CodeField),
				}))
			require.NoError(t, err)
			assert.NotEmpty(t, *vr.Payload.RedirectTo)
			http.Redirect(w, r, *vr.Payload.RedirectTo, http.StatusFound)
		}, testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "expect-reject-login")
	})

	t.Run("case=should fail because no cookie jar invalid csrf", func(t *testing.T) {
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewUI(t, acceptLoginHandler(t, "aeneas-rekkas", nil),
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		hc := new(http.Client)
		_, res := makeRequest(t, hc, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "No CSRF value available in the session cookie.")
	})

	t.Run("case=should fail because consent endpoints idles after login was granted - but consent endpoint should be called because cookie jar exists", func(t *testing.T) {
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewUI(t, acceptLoginHandler(t, "aeneas-rekkas", nil),
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		hc := new(http.Client)
		_, res := makeRequest(t, hc, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "No CSRF value available in the session cookie.")
	})

	t.Run("case=should because consent endpoints denies the request after login was granted", func(t *testing.T) {
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewUI(t,
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			func(w http.ResponseWriter, r *http.Request) {
				vr, err := adminClient.Admin.RejectConsentRequest(
					admin.NewRejectConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
						WithBody(&models.RejectRequest{
							Error:            fosite.ErrInteractionRequired.ErrorField,
							ErrorDescription: "expect-reject-consent",
							StatusCode:       int64(fosite.ErrInteractionRequired.CodeField)}))
				require.NoError(t, err)
				require.NotEmpty(t, *vr.Payload.RedirectTo)
				http.Redirect(w, r, *vr.Payload.RedirectTo, http.StatusFound)
			})

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), "expect-reject-consent")
	})

	t.Run("case=should pass if both login and consent are granted", func(t *testing.T) {
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewUI(t,
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("code"))
	})

	t.Run("case=should pass and set acr values properly", func(t *testing.T) {
		c := newClient(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewUI(t,
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		_, res := makeRequest(t, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("code"))
	})

	t.Run("case=should pass because login and consent have been granted, this time we remember the decision", func(t *testing.T) {
		c := &client.Client{OutfacingID: uuid.New(), TokenEndpointAuthMethod: "none",
			RedirectURIs: []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
				"custom://redirection-scheme/path", "custom://localhost/path"}}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		subject := "aeneas-rekkas"
		testhelpers.NewUI(t,
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true, RememberFor: 0}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true, RememberFor: 0}))

		hc := testhelpers.NewEmptyJarClient(t)

		t.Run("followup=perform first flow", func(t *testing.T) {
			_, res := makeRequest(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})
			assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
			assert.NotEmpty(t, res.Request.URL.Query().Get("code"))
		})

		// By not waiting here we ensure that there are no race conditions when it comes to authenticated_at and
		// requested_at time comparisons:
		//
		//	time.Sleep(time.Second)

		t.Run("followup=should pass because prompt=none, client is public, redirection scheme is HTTP and host is localhost", func(t *testing.T) {
			testhelpers.NewUI(t, func(w http.ResponseWriter, r *http.Request) {
				res, err := adminClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
				require.NoError(t, err)
				assert.True(t, *res.Payload.Skip)
				assert.Equal(t, subject, *res.Payload.Subject)
				acceptLoginHandler(t, subject, nil)(w, r)
			}, func(w http.ResponseWriter, r *http.Request) {
				res, err := adminClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
				require.NoError(t, err)
				assert.True(t, res.Payload.Skip)
				assert.Equal(t, subject, res.Payload.Subject)
				assert.NotEmpty(t, res.Payload.LoginSessionID)
				assert.NotEmpty(t, res.Payload.LoginChallenge)
				acceptConsentHandler(t, nil)(w, r)
			})

			_, res := makeRequest(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}, "prompt": {"none"}})
			assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
			assert.NotEmpty(t, res.Request.URL.Query().Get("code"), "%s", res.Request.URL.String())
		})

		return

		t.Run("followup=check remember values", func(t *testing.T) {
			for _, redir := range c.RedirectURIs[1:] {
				t.Run("redir=should pass because prompt=none, client is public, and redirection is "+redir, func(t *testing.T) {
					testhelpers.NewUI(t,
						acceptLoginHandler(t, subject, nil),
						acceptConsentHandler(t, nil))

					_, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(), url.Values{
						"response_type": {"code"},
						"state":         {uuid.New()},
						"redirect_uri":  {redir},
						"client_id":     {c.OutfacingID},
						"prompt":        {"none"},
					}).String())
					require.Error(t, err)
					assert.Contains(t, err.Error(), redir)
					assert.Contains(t, err.Error(), "error=login_required")
				})
			}
		})
	})

	/*
		{
			d:   "This should fail at login screen because subject from accept does not match subject from session",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					lr := res.Payload

					assert.True(t, *lr.Skip)
					assert.Equal(t, "user", *lr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("fooser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.Error(t, err)
					require.Empty(t, vr)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request},
			expectErr:             []bool{true},
		},

		{
			d:   "This should pass and confirm previous authentication and consent because it is a authorization_code",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", Secret: "should-not-be-included"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload

					assert.True(t, *rr.Skip)
					assert.Equal(t, "user", *rr.Subject)
					assert.Empty(t, rr.Client.ClientSecret)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.True(t, rr.Skip)
					assert.Equal(t, "client-id", rr.Client.ClientID)
					assert.Equal(t, "user", rr.Subject)
					assert.Empty(t, rr.Client.ClientSecret)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:      "This should pass and require re-authentication although session is set (because prompt=login)",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			prompt: "login+consent",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload

					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    true,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					body := `{"grant_scope": ["scope-a"], "remember": true}`
					require.NoError(t, err)
					req, err := http.NewRequest("PUT", api.URL+"/oauth2/auth/requests/consent/accept?challenge="+r.URL.Query().Get("consent_challenge"), strings.NewReader(body))
					req.Header.Add("Content-Type", "application/json")
					require.NoError(t, err)

					hres, err := http.DefaultClient.Do(req)
					require.NoError(t, err)
					defer hres.Body.Close()

					var v models.CompletedRequest
					require.NoError(t, json.NewDecoder(hres.Body).Decode(&v))
					require.EqualValues(t, http.StatusOK, hres.StatusCode)
					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session:        NewConsentRequestSessionData(),
			},
		},
		{
			d:      "This should pass and require re-authentication although session is set (because max_age=1)",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			maxAge: "1",
			setup: func() {
				time.Sleep(time.Second * 2)
			},
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    true,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.True(t, rr.Skip)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should fail because max_age=1 but prompt=none",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			setup: func() {
				time.Sleep(time.Second * 2)
			},
			maxAge:                "1",
			prompt:                "none",
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:                     "This should fail because prompt is none but no auth session exists",
			prompt:                "none",
			req:                   fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   newCookieJar(t),
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:      "This should fail because prompt is none and consent is missing a permission which requires re-authorization of the app",
			prompt: "none",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a", "this-scope-has-not-been-granted-before"}}},
			jar:    persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.True(t, *rr.Skip)
					assert.Equal(t, "user", *rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrConsentRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:      "This pass and properly require authentication as well as authorization because prompt is set to login and consent - although previous session exists",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			prompt: "login+consent",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:                     "This should fail because id_token_hint does not match authentication session and prompt is none",
			req:                   fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   persistentCJ,
			prompt:                "none",
			idTokenHint:           fooUserIDToken,
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:           "This should pass and require authentication because id_token_hint does not match subject from session",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         persistentCJ,
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("not-foouser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrLoginRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:           "This should pass and require authentication because id_token_hint does not match subject from session",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         persistentCJ,
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("foouser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
						WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
						WithBody(&models.AcceptConsentRequest{
							GrantScope:  []string{"scope-a"},
							Remember:    false,
							RememberFor: 0,
							Session: &models.ConsentRequestSession{
								AccessToken: map[string]interface{}{"foo": "bar"},
								IDToken:     map[string]interface{}{"bar": "baz"},
							},
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "foouser", SubjectIdentifier: "foouser"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should pass as regularly even though id_token_hint is expired",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: newCookieJar(t),
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "user",
				ExpiresAt: time.Now().Add(-time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph:                   passAuthentication(apiClient, false),
			cph:                   passAuthorization(apiClient, false),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},

		// Pairwise auth
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set by hydra",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd", // this is sha256("fooauth-user76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		}, // these tests depend on one another
		{
			d:   "This should pass as regularly and create a new session with pairwise subject and also with the ID token set",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd",
				ExpiresAt: time.Now().Add(time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		},
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set login request",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ4,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:                pointerx.String("auth-user"),
							ForceSubjectIdentifier: "forced-auth-user",
							Remember:               true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "forced-auth-user",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		}, // these tests depend on one another
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set on login request and also with the ID token set",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "forced-auth-user",
				ExpiresAt: time.Now().Add(time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:                pointerx.String("auth-user"),
							ForceSubjectIdentifier: "forced-auth-user",
							Remember:               false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "forced-auth-user",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		},

		// checks revoking sessions
		{
			d:   "This should pass as regularly and create a new session and forward data",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ2,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:      "This should pass and also revoke the session cookie",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ2,
			prompt: "login",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("not-auth-user"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		}, // these two tests depend on one another
		{
			d:   "This should require re-authentication because the session was revoked in the previous test",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ2,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.Empty(t, "", rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("foouser"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:   "This should require re-authentication because the session does not exist in the store",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: nonexistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.Empty(t, "", rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("foouser"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:           "This should fail because the user from the ID token does not match the user from the accept login request",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         newCookieJar(t),
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.EqualValues(t, "", *rr.Subject)
					assert.EqualValues(t, "foouser", rr.OidcContext.IDTokenHintClaims.(map[string]interface{})["sub"])

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("not-foouser"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrLoginRequired},
			expectErr:             []bool{true, true},
		},
	*/

	/*
		calls := -1
		aph = func(w http.ResponseWriter, r *http.Request) {
			calls++
			require.True(t, len(tc.expectErrType) >= calls+1, "%d (expect) < %d (got)", len(tc.expectErrType), calls+1)
			require.True(t, len(tc.expectErr) >= calls+1, "%d (expect) < %d (got)", len(tc.expectErr), calls+1)
			require.NoError(t, r.ParseForm())
			tc.req.Form = r.Form

			c, err := strategy.HandleOAuth2AuthorizationRequest(w, r, &tc.req)
			t.Logf("DefaultStrategy returned at call %d:\n\tresult: %+v\n\terr: %s", calls, c, err)

			if tc.expectErr[calls] {
				assert.Error(t, err)
				if tc.expectErrType[calls] != nil {
					assert.EqualError(t, err, tc.expectErrType[calls].Error(), "%+v", err)
				}
			} else {
				require.NoError(t, err)
				if tc.expectSession != nil {
					require.NotNil(t, c)
					assert.EqualValues(t, tc.expectSession.GrantedScope, c.GrantedScope)
					assert.EqualValues(t, tc.expectSession.Remember, c.Remember)
					assert.EqualValues(t, tc.expectSession.Session, c.Session)
					assert.EqualValues(t, tc.expectSession.RememberFor, c.RememberFor)
					assert.EqualValues(t, tc.expectSession.ConsentRequest.Subject, c.ConsentRequest.Subject)
					assert.EqualValues(t, tc.expectSession.ConsentRequest.SubjectIdentifier, c.ConsentRequest.SubjectIdentifier)
				}
			}

			if errors.Is(err, ErrAbortOAuth2Request) {
				// nothing to do, indicates redirect
			} else if err != nil {
				writer.WriteError(w, r, err)
			} else {
				writer.Write(w, r, c)
			}
		}

		cl := &http.Client{
			Jar: tc.jar,
		}
		resp, err := cl.Get(
			ap.URL + "?" +
				"login_verifier=" + tc.lv + "&" +
				"consent_verifier=" + tc.cv + "&" +
				"prompt=" + tc.prompt + "&" +
				"max_age=" + tc.maxAge + "&" +
				"id_token_hint=" + tc.idTokenHint + "&" + tc.other,
		)
		require.NoError(t, err)
		out, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		resp.Body.Close()
		assert.EqualValues(t, tc.expectFinalStatusCode, resp.StatusCode, "%s\n%s", resp.Request.URL.String(), out)
		// assert.Empty(t, resp.Request.URL.String())
	*/
}

func TestStrategyLoginConsent(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)

	var lph, cph, aph func(w http.ResponseWriter, r *http.Request)
	lp := mockProvider(&lph)
	cp := mockProvider(&cph)
	ap := mockProvider(&aph)

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)

	fooUserIDToken := genIDToken(t, reg, jwt.IDTokenClaims{Subject: "foouser", ExpiresAt: time.Now().Add(time.Hour), IssuedAt: time.Now()}.ToMapClaims())

	writer := reg.Writer()
	handler := reg.ConsentHandler()
	router := x.NewRouterAdmin()
	handler.SetRoutes(router)
	api := httptest.NewServer(router)
	defer api.Close()

	strategy := reg.ConsentStrategy()

	viper.Set(configuration.ViperKeyLoginURL, lp.URL)
	viper.Set(configuration.ViperKeyConsentURL, cp.URL)
	viper.Set(configuration.ViperKeyIssuerURL, ap.URL)
	viper.Set(configuration.ViperKeyConsentRequestMaxAge, time.Hour)
	viper.Set(configuration.ViperKeyScopeStrategy, "exact")
	viper.Set(configuration.ViperKeySubjectTypesSupported, []string{"pairwise", "public"})
	viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")

	apiClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(api.URL).Host})

	persistentCJ := newCookieJar(t)
	persistentCJ2 := newCookieJar(t)
	persistentCJ3 := newCookieJar(t)
	persistentCJ4 := newCookieJar(t)
	nonexistentCJ := newAuthCookieJar(t, reg, ap.URL, "i-do-not-exist")

	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), &client.Client{OutfacingID: "client-id"}))

	for k, tc := range []struct {
		setup                 func()
		d                     string
		lv                    string
		cv                    string
		lph, cph              func(t *testing.T) func(w http.ResponseWriter, r *http.Request)
		req                   fosite.AuthorizeRequest
		expectSession         *HandledConsentRequest
		expectErr             []bool
		expectErrType         []error
		expectFinalStatusCode int
		prompt                string
		maxAge                string
		idTokenHint           string
		other                 string
		jar                   http.CookieJar
	}{
		{
			d:                     "This should fail because a login verifier was given that doesn't exist in the store",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}}},
			lv:                    "invalid",
			expectErrType:         []error{fosite.ErrAccessDenied},
			expectErr:             []bool{true},
			expectFinalStatusCode: http.StatusForbidden,
		},
		{
			d:                     "This should fail because a consent verifier was given but no login verifier",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}}},
			lv:                    "",
			cv:                    "invalid",
			expectErrType:         []error{fosite.ErrAccessDenied},
			expectErr:             []bool{true},
			expectFinalStatusCode: http.StatusForbidden,
		},
		{
			d: "This should fail because the request was redirected but the login endpoint doesn't do anything (like redirecting back)",
			req: fosite.AuthorizeRequest{
				Request: fosite.Request{
					Client:         &client.Client{OutfacingID: "client-id"},
					RequestedScope: []string{"scope-a"},
				},
			},
			other: "display=page&ui_locales=de+en",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					lr := res.Payload

					assert.NotEmpty(t, lr.Challenge)
					assert.EqualValues(t, r.URL.Query().Get("login_challenge"), *lr.Challenge)
					assert.EqualValues(t, "client-id", lr.Client.ClientID)
					assert.EqualValues(t, []string{"scope-a"}, lr.RequestedScope)
					assert.Contains(t, *lr.RequestURL, "/oauth2/auth?login_verifier=&consent_verifier=&")
					assert.EqualValues(t, false, *lr.Skip)
					assert.EqualValues(t, "", *lr.Subject)
					assert.EqualValues(t, &models.OpenIDConnectContext{AcrValues: nil, Display: "page", UILocales: []string{"de", "en"}}, lr.OidcContext, "%s", res.Payload)
					w.WriteHeader(http.StatusNoContent)
				}
			},
			expectFinalStatusCode: http.StatusNoContent,
			expectErrType:         []error{ErrAbortOAuth2Request},
			expectErr:             []bool{true},
		},
		{
			d:   "This should fail because the request was redirected but the login endpoint rejected the request",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.RejectLoginRequest(admin.NewRejectLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.RejectRequest{
							Error:            fosite.ErrInteractionRequired.ErrorField,
							ErrorDebug:       fosite.ErrInteractionRequired.DebugField,
							ErrorDescription: fosite.ErrInteractionRequired.DescriptionField,
							ErrorHint:        fosite.ErrInteractionRequired.HintField,
							StatusCode:       int64(fosite.ErrInteractionRequired.CodeField),
						}))
					require.NoError(t, err)
					lr := vr.Payload

					assert.NotEmpty(t, *lr.RedirectTo)
					http.Redirect(w, r, *lr.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrInteractionRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:   "This should fail because no cookie jar / invalid csrf",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			lph: passAuthentication(apiClient, false),
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					// this should never be called because csrf doesn't make it that far
					require.True(t, false)
				}
			},
			expectFinalStatusCode: http.StatusForbidden,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrRequestForbidden},
			expectErr:             []bool{true, true},
		},
		{
			d:     "This should fail because consent endpoints idles after login was granted - but consent endpoint should be called because cookie jar exists",
			jar:   newCookieJar(t),
			req:   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			lph:   passAuthentication(apiClient, false),
			other: "display=page&ui_locales=de+en&acr_values=1+2",
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					lr := rrr.Payload

					assert.NotEmpty(t, *lr.Challenge)
					assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), *lr.Challenge)
					assert.EqualValues(t, "client-id", lr.Client.ClientID)
					assert.EqualValues(t, []string{"scope-a"}, lr.RequestedScope)
					assert.Contains(t, lr.RequestURL, "/oauth2/auth?login_verifier=&consent_verifier=&")
					assert.EqualValues(t, false, lr.Skip)
					assert.EqualValues(t, "user", lr.Subject)
					assert.NotEmpty(t, lr.LoginChallenge)
					assert.NotEmpty(t, lr.LoginSessionID)
					assert.EqualValues(t, &models.OpenIDConnectContext{AcrValues: []string{"1", "2"}, Display: "page", UILocales: []string{"de", "en"}}, lr.OidcContext)
					w.WriteHeader(http.StatusNoContent)
				}
			},
			expectFinalStatusCode: http.StatusNoContent,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request},
			expectErr:             []bool{true, true},
		},
		{
			d:                     "This should fail because consent verifier was set but does not exist",
			jar:                   newCookieJar(t),
			cv:                    "invalid",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			expectFinalStatusCode: http.StatusForbidden,
			expectErrType:         []error{fosite.ErrAccessDenied},
			expectErr:             []bool{true},
		},
		{
			d:   "This should fail because consent endpoints denies the request after login was granted",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: newCookieJar(t),
			lph: passAuthentication(apiClient, false),
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.RejectConsentRequest(
						admin.NewRejectConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
							WithBody(
								&models.RejectRequest{
									Error:            fosite.ErrInteractionRequired.ErrorField,
									ErrorDebug:       fosite.ErrInteractionRequired.DebugField,
									ErrorDescription: fosite.ErrInteractionRequired.DescriptionField,
									ErrorHint:        fosite.ErrInteractionRequired.HintField,
									StatusCode:       int64(fosite.ErrInteractionRequired.CodeField),
								}),
					)
					require.NoError(t, err)
					v := vr.Payload
					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, fosite.ErrInteractionRequired},
			expectErr:             []bool{true, true, true},
		},
		{
			d:                     "This should pass because login and consent have been granted",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   newCookieJar(t),
			lph:                   passAuthentication(apiClient, false),
			cph:                   passAuthorization(apiClient, false),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:                     "This should pass and set acr values properly",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   newCookieJar(t),
			lph:                   passAuthentication(apiClient, false),
			cph:                   passAuthorization(apiClient, false),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user", ACR: "1"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:                     "This should pass because login and consent have been granted, this time we remember the decision",
			req:                   fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   persistentCJ,
			lph:                   passAuthentication(apiClient, true),
			cph:                   passAuthorization(apiClient, true),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should pass because login and consent have been granted, this time we remember the decision",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					require.True(t, *res.Payload.Skip)
					passAuthentication(apiClient, true)(t)(w, r)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					require.True(t, rrr.Payload.Skip)
					passAuthorization(apiClient, true)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should pass because login and consent have been granted, this time we remember the decision",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					require.True(t, *res.Payload.Skip)
					passAuthentication(apiClient, true)(t)(w, r)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					require.True(t, rrr.Payload.Skip)
					passAuthorization(apiClient, true)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should pass because login was remembered and session id should be set and session context should also work",
			req: fosite.AuthorizeRequest{Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					lr := res.Payload

					assert.True(t, *lr.Skip)
					assert.NotEmpty(t, lr.SessionID)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
							Context:     map[string]interface{}{"foo": "bar"},
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					cr := rrr.Payload

					assert.True(t, cr.Skip)
					assert.NotEmpty(t, cr.LoginSessionID)
					assert.NotEmpty(t, cr.LoginChallenge)
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, cr.Context)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:                     "This should fail because prompt=none, client is public, and redirection scheme is not HTTPS but a custom scheme and acustom domain",
			req:                   fosite.AuthorizeRequest{RedirectURI: mustParseURL(t, "custom://redirection-scheme/path"), Request: fosite.Request{Client: &client.Client{TokenEndpointAuthMethod: "none", OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			prompt:                "none",
			jar:                   persistentCJ,
			lph:                   passAuthentication(apiClient, false),
			expectFinalStatusCode: fosite.ErrConsentRequired.StatusCode(),
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrConsentRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:                     "This should fail because prompt=none, client is public, and redirection scheme is not HTTPS but a custom scheme",
			req:                   fosite.AuthorizeRequest{RedirectURI: mustParseURL(t, "custom://localhost/path"), Request: fosite.Request{Client: &client.Client{TokenEndpointAuthMethod: "none", OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			prompt:                "none",
			jar:                   persistentCJ,
			lph:                   passAuthentication(apiClient, false),
			expectFinalStatusCode: fosite.ErrConsentRequired.StatusCode(),
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrConsentRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:                     "This should pass because prompt=none, client is public, redirection scheme is HTTP and host is localhost",
			req:                   fosite.AuthorizeRequest{RedirectURI: mustParseURL(t, "http://localhost/path"), Request: fosite.Request{Client: &client.Client{TokenEndpointAuthMethod: "none", OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			prompt:                "none",
			jar:                   persistentCJ,
			lph:                   passAuthentication(apiClient, true),
			cph:                   passAuthorization(apiClient, true),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should fail at login screen because subject from accept does not match subject from session",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					lr := res.Payload

					assert.True(t, *lr.Skip)
					assert.Equal(t, "user", *lr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("fooser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.Error(t, err)
					require.Empty(t, vr)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request},
			expectErr:             []bool{true},
		},
		{
			d:   "This should pass and confirm previous authentication and consent because it is a authorization_code",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", Secret: "should-not-be-included"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload

					assert.True(t, *rr.Skip)
					assert.Equal(t, "user", *rr.Subject)
					assert.Empty(t, rr.Client.ClientSecret)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.True(t, rr.Skip)
					assert.Equal(t, "client-id", rr.Client.ClientID)
					assert.Equal(t, "user", rr.Subject)
					assert.Empty(t, rr.Client.ClientSecret)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:      "This should pass and require re-authentication although session is set (because prompt=login)",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			prompt: "login+consent",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload

					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    true,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					body := `{"grant_scope": ["scope-a"], "remember": true}`
					require.NoError(t, err)
					req, err := http.NewRequest("PUT", api.URL+"/oauth2/auth/requests/consent/accept?challenge="+r.URL.Query().Get("consent_challenge"), strings.NewReader(body))
					req.Header.Add("Content-Type", "application/json")
					require.NoError(t, err)

					hres, err := http.DefaultClient.Do(req)
					require.NoError(t, err)
					defer hres.Body.Close()

					var v models.CompletedRequest
					require.NoError(t, json.NewDecoder(hres.Body).Decode(&v))
					require.EqualValues(t, http.StatusOK, hres.StatusCode)
					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       true,
				RememberFor:    0,
				Session:        NewConsentRequestSessionData(),
			},
		},
		{
			d:      "This should pass and require re-authentication although session is set (because max_age=1)",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			maxAge: "1",
			setup: func() {
				time.Sleep(time.Second * 2)
			},
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    true,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.True(t, rr.Skip)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "user", SubjectIdentifier: "user"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should fail because max_age=1 but prompt=none",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ,
			setup: func() {
				time.Sleep(time.Second * 2)
			},
			maxAge:                "1",
			prompt:                "none",
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:                     "This should fail because prompt is none but no auth session exists",
			prompt:                "none",
			req:                   fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   newCookieJar(t),
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:      "This should fail because prompt is none and consent is missing a permission which requires re-authorization of the app",
			prompt: "none",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a", "this-scope-has-not-been-granted-before"}}},
			jar:    persistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.True(t, *rr.Skip)
					assert.Equal(t, "user", *rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrConsentRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:      "This pass and properly require authentication as well as authorization because prompt is set to login and consent - although previous session exists",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ,
			prompt: "login+consent",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("user"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					passAuthorization(apiClient, false)(t)(w, r)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:                     "This should fail because id_token_hint does not match authentication session and prompt is none",
			req:                   fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:                   persistentCJ,
			prompt:                "none",
			idTokenHint:           fooUserIDToken,
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{fosite.ErrLoginRequired},
			expectErr:             []bool{true},
		},
		{
			d:           "This should pass and require authentication because id_token_hint does not match subject from session",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         persistentCJ,
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("not-foouser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: fosite.ErrLoginRequired.StatusCode(),
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrLoginRequired},
			expectErr:             []bool{true, true},
		},
		{
			d:           "This should pass and require authentication because id_token_hint does not match subject from session",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"code"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         persistentCJ,
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:     pointerx.String("foouser"),
							Remember:    false,
							RememberFor: 0,
							Acr:         "1",
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					rrr, err := apiClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))
					require.NoError(t, err)
					rr := rrr.Payload
					assert.False(t, rr.Skip)

					vr, err := apiClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
						WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
						WithBody(&models.AcceptConsentRequest{
							GrantScope:  []string{"scope-a"},
							Remember:    false,
							RememberFor: 0,
							Session: &models.ConsentRequestSession{
								AccessToken: map[string]interface{}{"foo": "bar"},
								IDToken:     map[string]interface{}{"bar": "baz"},
							},
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{Subject: "foouser", SubjectIdentifier: "foouser"},
				GrantedScope:   []string{"scope-a"},
				Remember:       false,
				RememberFor:    0,
				Session: &ConsentRequestSessionData{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
				},
			},
		},
		{
			d:   "This should pass as regularly even though id_token_hint is expired",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: newCookieJar(t),
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "user",
				ExpiresAt: time.Now().Add(-time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph:                   passAuthentication(apiClient, false),
			cph:                   passAuthorization(apiClient, false),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},

		// Pairwise auth
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set by hydra",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd", // this is sha256("fooauth-user76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		}, // these tests depend on one another
		{
			d:   "This should pass as regularly and create a new session with pairwise subject and also with the ID token set",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd",
				ExpiresAt: time.Now().Add(time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					require.NotEmpty(t, *v.RedirectTo)
					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "c737d5e1fec8896d096d49f6b1a73eb45ac7becb87de9ac3f0a350bad2a9c9fd",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		},
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set login request",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ4,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:                pointerx.String("auth-user"),
							ForceSubjectIdentifier: "forced-auth-user",
							Remember:               true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "forced-auth-user",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		}, // these tests depend on one another
		{
			d:   "This should pass as regularly and create a new session with pairwise subject set on login request and also with the ID token set",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id", SubjectType: "pairwise", SectorIdentifierURI: "foo"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ3,
			idTokenHint: genIDToken(t, reg, jwt.IDTokenClaims{
				Subject:   "forced-auth-user",
				ExpiresAt: time.Now().Add(time.Hour),
				IssuedAt:  time.Now(),
			}.ToMapClaims()),
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:                pointerx.String("auth-user"),
							ForceSubjectIdentifier: "forced-auth-user",
							Remember:               false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
			expectSession: &HandledConsentRequest{
				ConsentRequest: &ConsentRequest{
					Subject:           "auth-user",
					SubjectIdentifier: "forced-auth-user",
				},
				GrantedScope: []string{"scope-a"},
				Remember:     false,
				RememberFor:  0,
				Session:      NewConsentRequestSessionData(),
			},
		},

		// checks revoking sessions
		{
			d:   "This should pass as regularly and create a new session and forward data",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ2,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("auth-user"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:      "This should pass and also revoke the session cookie",
			req:    fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:    persistentCJ2,
			prompt: "login",
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("not-auth-user"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		}, // these two tests depend on one another
		{
			d:   "This should require re-authentication because the session was revoked in the previous test",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: persistentCJ2,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.Empty(t, "", rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("foouser"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:   "This should require re-authentication because the session does not exist in the store",
			req: fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar: nonexistentCJ,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.Empty(t, "", rr.Subject)

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("foouser"),
							Remember: true,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			cph:                   acceptRequest(apiClient, nil),
			expectFinalStatusCode: http.StatusOK,
			expectErrType:         []error{ErrAbortOAuth2Request, ErrAbortOAuth2Request, nil},
			expectErr:             []bool{true, true, false},
		},
		{
			d:           "This should fail because the user from the ID token does not match the user from the accept login request",
			req:         fosite.AuthorizeRequest{ResponseTypes: fosite.Arguments{"token", "code", "id_token"}, Request: fosite.Request{Client: &client.Client{OutfacingID: "client-id"}, RequestedScope: []string{"scope-a"}}},
			jar:         newCookieJar(t),
			idTokenHint: fooUserIDToken,
			lph: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					res, err := apiClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
					require.NoError(t, err)
					rr := res.Payload
					assert.False(t, *rr.Skip)
					assert.EqualValues(t, "", *rr.Subject)
					assert.EqualValues(t, "foouser", rr.OidcContext.IDTokenHintClaims.(map[string]interface{})["sub"])

					vr, err := apiClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
						WithLoginChallenge(r.URL.Query().Get("login_challenge")).
						WithBody(&models.AcceptLoginRequest{
							Subject:  pointerx.String("not-foouser"),
							Remember: false,
						}))
					require.NoError(t, err)
					v := vr.Payload

					http.Redirect(w, r, *v.RedirectTo, http.StatusFound)
				}
			},
			expectFinalStatusCode: http.StatusBadRequest,
			expectErrType:         []error{ErrAbortOAuth2Request, fosite.ErrLoginRequired},
			expectErr:             []bool{true, true},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			if tc.lph != nil {
				lph = tc.lph(t)
			} else {
				lph = noopHandler(t)
			}

			if tc.cph != nil {
				cph = tc.cph(t)
			} else {
				cph = noopHandler(t)
			}

			calls := -1
			aph = func(w http.ResponseWriter, r *http.Request) {
				calls++
				require.True(t, len(tc.expectErrType) >= calls+1, "%d (expect) < %d (got)", len(tc.expectErrType), calls+1)
				require.True(t, len(tc.expectErr) >= calls+1, "%d (expect) < %d (got)", len(tc.expectErr), calls+1)
				require.NoError(t, r.ParseForm())
				tc.req.Form = r.Form

				c, err := strategy.HandleOAuth2AuthorizationRequest(w, r, &tc.req)
				t.Logf("DefaultStrategy returned at call %d:\n\tresult: %+v\n\terr: %s", calls, c, err)

				if tc.expectErr[calls] {
					assert.Error(t, err)
					if tc.expectErrType[calls] != nil {
						assert.EqualError(t, err, tc.expectErrType[calls].Error(), "%+v", err)
					}
				} else {
					require.NoError(t, err)
					if tc.expectSession != nil {
						require.NotNil(t, c)
						assert.EqualValues(t, tc.expectSession.GrantedScope, c.GrantedScope)
						assert.EqualValues(t, tc.expectSession.Remember, c.Remember)
						assert.EqualValues(t, tc.expectSession.Session, c.Session)
						assert.EqualValues(t, tc.expectSession.RememberFor, c.RememberFor)
						assert.EqualValues(t, tc.expectSession.ConsentRequest.Subject, c.ConsentRequest.Subject)
						assert.EqualValues(t, tc.expectSession.ConsentRequest.SubjectIdentifier, c.ConsentRequest.SubjectIdentifier)
					}
				}

				if errors.Is(err, ErrAbortOAuth2Request) {
					// nothing to do, indicates redirect
				} else if err != nil {
					writer.WriteError(w, r, err)
				} else {
					writer.Write(w, r, c)
				}
			}

			cl := &http.Client{
				Jar: tc.jar,
			}
			resp, err := cl.Get(
				ap.URL + "?" +
					"login_verifier=" + tc.lv + "&" +
					"consent_verifier=" + tc.cv + "&" +
					"prompt=" + tc.prompt + "&" +
					"max_age=" + tc.maxAge + "&" +
					"id_token_hint=" + tc.idTokenHint + "&" + tc.other,
			)
			require.NoError(t, err)
			out, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			resp.Body.Close()
			assert.EqualValues(t, tc.expectFinalStatusCode, resp.StatusCode, "%s\n%s", resp.Request.URL.String(), out)
			// assert.Empty(t, resp.Request.URL.String())
		})
	}
}

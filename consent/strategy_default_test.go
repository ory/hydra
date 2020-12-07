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

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	"github.com/ory/x/sqlxx"

	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/urlx"

	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/x"
)

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

func newAuthCookieJar(t *testing.T, reg driver.Registry, u, sessionID string) http.CookieJar {
	cj, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	secrets := reg.Config().Source().Strings(config.KeyGetCookieSecrets)
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

	conf.Set(config.KeyIssuerURL, logoutServer.URL)
	conf.Set(config.KeyLogoutURL, logoutProviderServer.URL)
	conf.Set(config.KeyLogoutRedirectURL, defaultRedirServer.URL)

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

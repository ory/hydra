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
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
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

func TestAuthCodeSuite(t *testing.T) {
	// Regular request without any previous authentication
	// Regular request with previous authentication
	// Regular request with previous authentication and prompt login/consent
	// Regular request with previous authentication and prompt none
	// Regular request with previous authentication and prompt none and proper id_token_hint
	// Regular request without authentication and prompt none (fail)
	// Regular request with previous authentication and prompt none and very low max_age (fail)
	// Regular request with previous authentication and prompt none and mismatching id_token_hint (fail)
	// Regular request fails if login is denied
	// Regular request fails if consent is denied

	var callbackHandler *httprouter.Handle
	router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		(*callbackHandler)(w, r, ps)
	})
	m := sync.Mutex{}

	var code string
	for k, tc := range []struct {
		d                         string
		cb                        func(t *testing.T) httprouter.Handle
		authURL                   string
		shouldPassConsentStrategy bool
		expectOAuthAuthError      bool
		expectOAuthTokenError     bool
		authTime                  time.Time
		requestTime               time.Time
	}{
		{
			d:                         "should pass request if no previous authN or authZ exists",
			authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
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
			authTime:                  time.Now().UTC().Add(time.Minute),
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

			resp, err := (&http.Client{Jar: newCookieJar()}).Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if tc.expectOAuthAuthError {
				return
			}

			require.NotEmpty(t, code)

			token, err := oauthConfig.Exchange(oauth2.NoContext, code)

			if tc.expectOAuthTokenError {
				require.Error(t, err)
				return
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

			res, err := testRefresh(t, token)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			t.Run("duplicate code exchange fails", func(t *testing.T) {
				token, err := oauthConfig.Exchange(oauth2.NoContext, code)
				require.Error(t, err)
				require.Nil(t, token)
			})

			code = ""
		})
	}
}

func testRefresh(t *testing.T, token *oauth2.Token) (*http.Response, error) {
	req, err := http.NewRequest("POST", oauthClientConfig.TokenURL, strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{token.RefreshToken},
	}.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthClientConfig.ClientID, oauthClientConfig.ClientSecret)

	return http.DefaultClient.Do(req)
}

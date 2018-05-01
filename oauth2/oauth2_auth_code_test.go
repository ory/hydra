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
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthCode(t *testing.T) {
	var callbackHandler httprouter.Handle
	router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		callbackHandler(w, r, ps)
	})

	t.Run("case=test accept consent request", func(t *testing.T) {
		var code string

		callbackHandler = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			code = r.URL.Query().Get("code")
			w.Write([]byte(r.URL.Query().Get("code")))
		}

		cookieJar, _ := cookiejar.New(nil)
		req, err := http.NewRequest("GET", oauthConfig.AuthCodeURL("some-foo-state"), nil)
		require.NoError(t, err)

		resp, err := (&http.Client{Jar: cookieJar}).Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		require.NotEmpty(t, code)

		token, err := oauthConfig.Exchange(oauth2.NoContext, code)
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

		time.Sleep(time.Second * 5)

		res, err := testRefresh(t, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		t.Run("duplicate code exchange fails", func(t *testing.T) {
			token, err := oauthConfig.Exchange(oauth2.NoContext, code)
			require.Error(t, err)
			require.Nil(t, token)
		})
	})

	t.Run("case=test deny consent request", func(t *testing.T) {
		consentStrategy.deny = true

		callbackHandler = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			t.Logf("GOT URL: %s", r.URL.String())

			assert.Equal(t, "The request is not allowed", r.URL.Query().Get("error_description"))
			assert.Equal(t, "request_forbidden", r.URL.Query().Get("error"))
			w.WriteHeader(http.StatusNoContent)
		}

		cookieJar, _ := cookiejar.New(nil)
		req, err := http.NewRequest("GET", oauthConfig.AuthCodeURL("some-foo-state"), nil)
		require.NoError(t, err)

		resp, err := (&http.Client{Jar: cookieJar}).Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		consentStrategy.deny = false
	})
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

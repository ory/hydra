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

package consent

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

var passAuthentication = func(apiClient *swagger.AdminApi, remember bool) func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			v, res, err := apiClient.AcceptLoginRequest(r.URL.Query().Get("login_challenge"), swagger.AcceptLoginRequest{
				Subject:     "user",
				Remember:    remember,
				RememberFor: 0,
				Acr:         "1",
			})
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, res.StatusCode)
			require.NotEmpty(t, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}
}

var passAuthorization = func(apiClient *swagger.AdminApi, remember bool) func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			v, res, err := apiClient.AcceptConsentRequest(r.URL.Query().Get("consent_challenge"), swagger.AcceptConsentRequest{
				GrantScope:  []string{"scope-a"},
				Remember:    remember,
				RememberFor: 0,
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
	}
}

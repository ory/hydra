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
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	jwtgo "github.com/ory/fosite/token/jwt"

	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/urlx"

	"net/url"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/ioutilx"
)

func checkAndAcceptLoginHandler(t *testing.T, apiClient admin.ClientService, subject string, cb func(*testing.T, *admin.GetLoginRequestOK, error) *models.AcceptLoginRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := apiClient.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(r.URL.Query().Get("login_challenge")))
		payload := cb(t, res, err)
		payload.Subject = &subject

		v, err := apiClient.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(r.URL.Query().Get("login_challenge")).
			WithBody(payload))
		require.NoError(t, err)
		require.NotEmpty(t, *v.Payload.RedirectTo)
		http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
	}
}

func checkAndAcceptConsentHandler(t *testing.T, apiClient admin.ClientService, cb func(*testing.T, *admin.GetConsentRequestOK, error) *models.AcceptConsentRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := apiClient.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")))

		v, err := apiClient.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
			WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
			WithBody(cb(t, res, err)))
		require.NoError(t, err)
		require.NotEmpty(t, *v.Payload.RedirectTo)
		http.Redirect(w, r, *v.Payload.RedirectTo, http.StatusFound)
	}
}
func makeOAuth2Request(t *testing.T, reg driver.Registry, hc *http.Client, oc *client.Client, values url.Values) (gjson.Result, *http.Response) {
	if hc == nil {
		hc = testhelpers.NewEmptyJarClient(t)
	}

	values.Add("response_type", "code")
	values.Add("state", uuid.New().String())
	values.Add("client_id", oc.OutfacingID)
	res, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(), values).String())
	require.NoError(t, err)
	defer res.Body.Close()

	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body)), res
}

func createClient(t *testing.T, reg driver.Registry, c *client.Client) *client.Client {
	secret := uuid.New().String()
	c.Secret = secret
	c.Scope = "openid offline"
	c.OutfacingID = uuid.New().String()
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))
	c.Secret = secret
	return c
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
	cookie, _ := reg.CookieStore().Get(hr, CookieName(reg.Config().TLS(config.PublicInterface).Enabled(), CookieAuthenticationName))

	cookie.Values[CookieAuthenticationSIDName] = sessionID
	cookie.Options.HttpOnly = true

	rw := httptest.NewRecorder()
	require.NoError(t, cookie.Save(hr, rw))

	cj.SetCookies(urlx.ParseOrPanic(u), rw.Result().Cookies())
	return cj
}

func genIDToken(t *testing.T, reg driver.Registry, c jwtgo.MapClaims) string {
	r, _, err := reg.OpenIDJWTStrategy().Generate(context.TODO(), c, jwt.NewHeaders())
	require.NoError(t, err)
	return r
}

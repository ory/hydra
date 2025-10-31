// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	. "github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/urlx"
)

func checkAndAcceptDeviceHandler(t *testing.T, apiClient *hydra.APIClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCode := r.URL.Query().Get("user_code")
		payload := hydra.AcceptDeviceUserCodeRequest{
			UserCode: &userCode,
		}

		v, _, err := apiClient.OAuth2API.AcceptUserCodeRequest(context.Background()).
			DeviceChallenge(r.URL.Query().Get("device_challenge")).
			AcceptDeviceUserCodeRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}
}

func checkAndAcceptLoginHandler(t *testing.T, apiClient *hydra.APIClient, subject string, cb func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := apiClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
		payload := cb(t, res, err)
		payload.Subject = subject

		v, _, err := apiClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
			LoginChallenge(r.URL.Query().Get("login_challenge")).
			AcceptOAuth2LoginRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}
}

func checkAndAcceptConsentHandler(t *testing.T, apiClient *hydra.APIClient, cb func(*testing.T, *hydra.OAuth2ConsentRequest, error) hydra.AcceptOAuth2ConsentRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := apiClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
		payload := cb(t, res, err)

		v, _, err := apiClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
			ConsentChallenge(res.Challenge).
			AcceptOAuth2ConsentRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}
}

func makeOAuth2Request(t *testing.T, reg *driver.RegistrySQL, hc *http.Client, oc *client.Client, values url.Values) (gjson.Result, *http.Response) {
	ctx := context.Background()
	if hc == nil {
		hc = testhelpers.NewEmptyJarClient(t)
	}

	values.Add("response_type", "code")
	values.Add("state", uuid.Must(uuid.NewV4()).String())
	values.Add("client_id", oc.GetID())
	values.Add("redirect_uri", oc.GetRedirectURIs()[0])
	res, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(ctx), values).String())
	require.NoError(t, err)
	defer res.Body.Close() //nolint:errcheck

	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body)), res
}

func makeOAuth2DeviceAuthRequest(t *testing.T, reg *driver.RegistrySQL, hc *http.Client, oc *client.Client, scope string) (gjson.Result, *http.Response) {
	ctx := context.Background()
	if hc == nil {
		hc = testhelpers.NewEmptyJarClient(t)
	}

	data := url.Values{}
	data.Set("scope", scope)
	data.Set("client_id", oc.GetID())
	req, err := http.NewRequest(
		http.MethodPost,
		reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
		strings.NewReader(data.Encode()),
	)
	require.NoError(t, err)
	req.SetBasicAuth(oc.GetID(), oc.Secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := hc.Do(req)
	require.NoError(t, err)

	defer res.Body.Close() //nolint:errcheck

	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body)), res
}

func makeOAuth2DeviceVerificationRequest(t *testing.T, reg *driver.RegistrySQL, hc *http.Client, oc *client.Client, values url.Values) (gjson.Result, *http.Response) {
	ctx := context.Background()
	if hc == nil {
		hc = testhelpers.NewEmptyJarClient(t)
	}

	values.Add("client_id", oc.GetID())
	res, err := hc.Get(urlx.CopyWithQuery(urlx.AppendPaths(reg.Config().PublicURL(ctx), oauth2.DeviceVerificationPath), values).String())
	require.NoError(t, err)
	defer res.Body.Close() //nolint:errcheck

	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body)), res
}

func createClient(t *testing.T, reg *driver.RegistrySQL, c *client.Client) *client.Client {
	secret := uuid.Must(uuid.NewV4()).String()
	c.Secret = secret
	c.Scope = "openid offline"
	c.ID = uuid.Must(uuid.NewV4()).String()
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))
	c.Secret = secret
	return c
}

func newAuthCookieJar(t *testing.T, reg *driver.RegistrySQL, u, sessionID string) http.CookieJar {
	ctx := context.Background()
	cj, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)

	hr := &http.Request{Header: map[string][]string{}, URL: urlx.ParseOrPanic(u), RequestURI: u}
	s, err := reg.CookieStore(ctx)
	require.NoError(t, err)
	cookie, _ := s.Get(hr, reg.Config().SessionCookieName(ctx))

	cookie.Values[CookieAuthenticationSIDName] = sessionID
	cookie.Options.HttpOnly = true

	rw := httptest.NewRecorder()
	require.NoError(t, cookie.Save(hr, rw))

	cj.SetCookies(urlx.ParseOrPanic(u), rw.Result().Cookies())
	return cj
}

func genIDToken(t *testing.T, reg *driver.RegistrySQL, c jwt.MapClaims) string {
	r, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), c, jwt.NewHeaders())
	require.NoError(t, err)
	return r
}

func checkAndDuplicateAcceptLoginHandler(t *testing.T, apiClient *hydra.APIClient, subject string, cb func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := apiClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
		payload := cb(t, res, err)
		payload.Subject = subject

		v, _, err := apiClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
			LoginChallenge(r.URL.Query().Get("login_challenge")).
			AcceptOAuth2LoginRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)

		v2, _, err := apiClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
			LoginChallenge(r.URL.Query().Get("login_challenge")).
			AcceptOAuth2LoginRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v2.RedirectTo)
		http.Redirect(w, r, v2.RedirectTo, http.StatusFound)
	}
}

func checkAndDuplicateAcceptConsentHandler(t *testing.T, apiClient *hydra.APIClient, cb func(*testing.T, *hydra.OAuth2ConsentRequest, error) hydra.AcceptOAuth2ConsentRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := apiClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).
			ConsentChallenge(r.URL.Query().Get("consent_challenge")).
			Execute()
		payload := cb(t, res, err)

		v, _, err := apiClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
			ConsentChallenge(r.URL.Query().Get("consent_challenge")).
			AcceptOAuth2ConsentRequest(payload).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)

		res2, _, err := apiClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
		payload2 := cb(t, res2, err)

		v2, _, err := apiClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
			ConsentChallenge(r.URL.Query().Get("consent_challenge")).
			AcceptOAuth2ConsentRequest(payload2).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v2.RedirectTo)
		http.Redirect(w, r, v2.RedirectTo, http.StatusFound)
	}
}

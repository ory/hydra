// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ory/hydra/v2/driver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	jwtgo "github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/configx"
	"github.com/ory/x/ioutilx"
)

func makeDeps(t *testing.T, defaultLogoutURL string) (*kratos.FakeKratos, *driver.RegistrySQL, *httptest.Server, *httptest.Server, *hydra.APIClient) {
	fakeKratos := kratos.NewFake()
	reg := testhelpers.NewRegistryMemory(t,
		driver.WithConfigOptions(configx.WithValues(map[string]any{
			config.KeyAccessTokenStrategy:  "opaque",
			config.KeyConsentRequestMaxAge: time.Hour,
		})),
		driver.WithKratosClient(fakeKratos))
	reg.Config().MustSet(t.Context(), config.KeyLogoutRedirectURL, defaultLogoutURL)
	publicTS, adminTS := testhelpers.NewOAuth2Server(t.Context(), t, reg)
	t.Cleanup(publicTS.Close)
	t.Cleanup(adminTS.Close)

	adminApi := hydra.NewAPIClient(hydra.NewConfiguration())
	adminApi.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	return fakeKratos, reg, publicTS, adminTS, adminApi
}

func createBrowserWithSession(t *testing.T, c *client.Client, reg *driver.RegistrySQL) *http.Client {
	hc := testhelpers.NewEmptyJarClient(t)
	makeOAuth2Request(t, reg, hc, c, url.Values{})
	return hc
}

func createSampleClient(t *testing.T, reg *driver.RegistrySQL, customPostLogoutURL string) *client.Client {
	return createClient(t, reg, &client.Client{
		RedirectURIs:           []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)},
		PostLogoutRedirectURIs: []string{customPostLogoutURL},
	})
}

func createClientWithBackchannelLogout(t *testing.T, reg *driver.RegistrySQL, customPostLogoutURL string, wg *sync.WaitGroup, cb func(t *testing.T, logoutToken gjson.Result)) *client.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		require.NoError(t, r.ParseForm())
		lt := r.PostFormValue("logout_token")
		assert.NotEmpty(t, lt)
		token, err := reg.OpenIDJWTSigner().Decode(r.Context(), lt)
		require.NoError(t, err)

		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(token.Claims))
		cb(t, gjson.Parse(b.String()))
	}))
	t.Cleanup(server.Close)

	return createClient(t, reg, &client.Client{
		BackChannelLogoutURI:   server.URL,
		RedirectURIs:           []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)},
		PostLogoutRedirectURIs: []string{customPostLogoutURL},
	})
}

func makeLogoutRequest(t *testing.T, publicTSURL string, hc *http.Client, method string, values url.Values) (body string, resp *http.Response) {
	var err error
	if method == http.MethodGet {
		resp, err = hc.Get(publicTSURL + "/oauth2/sessions/logout?" + values.Encode())
	} else if method == http.MethodPost {
		resp, err = hc.PostForm(publicTSURL+"/oauth2/sessions/logout", values)
	}
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	return string(ioutilx.MustReadAll(resp.Body)), resp
}

func makeHeadlessLogoutRequest(t *testing.T, adminTSURL string, hc *http.Client, values url.Values) (body string, resp *http.Response) {
	var err error
	req, err := http.NewRequest(http.MethodDelete, adminTSURL+"/admin/oauth2/auth/sessions/login?"+values.Encode(), nil)
	require.NoError(t, err)

	resp, err = hc.Do(req)

	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	return string(ioutilx.MustReadAll(resp.Body)), resp
}

func logoutViaHeadlessAndExpectNoContent(t *testing.T, adminTSURL string, browser *http.Client, values url.Values) {
	_, res := makeHeadlessLogoutRequest(t, adminTSURL, browser, values)
	assert.EqualValues(t, http.StatusNoContent, res.StatusCode)
}

func logoutViaHeadlessAndExpectError(t *testing.T, adminTSURL string, browser *http.Client, values url.Values, expectedErrorMessage string) {
	body, res := makeHeadlessLogoutRequest(t, adminTSURL, browser, values)
	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, body, expectedErrorMessage)
}

func logoutAndExpectErrorPage(t *testing.T, publicTSURL string, browser *http.Client, method string, values url.Values, expectedErrorMessage string) {
	body, res := makeLogoutRequest(t, publicTSURL, browser, method, values)
	assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)
	assert.Contains(t, body, expectedErrorMessage)
}

func logoutAndExpectPostLogoutPage(t *testing.T, publicTSURL string, browser *http.Client, method string, values url.Values, expectedMessage string) {
	body, res := makeLogoutRequest(t, publicTSURL, browser, method, values)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Contains(t, body, expectedMessage)
}

func setupCheckAndAcceptLogoutHandler(t *testing.T, reg *driver.RegistrySQL, adminApi *hydra.APIClient, wg *sync.WaitGroup, cb func(*testing.T, *hydra.OAuth2LogoutRequest, error)) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wg != nil {
			defer wg.Done()
		}

		challenge := r.URL.Query().Get("logout_challenge")
		res, _, err := adminApi.OAuth2API.GetOAuth2LogoutRequest(t.Context()).LogoutChallenge(challenge).Execute()
		if cb != nil {
			cb(t, res, err)
		} else {
			require.NoError(t, err)
		}
		require.NotNil(t, res)
		require.NotNil(t, res.Challenge)

		v, _, err := adminApi.OAuth2API.AcceptOAuth2LogoutRequest(t.Context()).LogoutChallenge(*res.Challenge).Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}))

	t.Cleanup(server.Close)

	reg.Config().MustSet(t.Context(), config.KeyLogoutURL, server.URL)
}

func acceptLoginAsAndWatchSidForConsumers(t *testing.T, reg *driver.RegistrySQL, adminApi *hydra.APIClient, subject string, sid chan<- string, remember bool, numSidConsumers int) {
	testhelpers.NewLoginConsentUI(t, reg.Config(),
		checkAndAcceptLoginHandler(t, adminApi, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
			require.NoError(t, err)
			return hydra.AcceptOAuth2LoginRequest{
				Remember:                  new(true),
				IdentityProviderSessionId: new(kratos.FakeSessionID),
			}
		}),
		checkAndAcceptConsentHandler(t, adminApi, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
			require.NoError(t, err)
			if sid != nil {
				for range numSidConsumers {
					sid <- *res.LoginSessionId
				}
			}
			return hydra.AcceptOAuth2ConsentRequest{Remember: new(remember)}
		}))
}

func acceptLoginAsAndWatchSid(t *testing.T, reg *driver.RegistrySQL, adminApi *hydra.APIClient, subject string) <-chan string {
	sid := make(chan string, 1)
	acceptLoginAsAndWatchSidForConsumers(t, reg, adminApi, subject, sid, true, 1)
	return sid
}

func acceptLoginAs(t *testing.T, reg *driver.RegistrySQL, adminApi *hydra.APIClient, subject string) {
	acceptLoginAsAndWatchSidForConsumers(t, reg, adminApi, subject, nil, true, 0)
}

func TestLogoutFlows(t *testing.T) {
	t.Parallel()
	// Only truly immutable values are shared across subtests.
	defaultRedirectedMessage := "redirected to default server"
	postLogoutCallback := func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		_, _ = fmt.Fprintf(w, "%s%s%s", defaultRedirectedMessage, r.Form.Get("state"), strings.TrimLeft(r.URL.Path, "/"))
	}
	defaultLogoutURL := testhelpers.NewCallbackURL(t, "logged-out", postLogoutCallback)
	customPostLogoutURL := testhelpers.NewCallbackURL(t, "logged-out/custom", postLogoutCallback)
	subject := "aeneas-rekkas"

	newWg := func(add int) *sync.WaitGroup {
		var wg sync.WaitGroup
		wg.Add(add)
		return &wg
	}

	t.Run("case=should ignore / redirect non-rp initiated logout if no session exists", func(t *testing.T) {
		t.Parallel()
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectPostLogoutPage(t, publicTS.URL, new(http.Client), http.MethodGet, url.Values{}, defaultRedirectedMessage)
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectPostLogoutPage(t, publicTS.URL, new(http.Client), http.MethodPost, url.Values{}, defaultRedirectedMessage)
		})
	})

	t.Run("case=should fail if non-rp initiated logout is initiated with state (indicating rp-flow)", func(t *testing.T) {
		t.Parallel()
		expectedMessage := "Logout failed because query parameter state is set but id_token_hint is missing"
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectErrorPage(t, publicTS.URL, new(http.Client), http.MethodGet, url.Values{"state": {"foobar"}}, expectedMessage)
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectErrorPage(t, publicTS.URL, new(http.Client), http.MethodPost, url.Values{"state": {"foobar"}}, expectedMessage)
		})
	})

	t.Run("case=should fail if non-rp initiated logout is initiated with post_logout_redirect_uri (indicating rp-flow)", func(t *testing.T) {
		t.Parallel()
		expectedMessage := "Logout failed because query parameter post_logout_redirect_uri is set but id_token_hint is missing"
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectErrorPage(t, publicTS.URL, new(http.Client), http.MethodGet, url.Values{"post_logout_redirect_uri": {"foobar"}}, expectedMessage)
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, _, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			logoutAndExpectErrorPage(t, publicTS.URL, new(http.Client), http.MethodPost, url.Values{"post_logout_redirect_uri": {"foobar"}}, expectedMessage)
		})
	})

	t.Run("case=should ignore / redirect non-rp initiated logout if a session cookie exists but the session itself is no longer active / invalid", func(t *testing.T) {
		t.Parallel()
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			browser := &http.Client{Jar: newAuthCookieJar(t, reg, publicTS.URL, "i-do-not-exist")}
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, _ := makeDeps(t, defaultLogoutURL)
			browser := &http.Client{Jar: newAuthCookieJar(t, reg, publicTS.URL, "i-do-not-exist")}
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodPost, url.Values{}, defaultRedirectedMessage)
		})
	})

	t.Run("case=should redirect to logout provider if session exists and it's not rp-flow", func(t *testing.T) {
		t.Parallel()
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			acceptLoginAs(t, reg, adminApi, subject)
			wg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, func(t *testing.T, res *hydra.OAuth2LogoutRequest, err error) {
				require.NoError(t, err)
				assert.EqualValues(t, subject, *res.Subject)
				assert.NotEmpty(t, subject, res.Sid)
			})
			logoutAndExpectPostLogoutPage(t, publicTS.URL, createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg), http.MethodGet, url.Values{}, defaultRedirectedMessage)
			wg.Wait()
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			acceptLoginAs(t, reg, adminApi, subject)
			wg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, func(t *testing.T, res *hydra.OAuth2LogoutRequest, err error) {
				require.NoError(t, err)
				assert.EqualValues(t, subject, *res.Subject)
				assert.NotEmpty(t, subject, res.Sid)
			})
			logoutAndExpectPostLogoutPage(t, publicTS.URL, createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg), http.MethodPost, url.Values{}, defaultRedirectedMessage)
			wg.Wait()
		})
	})

	t.Run("case=should redirect to post logout url because logout was already done before", func(t *testing.T) {
		t.Parallel()
		// Formerly: should redirect to logout provider because the session has been removed previously.
		// Each inner subtest invalidates its own session first, then confirms subsequent calls redirect
		// without triggering the logout handler again.
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			acceptLoginAs(t, reg, adminApi, subject)
			browser := createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg)
			wg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, nil)
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)
			wg.Wait() // ensure logout ui was called exactly once
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			acceptLoginAs(t, reg, adminApi, subject)
			browser := createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg)
			wg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, nil)
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)
			wg.Wait() // ensure logout ui was called exactly once
			logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodPost, url.Values{}, defaultRedirectedMessage)
		})
	})

	t.Run("case=should handle double-submit of the logout challenge gracefully", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		acceptLoginAs(t, reg, adminApi, subject)
		browser := createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg)

		var logoutReq *hydra.OAuth2LogoutRequest
		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, func(t *testing.T, req *hydra.OAuth2LogoutRequest, err error) {
			require.NoError(t, err)
			require.NotNil(t, req.Challenge)
			logoutReq = req
		})

		// run once to log out
		logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)

		require.NotNil(t, logoutReq)
		require.NotEmpty(t, logoutReq.GetChallenge())

		// run again: with the stateless logout flow the challenge is an opaque
		// AEAD blob with no server-side state, so GetOAuth2LogoutRequest must
		// return the same values on a repeated call.
		repeated, _, err := adminApi.OAuth2API.GetOAuth2LogoutRequest(t.Context()).LogoutChallenge(logoutReq.GetChallenge()).Execute()
		require.NoError(t, err)
		require.NotNil(t, repeated)
		assert.Equal(t, logoutReq.GetChallenge(), repeated.GetChallenge())
		assert.Equal(t, logoutReq.GetSubject(), repeated.GetSubject())
		assert.Equal(t, logoutReq.GetSid(), repeated.GetSid())
		assert.Equal(t, logoutReq.GetRequestUrl(), repeated.GetRequestUrl())
		assert.Equal(t, logoutReq.GetRpInitiated(), repeated.GetRpInitiated())

		v, _, err := adminApi.OAuth2API.AcceptOAuth2LogoutRequest(t.Context()).LogoutChallenge(logoutReq.GetChallenge()).Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)

		res, err := browser.Get(v.RedirectTo)
		require.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("case=should handle double-submit of the logout verifier", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		acceptLoginAs(t, reg, adminApi, subject)
		c := createSampleClient(t, reg, customPostLogoutURL)
		browser := createBrowserWithSession(t, c, reg)

		sessionCookie := getSessionCookie(t, browser, publicTS.URL)
		require.NotNil(t, sessionCookie)

		// capture the request with the logout verifier
		var verifiedLogoutReq *http.Request
		browser.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if lv := req.FormValue("logout_verifier"); lv != "" {
				verifiedLogoutReq = req.Clone(t.Context())
			}
			return nil
		}

		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)

		// run once to log out
		logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, http.MethodGet, url.Values{}, defaultRedirectedMessage)
		require.Nil(t, getSessionCookie(t, browser, publicTS.URL), "Session cookie should be gone after logout")

		// Attempt to re-use the logout verifier from a different browser
		// session. This checks that the logout verifier can be used only once,
		// and that a session cannot be invalidated by tricking another user
		// into visiting the /oauth2/sessions/logout?logout_verifier=<verifier>
		// URL.
		browser2 := createBrowserWithSession(t, c, reg)
		sessionCookie2 := getSessionCookie(t, browser2, publicTS.URL)
		require.NotNil(t, sessionCookie2)
		require.NotEqual(t, sessionCookie.Value, sessionCookie2.Value, "Should have two different session cookies after logout + login again")

		// Re-use the logout verifier from the first browser session
		res, err := browser2.Do(verifiedLogoutReq)
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		require.Equal(t, 200, res.StatusCode)

		sessionCookie2AfterReuse := getSessionCookie(t, browser2, publicTS.URL)
		require.NotNil(t, sessionCookie2AfterReuse, "Session cookie should still exist after re-using the logout verifier")
		require.Equal(t, sessionCookie2.Value, sessionCookie2AfterReuse.Value, "Session cookie should not have changed after re-using the logout verifier")

		// check that the session in browser2 is still valid and recognized
		var wg sync.WaitGroup
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminApi, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				defer wg.Done()
				require.NoError(t, err)
				assert.True(t, res.Skip)
				return hydra.AcceptOAuth2LoginRequest{Remember: new(true)}
			}),
			checkAndAcceptConsentHandler(t, adminApi, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				return hydra.AcceptOAuth2ConsentRequest{Remember: new(true)}
			}))

		// Make an oauth 2 request to trigger the login check.
		wg.Add(1)
		_, res = makeOAuth2Request(t, reg, browser2, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("code"))
		wg.Wait()
	})

	t.Run("case=valid logout verifiers cannot be used to log out other users", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		acceptLoginAs(t, reg, adminApi, subject)
		c := createSampleClient(t, reg, customPostLogoutURL)
		attacker := createBrowserWithSession(t, c, reg)
		// Capture the request with the logout verifier and stop navigation just before.
		var verifiedLogoutReq *http.Request
		attacker.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if lv := req.FormValue("logout_verifier"); lv != "" {
				verifiedLogoutReq = req.Clone(t.Context())
				return http.ErrUseLastResponse
			}
			return nil
		}

		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)

		_, res := makeLogoutRequest(t, publicTS.URL, attacker, http.MethodGet, url.Values{})
		require.Equal(t, http.StatusFound, res.StatusCode)
		require.NotNil(t, verifiedLogoutReq)

		victim := createBrowserWithSession(t, c, reg)
		res, err := victim.Do(verifiedLogoutReq) // we've tricked the victim into clicking a link with a valid login_verifier
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), defaultRedirectedMessage) // we've been redirected to the post-logout-redirect-uri

		// Check if victim's session cookie is still valid
		var wg sync.WaitGroup
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminApi, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				defer wg.Done()
				require.NoError(t, err)
				assert.True(t, res.Skip)
				return hydra.AcceptOAuth2LoginRequest{Remember: new(true)}
			}),
			checkAndAcceptConsentHandler(t, adminApi, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				return hydra.AcceptOAuth2ConsentRequest{Remember: new(true)}
			}))

		// Make an oauth 2 request to trigger the login check.
		wg.Add(1)
		_, res = makeOAuth2Request(t, reg, victim, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("code"))
		wg.Wait()
	})

	t.Run("case=should handle an invalid logout challenge", func(t *testing.T) {
		t.Parallel()
		_, _, _, _, adminApi := makeDeps(t, defaultLogoutURL)
		_, res, err := adminApi.OAuth2API.GetOAuth2LogoutRequest(t.Context()).LogoutChallenge("some-invalid-challenge").Execute()
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		_, res, err = adminApi.OAuth2API.AcceptOAuth2LogoutRequest(t.Context()).LogoutChallenge("some-invalid-challenge").Execute()
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		res, err = adminApi.OAuth2API.RejectOAuth2LogoutRequest(t.Context()).LogoutChallenge("some-invalid-challenge").Execute()
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("case=should handle an invalid logout verifier", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
		logoutAndExpectErrorPage(t, publicTS.URL, http.DefaultClient, http.MethodGet, url.Values{
			"logout_verifier": {"an-invalid-verifier"},
		}, "Description: Unable to locate the requested resource")
	})

	t.Run("case=should execute backchannel logout if issued without rp-involvement", func(t *testing.T) {
		t.Parallel()
		t.Run("method=get", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
			logoutWg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, logoutWg, nil)
			backChannelWG := newWg(1)
			c := createClientWithBackchannelLogout(t, reg, customPostLogoutURL, backChannelWG, func(t *testing.T, logoutToken gjson.Result) {
				assert.EqualValues(t, <-sid, logoutToken.Get("sid").String(), logoutToken.Raw)
				assert.Empty(t, logoutToken.Get("sub").String(), logoutToken.Raw) // The sub claim should be empty because it doesn't work with forced obfuscation and thus we can't easily recover it.
				assert.Empty(t, logoutToken.Get("nonce").String(), logoutToken.Raw)
			})
			logoutAndExpectPostLogoutPage(t, publicTS.URL, createBrowserWithSession(t, c, reg), http.MethodGet, url.Values{}, defaultRedirectedMessage)
			logoutWg.Wait()
			backChannelWG.Wait()
		})
		t.Run("method=post", func(t *testing.T) {
			t.Parallel()
			_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
			sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
			logoutWg := newWg(1)
			setupCheckAndAcceptLogoutHandler(t, reg, adminApi, logoutWg, nil)
			backChannelWG := newWg(1)
			c := createClientWithBackchannelLogout(t, reg, customPostLogoutURL, backChannelWG, func(t *testing.T, logoutToken gjson.Result) {
				assert.EqualValues(t, <-sid, logoutToken.Get("sid").String(), logoutToken.Raw)
				assert.Empty(t, logoutToken.Get("sub").String(), logoutToken.Raw) // The sub claim should be empty because it doesn't work with forced obfuscation and thus we can't easily recover it.
				assert.Empty(t, logoutToken.Get("nonce").String(), logoutToken.Raw)
			})
			logoutAndExpectPostLogoutPage(t, publicTS.URL, createBrowserWithSession(t, c, reg), http.MethodPost, url.Values{}, defaultRedirectedMessage)
			logoutWg.Wait()
			backChannelWG.Wait()
		})
	})

	// Only do GET requests from here on out, POST should be tested enough to ensure that it is working fine already.

	t.Run("case=should fail several flows when id_token_hint is invalid", func(t *testing.T) {
		t.Parallel()
		t.Run("case=should error when rp-flow without valid id token", func(t *testing.T) {
			t.Parallel()
			t.Run("method=get", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				acceptLoginAs(t, reg, adminApi, "aeneas-rekkas")
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg)
				values := url.Values{"state": {"1234"}, "post_logout_redirect_uri": {customPostLogoutURL}, "id_token_hint": {"i am not valid"}}
				logoutAndExpectErrorPage(t, publicTS.URL, browser, http.MethodGet, values, "compact JWS format must have three parts")
			})
			t.Run("method=post", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				acceptLoginAs(t, reg, adminApi, "aeneas-rekkas")
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, createSampleClient(t, reg, customPostLogoutURL), reg)
				values := url.Values{"state": {"1234"}, "post_logout_redirect_uri": {customPostLogoutURL}, "id_token_hint": {"i am not valid"}}
				logoutAndExpectErrorPage(t, publicTS.URL, browser, http.MethodPost, values, "compact JWS format must have three parts")
			})
		})

		for _, tc := range []struct {
			d                  string
			claims             jwtgo.MapClaims
			useRealIssuer      bool
			expectedErrMessage string
		}{
			{
				d: "should fail rp-initiated flow because id token hint is missing issuer",
				claims: jwtgo.MapClaims{
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				},
				expectedErrMessage: "Logout failed because issuer claim value &#39;&#39; from query parameter id_token_hint does not match with issuer value from configuration",
			},
			{
				d: "should fail rp-initiated flow because id token hint is using wrong issuer",
				claims: jwtgo.MapClaims{
					"iss": "some-issuer",
					"iat": time.Now().Add(-time.Hour * 2).Unix(),
				},
				expectedErrMessage: "Logout failed because issuer claim value &#39;some-issuer&#39; from query parameter id_token_hint does not match with issuer value from configuration",
			},
			{
				d: "should fail rp-initiated flow because iat is in the future",
				claims: jwtgo.MapClaims{
					"iat": time.Now().Add(time.Hour * 2).Unix(),
				},
				useRealIssuer:      true,
				expectedErrMessage: "Token used before issued",
			},
		} {
			t.Run("case="+tc.d, func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				c := createSampleClient(t, reg, customPostLogoutURL)
				sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
				browser := createBrowserWithSession(t, c, reg)

				wg := newWg(1)
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, nil)
				// Clone to avoid mutating the shared table entry across iterations.
				claims := maps.Clone(tc.claims)
				if tc.useRealIssuer {
					claims["iss"] = reg.Config().IssuerURL(t.Context()).String()
				}
				claims["sub"] = subject
				claims["sid"] = <-sid
				claims["aud"] = c.GetID()
				claims["exp"] = time.Now().Add(-time.Hour).Unix()

				logoutAndExpectErrorPage(t, publicTS.URL, browser, http.MethodGet, url.Values{
					"state":                    {"1234"},
					"post_logout_redirect_uri": {customPostLogoutURL},
					"id_token_hint":            {testhelpers.NewIDTokenWithClaims(t, reg, claims)},
				}, tc.expectedErrMessage)

				wg.Done()
			})
		}
	})

	t.Run("case=should fail because post-logout url is not registered", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		c := createSampleClient(t, reg, customPostLogoutURL)
		acceptLoginAs(t, reg, adminApi, subject)

		browser := createBrowserWithSession(t, c, reg)
		values := url.Values{
			"state":                    {"1234"},
			"post_logout_redirect_uri": {"https://this-is-not-a-valid-redirect-url/custom"},
			"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
				"aud": c.GetID(),
				"iss": reg.Config().IssuerURL(t.Context()).String(),
				"sub": subject,
				"sid": "logout-session-temp4",
				"exp": time.Now().Add(-time.Hour).Unix(),
				"iat": time.Now().Add(-time.Hour * 2).Unix(),
			})},
		}

		logoutAndExpectErrorPage(t, publicTS.URL, browser, http.MethodGet, values, "Logout failed because query parameter post_logout_redirect_uri is not a whitelisted as a post_logout_redirect_uri for the client")
	})

	t.Run("case=should pass rp-initiated flows", func(t *testing.T) {
		t.Parallel()
		t.Run("case=should pass even if expiry is in the past", func(t *testing.T) {
			t.Parallel()
			// formerly: should pass rp-initiated even when expiry is in the past
			t.Run("method=GET", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				c := createSampleClient(t, reg, customPostLogoutURL)
				sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, c, reg)
				body, res := makeLogoutRequest(t, publicTS.URL, browser, "GET", url.Values{
					"state":                    {"1234"},
					"post_logout_redirect_uri": {customPostLogoutURL},
					"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
						"iss": reg.Config().IssuerURL(t.Context()).String(),
						"aud": c.GetID(), "sid": <-sid, "sub": subject,
						"exp": time.Now().Add(-time.Hour).Unix(),
						"iat": time.Now().Add(-time.Hour).Unix(),
					})},
				})
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.Contains(t, body, "redirected to default server1234logged-out/custom")
				assert.Contains(t, res.Request.URL.String(), "/logged-out/custom?state=1234")
			})
			t.Run("method=POST", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				c := createSampleClient(t, reg, customPostLogoutURL)
				sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, c, reg)
				body, res := makeLogoutRequest(t, publicTS.URL, browser, "POST", url.Values{
					"state":                    {"1234"},
					"post_logout_redirect_uri": {customPostLogoutURL},
					"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
						"iss": reg.Config().IssuerURL(t.Context()).String(),
						"aud": c.GetID(), "sid": <-sid, "sub": subject,
						"exp": time.Now().Add(-time.Hour).Unix(),
						"iat": time.Now().Add(-time.Hour).Unix(),
					})},
				})
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.Contains(t, body, "redirected to default server1234logged-out/custom")
				assert.Contains(t, res.Request.URL.String(), "/logged-out/custom?state=1234")
			})
		})

		t.Run("case=should pass even if audience is an array not a string", func(t *testing.T) {
			t.Parallel()
			// formerly: should pass rp-initiated flow
			t.Run("method=GET", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				c := createSampleClient(t, reg, customPostLogoutURL)
				sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, c, reg)
				body, res := makeLogoutRequest(t, publicTS.URL, browser, "GET", url.Values{
					"state":                    {"1234"},
					"post_logout_redirect_uri": {customPostLogoutURL},
					"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
						"iss": reg.Config().IssuerURL(t.Context()).String(),
						"aud": []string{c.GetID()}, "sid": <-sid, "sub": subject,
						"exp": time.Now().Add(time.Hour).Unix(),
						"iat": time.Now().Add(-time.Hour).Unix(),
					})},
				})
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.Contains(t, body, "redirected to default server1234logged-out/custom")
				assert.Contains(t, res.Request.URL.String(), "/logged-out/custom?state=1234")
			})
			t.Run("method=POST", func(t *testing.T) {
				t.Parallel()
				_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
				c := createSampleClient(t, reg, customPostLogoutURL)
				sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
				setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
				browser := createBrowserWithSession(t, c, reg)
				body, res := makeLogoutRequest(t, publicTS.URL, browser, "POST", url.Values{
					"state":                    {"1234"},
					"post_logout_redirect_uri": {customPostLogoutURL},
					"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
						"iss": reg.Config().IssuerURL(t.Context()).String(),
						"aud": []string{c.GetID()}, "sid": <-sid, "sub": subject,
						"exp": time.Now().Add(time.Hour).Unix(),
						"iat": time.Now().Add(-time.Hour).Unix(),
					})},
				})
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.Contains(t, body, "redirected to default server1234logged-out/custom")
				assert.Contains(t, res.Request.URL.String(), "/logged-out/custom?state=1234")
			})
		})
	})

	t.Run("case=should pass rp-initiated flow without any action because SID is unknown", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		c := createSampleClient(t, reg, customPostLogoutURL)
		acceptLoginAs(t, reg, adminApi, subject)

		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, func(t *testing.T, _ *hydra.OAuth2LogoutRequest, _ error) {
			t.Fatalf("Logout should not have been called")
		})
		browser := createBrowserWithSession(t, c, reg)

		logoutAndExpectPostLogoutPage(t, publicTS.URL, browser, "GET", url.Values{
			"state":                    {"1234"},
			"post_logout_redirect_uri": {customPostLogoutURL},
			"id_token_hint": {genIDToken(t, reg, jwtgo.MapClaims{
				"aud": []string{c.GetID()}, // make sure this works with string slices too
				"iss": reg.Config().IssuerURL(t.Context()).String(),
				"sub": subject,
				"sid": "i-do-not-exist",
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Add(-time.Hour).Unix(),
			})},
		}, defaultRedirectedMessage+"1234logged-out/custom")
	})

	t.Run("case=should not append a state param if no state was passed to logout server", func(t *testing.T) {
		t.Parallel()
		_, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		c := createSampleClient(t, reg, customPostLogoutURL)
		sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)

		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, nil, nil)
		browser := createBrowserWithSession(t, c, reg)

		body, res := makeLogoutRequest(t, publicTS.URL, browser, "GET", url.Values{
			"post_logout_redirect_uri": {customPostLogoutURL},
			"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
				"iss": reg.Config().IssuerURL(t.Context()).String(),
				"aud": c.GetID(),
				"sid": <-sid,
				"sub": subject,
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Add(-time.Hour).Unix(),
			})},
		})

		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "redirected to default serverlogged-out/custom")
		assert.Contains(t, res.Request.URL.String(), "/logged-out/custom")
		assert.NotContains(t, res.Request.URL.String(), "state=1234")
	})

	t.Run("case=should return to default post logout because session was revoked in browser context", func(t *testing.T) {
		t.Parallel()
		fakeKratos, reg, publicTS, _, adminApi := makeDeps(t, defaultLogoutURL)
		sid := acceptLoginAsAndWatchSid(t, reg, adminApi, subject)
		var SID string
		wg := newWg(4)
		fakeKratos.DisableSessionCB = wg.Done
		setupCheckAndAcceptLogoutHandler(t, reg, adminApi, wg, nil)
		var bcURLCalled atomic.Bool
		c := createClientWithBackchannelLogout(t, reg, customPostLogoutURL, wg, func(t *testing.T, logoutToken gjson.Result) {
			assert.Equal(t, SID, logoutToken.Get("sid").String(), logoutToken.Raw)
			assert.False(t, bcURLCalled.Swap(true))
		})

		browser := createBrowserWithSession(t, c, reg)
		SID = <-sid

		// Use another browser (without session cookie) to make the logout request:
		otherBrowser := &http.Client{
			Jar: testhelpers.NewEmptyCookieJar(t),
		}
		// Capture the request with the logout verifier to test reuse detection later.
		var verifiedLogoutReq *http.Request
		otherBrowser.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if lv := req.FormValue("logout_verifier"); lv != "" {
				verifiedLogoutReq = req.Clone(t.Context())
			}
			return nil
		}
		logoutAndExpectPostLogoutPage(t, publicTS.URL, otherBrowser, "GET", url.Values{
			"post_logout_redirect_uri": {customPostLogoutURL},
			"id_token_hint": {testhelpers.NewIDTokenWithClaims(t, reg, jwtgo.MapClaims{
				"iss": reg.Config().IssuerURL(t.Context()).String(),
				"aud": c.GetID(),
				"sid": SID,
				"sub": subject,
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Add(-time.Hour).Unix(),
			})},
		}, "redirected to default serverlogged-out/custom") // this means RP-initiated flow worked!

		// Set up login / consent and check if skip is set to false (because logout happened), but use
		// the original login browser which still has the session.
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminApi, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				defer wg.Done()
				require.NoError(t, err)
				assert.False(t, res.Skip)
				return hydra.AcceptOAuth2LoginRequest{Remember: new(true)}
			}),
			checkAndAcceptConsentHandler(t, adminApi, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				return hydra.AcceptOAuth2ConsentRequest{Remember: new(true)}
			}))

		// Make an oauth 2 request to trigger the login check.
		_, res := makeOAuth2Request(t, reg, browser, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("code"))

		wg.Wait()

		assert.True(t, fakeKratos.DisableSessionWasCalled)
		assert.Equal(t, fakeKratos.LastDisabledSession, kratos.FakeSessionID)

		wg.Add(1) // in case the backchannel logout callback is incorrectly called again, we want to avoid a panic in a different goroutine
		res, err := otherBrowser.Do(verifiedLogoutReq)
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		require.NotPanics(t, wg.Done, "The backchannel logout callback should not have been called again after the session was already disabled.")
		wg.Wait()
		require.Equal(t, 200, res.StatusCode)

		wg.Add(1) // in case the backchannel logout callback is incorrectly called again, we want to avoid a panic in a different goroutine
		res, err = browser.Do(verifiedLogoutReq)
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		require.NotPanics(t, wg.Done, "The backchannel logout callback should not have been called again after the session was already disabled.")
		wg.Wait()
		require.Equal(t, 200, res.StatusCode)
	})

	t.Run("case=should execute backchannel logout in headless flow with sid", func(t *testing.T) {
		t.Parallel()
		fakeKratos, reg, _, adminTS, adminApi := makeDeps(t, defaultLogoutURL)
		numSidConsumers := 2
		sid := make(chan string, numSidConsumers)
		acceptLoginAsAndWatchSidForConsumers(t, reg, adminApi, subject, sid, true, numSidConsumers)

		backChannelWG := newWg(2)
		fakeKratos.DisableSessionCB = backChannelWG.Done

		c := createClientWithBackchannelLogout(t, reg, customPostLogoutURL, backChannelWG, func(t *testing.T, logoutToken gjson.Result) {
			assert.EqualValues(t, <-sid, logoutToken.Get("sid").String(), logoutToken.Raw)
			assert.Empty(t, logoutToken.Get("sub").String(), logoutToken.Raw) // The sub claim should be empty because it doesn't work with forced obfuscation and thus we can't easily recover it.
			assert.Empty(t, logoutToken.Get("nonce").String(), logoutToken.Raw)
		})

		logoutViaHeadlessAndExpectNoContent(t, adminTS.URL, createBrowserWithSession(t, c, reg), url.Values{"sid": {<-sid}})

		backChannelWG.Wait() // we want to ensure that all back channels have been called!
		assert.True(t, fakeKratos.DisableSessionWasCalled)
		assert.Equal(t, fakeKratos.LastDisabledSession, kratos.FakeSessionID)
	})

	t.Run("case=should logout in headless flow with non-existing sid", func(t *testing.T) {
		t.Parallel()
		fakeKratos, _, _, adminTS, _ := makeDeps(t, defaultLogoutURL)
		logoutViaHeadlessAndExpectNoContent(t, adminTS.URL, new(http.Client), url.Values{"sid": {"non-existing-sid"}})
		assert.False(t, fakeKratos.DisableSessionWasCalled)
	})

	t.Run("case=should logout in headless flow with session that has remember=false", func(t *testing.T) {
		t.Parallel()
		fakeKratos, reg, _, adminTS, adminApi := makeDeps(t, defaultLogoutURL)
		sid := make(chan string, 1)
		acceptLoginAsAndWatchSidForConsumers(t, reg, adminApi, subject, sid, false, 1)

		wg := newWg(1)
		fakeKratos.DisableSessionCB = wg.Done

		c := createSampleClient(t, reg, customPostLogoutURL)

		logoutViaHeadlessAndExpectNoContent(t, adminTS.URL, createBrowserWithSession(t, c, reg), url.Values{"sid": {<-sid}})
		wg.Wait()
		assert.True(t, fakeKratos.DisableSessionWasCalled)
		assert.Equal(t, fakeKratos.LastDisabledSession, kratos.FakeSessionID)
	})

	t.Run("case=should fail headless logout because neither sid nor subject were provided", func(t *testing.T) {
		t.Parallel()
		fakeKratos, _, _, adminTS, _ := makeDeps(t, defaultLogoutURL)
		logoutViaHeadlessAndExpectError(t, adminTS.URL, new(http.Client), url.Values{}, `Either 'subject' or 'sid' query parameters need to be defined.`)
		assert.False(t, fakeKratos.DisableSessionWasCalled)
	})
}

func getSessionCookie(t *testing.T, browser *http.Client, publicTSURL string) *http.Cookie {
	u, err := url.Parse(publicTSURL)
	require.NoError(t, err)

	cookies := browser.Jar.Cookies(u)
	idx := slices.IndexFunc(cookies, func(c *http.Cookie) bool {
		return strings.HasPrefix(c.Name, "ory_hydra_session")
	})
	if idx >= 0 {
		return cookies[idx]
	}
	return nil
}

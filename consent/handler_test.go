// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/v2/internal"

	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	. "github.com/ory/hydra/v2/consent"
)

func TestGetLogoutRequest(t *testing.T) {
	for k, tc := range []struct {
		exists  bool
		handled bool
		status  int
	}{
		{false, false, http.StatusNotFound},
		{true, false, http.StatusOK},
		{true, true, http.StatusGone},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := internal.NewConfigurationWithDefaults()
			reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{LegacyClientID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				require.NoError(t, reg.ConsentManager().CreateLogoutRequest(context.TODO(), &LogoutRequest{
					Client:     cl,
					ID:         challenge,
					WasHandled: tc.handled,
					RequestURL: requestURL,
				}))
			}

			h := NewHandler(reg, conf)
			r := x.NewRouterAdmin(conf.AdminURL)
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + "/admin" + LogoutPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result LogoutRequest
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, challenge, result.ID)
				require.Equal(t, requestURL, result.RequestURL)
			}
		})
	}
}

func TestGetLoginRequest(t *testing.T) {
	for k, tc := range []struct {
		exists  bool
		handled bool
		status  int
	}{
		{false, false, http.StatusNotFound},
		{true, false, http.StatusOK},
		{true, true, http.StatusGone},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := internal.NewConfigurationWithDefaults()
			reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{LegacyClientID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), &LoginRequest{
					Client:     cl,
					ID:         challenge,
					RequestURL: requestURL,
				}))

				if tc.handled {
					_, err := reg.ConsentManager().HandleLoginRequest(context.Background(), challenge, &HandledLoginRequest{ID: challenge, WasHandled: true})
					require.NoError(t, err)
				}
			}

			h := NewHandler(reg, conf)
			r := x.NewRouterAdmin(conf.AdminURL)
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + "/admin" + LoginPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result LoginRequest
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, challenge, result.ID)
				require.Equal(t, requestURL, result.RequestURL)
				require.NotNil(t, result.Client)
			}
		})
	}
}

func TestGetConsentRequest(t *testing.T) {
	for k, tc := range []struct {
		exists  bool
		handled bool
		status  int
	}{
		{false, false, http.StatusNotFound},
		{true, false, http.StatusOK},
		{true, true, http.StatusGone},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := internal.NewConfigurationWithDefaults()
			reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{LegacyClientID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				lr := &LoginRequest{ID: "login-" + challenge, Client: cl, RequestURL: requestURL}
				require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), lr))
				_, err := reg.ConsentManager().HandleLoginRequest(context.Background(), lr.ID, &HandledLoginRequest{
					ID: lr.ID,
				})
				require.NoError(t, err)
				require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), &OAuth2ConsentRequest{
					Client:         cl,
					ID:             challenge,
					Verifier:       challenge,
					CSRF:           challenge,
					LoginChallenge: sqlxx.NullString(lr.ID),
				}))

				if tc.handled {
					_, err := reg.ConsentManager().HandleConsentRequest(context.Background(), &AcceptOAuth2ConsentRequest{
						ID:         challenge,
						WasHandled: true,
						HandledAt:  sqlxx.NullTime(time.Now()),
					})
					require.NoError(t, err)
				}
			}

			h := NewHandler(reg, conf)

			r := x.NewRouterAdmin(conf.AdminURL)
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + "/admin" + ConsentPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result OAuth2ConsentRequest
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, challenge, result.ID)
				require.Equal(t, requestURL, result.RequestURL)
				require.NotNil(t, result.Client)
			}
		})
	}
}

func TestExtendConsentRequest(t *testing.T) {
	t.Run("case=extend consent expiry time", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		h := NewHandler(reg, conf)
		r := x.NewRouterAdmin(conf.AdminURL)
		h.SetRoutes(r)
		ts := httptest.NewServer(r)
		defer ts.Close()

		c := &http.Client{}
		cl := &client.Client{LegacyClientID: "client-1"}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))

		var initialRememberFor time.Duration = 300
		var remainingValidTime time.Duration = 100

		require.NoError(t, reg.ConsentManager().CreateLoginSession(context.Background(), &LoginSession{
			ID:      makeID("fk-login-session", "1", "1"),
			Subject: "subject-1",
		}))
		requestedTimeInPast := time.Now().UTC().Add(-(initialRememberFor - remainingValidTime) * time.Second)
		require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), &LoginRequest{
			ID:          makeID("challenge", "1", "1"),
			SessionID:   sqlxx.NullString(makeID("fk-login-session", "1", "1")),
			Client:      cl,
			Subject:     "subject-1",
			RequestedAt: requestedTimeInPast,
		}))
		require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), &OAuth2ConsentRequest{
			ID:             makeID("challenge", "1", "1"),
			Subject:        "subject-1",
			Client:         cl,
			LoginSessionID: sqlxx.NullString(makeID("fk-login-session", "1", "1")),
			LoginChallenge: sqlxx.NullString(makeID("challenge", "1", "1")),
			Verifier:       makeID("verifier", "1", "1"),
			CSRF:           "csrf1",
			Skip:           false,
			ACR:            "1",
		}))
		_, err := reg.ConsentManager().HandleConsentRequest(context.Background(), &AcceptOAuth2ConsentRequest{
			ID:          makeID("challenge", "1", "1"),
			Remember:    true,
			RememberFor: int(initialRememberFor),
			WasHandled:  true,
			HandledAt:   sqlxx.NullTime(time.Now().UTC()),
		})
		require.NoError(t, err)

		require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), &LoginRequest{
			ID:          makeID("challenge", "1", "2"),
			SessionID:   sqlxx.NullString(makeID("fk-login-session", "1", "1")),
			Verifier:    makeID("verifier", "1", "1"),
			Client:      cl,
			RequestedAt: time.Now().UTC(),
			Subject:     "subject-1",
		}))
		require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), &OAuth2ConsentRequest{
			ID:             makeID("challenge", "1", "2"),
			Subject:        "subject-1",
			Client:         cl,
			LoginSessionID: sqlxx.NullString(makeID("fk-login-session", "1", "1")),
			LoginChallenge: sqlxx.NullString(makeID("challenge", "1", "2")),
			Verifier:       makeID("verifier", "1", "2"),
			CSRF:           "csrf2",
			Skip:           true,
		}))

		var b bytes.Buffer
		var extendRememberFor time.Duration = 300
		require.NoError(t, json.NewEncoder(&b).Encode(&AcceptOAuth2ConsentRequest{
			Remember:    true,
			RememberFor: int(extendRememberFor),
		}))

		req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+ConsentPath+"/accept?challenge=challenge-1-2", &b)
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		require.EqualValues(t, 200, resp.StatusCode)

		crs, err := reg.ConsentManager().FindSubjectsGrantedConsentRequests(context.Background(), "subject-1", 100, 0)
		require.NoError(t, err)
		require.NotNil(t, crs)
		require.EqualValues(t, 1, len(crs))
		expectedRememberFor := int(initialRememberFor + extendRememberFor - remainingValidTime)
		cr := crs[0]
		require.EqualValues(t, "challenge-1-1", cr.ID)
		require.InDelta(t, expectedRememberFor, cr.RememberFor, 1)
	})
}

func TestGetLoginRequestWithDuplicateAccept(t *testing.T) {
	t.Run("Test get login request with duplicate accept", func(t *testing.T) {
		challenge := "challenge"
		requestURL := "http://192.0.2.1"

		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

		cl := &client.Client{LegacyClientID: "client"}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
		require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), &LoginRequest{
			Client:     cl,
			ID:         challenge,
			RequestURL: requestURL,
		}))

		h := NewHandler(reg, conf)
		r := x.NewRouterAdmin(conf.AdminURL)
		h.SetRoutes(r)
		ts := httptest.NewServer(r)
		defer ts.Close()

		c := &http.Client{}

		sub := "sub123"
		acceptLogin := &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true), Subject: sub}

		// marshal User to json
		acceptLoginJson, err := json.Marshal(acceptLogin)
		if err != nil {
			panic(err)
		}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+LoginPath+"/accept?challenge="+challenge, bytes.NewBuffer(acceptLoginJson))
		if err != nil {
			panic(err)
		}

		resp, err := c.Do(req)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.RedirectTo)
		require.Contains(t, result.RedirectTo, "login_verifier")

		req2, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+LoginPath+"/accept?challenge="+challenge, bytes.NewBuffer(acceptLoginJson))
		if err != nil {
			panic(err)
		}

		resp2, err := c.Do(req2)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp2.StatusCode)

		var result2 OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))
		require.NotNil(t, result2.RedirectTo)
		require.Contains(t, result2.RedirectTo, "login_verifier")
	})
}

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
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/driver"

	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
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

func TestRevokeConsentSession(t *testing.T) {
	newWg := func(add int) *sync.WaitGroup {
		var wg sync.WaitGroup
		wg.Add(add)
		return &wg
	}

	t.Run("case=subject=subject-1,client=client-1,session=session-1,trigger_back_channel_logout=true", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(1)
		cl := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{"login-session-1"}, backChannelWG)
		performLoginFlow(t, reg, "1", cl)
		performLoginFlow(t, reg, "2", cl)
		performDeleteConsentSession(t, reg, "client-1", "login-session-1", true)
		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.NoError(t, err)
		require.NotNil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,client=client-1,session=session-1,trigger_back_channel_logout=false", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(0)
		cl := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{}, backChannelWG)
		performLoginFlow(t, reg, "1", cl)
		performLoginFlow(t, reg, "2", cl)
		performDeleteConsentSession(t, reg, "client-1", "login-session-1", false)
		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.NoError(t, err)
		require.NotNil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,client=client-1,trigger_back_channel_logout=true", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(2)
		cl := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{"login-session-1", "login-session-2"}, backChannelWG)
		performLoginFlow(t, reg, "1", cl)
		performLoginFlow(t, reg, "2", cl)

		performDeleteConsentSession(t, reg, "client-1", nil, true)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,client=client-1,trigger_back_channel_logout=false", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(0)
		cl := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{}, backChannelWG)
		performLoginFlow(t, reg, "1", cl)
		performLoginFlow(t, reg, "2", cl)

		performDeleteConsentSession(t, reg, "client-1", nil, false)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,all=true,session=session-1,trigger_back_channel_logout=true", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(1)
		cl1 := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{"login-session-1"}, backChannelWG)
		cl2 := createClientWithBackChannelEndpoint(t, reg, "client-2", []string{}, backChannelWG)
		performLoginFlow(t, reg, "1", cl1)
		performLoginFlow(t, reg, "2", cl2)

		performDeleteConsentSession(t, reg, nil, "login-session-1", true)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.NoError(t, err)
		require.NotNil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,all=true,session=session-1,trigger_back_channel_logout=false", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(0)
		cl1 := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{}, backChannelWG)
		cl2 := createClientWithBackChannelEndpoint(t, reg, "client-2", []string{}, backChannelWG)
		performLoginFlow(t, reg, "1", cl1)
		performLoginFlow(t, reg, "2", cl2)

		performDeleteConsentSession(t, reg, nil, "login-session-1", false)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.NoError(t, err)
		require.NotNil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,all=true,trigger_back_channel_logout=true", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(2)
		cl1 := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{"login-session-1"}, backChannelWG)
		cl2 := createClientWithBackChannelEndpoint(t, reg, "client-2", []string{"login-session-2"}, backChannelWG)
		performLoginFlow(t, reg, "1", cl1)
		performLoginFlow(t, reg, "2", cl2)

		performDeleteConsentSession(t, reg, nil, nil, true)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c2)
		backChannelWG.Wait()
	})

	t.Run("case=subject=subject-1,all=true,trigger_back_channel_logout=false", func(t *testing.T) {
		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})
		backChannelWG := newWg(0)
		cl1 := createClientWithBackChannelEndpoint(t, reg, "client-1", []string{}, backChannelWG)
		cl2 := createClientWithBackChannelEndpoint(t, reg, "client-2", []string{}, backChannelWG)
		performLoginFlow(t, reg, "1", cl1)
		performLoginFlow(t, reg, "2", cl2)

		performDeleteConsentSession(t, reg, nil, nil, false)

		c1, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-1")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c1)
		c2, err := reg.ConsentManager().GetConsentRequest(context.Background(), "consent-challenge-2")
		require.Error(t, x.ErrNotFound, err)
		require.Nil(t, c2)
		backChannelWG.Wait()
	})
}

func performDeleteConsentSession(t *testing.T, reg driver.Registry, client, loginSessionId interface{}, triggerBackChannelLogout bool) {
	conf := internal.NewConfigurationWithDefaults()
	h := NewHandler(reg, conf)
	r := x.NewRouterAdmin(conf.AdminURL)
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()
	c := &http.Client{}

	u, _ := url.Parse(ts.URL + "/admin" + SessionsPath + "/consent")
	q := u.Query()
	q.Set("subject", "subject-1")
	if client != nil && len(client.(string)) != 0 {
		q.Set("client", client.(string))
	} else {
		q.Set("all", "true")
	}
	if loginSessionId != nil && len(loginSessionId.(string)) != 0 {
		q.Set("login_session_id", loginSessionId.(string))
	}
	if triggerBackChannelLogout {
		q.Set("trigger_back_channel_logout", "true")
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)

	require.NoError(t, err)
	_, err = c.Do(req)
	require.NoError(t, err)
}

func performLoginFlow(t *testing.T, reg driver.Registry, flowId string, cl *client.Client) {
	subject := "subject-1"
	loginSessionId := "login-session-" + flowId
	loginChallenge := "login-challenge-" + flowId
	consentChallenge := "consent-challenge-" + flowId
	requestURL := "http://192.0.2.1"

	ls := &LoginSession{
		ID:      loginSessionId,
		Subject: subject,
	}
	lr := &LoginRequest{
		ID:         loginChallenge,
		Subject:    subject,
		Client:     cl,
		RequestURL: requestURL,
		Verifier:   "login-verifier-" + flowId,
		SessionID:  sqlxx.NullString(loginSessionId),
	}
	cr := &OAuth2ConsentRequest{
		Client:         cl,
		ID:             consentChallenge,
		Verifier:       consentChallenge,
		CSRF:           consentChallenge,
		Subject:        subject,
		LoginChallenge: sqlxx.NullString(loginChallenge),
		LoginSessionID: sqlxx.NullString(loginSessionId),
	}
	hcr := &AcceptOAuth2ConsentRequest{
		ConsentRequest: cr,
		ID:             consentChallenge,
		WasHandled:     true,
		HandledAt:      sqlxx.NullTime(time.Now().UTC()),
	}

	require.NoError(t, reg.ConsentManager().CreateLoginSession(context.Background(), ls))
	require.NoError(t, reg.ConsentManager().CreateLoginRequest(context.Background(), lr))
	require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), cr))
	_, err := reg.ConsentManager().HandleConsentRequest(context.Background(), hcr)
	require.NoError(t, err)
}

func createClientWithBackChannelEndpoint(t *testing.T, reg driver.Registry, clientId string, expectedBackChannelLogoutFlowIds []string, wg *sync.WaitGroup) *client.Client {
	return func(t *testing.T, key string, wg *sync.WaitGroup, cb func(t *testing.T, logoutToken gjson.Result)) *client.Client {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer wg.Done()
			require.NoError(t, r.ParseForm())
			lt := r.PostFormValue("logout_token")
			assert.NotEmpty(t, lt)
			token, err := reg.OpenIDJWTStrategy().Decode(r.Context(), lt)
			require.NoError(t, err)
			var b bytes.Buffer
			require.NoError(t, json.NewEncoder(&b).Encode(token.Claims))
			cb(t, gjson.Parse(b.String()))
		}))
		t.Cleanup(server.Close)
		c := &client.Client{
			LegacyClientID:       clientId,
			BackChannelLogoutURI: server.URL,
		}
		err := reg.ClientManager().CreateClient(context.Background(), c)
		require.NoError(t, err)
		return c
	}(t, clientId, wg, func(t *testing.T, logoutToken gjson.Result) {
		sid := logoutToken.Get("sid").String()
		assert.Contains(t, expectedBackChannelLogoutFlowIds, sid)
		for i, v := range expectedBackChannelLogoutFlowIds {
			if v == sid {
				expectedBackChannelLogoutFlowIds = append(expectedBackChannelLogoutFlowIds[:i], expectedBackChannelLogoutFlowIds[i+1:]...)
				break
			}
		}
	})
}

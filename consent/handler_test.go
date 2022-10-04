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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go"
	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
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
				var result HandledOAuth2ConsentRequest
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
				var result HandledOAuth2ConsentRequest
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
				_, err := reg.ConsentManager().HandleLoginRequest(context.Background(), lr.ID, &consent.HandledLoginRequest{
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
				var result HandledOAuth2ConsentRequest
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

		var result RequestHandlerResponse
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

		var result2 RequestHandlerResponse
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))
		require.NotNil(t, result2.RedirectTo)
		require.Contains(t, result2.RedirectTo, "login_verifier")
	})
}

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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/pborman/uuid"

	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/require"

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
			reg := internal.NewRegistryMemory(t, conf)

			if tc.exists {
				cl := &client.Client{OutfacingID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				require.NoError(t, reg.ConsentManager().CreateLogoutRequest(context.TODO(), &LogoutRequest{
					Client:     cl,
					ID:         challenge,
					WasHandled: tc.handled,
					RequestURL: requestURL,
				}))
			}

			h := NewHandler(reg, conf)
			r := x.NewRouterAdmin()
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + LogoutPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result RequestWasHandledResponse
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
			reg := internal.NewRegistryMemory(t, conf)

			if tc.exists {
				cl := &client.Client{OutfacingID: "client" + key}
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
			r := x.NewRouterAdmin()
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + LoginPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result RequestWasHandledResponse
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
			reg := internal.NewRegistryMemory(t, conf)

			if tc.exists {
				cl := &client.Client{OutfacingID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				require.NoError(t, reg.ConsentManager().CreateConsentRequest(context.Background(), &ConsentRequest{
					Client:     cl,
					ID:         challenge,
					RequestURL: requestURL,
				}))

				if tc.handled {
					_, err := reg.ConsentManager().HandleConsentRequest(context.Background(), challenge, &HandledConsentRequest{
						ID:         challenge,
						WasHandled: true,
					})
					require.NoError(t, err)
				}
			}

			h := NewHandler(reg, conf)

			r := x.NewRouterAdmin()
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			c := &http.Client{}
			resp, err := c.Get(ts.URL + ConsentPath + "?challenge=" + challenge)
			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)

			if tc.handled {
				var result RequestWasHandledResponse
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result ConsentRequest
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, challenge, result.ID)
				require.Equal(t, requestURL, result.RequestURL)
				require.NotNil(t, result.Client)
			}
		})
	}
}

func TestGetDeviceLoginRequest(t *testing.T) {
	for k, tc := range []struct {
		createUserSession   bool
		createDeviceSession bool
		handled             bool
		status              int
		user_code           string
		device_challenge    string
	}{
		{
			createUserSession:   false,
			createDeviceSession: false,
			handled:             false,
			status:              http.StatusBadRequest,
			user_code:           "",
			device_challenge:    "",
		},
		{
			createUserSession:   false,
			createDeviceSession: false,
			handled:             false,
			status:              http.StatusNotFound,
			user_code:           "AAABBBCCC",
			device_challenge:    "muyjbkdhjsbvc8",
		},
		{
			createUserSession:   true,
			createDeviceSession: false,
			handled:             false,
			status:              http.StatusNotFound,
			user_code:           "AAABBBCCC",
			device_challenge:    "muyjbkdhjsbvc8",
		},
		{
			createUserSession:   true,
			createDeviceSession: true,
			handled:             false,
			status:              http.StatusFound,
			user_code:           "AAABBBCCC",
			device_challenge:    "muyjbkdhjsbvc8",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {

			conf := internal.NewConfigurationWithDefaults()
			reg := internal.NewRegistryMemory(t, conf)

			h := NewHandler(reg, conf)
			r := x.NewRouterAdmin()
			h.SetRoutes(r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			cl := &client.Client{OutfacingID: "test"}
			reg.ClientManager().CreateClient(context.Background(), cl)

			params := "?state=abc12345"
			if tc.device_challenge != "" {
				params = params + "&device_challenge=" + tc.device_challenge
			}

			if tc.user_code != "" {

				verifier := strings.Replace(uuid.New(), "-", "", -1)
				csrf := strings.Replace(uuid.New(), "-", "", -1)

				if tc.createDeviceSession {
					reg.ConsentManager().CreateDeviceGrantRequest(context.TODO(), &DeviceGrantRequest{
						ID:       tc.device_challenge,
						Verifier: verifier,
						CSRF:     csrf,
					})
				}

				userCodeHash := reg.OAuth2HMACStrategy().UserCodeSignature(tc.user_code)
				deviceCodeHash := reg.OAuth2HMACStrategy().DeviceCodeSignature("AAABBBCCCDDD")

				req := &fosite.AccessRequest{
					GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
					Request: fosite.Request{
						Client:      cl,
						Session:     &fosite.DefaultSession{Subject: "A"},
						RequestedAt: time.Now().UTC(),
						Form:        url.Values{"device_code": {"ABC1234"}},
					},
				}
				req.SetID(deviceCodeHash)
				req.Session = &oauth2.Session{}
				if tc.createUserSession {
					reg.OAuth2Storage().CreateUserCodeSession(context.TODO(), userCodeHash, req)
					reg.OAuth2Storage().CreateDeviceCodeSession(context.TODO(), tc.device_challenge, req)
				}
				params = params + "&user_code=" + tc.user_code
			}

			req, err := http.NewRequest("GET", ts.URL+DevicePath+params, nil)
			if err != nil {
				t.Fatal(err)
			}

			transport := http.Transport{}
			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}

			require.NoError(t, err)
			require.EqualValues(t, tc.status, resp.StatusCode)
		})
	}
}

func TestGetDeviceSessionCreateDelete(t *testing.T) {
	t.Run("case=should pass creating / deleting device sessions", func(t *testing.T) {

		conf := internal.NewConfigurationWithDefaults()
		reg := internal.NewRegistryMemory(t, conf)

		cl := &client.Client{OutfacingID: "test"}
		reg.ClientManager().CreateClient(context.Background(), cl)

		userCodeHash := reg.OAuth2HMACStrategy().UserCodeSignature("ABCD12345")
		deviceCodeHash := reg.OAuth2HMACStrategy().DeviceCodeSignature("AAABBB.CCCDDD")

		req := &fosite.AccessRequest{
			GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
			Request: fosite.Request{
				Client:      cl,
				Session:     &fosite.DefaultSession{Subject: "A"},
				RequestedAt: time.Now().UTC(),
				Form:        url.Values{"device_code": {"ABC1234"}},
			},
		}
		req.SetID(deviceCodeHash)
		req.Session = &oauth2.Session{}
		require.NoError(t, reg.OAuth2Storage().CreateUserCodeSession(context.TODO(), userCodeHash, req))
		require.NoError(t, reg.OAuth2Storage().CreateDeviceCodeSession(context.TODO(), deviceCodeHash, req))

		_, err := reg.OAuth2Storage().GetDeviceCodeSession(context.TODO(), deviceCodeHash, req.Session.Clone())

		require.NoError(t, err)

		require.NoError(t, reg.OAuth2Storage().DeleteUserCodeSession(context.TODO(), userCodeHash))
		require.NoError(t, reg.OAuth2Storage().DeleteDeviceCodeSession(context.TODO(), deviceCodeHash))
	})
}

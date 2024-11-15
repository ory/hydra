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

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	. "github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlxx"
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
			ctx := context.Background()
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := testhelpers.NewConfigurationWithDefaults()
			reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{ID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
				require.NoError(t, reg.ConsentManager().CreateLogoutRequest(context.TODO(), &flow.LogoutRequest{
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
				var result flow.OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result flow.LogoutRequest
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
		t.Run(fmt.Sprintf("exists=%v/handled=%v", tc.exists, tc.handled), func(t *testing.T) {
			ctx := context.Background()
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := testhelpers.NewConfigurationWithDefaults()
			reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{ID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cl))
				f, err := reg.ConsentManager().CreateLoginRequest(context.Background(), nil, &flow.LoginRequest{
					Client:      cl,
					ID:          challenge,
					RequestURL:  requestURL,
					RequestedAt: time.Now(),
				})
				require.NoError(t, err)
				challenge, err = f.ToLoginChallenge(ctx, reg)
				require.NoError(t, err)

				if tc.handled {
					_, err := reg.ConsentManager().HandleLoginRequest(ctx, f, challenge, &flow.HandledLoginRequest{ID: challenge, WasHandled: true})
					require.NoError(t, err)
					challenge, err = f.ToLoginChallenge(ctx, reg)
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
				var result flow.OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result flow.LoginRequest
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
			ctx := context.Background()
			key := fmt.Sprint(k)
			challenge := "challenge" + key
			requestURL := "http://192.0.2.1"

			conf := testhelpers.NewConfigurationWithDefaults()
			reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

			if tc.exists {
				cl := &client.Client{ID: "client" + key}
				require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
				lr := &flow.LoginRequest{
					ID:          "login-" + challenge,
					Client:      cl,
					RequestURL:  requestURL,
					RequestedAt: time.Now(),
				}
				f, err := reg.ConsentManager().CreateLoginRequest(ctx, nil, lr)
				require.NoError(t, err)
				challenge, err = f.ToLoginChallenge(ctx, reg)
				require.NoError(t, err)
				_, err = reg.ConsentManager().HandleLoginRequest(ctx, f, challenge, &flow.HandledLoginRequest{
					ID: challenge,
				})
				require.NoError(t, err)
				challenge, err = f.ToConsentChallenge(ctx, reg)
				require.NoError(t, err)
				require.NoError(t, reg.ConsentManager().CreateConsentRequest(ctx, f, &flow.OAuth2ConsentRequest{
					Client:         cl,
					ID:             challenge,
					Verifier:       challenge,
					CSRF:           challenge,
					LoginChallenge: sqlxx.NullString(lr.ID),
				}))

				if tc.handled {
					_, err := reg.ConsentManager().HandleConsentRequest(ctx, f, &flow.AcceptOAuth2ConsentRequest{
						ID:         challenge,
						WasHandled: true,
						HandledAt:  sqlxx.NullTime(time.Now()),
					})
					require.NoError(t, err)
					challenge, err = f.ToConsentChallenge(ctx, reg)
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
				var result flow.OAuth2RedirectTo
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, requestURL, result.RedirectTo)
			} else if tc.exists {
				var result flow.OAuth2ConsentRequest
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
		ctx := context.Background()
		challenge := "challenge"
		requestURL := "http://192.0.2.1"

		conf := testhelpers.NewConfigurationWithDefaults()
		reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

		cl := &client.Client{ID: "client"}
		require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
		f, err := reg.ConsentManager().CreateLoginRequest(ctx, nil, &flow.LoginRequest{
			Client:      cl,
			ID:          challenge,
			RequestURL:  requestURL,
			RequestedAt: time.Now(),
		})
		require.NoError(t, err)
		challenge, err = f.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)

		h := NewHandler(reg, conf)
		r := x.NewRouterAdmin(conf.AdminURL)
		h.SetRoutes(r)
		ts := httptest.NewServer(r)
		defer ts.Close()

		c := &http.Client{}

		sub := "sub123"
		acceptLogin := &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Ptr(true), Subject: sub}

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

		var result flow.OAuth2RedirectTo
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

		var result2 flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))
		require.NotNil(t, result2.RedirectTo)
		require.Contains(t, result2.RedirectTo, "login_verifier")
	})
}

func TestAcceptDeviceRequest(t *testing.T) {
	ctx := context.Background()
	challenge := "challenge"
	requestURL := "https://hydra.example.com/" + oauth2.DeviceVerificationPath

	conf := testhelpers.NewConfigurationWithDefaults()
	reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

	cl := &client.Client{ID: "client"}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
	f, err := reg.ConsentManager().CreateDeviceUserAuthRequest(ctx, &flow.DeviceUserAuthRequest{
		Client:      cl,
		ID:          challenge,
		RequestURL:  requestURL,
		RequestedAt: time.Now(),
	})
	require.NoError(t, err)
	challenge, err = f.ToDeviceChallenge(ctx, reg)
	require.NoError(t, err)

	h := NewHandler(reg, conf)
	r := x.NewRouterAdmin(conf.AdminURL)
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	c := &http.Client{}

	deviceRequest := fosite.NewDeviceRequest()
	deviceRequest.Client = cl
	deviceRequest.SetSession(
		&oauth2.Session{
			DefaultSession: &openid.DefaultSession{
				Headers: &jwt.Headers{},
			},
		},
	)
	_, deviceCodesig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(ctx)
	require.NoError(t, err)
	userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
	require.NoError(t, err)
	reg.OAuth2Storage().CreateDeviceAuthSession(ctx, deviceCodesig, sig, deviceRequest)
	require.NoError(t, err)

	acceptUserCode := &hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode}

	// marshal User to json
	acceptUserCodeJson, err := json.Marshal(acceptUserCode)
	require.NoError(t, err)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+DevicePath+"/accept?device_challenge="+challenge, bytes.NewBuffer(acceptUserCodeJson))
	require.NoError(t, err)

	resp, err := c.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)

	var result flow.OAuth2RedirectTo
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(t, result.RedirectTo)
	require.Contains(t, result.RedirectTo, requestURL)
	require.Contains(t, result.RedirectTo, "device_verifier")
}

func TestAcceptDuplicateDeviceRequest(t *testing.T) {
	ctx := context.Background()
	challenge := "challenge"
	requestURL := "https://hydra.example.com/" + oauth2.DeviceVerificationPath

	conf := testhelpers.NewConfigurationWithDefaults()
	reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

	cl := &client.Client{ID: "client"}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
	f, err := reg.ConsentManager().CreateDeviceUserAuthRequest(ctx, &flow.DeviceUserAuthRequest{
		Client:      cl,
		ID:          challenge,
		RequestURL:  requestURL,
		RequestedAt: time.Now(),
	})
	require.NoError(t, err)
	challenge, err = f.ToDeviceChallenge(ctx, reg)
	require.NoError(t, err)

	h := NewHandler(reg, conf)
	r := x.NewRouterAdmin(conf.AdminURL)
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	c := &http.Client{}

	deviceRequest := fosite.NewDeviceRequest()
	deviceRequest.Client = cl
	deviceRequest.SetSession(
		&oauth2.Session{
			DefaultSession: &openid.DefaultSession{
				Headers: &jwt.Headers{},
			},
		},
	)
	_, deviceCodesig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(ctx)
	require.NoError(t, err)
	userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
	require.NoError(t, err)
	reg.OAuth2Storage().CreateDeviceAuthSession(ctx, deviceCodesig, sig, deviceRequest)
	require.NoError(t, err)

	acceptUserCode := &hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode}

	// marshal User to json
	acceptUserCodeJson, err := json.Marshal(acceptUserCode)
	require.NoError(t, err)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+DevicePath+"/accept?device_challenge="+challenge, bytes.NewBuffer(acceptUserCodeJson))
	require.NoError(t, err)

	resp, err := c.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)

	var result flow.OAuth2RedirectTo
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(t, result.RedirectTo)
	require.Contains(t, result.RedirectTo, requestURL)
	require.Contains(t, result.RedirectTo, "device_verifier")

	req2, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+DevicePath+"/accept?device_challenge="+challenge, bytes.NewBuffer(acceptUserCodeJson))
	require.NoError(t, err)
	resp2, err := c.Do(req2)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp2.StatusCode)

	var result2 flow.OAuth2RedirectTo
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))
	require.NotNil(t, result2.RedirectTo)
	require.Contains(t, result2.RedirectTo, requestURL)
	require.Contains(t, result2.RedirectTo, "device_verifier")
}

func TestAcceptCodeDeviceRequestFailure(t *testing.T) {
	ctx := context.Background()
	challenge := "challenge"
	requestURL := "https://hydra.example.com/" + oauth2.DeviceVerificationPath

	conf := testhelpers.NewConfigurationWithDefaults()
	reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

	cl := &client.Client{ID: "client"}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, cl))
	f, err := reg.ConsentManager().CreateDeviceUserAuthRequest(ctx, &flow.DeviceUserAuthRequest{
		Client:      cl,
		ID:          challenge,
		RequestURL:  requestURL,
		RequestedAt: time.Now(),
	})
	require.NoError(t, err)
	challenge, err = f.ToDeviceChallenge(ctx, reg)
	require.NoError(t, err)

	h := NewHandler(reg, conf)
	r := x.NewRouterAdmin(conf.AdminURL)
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	c := &http.Client{}

	for _, tc := range []struct {
		desc             string
		getBody          func() ([]byte, error)
		getURL           func() string
		validateResponse func(*http.Response)
	}{
		{
			desc: "random user_code, not persisted in the database",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
				require.NoError(t, err)
				return json.Marshal(&hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode})
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept?device_challenge=" + challenge
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			desc: "empty user_code",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				userCode := ""
				return json.Marshal(&hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode})
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept?device_challenge=" + challenge
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			desc: "empty challenge",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
				require.NoError(t, err)
				return json.Marshal(&hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode})
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept"
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			desc: "random challenge",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
				require.NoError(t, err)
				return json.Marshal(&hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode})
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept?device_challenge=abc"
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			desc: "expired user_code",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				_, deviceCodesig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(ctx)
				require.NoError(t, err)
				userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
				require.NoError(t, err)
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				exp := time.Now().UTC()
				deviceRequest.Session.SetExpiresAt(fosite.UserCode, exp)
				err = reg.OAuth2Storage().CreateDeviceAuthSession(ctx, deviceCodesig, sig, deviceRequest)
				require.NoError(t, err)
				return json.Marshal(&hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode})
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept?device_challenge=" + challenge
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusUnauthorized, resp.StatusCode)
				result := &fosite.RFC6749Error{}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.EqualValues(t, result.ErrorField, fosite.ErrTokenExpired.ErrorField)
			},
		},
		{
			desc: "extra fields",
			getBody: func() ([]byte, error) {
				deviceRequest := fosite.NewDeviceRequest()
				deviceRequest.Client = cl
				deviceRequest.SetSession(
					&oauth2.Session{
						DefaultSession: &openid.DefaultSession{
							Headers: &jwt.Headers{},
						},
					},
				)
				userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(ctx)
				require.NoError(t, err)
				ret := struct {
					UserCode *string
					Extra    string
				}{
					UserCode: &userCode,
					Extra:    "extra",
				}
				return json.Marshal(ret)
			},
			getURL: func() string {
				return ts.URL + "/admin" + DevicePath + "/accept?device_challenge=" + challenge
			},
			validateResponse: func(resp *http.Response) {
				require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
	} {
		tc := tc
		t.Run("case="+tc.desc, func(t *testing.T) {
			acceptUserCodeJson, err := tc.getBody()
			require.NoError(t, err)

			// set the HTTP method, url, and request body
			req, err := http.NewRequest(http.MethodPut, tc.getURL(), bytes.NewBuffer(acceptUserCodeJson))
			require.NoError(t, err)

			resp, err := c.Do(req)
			require.NoError(t, err)
			tc.validateResponse(resp)
		})
	}

}

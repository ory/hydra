// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	. "github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/uuidx"
)

func TestGetLogoutRequest(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t)

	h := NewHandler(reg)
	r := x.NewRouterAdmin(prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date))
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	cl := &client.Client{
		ID:   "test client id",
		Name: "test client name",
	}
	require.NoError(t, reg.ClientManager().CreateClient(t.Context(), cl))

	requestURL := "http://192.0.2.1"

	t.Run("unhandled logout request", func(t *testing.T) {
		challenge := "test-challenge-unhandled"
		require.NoError(t, reg.LogoutManager().CreateLogoutRequest(t.Context(), &flow.LogoutRequest{
			Client:     cl,
			ID:         challenge,
			RequestURL: requestURL,
			Verifier:   uuidx.NewV4().String(),
			SessionID:  "test-session-id",
			Subject:    "test-subject",
		}))

		resp, err := ts.Client().Get(ts.URL + "/admin" + LogoutPath + "?challenge=" + challenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result flow.LogoutRequest
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, challenge, result.ID)
		assert.Equal(t, requestURL, result.RequestURL)
	})

	t.Run("handled logout request", func(t *testing.T) {
		challenge := "test-challenge-handled"
		require.NoError(t, reg.LogoutManager().CreateLogoutRequest(t.Context(), &flow.LogoutRequest{
			Client:     cl,
			ID:         challenge,
			RequestURL: requestURL,
			Verifier:   uuidx.NewV4().String(),
			WasHandled: true,
		}))

		resp, err := ts.Client().Get(ts.URL + "/admin" + LogoutPath + "?challenge=" + challenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusGone, resp.StatusCode)

		var result flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, requestURL, result.RedirectTo)
	})

	t.Run("unknown challenge", func(t *testing.T) {
		resp, err := ts.Client().Get(ts.URL + "/admin" + LogoutPath + "?challenge=unknown-challenge")
		require.NoError(t, err)
		assert.EqualValuesf(t, http.StatusNotFound, resp.StatusCode, "%s", ioutilx.MustReadAll(resp.Body))
	})
}

func TestGetLoginRequest(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t)

	h := NewHandler(reg)
	r := x.NewRouterAdmin(prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date))
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	cl := &client.Client{
		ID:   "test client id",
		Name: "test client name",
	}

	requestURL := "http://192.0.2.1"

	f := &flow.Flow{
		Client:            cl,
		RequestURL:        requestURL,
		RequestedAt:       time.Now(),
		State:             flow.FlowStateLoginUnused,
		NID:               reg.Persister().NetworkID(t.Context()),
		RequestedAudience: []string{"audience1", "audience2"},
		RequestedScope:    []string{"scope1", "scope2"},
		Subject:           "test subject",
		SessionID:         "test session id",
	}

	unhandledChallenge, err := f.ToLoginChallenge(t.Context(), reg)
	require.NoError(t, err)

	t.Run("unhandled flow", func(t *testing.T) {
		resp, err := ts.Client().Get(ts.URL + "/admin" + LoginPath + "?challenge=" + unhandledChallenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result flow.LoginRequest
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, unhandledChallenge, result.ID)
		assert.Equal(t, requestURL, result.RequestURL)
		assert.NotNil(t, result.Client)
	})

	t.Run("handled flow", func(t *testing.T) {
		f.State = flow.FlowStateLoginUnused
		require.NoError(t, f.InvalidateLoginRequest())
		handledChallenge, err := f.ToLoginChallenge(t.Context(), reg)
		require.NoError(t, err)

		resp, err := ts.Client().Get(ts.URL + "/admin" + LoginPath + "?challenge=" + handledChallenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusGone, resp.StatusCode)

		var result flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, requestURL, result.RedirectTo)
	})

	t.Run("unknown challenge", func(t *testing.T) {
		resp, err := ts.Client().Get(ts.URL + "/admin" + LoginPath + "?challenge=unknown-challenge")
		require.NoError(t, err)
		assert.EqualValuesf(t, http.StatusNotFound, resp.StatusCode, "%s", ioutilx.MustReadAll(resp.Body))
	})
}

func TestGetConsentRequest(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t)

	h := NewHandler(reg)
	r := x.NewRouterAdmin(prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date))
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	cl := &client.Client{
		ID:   "test client id",
		Name: "test client name",
	}

	requestURL := "http://192.0.2.1"
	consentRequestID := "test consent request id"

	f := &flow.Flow{
		Client:           cl,
		RequestURL:       requestURL,
		RequestedAt:      time.Now(),
		State:            flow.FlowStateConsentUnused,
		NID:              reg.Persister().NetworkID(t.Context()),
		ConsentRequestID: sqlxx.NullString(consentRequestID),
	}

	unhandledChallenge, err := f.ToConsentChallenge(t.Context(), reg)
	require.NoError(t, err)

	t.Run("unhandled flow", func(t *testing.T) {
		resp, err := ts.Client().Get(ts.URL + "/admin" + ConsentPath + "?challenge=" + unhandledChallenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result flow.OAuth2ConsentRequest
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, unhandledChallenge, result.Challenge)
		assert.Equal(t, requestURL, result.RequestURL)
		assert.NotNil(t, result.Client)
	})

	t.Run("handled flow", func(t *testing.T) {
		f.State = flow.FlowStateConsentUnused
		require.NoError(t, f.InvalidateConsentRequest())
		handledChallenge, err := f.ToConsentChallenge(t.Context(), reg)
		require.NoError(t, err)

		resp, err := ts.Client().Get(ts.URL + "/admin" + ConsentPath + "?challenge=" + handledChallenge)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusGone, resp.StatusCode)

		var result flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, requestURL, result.RedirectTo)
	})

	t.Run("unknown challenge", func(t *testing.T) {
		resp, err := ts.Client().Get(ts.URL + "/admin" + ConsentPath + "?challenge=unknown-challenge")
		require.NoError(t, err)
		assert.EqualValuesf(t, http.StatusNotFound, resp.StatusCode, "%s", ioutilx.MustReadAll(resp.Body))
	})
}

func TestAcceptLoginRequestDouble(t *testing.T) {
	t.Parallel()

	requestURL := "http://192.0.2.1"

	reg := testhelpers.NewRegistryMemory(t)

	f := flow.Flow{
		Client:      &client.Client{ID: "client"},
		RequestURL:  requestURL,
		RequestedAt: time.Now(),
		NID:         reg.Persister().NetworkID(t.Context()),
		State:       flow.FlowStateLoginUnused,
	}
	challenge, err := f.ToLoginChallenge(t.Context(), reg)
	require.NoError(t, err)

	h := NewHandler(reg)
	r := x.NewRouterAdmin(prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date))
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// marshal User to json
	acceptLoginJson, err := json.Marshal(&flow.HandledLoginRequest{Subject: "sub123"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+LoginPath+"/accept?challenge="+challenge, bytes.NewReader(acceptLoginJson))
	require.NoError(t, err)

	for range 2 {
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Contains(t, result.RedirectTo, "login_verifier")
	}
}

func TestAcceptCodeDeviceRequest(t *testing.T) {
	requestURL := "https://hydra.example.com/" + oauth2.DeviceVerificationPath

	reg := testhelpers.NewRegistryMemory(t)

	cl := &client.Client{ID: "client"}
	require.NoError(t, reg.ClientManager().CreateClient(t.Context(), cl))
	f := &flow.Flow{
		Client:      cl,
		RequestURL:  requestURL,
		RequestedAt: time.Now(),
		State:       flow.DeviceFlowStateUnused,
	}
	f.NID = reg.Networker().NetworkID(t.Context())
	challenge, err := f.ToDeviceChallenge(t.Context(), reg)
	require.NoError(t, err)

	h := NewHandler(reg)
	r := x.NewRouterAdmin(prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date))
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	submitCode := func(t *testing.T, reqBody any, challenge string) *http.Response {
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, ts.URL+"/admin"+DevicePath+"/accept?device_challenge="+challenge, bytes.NewReader(body))
		require.NoError(t, err)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err)

		return resp
	}

	t.Run("case=successful user_code submission", func(t *testing.T) {
		deviceRequest := fosite.NewDeviceRequest()
		deviceRequest.Client = cl
		deviceRequest.SetSession(oauth2.NewTestSession(t, "test-subject"))

		_, deviceCodeSig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(t.Context())
		require.NoError(t, err)
		userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(t.Context(), deviceCodeSig, sig, deviceRequest))

		resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: userCode}, challenge)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		var result flow.OAuth2RedirectTo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.RedirectTo, requestURL)
		assert.Contains(t, result.RedirectTo, "device_verifier")

		t.Run("double submit", func(t *testing.T) {
			resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: userCode}, challenge)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, resp.StatusCode)

			var result flow.OAuth2RedirectTo
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
			assert.Contains(t, result.RedirectTo, requestURL)
			assert.Contains(t, result.RedirectTo, "device_verifier")
		})
	})

	t.Run("case=random user_code, not persisted in the database", func(t *testing.T) {
		userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)

		resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: userCode}, challenge)
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "The 'user_code' session could not be found or has expired or is otherwise malformed.")
	})

	t.Run("case=empty user_code", func(t *testing.T) {
		resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: ""}, challenge)
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "'user_code' must not be empty")
	})

	t.Run("case=empty challenge", func(t *testing.T) {
		userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: userCode}, "")
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "'device_challenge' is not defined but should have been")
	})

	t.Run("case=invalid challenge", func(t *testing.T) {
		userCode, _, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		resp := submitCode(t, &hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode}, "invalid-challenge")
		require.EqualValues(t, http.StatusNotFound, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, x.ErrNotFound.DescriptionField)
	})

	t.Run("case=expired user_code", func(t *testing.T) {
		deviceRequest := fosite.NewDeviceRequest()
		deviceRequest.Client = cl
		deviceRequest.SetSession(oauth2.NewTestSession(t, "test-subject"))
		deviceRequest.Session.SetExpiresAt(fosite.UserCode, time.Now().Add(-time.Hour).UTC())

		_, deviceCodeSig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(t.Context())
		require.NoError(t, err)
		userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(t.Context(), deviceCodeSig, sig, deviceRequest))

		resp := submitCode(t, &flow.AcceptDeviceUserCodeRequest{UserCode: userCode}, challenge)
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "The 'user_code' session could not be found or has expired or is otherwise malformed.")
	})

	t.Run("case=accepted user_code", func(t *testing.T) {
		deviceRequest := fosite.NewDeviceRequest()
		deviceRequest.Client = cl
		deviceRequest.SetSession(oauth2.NewTestSession(t, "test-subject"))
		deviceRequest.UserCodeState = fosite.UserCodeAccepted

		_, deviceCodeSig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(t.Context())
		require.NoError(t, err)
		userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(t.Context(), deviceCodeSig, sig, deviceRequest))

		resp := submitCode(t, &hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode}, challenge)
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "The 'user_code' session could not be found or has expired or is otherwise malformed.")
	})

	t.Run("case=rejected user_code", func(t *testing.T) {
		deviceRequest := fosite.NewDeviceRequest()
		deviceRequest.Client = cl
		deviceRequest.SetSession(oauth2.NewTestSession(t, "test-subject"))
		deviceRequest.UserCodeState = fosite.UserCodeRejected

		_, deviceCodesig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(t.Context())
		require.NoError(t, err)
		userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(t.Context(), deviceCodesig, sig, deviceRequest))

		resp := submitCode(t, &hydra.AcceptDeviceUserCodeRequest{UserCode: &userCode}, challenge)
		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

		result := fosite.RFC6749Error{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Contains(t, result.DescriptionField, fosite.ErrInvalidRequest.DescriptionField)
		assert.Contains(t, result.DescriptionField, "The 'user_code' session could not be found or has expired or is otherwise malformed.")
	})

	t.Run("case=extra fields", func(t *testing.T) {
		deviceRequest := fosite.NewDeviceRequest()
		deviceRequest.Client = cl
		deviceRequest.SetSession(oauth2.NewTestSession(t, "test-subject"))

		_, deviceCodeSig, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(t.Context())
		require.NoError(t, err)
		userCode, sig, err := reg.RFC8628HMACStrategy().GenerateUserCode(t.Context())
		require.NoError(t, err)
		require.NoError(t, reg.OAuth2Storage().CreateDeviceAuthSession(t.Context(), deviceCodeSig, sig, deviceRequest))

		resp := submitCode(t, map[string]string{
			"user_code": userCode,
			"extra":     "extra",
		}, challenge)
		assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	})
}

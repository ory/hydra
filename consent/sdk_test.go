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
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSDK(t *testing.T) {
	m := NewMemoryManager()
	router := httprouter.New()
	h := NewHandler(herodot.NewJSONWriter(logrus.New()), m)

	h.SetRoutes(router)
	ts := httptest.NewServer(router)

	sdk, err := hydra.NewSDK(&hydra.Configuration{
		EndpointURL: ts.URL,
	})
	require.NoError(t, err)

	require.NoError(t, m.CreateAuthenticationSession(&AuthenticationSession{
		ID:      "session1",
		Subject: "subject1",
	}))

	ar1, _ := mockAuthRequest("1", false)
	ar2, _ := mockAuthRequest("2", false)
	require.NoError(t, m.CreateAuthenticationRequest(ar1))
	require.NoError(t, m.CreateAuthenticationRequest(ar2))

	cr1, hcr1 := mockConsentRequest("1", false, 0, false, false, false)
	cr2, hcr2 := mockConsentRequest("2", false, 0, false, false, false)
	require.NoError(t, m.CreateConsentRequest(cr1))
	require.NoError(t, m.CreateConsentRequest(cr2))
	_, err = m.HandleConsentRequest("challenge1", hcr1)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest("challenge2", hcr2)
	require.NoError(t, err)

	crGot, res, err := sdk.GetConsentRequest("challenge1")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	compareSDKConsentRequest(t, cr1, crGot)

	crGot, res, err = sdk.GetConsentRequest("challenge2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	compareSDKConsentRequest(t, cr2, crGot)

	arGot, res, err := sdk.GetLoginRequest("challenge1")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	compareSDKLoginRequest(t, ar1, arGot)

	arGot, res, err = sdk.GetLoginRequest("challenge2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	compareSDKLoginRequest(t, ar2, arGot)

	res, err = sdk.RevokeAuthenticationSession("subject1")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNoContent, res.StatusCode)

	res, err = sdk.RevokeAllUserConsentSessions("subject1")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNoContent, res.StatusCode)

	_, res, err = sdk.GetConsentRequest("challenge1")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNotFound, res.StatusCode)

	crGot, res, err = sdk.GetConsentRequest("challenge2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	compareSDKConsentRequest(t, cr2, crGot)

	res, err = sdk.RevokeUserClientConsentSessions("subject2", "client2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNoContent, res.StatusCode)

	_, res, err = sdk.GetConsentRequest("challenge2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNotFound, res.StatusCode)
}

func compareSDKLoginRequest(t *testing.T, expected *AuthenticationRequest, got *swagger.LoginRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.ID, got.Client.Id)
}

func compareSDKConsentRequest(t *testing.T, expected *ConsentRequest, got *swagger.ConsentRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.ID, got.Client.Id)
}

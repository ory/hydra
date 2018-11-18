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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/herodot"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

func TestSDK(t *testing.T) {
	m := NewMemoryManager(oauth2.NewFositeMemoryStore(nil, time.Minute))
	router := httprouter.New()
	h := NewHandler(herodot.NewJSONWriter(logrus.New()), m, sessions.NewCookieStore([]byte("secret")), "https://www.ory.sh")

	h.SetRoutes(router, router)
	ts := httptest.NewServer(router)

	sdk, err := hydra.NewSDK(&hydra.Configuration{
		AdminURL: ts.URL,
	})
	require.NoError(t, err)

	require.NoError(t, m.CreateAuthenticationSession(context.TODO(), &AuthenticationSession{
		ID:      "session1",
		Subject: "subject1",
	}))

	ar1, _ := MockAuthRequest("1", false)
	ar2, _ := MockAuthRequest("2", false)
	require.NoError(t, m.CreateAuthenticationRequest(context.TODO(), ar1))
	require.NoError(t, m.CreateAuthenticationRequest(context.TODO(), ar2))

	cr1, hcr1 := MockConsentRequest("1", false, 0, false, false, false)
	cr2, hcr2 := MockConsentRequest("2", false, 0, false, false, false)
	cr3, hcr3 := MockConsentRequest("3", true, 3600, false, false, false)
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr1))
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr2))
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr3))
	_, err = m.HandleConsentRequest(context.TODO(), "challenge1", hcr1)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge2", hcr2)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge3", hcr3)
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

	res, err = sdk.RevokeUserClientConsentSessions("subject2", "fk-client-2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNoContent, res.StatusCode)

	_, res, err = sdk.GetConsentRequest("challenge2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNotFound, res.StatusCode)

	csGot, res, err := sdk.ListUserConsentSessions("subject3")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, len(csGot))
	cs := csGot[0]
	assert.Equal(t, "challenge3", cs.ConsentRequest.Challenge)

	csGot, res, err = sdk.ListUserConsentSessions("subject2")
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 0, len(csGot))
}

func compareSDKLoginRequest(t *testing.T, expected *AuthenticationRequest, got *swagger.LoginRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), got.Client.ClientId)
}

func compareSDKConsentRequest(t *testing.T, expected *ConsentRequest, got *swagger.ConsentRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), got.Client.ClientId)
}

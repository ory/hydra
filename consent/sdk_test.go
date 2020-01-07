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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/x"

	"github.com/ory/viper"

	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/consent"
)

func TestSDK(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyIssuerURL, "https://www.ory.sh")
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Minute)
	reg := internal.NewRegistry(conf)

	router := x.NewRouterPublic()
	h := NewHandler(reg, conf)

	h.SetRoutes(router.RouterAdmin())
	ts := httptest.NewServer(router)

	sdk := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(ts.URL).Host})

	m := reg.ConsentManager()

	require.NoError(t, m.CreateLoginSession(context.TODO(), &LoginSession{
		ID:      "session1",
		Subject: "subject1",
	}))

	ar1, _ := MockAuthRequest("1", false)
	ar2, _ := MockAuthRequest("2", false)
	require.NoError(t, m.CreateLoginRequest(context.TODO(), ar1))
	require.NoError(t, m.CreateLoginRequest(context.TODO(), ar2))

	cr1, hcr1 := MockConsentRequest("1", false, 0, false, false, false)
	cr2, hcr2 := MockConsentRequest("2", false, 0, false, false, false)
	cr3, hcr3 := MockConsentRequest("3", true, 3600, false, false, false)
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr1))
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr2))
	require.NoError(t, m.CreateConsentRequest(context.TODO(), cr3))
	_, err := m.HandleConsentRequest(context.TODO(), "challenge1", hcr1)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge2", hcr2)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge3", hcr3)
	require.NoError(t, err)

	lur1 := MockLogoutRequest("testsdk-1", true)
	require.NoError(t, m.CreateLogoutRequest(context.TODO(), lur1))

	lur2 := MockLogoutRequest("testsdk-2", false)
	require.NoError(t, m.CreateLogoutRequest(context.TODO(), lur2))

	crGot, err := sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge("challenge1"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr1, *crGot.Payload)

	crGot, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot.Payload)

	arGot, err := sdk.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge("challenge1"))
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar1, *arGot.Payload)

	arGot, err = sdk.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar2, *arGot.Payload)

	_, err = sdk.Admin.RevokeAuthenticationSession(admin.NewRevokeAuthenticationSessionParams().WithSubject("subject1"))
	require.NoError(t, err)

	_, err = sdk.Admin.RevokeConsentSessions(admin.NewRevokeConsentSessionsParams().WithSubject("subject1"))
	require.NoError(t, err)

	_, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge("challenge1"))
	require.Error(t, err)

	crGot, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot.Payload)

	_, err = sdk.Admin.RevokeConsentSessions(admin.NewRevokeConsentSessionsParams().WithSubject("subject1").WithSubject("subject2").WithClient(pointerx.String("fk-client-2")))
	require.NoError(t, err)

	_, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge("challenge2"))
	require.Error(t, err)

	csGot, err := sdk.Admin.ListSubjectConsentSessions(admin.NewListSubjectConsentSessionsParams().WithSubject("subject3"))
	require.NoError(t, err)
	assert.Equal(t, 1, len(csGot.Payload))
	cs := csGot.Payload[0]
	assert.Equal(t, "challenge3", cs.ConsentRequest.Challenge)

	csGot, err = sdk.Admin.ListSubjectConsentSessions(admin.NewListSubjectConsentSessionsParams().WithSubject("subject2"))
	require.NoError(t, err)
	assert.Equal(t, 0, len(csGot.Payload))

	luGot, err := sdk.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().WithLogoutChallenge("challengetestsdk-1"))
	require.NoError(t, err)
	compareSDKLogoutRequest(t, lur1, luGot.Payload)

	luaGot, err := sdk.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().WithLogoutChallenge("challengetestsdk-1"))
	require.NoError(t, err)
	assert.EqualValues(t, "https://www.ory.sh/oauth2/sessions/logout?logout_verifier=verifiertestsdk-1", luaGot.Payload.RedirectTo)

	_, err = sdk.Admin.RejectLogoutRequest(admin.NewRejectLogoutRequestParams().WithLogoutChallenge("challengetestsdk-2"))
	require.NoError(t, err)

	_, err = sdk.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().WithLogoutChallenge("challengetestsdk-2"))
	require.Error(t, err)
}

func compareSDKLoginRequest(t *testing.T, expected *LoginRequest, got models.LoginRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), got.Client.ClientID)
}

func compareSDKConsentRequest(t *testing.T, expected *ConsentRequest, got models.ConsentRequest) {
	assert.EqualValues(t, expected.Challenge, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), got.Client.ClientID)
}

func compareSDKLogoutRequest(t *testing.T, expected *LogoutRequest, got *models.LogoutRequest) {
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.SessionID, got.Sid)
	assert.EqualValues(t, expected.SessionID, got.Sid)
	assert.EqualValues(t, expected.RequestURL, got.RequestURL)
	assert.EqualValues(t, expected.RPInitiated, got.RpInitiated)
}

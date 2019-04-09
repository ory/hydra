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

	"github.com/ory/hydra/sdk/go/hydra/client"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/x"

	"github.com/spf13/viper"

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

	h.SetRoutes(router.RouterAdmin(), router)
	ts := httptest.NewServer(router)

	sdk := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(ts.URL).Host})

	m := reg.ConsentManager()

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
	_, err := m.HandleConsentRequest(context.TODO(), "challenge1", hcr1)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge2", hcr2)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.TODO(), "challenge3", hcr3)
	require.NoError(t, err)

	crGot, err := sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithChallenge("challenge1"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr1, *crGot.Payload)

	crGot, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot.Payload)

	arGot, err := sdk.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithChallenge("challenge1"))
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar1, *arGot.Payload)

	arGot, err = sdk.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar2, *arGot.Payload)

	_, err = sdk.Admin.RevokeAuthenticationSession(admin.NewRevokeAuthenticationSessionParams().WithUser("subject1"))
	require.NoError(t, err)

	_, err = sdk.Admin.RevokeAllUserConsentSessions(admin.NewRevokeAllUserConsentSessionsParams().WithUser("subject1"))
	require.NoError(t, err)

	_, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithChallenge("challenge1"))
	require.NoError(t, err)

	crGot, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithChallenge("challenge2"))
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot.Payload)

	_, err = sdk.Admin.RevokeUserClientConsentSessions(admin.NewRevokeUserClientConsentSessionsParams().WithUser("subject2").WithClient("fk-client-2"))
	require.NoError(t, err)

	_, err = sdk.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithChallenge("challenge2"))
	require.Error(t, err)

	csGot, err := sdk.Admin.ListUserConsentSessions(admin.NewListUserConsentSessionsParams().WithUser("subject3"))
	require.NoError(t, err)
	assert.Equal(t, 1, len(csGot.Payload))
	cs := csGot.Payload[0]
	assert.Equal(t, "challenge3", cs.ConsentRequest.Challenge)

	csGot, err = sdk.Admin.ListUserConsentSessions(admin.NewListUserConsentSessionsParams().WithUser("subject2"))
	require.NoError(t, err)
	assert.Equal(t, 0, len(csGot.Payload))
}

func compareSDKLoginRequest(t *testing.T, expected *AuthenticationRequest, got models.AuthenticationRequest) {
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

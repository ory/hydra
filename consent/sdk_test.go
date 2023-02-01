// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	hydra "github.com/ory/hydra-client-go/v2"

	"github.com/ory/x/httprouterx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
)

func makeID(base string, network string, key string) string {
	return fmt.Sprintf("%s-%s-%s", base, network, key)
}

func TestSDK(t *testing.T) {
	ctx := context.Background()
	network := "t1"
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(ctx, config.KeyIssuerURL, "https://www.ory.sh")
	conf.MustSet(ctx, config.KeyAccessTokenLifespan, time.Minute)
	reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

	router := x.NewRouterPublic()
	h := NewHandler(reg, conf)

	h.SetRoutes(httprouterx.NewRouterAdminWithPrefixAndRouter(router.Router, "/admin", conf.AdminURL))
	ts := httptest.NewServer(router)

	sdk := hydra.NewAPIClient(hydra.NewConfiguration())
	sdk.GetConfig().Servers = hydra.ServerConfigurations{{URL: ts.URL}}

	m := reg.ConsentManager()

	require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
		ID:      "session1",
		Subject: "subject1",
	}))

	ar1, _ := MockAuthRequest("ar-1", false, network)
	ar2, _ := MockAuthRequest("ar-2", false, network)
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), ar1.Client))
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), ar2.Client))
	require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
		ID:      ar1.SessionID.String(),
		Subject: ar1.Subject,
	}))
	require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
		ID:      ar2.SessionID.String(),
		Subject: ar2.Subject,
	}))
	require.NoError(t, m.CreateLoginRequest(context.Background(), ar1))
	require.NoError(t, m.CreateLoginRequest(context.Background(), ar2))

	cr1, hcr1 := MockConsentRequest("1", false, 0, false, false, false, "fk-login-challenge", network)
	cr2, hcr2 := MockConsentRequest("2", false, 0, false, false, false, "fk-login-challenge", network)
	cr3, hcr3 := MockConsentRequest("3", true, 3600, false, false, false, "fk-login-challenge", network)
	cr4, hcr4 := MockConsentRequest("4", true, 3600, false, false, false, "fk-login-challenge", network)
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cr1.Client))
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cr2.Client))
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cr3.Client))
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), cr4.Client))
	require.NoError(t, m.CreateLoginRequest(context.Background(), &LoginRequest{ID: cr1.LoginChallenge.String(), Subject: cr1.Subject, Client: cr1.Client, Verifier: cr1.ID}))
	require.NoError(t, m.CreateLoginRequest(context.Background(), &LoginRequest{ID: cr2.LoginChallenge.String(), Subject: cr2.Subject, Client: cr2.Client, Verifier: cr2.ID}))
	require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{ID: cr3.LoginSessionID.String()}))
	require.NoError(t, m.CreateLoginRequest(context.Background(), &LoginRequest{ID: cr3.LoginChallenge.String(), Subject: cr3.Subject, Client: cr3.Client, Verifier: cr3.ID, RequestedAt: hcr3.RequestedAt, SessionID: cr3.LoginSessionID}))
	require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{ID: cr4.LoginSessionID.String()}))
	require.NoError(t, m.CreateLoginRequest(context.Background(), &LoginRequest{ID: cr4.LoginChallenge.String(), Client: cr4.Client, Verifier: cr4.ID, SessionID: cr4.LoginSessionID}))
	require.NoError(t, m.CreateConsentRequest(context.Background(), cr1))
	require.NoError(t, m.CreateConsentRequest(context.Background(), cr2))
	require.NoError(t, m.CreateConsentRequest(context.Background(), cr3))
	require.NoError(t, m.CreateConsentRequest(context.Background(), cr4))
	_, err := m.HandleConsentRequest(context.Background(), hcr1)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.Background(), hcr2)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.Background(), hcr3)
	require.NoError(t, err)
	_, err = m.HandleConsentRequest(context.Background(), hcr4)
	require.NoError(t, err)

	lur1 := MockLogoutRequest("testsdk-1", true, network)
	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), lur1.Client))
	require.NoError(t, m.CreateLogoutRequest(context.Background(), lur1))

	lur2 := MockLogoutRequest("testsdk-2", false, network)
	require.NoError(t, m.CreateLogoutRequest(context.Background(), lur2))

	crGot, _, err := sdk.OAuth2Api.GetOAuth2ConsentRequest(ctx).ConsentChallenge(makeID("challenge", network, "1")).Execute()
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr1, *crGot)

	crGot, _, err = sdk.OAuth2Api.GetOAuth2ConsentRequest(ctx).ConsentChallenge(makeID("challenge", network, "2")).Execute()
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot)

	arGot, _, err := sdk.OAuth2Api.GetOAuth2LoginRequest(ctx).LoginChallenge(makeID("challenge", network, "ar-1")).Execute()
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar1, *arGot)

	arGot, _, err = sdk.OAuth2Api.GetOAuth2LoginRequest(ctx).LoginChallenge(makeID("challenge", network, "ar-2")).Execute()
	require.NoError(t, err)
	compareSDKLoginRequest(t, ar2, *arGot)

	_, err = sdk.OAuth2Api.RevokeOAuth2LoginSessions(ctx).Subject("subject1").Execute()
	require.NoError(t, err)

	_, err = sdk.OAuth2Api.RevokeOAuth2ConsentSessions(ctx).Subject("subject1").Execute()
	require.Error(t, err)

	_, err = sdk.OAuth2Api.RevokeOAuth2ConsentSessions(ctx).Subject(cr4.Subject).Client(cr4.Client.GetID()).Execute()
	require.NoError(t, err)

	_, err = sdk.OAuth2Api.RevokeOAuth2ConsentSessions(ctx).Subject("subject1").All(true).Execute()
	require.NoError(t, err)

	_, _, err = sdk.OAuth2Api.GetOAuth2ConsentRequest(ctx).ConsentChallenge(makeID("challenge", network, "1")).Execute()
	require.Error(t, err)

	crGot, _, err = sdk.OAuth2Api.GetOAuth2ConsentRequest(ctx).ConsentChallenge(makeID("challenge", network, "2")).Execute()
	require.NoError(t, err)
	compareSDKConsentRequest(t, cr2, *crGot)

	_, err = sdk.OAuth2Api.RevokeOAuth2ConsentSessions(ctx).Subject("subject2").Client("fk-client-2").Execute()
	require.NoError(t, err)

	_, _, err = sdk.OAuth2Api.GetOAuth2ConsentRequest(ctx).ConsentChallenge(makeID("challenge", network, "2")).Execute()
	require.Error(t, err)

	csGot, _, err := sdk.OAuth2Api.ListOAuth2ConsentSessions(ctx).Subject("subject3").Execute()
	require.NoError(t, err)
	assert.Equal(t, 1, len(csGot))
	cs := csGot[0]
	assert.Equal(t, makeID("challenge", network, "3"), cs.ConsentRequest.Challenge)

	csGot, _, err = sdk.OAuth2Api.ListOAuth2ConsentSessions(ctx).Subject("subject2").Execute()
	require.NoError(t, err)
	assert.Equal(t, 0, len(csGot))

	csGot, _, err = sdk.OAuth2Api.ListOAuth2ConsentSessions(ctx).Subject("subject3").LoginSessionId("fk-login-session-t1-3").Execute()
	require.NoError(t, err)
	assert.Equal(t, 1, len(csGot))
	cs = csGot[0]
	assert.Equal(t, makeID("challenge", network, "3"), cs.ConsentRequest.Challenge)

	csGot, _, err = sdk.OAuth2Api.ListOAuth2ConsentSessions(ctx).Subject("subject3").LoginSessionId("fk-login-session-t1-X").Execute()
	require.NoError(t, err)
	assert.Equal(t, 0, len(csGot))

	luGot, _, err := sdk.OAuth2Api.GetOAuth2LogoutRequest(ctx).LogoutChallenge(makeID("challenge", network, "testsdk-1")).Execute()
	require.NoError(t, err)
	compareSDKLogoutRequest(t, lur1, luGot)

	luaGot, _, err := sdk.OAuth2Api.AcceptOAuth2LogoutRequest(ctx).LogoutChallenge(makeID("challenge", network, "testsdk-1")).Execute()
	require.NoError(t, err)
	assert.EqualValues(t, "https://www.ory.sh/oauth2/sessions/logout?logout_verifier="+makeID("verifier", network, "testsdk-1"), luaGot.RedirectTo)

	_, err = sdk.OAuth2Api.RejectOAuth2LogoutRequest(ctx).LogoutChallenge(lur2.ID).Execute()
	require.NoError(t, err)

	_, _, err = sdk.OAuth2Api.GetOAuth2LogoutRequest(ctx).LogoutChallenge(lur2.ID).Execute()
	require.Error(t, err)
}

func compareSDKLoginRequest(t *testing.T, expected *LoginRequest, got hydra.OAuth2LoginRequest) {
	assert.EqualValues(t, expected.ID, got.Challenge)
	assert.EqualValues(t, expected.Subject, got.Subject)
	assert.EqualValues(t, expected.Skip, got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), *got.Client.ClientId)
}

func compareSDKConsentRequest(t *testing.T, expected *OAuth2ConsentRequest, got hydra.OAuth2ConsentRequest) {
	assert.EqualValues(t, expected.ID, got.Challenge)
	assert.EqualValues(t, expected.Subject, *got.Subject)
	assert.EqualValues(t, expected.Skip, *got.Skip)
	assert.EqualValues(t, expected.Client.GetID(), *got.Client.ClientId)
}

func compareSDKLogoutRequest(t *testing.T, expected *LogoutRequest, got *hydra.OAuth2LogoutRequest) {
	assert.EqualValues(t, expected.Subject, *got.Subject)
	assert.EqualValues(t, expected.SessionID, *got.Sid)
	assert.EqualValues(t, expected.SessionID, *got.Sid)
	assert.EqualValues(t, expected.RequestURL, *got.RequestUrl)
	assert.EqualValues(t, expected.RPInitiated, *got.RpInitiated)
}

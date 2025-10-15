// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

func mockConsentRequest(key string, remember bool, rememberFor int, skip bool, loginChallengeBase, network string) (c *flow.OAuth2ConsentRequest, h *flow.AcceptOAuth2ConsentRequest, f *flow.Flow) {
	c = &flow.OAuth2ConsentRequest{
		ConsentRequestID:  makeID("challenge", network, key),
		RequestedScope:    []string{"scopea" + key, "scopeb" + key},
		RequestedAudience: []string{"auda" + key, "audb" + key},
		Skip:              skip,
		Subject:           "subject" + key,
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		Client:         &client.Client{ID: "fk-client-" + key},
		RequestURL:     "https://request-url/path" + key,
		LoginChallenge: sqlxx.NullString(makeID(loginChallengeBase, network, key)),
		LoginSessionID: sqlxx.NullString(makeID("fk-login-session", network, key)),
		ACR:            "1",
		Context:        sqlxx.JSONRawMessage(`{"foo": "bar` + key + `"}`),
	}

	f = &flow.Flow{
		ID:                   c.LoginChallenge.String(),
		LoginVerifier:        makeID("login-verifier", network, key),
		SessionID:            c.LoginSessionID,
		Client:               c.Client,
		State:                flow.FlowStateConsentUnused,
		ConsentRequestID:     sqlxx.NullString(c.ConsentRequestID),
		ConsentSkip:          c.Skip,
		ConsentVerifier:      sqlxx.NullString(makeID("verifier", network, key)),
		ConsentCSRF:          sqlxx.NullString("csrf" + key),
		OpenIDConnectContext: c.OpenIDConnectContext,
		Subject:              c.Subject,
		RequestedScope:       c.RequestedScope,
		RequestedAudience:    c.RequestedAudience,
		RequestURL:           c.RequestURL,
		RequestedAt:          time.Now().UTC(),
	}

	h = &flow.AcceptOAuth2ConsentRequest{
		RememberFor:     rememberFor,
		Remember:        remember,
		GrantedScope:    []string{"scopea" + key, "scopeb" + key},
		GrantedAudience: []string{"auda" + key, "audb" + key},
	}

	return c, h, f
}

func mockConsentError(key string) *flow.RequestDeniedError {
	return &flow.RequestDeniedError{
		Name:        "error_name" + key,
		Description: "error_description" + key,
		Hint:        "error_hint,omitempty" + key,
		Code:        100,
		Debug:       "error_debug,omitempty" + key,
		Valid:       true,
	}
}

func mockLogoutRequest(key string, withClient bool, network string) (c *flow.LogoutRequest) {
	var cl *client.Client
	if withClient {
		cl = &client.Client{
			ID: "fk-client-" + key,
		}
	}
	return &flow.LogoutRequest{
		Subject:               "subject" + key,
		ID:                    makeID("challenge", network, key),
		Verifier:              makeID("verifier", network, key),
		SessionID:             makeID("session", network, key),
		RPInitiated:           true,
		RequestURL:            "http://request-me/",
		PostLogoutRedirectURI: "http://redirect-me/",
		WasHandled:            false,
		Accepted:              false,
		Client:                cl,
	}
}

func mockDeviceRequest(key, network string) (h *flow.HandledDeviceUserAuthRequest, f *flow.Flow) {
	cl := &client.Client{ID: "fk-client-" + key}

	f = &flow.Flow{
		DeviceChallengeID: sqlxx.NullString(makeID("challenge", network, key)),
		DeviceVerifier:    sqlxx.NullString(makeID("verifier", network, key)),
		DeviceCSRF:        sqlxx.NullString("csrf" + key),
		RequestedAt:       time.Now().UTC().Add(-time.Minute),
		Client:            cl,
		RequestURL:        "https://request-url/path" + key,
		State:             flow.DeviceFlowStateUnused,
		Context:           sqlxx.JSONRawMessage("{}"),
	}

	h = &flow.HandledDeviceUserAuthRequest{Client: cl}

	return h, f
}

func mockAuthRequest(key, network string) (h *flow.HandledLoginRequest, f *flow.Flow) {
	f = &flow.Flow{
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		RequestedAt:    time.Now().UTC().Add(-time.Minute),
		Client:         &client.Client{ID: "fk-client-" + key},
		Subject:        "subject" + key,
		RequestURL:     "https://request-url/path" + key,
		LoginSkip:      true,
		ID:             makeID("challenge", network, key),
		LoginVerifier:  makeID("verifier", network, key),
		RequestedScope: []string{"scopea" + key, "scopeb" + key},
		LoginCSRF:      "csrf" + key,
		SessionID:      sqlxx.NullString(makeID("fk-login-session", network, key)),
		NID:            uuid.FromStringOrNil(network),
		State:          flow.FlowStateLoginUnused,
	}

	h = &flow.HandledLoginRequest{
		RememberFor:            120,
		Remember:               true,
		Subject:                f.Subject,
		ACR:                    "acr",
		ForceSubjectIdentifier: "forced-subject",
	}

	return h, f
}

func saneMockHandleConsentRequest(t *testing.T, f *flow.Flow, rememberFor int, remember bool) *flow.AcceptOAuth2ConsentRequest {
	h := &flow.AcceptOAuth2ConsentRequest{
		RememberFor:     rememberFor,
		Remember:        remember,
		GrantedScope:    []string{"scopea", "scopeb"},
		GrantedAudience: []string{"auda", "audb"},
	}

	require.NoError(t, f.HandleConsentRequest(h))

	return h
}

func makeID(base, network, key string) string {
	return fmt.Sprintf("%s-%s-%s", base, network, key)
}

func TestHelperNID(t1ValidNID, t2InvalidNID consent.Manager) func(t *testing.T) {
	testLS := flow.LoginSession{
		ID:      "2022-03-11-ls-nid-test-1",
		Subject: "2022-03-11-test-1-sub",
	}
	return func(t *testing.T) {
		ctx := t.Context()

		require.Error(t, t2InvalidNID.ConfirmLoginSession(ctx, &testLS))
		require.NoError(t, t1ValidNID.ConfirmLoginSession(ctx, &testLS))
		ls, err := t2InvalidNID.DeleteLoginSession(ctx, testLS.ID)
		require.Error(t, err)
		assert.Nil(t, ls)
		ls, err = t1ValidNID.DeleteLoginSession(ctx, testLS.ID)
		require.NoError(t, err)
		assert.EqualValues(t, testLS.ID, ls.ID)
	}
}

type Deps interface {
	FlowCipher() *aead.XChaCha20Poly1305
	contextx.Provider
	x.TracingProvider
	x.NetworkProvider
	config.Provider
}

func ManagerTests(deps Deps, m consent.Manager, clientManager client.Manager, fositeManager x.FositeStorer, network string, parallel bool) func(t *testing.T) {
	lr := make(map[string]*flow.Flow)

	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		ctx := context.Background()
		t.Run("case=init-fks", func(t *testing.T) {
			for _, k := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "rv1", "rv2"} {
				require.NoError(t, clientManager.CreateClient(ctx, &client.Client{ID: fmt.Sprintf("fk-client-%s", k)}))

				loginSession := &flow.LoginSession{
					ID:              makeID("fk-login-session", network, k),
					NID:             m.NetworkID(ctx),
					AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
					Subject:         fmt.Sprintf("subject-%s", k),
				}
				require.NoError(t, m.ConfirmLoginSession(ctx, loginSession))

				lr[k] = &flow.Flow{
					ID:                   makeID("fk-login-challenge", network, k),
					Subject:              fmt.Sprintf("subject%s", k),
					SessionID:            sqlxx.NullString(makeID("fk-login-session", network, k)),
					LoginVerifier:        makeID("fk-login-verifier", network, k),
					Client:               &client.Client{ID: fmt.Sprintf("fk-client-%s", k)},
					LoginAuthenticatedAt: sqlxx.NullTime(time.Now()),
					RequestedAt:          time.Now(),
					State:                flow.FlowStateLoginUnused,
				}
			}
		})

		t.Run("case=auth-session", func(t *testing.T) {
			for _, s := range []flow.LoginSession{
				{
					ID:              makeID("session", network, "1"),
					NID:             m.NetworkID(ctx),
					AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-time.Minute).UTC()),
					Subject:         "subject1",
				},
				{
					ID:              makeID("session", network, "2"),
					NID:             m.NetworkID(ctx),
					AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
					Subject:         "subject2",
				},
			} {
				t.Run("case=create-get-"+s.ID, func(t *testing.T) {
					_, err := m.GetRememberedLoginSession(ctx, s.ID)
					require.ErrorIs(t, err, x.ErrNotFound)

					updatedAuth := time.Time(s.AuthenticatedAt).Add(time.Second)
					s.AuthenticatedAt = sqlxx.NullTime(updatedAuth)
					require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
						ID:              s.ID,
						AuthenticatedAt: sqlxx.NullTime(updatedAuth),
						Subject:         s.Subject,
						Remember:        true,
					}))

					got, err := m.GetRememberedLoginSession(ctx, s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, s.ID, got.ID)
					assert.EqualValues(t, s.AuthenticatedAt, got.AuthenticatedAt) // this was updated from confirm...
					assert.EqualValues(t, s.Subject, got.Subject)

					// Make sure AuthAt does not equal...
					updatedAuth2 := updatedAuth.Add(1 * time.Second).UTC()
					require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
						ID:              s.ID,
						AuthenticatedAt: sqlxx.NullTime(updatedAuth2),
						Subject:         "some-other-subject",
						Remember:        true,
					}))

					got2, err := m.GetRememberedLoginSession(ctx, s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, s.ID, got2.ID)
					assert.EqualValues(t, updatedAuth2.Unix(), time.Time(got2.AuthenticatedAt).Unix()) // this was updated from confirm...
					assert.EqualValues(t, "some-other-subject", got2.Subject)
				})
			}
			for _, tc := range []struct {
				id string
			}{
				{
					id: makeID("session", network, "1"),
				},
				{
					id: makeID("session", network, "2"),
				},
			} {
				t.Run("case=delete-get-"+tc.id, func(t *testing.T) {
					ls, err := m.DeleteLoginSession(ctx, tc.id)
					require.NoError(t, err)
					assert.EqualValues(t, tc.id, ls.ID)

					_, err = m.GetRememberedLoginSession(ctx, tc.id)
					require.Error(t, err)
				})
			}
		})

		t.Run("case=device-request", func(t *testing.T) {
			challenges := make([]string, 0)

			h, f := mockDeviceRequest("0", network)
			_ = clientManager.CreateClient(ctx, f.Client) // Ignore errors that are caused by duplication
			deviceChallenge := x.Must(f.ToDeviceChallenge(ctx, deps))

			_, err := flow.DecodeFromDeviceChallenge(ctx, deps, deviceChallenge)
			require.ErrorIs(t, err, x.ErrNotFound)

			f.NID = m.NetworkID(ctx)
			deviceChallenge = x.Must(f.ToDeviceChallenge(ctx, deps))

			got1, err := flow.DecodeFromDeviceChallenge(ctx, deps, deviceChallenge)
			require.NoError(t, err)
			got1.Client.NID = f.NID // the client nid is not encoded in the challenge
			assert.Equal(t, f, got1)

			require.NoError(t, f.HandleDeviceUserAuthRequest(h))

			for _, key := range []string{"1", "2", "3", "4", "5", "6", "7"} {
				h, f := mockDeviceRequest(key, network)
				deviceChallenge := x.Must(f.ToDeviceChallenge(ctx, deps))

				_, err := flow.DecodeFromDeviceChallenge(ctx, deps, deviceChallenge)
				require.ErrorIs(t, err, x.ErrNotFound)

				f.NID = m.NetworkID(ctx)
				deviceChallenge = x.Must(f.ToDeviceChallenge(ctx, deps))
				challenges = append(challenges, deviceChallenge)

				got1, err := flow.DecodeFromDeviceChallenge(ctx, deps, deviceChallenge)
				require.NoError(t, err)
				assert.Equal(t, flow.DeviceFlowStateUnused, got1.State)
				assert.Equal(t, f, got1)

				require.NoError(t, f.HandleDeviceUserAuthRequest(h))
			}

			deviceVerifier := x.Must(f.ToDeviceVerifier(ctx, deps))

			{
				f := *f
				got2, err := flow.DecodeAndInvalidateDeviceVerifier(ctx, deps, deviceVerifier)
				require.NoError(t, err)
				got2.Client.NID = f.NID // the client nid is not encoded in the challenge
				f.State = flow.DeviceFlowStateUsed
				assert.Equal(t, &f, got2)
			}
			{
				f := *f
				deviceChallenge = x.Must(f.ToDeviceChallenge(ctx, deps))
				authReq, err := flow.DecodeFromDeviceChallenge(ctx, deps, deviceChallenge)
				require.NoError(t, err)
				authReq.Client.NID = f.NID // the client nid is not encoded in the challenge
				assert.Equal(t, &f, authReq)
			}
		})

		t.Run("case=auth-request", func(t *testing.T) {
			for _, key := range []string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
			} {
				t.Run("key="+key, func(t *testing.T) {
					h, f := mockAuthRequest(key, network)
					_ = clientManager.CreateClient(ctx, f.Client) // Ignore errors that are caused by duplication
					loginChallenge := x.Must(f.ToLoginChallenge(ctx, deps))

					_, err := flow.DecodeFromLoginChallenge(ctx, deps, loginChallenge)
					require.Error(t, err)

					f.NID = m.NetworkID(ctx)
					loginChallenge = x.Must(f.ToLoginChallenge(ctx, deps))

					got1, err := flow.DecodeFromLoginChallenge(ctx, deps, loginChallenge)
					require.NoError(t, err)
					assert.Equal(t, flow.FlowStateLoginUnused, got1.State)
					compareAuthenticationRequest(t, f.GetLoginRequest(), got1.GetLoginRequest())

					err = f.HandleLoginRequest(h)
					require.NoError(t, err)

					require.NoError(t, f.InvalidateLoginRequest())
					compareAuthenticationRequestFlow(t, f.GetLoginRequest(), f)

					loginChallenge = x.Must(f.ToLoginChallenge(ctx, deps))
					got1, err = flow.DecodeFromLoginChallenge(ctx, deps, loginChallenge)
					require.NoError(t, err)
				})
			}
		})

		t.Run("case=consent-request", func(t *testing.T) {
			for _, tc := range []struct {
				key         string
				remember    bool
				rememberFor int
				skip        bool
				authAt      bool
			}{
				{"1", true, 0, false, true},
				// {"2", true, 0, false, true}, // error case, moved to its own test below
				{"3", true, 1, false, true},
				{"4", false, 0, false, true},
				{"5", true, 120, false, true},
				{"6", true, 120, true, true},
				{"7", false, 0, false, false},
			} {
				t.Run("key="+tc.key, func(t *testing.T) {
					consentRequest, h, f := mockConsentRequest(tc.key, tc.remember, tc.rememberFor, tc.skip, "challenge", network)
					_ = clientManager.CreateClient(ctx, consentRequest.Client) // Ignore errors that are caused by duplication
					f.NID = m.NetworkID(ctx)

					require.NoError(t, f.HandleConsentRequest(h))
					assert.WithinDuration(t, time.Now(), time.Time(f.ConsentHandledAt), 5*time.Second)

					h.GrantedAudience = sqlxx.StringSliceJSONFormat{"new-audience"}
					require.NoError(t, f.HandleConsentRequest(h))

					require.NoError(t, f.InvalidateConsentRequest())
					require.NoError(t, m.CreateConsentSession(ctx, f))
					compareConsentRequestFlow(t, consentRequest, f)
					assert.EqualValues(t, consentRequest.ConsentRequestID, f.ConsentRequestID)
					assert.EqualValues(t, h.GrantedAudience, f.GrantedAudience)

					t.Run("sub=detect double-submit for consent verifier", func(t *testing.T) {
						require.ErrorIs(t, m.CreateConsentSession(ctx, f), sqlcon.ErrUniqueViolation)
					})

					assert.EqualValues(t, tc.remember, f.ConsentRemember)
					assert.EqualValues(t, tc.rememberFor, *f.ConsentRememberFor)
				})
			}

			t.Run("case=consent-error", func(t *testing.T) {
				err := mockConsentError("2")
				_, _, f := mockConsentRequest("2", true, 0, false, "challenge", network)
				f.NID = m.NetworkID(ctx)

				_ = clientManager.CreateClient(ctx, f.Client) // Ignore errors that are caused by duplication
				require.NoError(t, f.HandleConsentError(err))
				assert.Equal(t, flow.FlowStateConsentError, f.State)
				assert.True(t, f.ConsentError.IsError())
			})

			for _, tc := range []struct {
				keyC          string
				keyS          string
				expectedFound bool
			}{
				{"1", "1", true},
				{"2", "2", false},
				// {"3", "3", false},  // Some consent is given in some other test case. Yay global fixtues :)
				{"4", "4", false},
				{"1", "2", false},
				{"2", "1", false},
				{"5", "5", true},
				{"6", "6", false},
			} {
				t.Run("key="+tc.keyC+"-"+tc.keyS, func(t *testing.T) {
					rs, err := m.FindGrantedAndRememberedConsentRequest(ctx, "fk-client-"+tc.keyC, "subject"+tc.keyS)
					if !tc.expectedFound {
						assert.Nil(t, rs)
						assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
					} else {
						require.NoError(t, err)
						assert.NotNil(t, rs)
					}
				})
			}
		})

		t.Run("case=revoke-auth-request", func(t *testing.T) {
			require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
				ID:              makeID("rev-session", network, "-1"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-1",
			}))

			require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
				ID:              makeID("rev-session", network, "-2"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-2",
			}))

			require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
				ID:              makeID("rev-session", network, "-3"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-1",
			}))

			for i, tc := range []struct {
				subject string
				ids     []string
			}{
				{
					subject: "subject-1",
					ids:     []string{makeID("rev-session", network, "-1"), makeID("rev-session", network, "-3")},
				},
				{
					subject: "subject-2",
					ids:     []string{makeID("rev-session", network, "-1"), makeID("rev-session", network, "-3"), makeID("rev-session", network, "-2")},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					require.NoError(t, m.RevokeSubjectLoginSession(ctx, tc.subject))

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetRememberedLoginSession(ctx, id)
							assert.EqualError(t, err, x.ErrNotFound.Error())
						})
					}
				})
			}
		})

		challengerv1 := makeID("challenge", network, "rv1")
		challengerv2 := makeID("challenge", network, "rv2")
		t.Run("case=revoke-used-consent-request", func(t *testing.T) {

			cr1, hcr1, f1 := mockConsentRequest("rv1", false, 0, false, "fk-login-challenge", network)
			cr2, hcr2, f2 := mockConsentRequest("rv2", false, 0, false, "fk-login-challenge", network)
			f1.NID = m.NetworkID(ctx)
			f2.NID = m.NetworkID(ctx)

			// Ignore duplication errors
			_ = clientManager.CreateClient(ctx, cr1.Client)
			_ = clientManager.CreateClient(ctx, cr2.Client)

			require.NoError(t, f1.HandleConsentRequest(hcr1))
			require.NoError(t, f2.HandleConsentRequest(hcr2))

			require.NoError(t, f1.InvalidateConsentRequest())
			require.NoError(t, f2.InvalidateConsentRequest())

			require.NoError(t, m.CreateConsentSession(ctx, f1))
			require.NoError(t, m.CreateConsentSession(ctx, f2))

			require.NoError(t, fositeManager.CreateAccessTokenSession(
				ctx,
				makeID("", network, "trva1"),
				&fosite.Request{Client: cr1.Client, ID: f1.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(
				ctx,
				makeID("", network, "rrva1"),
				"",
				&fosite.Request{Client: cr1.Client, ID: f1.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateAccessTokenSession(
				ctx,
				makeID("", network, "trva2"),
				&fosite.Request{Client: cr2.Client, ID: f2.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(
				ctx,
				makeID("", network, "rrva2"),
				"",
				&fosite.Request{Client: cr2.Client, ID: f2.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))

			for i, tc := range []struct {
				subject string
				client  string
				at      string
				rt      string
				ids     []string
			}{
				{
					at:      makeID("", network, "trva1"),
					rt:      makeID("", network, "rrva1"),
					subject: "subjectrv1",
					client:  "",
					ids:     []string{challengerv1},
				},
				{
					at:      makeID("", network, "trva2"),
					rt:      makeID("", network, "rrva2"),
					subject: "subjectrv2",
					client:  "fk-client-rv2",
					ids:     []string{challengerv2},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					_, err := fositeManager.GetAccessTokenSession(ctx, tc.at, nil)
					assert.NoError(t, err)
					_, err = fositeManager.GetRefreshTokenSession(ctx, tc.rt, nil)
					assert.NoError(t, err)

					if tc.client == "" {
						require.NoError(t, m.RevokeSubjectConsentSession(ctx, tc.subject))
					} else {
						require.NoError(t, m.RevokeSubjectClientConsentSession(ctx, tc.subject, tc.client))
					}

					r, err := fositeManager.GetAccessTokenSession(ctx, tc.at, nil)
					assert.Error(t, err, "%+v", r)
					r, err = fositeManager.GetRefreshTokenSession(ctx, tc.rt, nil)
					assert.Error(t, err, "%+v", r)
				})
			}

			require.NoError(t, m.RevokeSubjectConsentSession(ctx, "i-do-not-exist"))
			require.NoError(t, m.RevokeSubjectClientConsentSession(ctx, "i-do-not-exist", "i-do-not-exist"))
		})

		t.Run("case=list-used-consent-requests", func(t *testing.T) {
			_, hcr1, f1 := mockConsentRequest("rv1", true, 0, false, "fk-login-challenge", network)
			f1.NID = m.NetworkID(ctx)
			_, hcr2, f2 := mockConsentRequest("rv2", false, 0, false, "fk-login-challenge", network)
			f2.NID = m.NetworkID(ctx)

			// Ignore duplicate errors
			_ = clientManager.CreateClient(ctx, f1.Client)
			_ = clientManager.CreateClient(ctx, f2.Client)

			require.NoError(t, f1.HandleConsentRequest(hcr1))
			require.NoError(t, f2.HandleConsentRequest(hcr2))

			require.NoError(t, f1.InvalidateConsentRequest())
			require.NoError(t, f2.InvalidateConsentRequest())

			require.NoError(t, m.CreateConsentSession(ctx, f1))
			require.NoError(t, m.CreateConsentSession(ctx, f2))

			for i, tc := range []struct {
				subject    string
				sid        string
				consentIDs []string
				clients    []string
			}{
				{
					subject:    f1.Subject,
					sid:        makeID("fk-login-session", network, "rv1"),
					consentIDs: []string{f1.ConsentRequestID.String()},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    f2.Subject,
					sid:        makeID("fk-login-session", network, "rv2"),
					consentIDs: []string{f2.ConsentRequestID.String()},
					clients:    []string{"fk-client-rv2"},
				},
				{
					subject:    "subjectrv3",
					sid:        makeID("fk-login-session", network, "rv2"),
					consentIDs: []string{},
					clients:    []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s/session=%s", i, tc.subject, tc.sid), func(t *testing.T) {
					consents, nextPage, err := m.FindSubjectsSessionGrantedConsentRequests(ctx, tc.subject, tc.sid)
					if len(tc.consentIDs) == 0 {
						assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
					} else {
						require.NoError(t, err)
						require.Len(t, consents, len(tc.consentIDs))
						assert.True(t, nextPage.IsLast())
						for _, cs := range consents {
							assert.Contains(t, tc.consentIDs, cs.ConsentRequestID.String())
							assert.Contains(t, tc.clients, cs.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(ctx, tc.subject)
					require.NoError(t, err)
					assert.EqualValues(t, n, len(tc.consentIDs))

				})
			}

			for i, tc := range []struct {
				subject           string
				consentRequestIDs []string
				clients           []string
			}{
				{
					subject:           "subjectrv1",
					consentRequestIDs: []string{f1.ConsentRequestID.String()},
					clients:           []string{"fk-client-rv1"},
				},
				{
					subject:           "subjectrv2",
					consentRequestIDs: []string{f2.ConsentRequestID.String()},
					clients:           []string{"fk-client-rv2"},
				},
				{
					subject:           "subjectrv3",
					consentRequestIDs: []string{},
					clients:           []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					consents, nextPage, err := m.FindSubjectsGrantedConsentRequests(ctx, tc.subject)
					if len(tc.consentRequestIDs) == 0 {
						assert.EqualError(t, err, consent.ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						require.Len(t, consents, len(tc.consentRequestIDs))
						assert.True(t, nextPage.IsLast())
						for _, cs := range consents {
							assert.Contains(t, tc.consentRequestIDs, cs.ConsentRequestID.String())
							assert.Contains(t, tc.clients, cs.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(ctx, tc.subject)
					require.NoError(t, err)
					assert.EqualValues(t, n, len(tc.consentRequestIDs))

				})
			}

			t.Run("case=obfuscated", func(t *testing.T) {
				_, err := m.GetForcedObfuscatedLoginSession(ctx, "fk-client-1", "obfuscated-1")
				require.True(t, errors.Is(err, x.ErrNotFound))

				expect := &consent.ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-1",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(ctx, expect))

				got, err := m.GetForcedObfuscatedLoginSession(ctx, "fk-client-1", "obfuscated-1")
				require.NoError(t, err)
				require.NotEqualValues(t, got.NID, uuid.Nil)
				got.NID = uuid.Nil
				assert.EqualValues(t, expect, got)

				expect = &consent.ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-2",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(ctx, expect))

				got, err = m.GetForcedObfuscatedLoginSession(ctx, "fk-client-1", "obfuscated-2")
				require.NotEqualValues(t, got.NID, uuid.Nil)
				got.NID = uuid.Nil
				require.NoError(t, err)
				assert.EqualValues(t, expect, got)

				_, err = m.GetForcedObfuscatedLoginSession(ctx, "fk-client-1", "obfuscated-1")
				require.True(t, errors.Is(err, x.ErrNotFound))
			})

			t.Run("case=ListUserAuthenticatedClientsWithFrontAndBackChannelLogout", func(t *testing.T) {
				// The idea of this test is to create two identities (subjects) with 4 sessions each, where
				// only some sessions have been associated with a client that has a front channel logout url

				subjects := make([]string, 1)
				for k := range subjects {
					subjects[k] = fmt.Sprintf("subject-ListUserAuthenticatedClientsWithFrontAndBackChannelLogout-%d", k)
				}

				sessions := make([]flow.LoginSession, len(subjects)*1)
				frontChannels := map[string][]client.Client{}
				backChannels := map[string][]client.Client{}
				for k := range sessions {
					id := uuid.Must(uuid.NewV4()).String()
					subject := subjects[k%len(subjects)]
					t.Run(fmt.Sprintf("create/session=%s/subject=%s", id, subject), func(t *testing.T) {
						ls := &flow.LoginSession{
							ID:              id,
							NID:             m.NetworkID(ctx),
							AuthenticatedAt: sqlxx.NullTime(time.Now()),
							Subject:         subject,
							Remember:        true,
						}
						require.NoError(t, m.ConfirmLoginSession(ctx, ls))

						cl := &client.Client{ID: uuid.Must(uuid.NewV4()).String()}
						switch k % 4 {
						case 0:
							cl.FrontChannelLogoutURI = "http://some-url.com/"
							frontChannels[id] = append(frontChannels[id], *cl)
						case 1:
							cl.BackChannelLogoutURI = "http://some-url.com/"
							backChannels[id] = append(backChannels[id], *cl)
						case 2:
							cl.FrontChannelLogoutURI = "http://some-url.com/"
							cl.BackChannelLogoutURI = "http://some-url.com/"
							frontChannels[id] = append(frontChannels[id], *cl)
							backChannels[id] = append(backChannels[id], *cl)
						}
						require.NoError(t, clientManager.CreateClient(ctx, cl))

						f := &flow.Flow{
							NID: m.NetworkID(ctx),
							OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
								ACRValues: []string{"1", "2"},
								UILocales: []string{"fr", "de"},
								Display:   "popup",
							},
							RequestedAt:      time.Now().UTC().Add(-time.Hour),
							Client:           cl,
							Subject:          ls.Subject,
							RequestURL:       "https://request-url/path",
							LoginSkip:        true,
							RequestedScope:   []string{"scopea", "scopeb"},
							SessionID:        sqlxx.NullString(ls.ID),
							ConsentRequestID: sqlxx.NullString(uuid.Must(uuid.NewV4()).String()),
							ConsentVerifier:  sqlxx.NullString(uuid.Must(uuid.NewV4()).String()),
							ConsentCSRF:      sqlxx.NullString(uuid.Must(uuid.NewV4()).String()),

							LoginCSRF:     uuid.Must(uuid.NewV4()).String(),
							ID:            uuid.Must(uuid.NewV4()).String(),
							LoginVerifier: uuid.Must(uuid.NewV4()).String(),
							State:         flow.FlowStateConsentUnused,
						}
						_ = saneMockHandleConsentRequest(t, f, 0, false)

						require.NoError(t, f.InvalidateConsentRequest())
						require.NoError(t, m.CreateConsentSession(ctx, f))

						sessions[k] = *ls
					})
				}

				for _, ls := range sessions {
					check := func(t *testing.T, expected map[string][]client.Client, actual []client.Client) {
						es, ok := expected[ls.ID]
						if !ok {
							require.Len(t, actual, 0)
							return
						}
						require.Len(t, actual, len(es))

						for _, e := range es {
							var found bool
							for _, a := range actual {
								if e.GetID() == a.GetID() {
									found = true
								}
								assert.EqualValues(t, e.GetID(), a.GetID())
								assert.EqualValues(t, e.FrontChannelLogoutURI, a.FrontChannelLogoutURI)
								assert.EqualValues(t, e.BackChannelLogoutURI, a.BackChannelLogoutURI)
							}
							require.True(t, found)
						}
					}

					t.Run(fmt.Sprintf("method=ListUserAuthenticatedClientsWithFrontChannelLogout/session=%s/subject=%s", ls.ID, ls.Subject), func(t *testing.T) {
						actual, err := m.ListUserAuthenticatedClientsWithFrontChannelLogout(ctx, ls.Subject, ls.ID)
						require.NoError(t, err)
						check(t, frontChannels, actual)
					})

					t.Run(fmt.Sprintf("method=ListUserAuthenticatedClientsWithBackChannelLogout/session=%s", ls.ID), func(t *testing.T) {
						actual, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(ctx, ls.Subject, ls.ID)
						require.NoError(t, err)
						check(t, backChannels, actual)
					})
				}
			})

			t.Run("case=LogoutRequest", func(t *testing.T) {
				for k, tc := range []struct {
					key        string
					authAt     bool
					withClient bool
				}{
					{"LogoutRequest-1", true, true},
					{"LogoutRequest-2", true, true},
					{"LogoutRequest-3", true, true},
					{"LogoutRequest-4", true, true},
					{"LogoutRequest-5", true, false},
					{"LogoutRequest-6", false, false},
				} {
					t.Run("key="+tc.key, func(t *testing.T) {
						challenge := makeID("challenge", network, tc.key)
						verifier := makeID("verifier", network, tc.key)
						c := mockLogoutRequest(tc.key, tc.withClient, network)
						if tc.withClient {
							require.NoError(t, clientManager.CreateClient(ctx, c.Client)) // Ignore errors that are caused by duplication
						}

						_, err := m.GetLogoutRequest(ctx, challenge)
						require.Error(t, err)

						require.NoError(t, m.CreateLogoutRequest(ctx, c))

						got2, err := m.GetLogoutRequest(ctx, challenge)
						require.NoError(t, err)
						assert.False(t, got2.WasHandled)
						assert.False(t, got2.Accepted)
						compareLogoutRequest(t, c, got2)

						if k%2 == 0 {
							got2, err = m.AcceptLogoutRequest(ctx, challenge)
							require.NoError(t, err)
							assert.True(t, got2.Accepted)
							compareLogoutRequest(t, c, got2)

							got3, err := m.VerifyAndInvalidateLogoutRequest(ctx, verifier)
							require.NoError(t, err)
							assert.True(t, got3.Accepted)
							assert.True(t, got3.WasHandled)
							compareLogoutRequest(t, c, got3)

							_, err = m.VerifyAndInvalidateLogoutRequest(ctx, verifier)
							require.NoError(t, err)

							got2, err = m.GetLogoutRequest(ctx, challenge)
							require.NoError(t, err)
							compareLogoutRequest(t, got3, got2)
							assert.True(t, got2.WasHandled)
						} else {
							require.NoError(t, m.RejectLogoutRequest(ctx, challenge))
							_, err = m.GetLogoutRequest(ctx, challenge)
							require.Error(t, err)
						}
					})
				}
			})
		})

		t.Run("case=foreign key regression", func(t *testing.T) {
			cl := &client.Client{ID: uuid.Must(uuid.NewV4()).String()}
			require.NoError(t, clientManager.CreateClient(ctx, cl))

			subject := uuid.Must(uuid.NewV4()).String()
			s := flow.LoginSession{
				ID:              uuid.Must(uuid.NewV4()).String(),
				AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
				Subject:         subject,
			}

			require.NoError(t, m.ConfirmLoginSession(ctx, &s))

			_, err := m.DeleteLoginSession(ctx, s.ID)
			require.NoError(t, err)
		})
	}
}

func compareLogoutRequest(t *testing.T, a, b *flow.LogoutRequest) {
	require.True(t, (a.Client != nil && b.Client != nil) || (a.Client == nil && b.Client == nil))
	if a.Client != nil {
		assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	}

	assert.EqualValues(t, a.ID, b.ID)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.PostLogoutRedirectURI, b.PostLogoutRedirectURI)
	assert.EqualValues(t, a.RPInitiated, b.RPInitiated)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareAuthenticationRequest(t *testing.T, a, b *flow.LoginRequest) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.CSRF, b.CSRF)
	assert.EqualValues(t, a.Skip, b.Skip)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareAuthenticationRequestFlow(t *testing.T, a *flow.LoginRequest, b *flow.Flow) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.Verifier, b.LoginVerifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.CSRF, b.LoginCSRF)
	assert.EqualValues(t, a.Skip, b.LoginSkip)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareDeviceRequestFlow(t *testing.T, a *flow.DeviceUserAuthRequest, b *flow.Flow) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.CSRF, b.DeviceCSRF.String())
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.Verifier, b.DeviceVerifier.String())
	assert.EqualValues(t, a.HandledAt, b.DeviceHandledAt)
	assert.EqualValues(t, a.RequestedAudience, b.RequestedAudience)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
}

func compareConsentRequestFlow(t *testing.T, a *flow.OAuth2ConsentRequest, b *flow.Flow) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.ConsentRequestID, b.ConsentRequestID)
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.Skip, b.ConsentSkip)
	assert.EqualValues(t, a.LoginChallenge, b.ID)
	assert.EqualValues(t, a.LoginSessionID, b.SessionID)
	assert.EqualValues(t, a.DeviceChallenge, b.DeviceChallengeID.String())
}

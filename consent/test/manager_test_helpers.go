// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/oauth2"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/x/assertx"
	"github.com/ory/x/contextx"

	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlxx"

	"github.com/ory/fosite"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/x"
)

func MockConsentRequest(key string, remember bool, rememberFor int, hasError bool, skip bool, authAt bool, loginChallengeBase string, network string) (c *flow.OAuth2ConsentRequest, h *flow.AcceptOAuth2ConsentRequest, f *flow.Flow) {
	c = &flow.OAuth2ConsentRequest{
		ID:                makeID("challenge", network, key),
		RequestedScope:    []string{"scopea" + key, "scopeb" + key},
		RequestedAudience: []string{"auda" + key, "audb" + key},
		Skip:              skip,
		Subject:           "subject" + key,
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		Client:                 &client.Client{ID: "fk-client-" + key},
		RequestURL:             "https://request-url/path" + key,
		LoginChallenge:         sqlxx.NullString(makeID(loginChallengeBase, network, key)),
		LoginSessionID:         sqlxx.NullString(makeID("fk-login-session", network, key)),
		ForceSubjectIdentifier: "forced-subject",
		Verifier:               makeID("verifier", network, key),
		CSRF:                   "csrf" + key,
		ACR:                    "1",
		AuthenticatedAt:        sqlxx.NullTime(time.Now().UTC().Add(-time.Hour)),
		RequestedAt:            time.Now().UTC(),
		Context:                sqlxx.JSONRawMessage(`{"foo": "bar` + key + `"}`),
	}

	f = &flow.Flow{
		ID:                   c.LoginChallenge.String(),
		LoginVerifier:        makeID("login-verifier", network, key),
		SessionID:            c.LoginSessionID,
		Client:               c.Client,
		State:                flow.FlowStateConsentInitialized,
		ConsentChallengeID:   sqlxx.NullString(c.ID),
		ConsentSkip:          c.Skip,
		ConsentVerifier:      sqlxx.NullString(c.Verifier),
		ConsentCSRF:          sqlxx.NullString(c.CSRF),
		OpenIDConnectContext: c.OpenIDConnectContext,
		Subject:              c.Subject,
		RequestedScope:       c.RequestedScope,
		RequestedAudience:    c.RequestedAudience,
		RequestURL:           c.RequestURL,
		RequestedAt:          c.RequestedAt,
	}

	var err *flow.RequestDeniedError
	if hasError {
		err = &flow.RequestDeniedError{
			Name:        "error_name" + key,
			Description: "error_description" + key,
			Hint:        "error_hint,omitempty" + key,
			Code:        100,
			Debug:       "error_debug,omitempty" + key,
			Valid:       true,
		}
	}

	var authenticatedAt sqlxx.NullTime
	if authAt {
		authenticatedAt = sqlxx.NullTime(time.Now().UTC().Add(-time.Minute))
	}

	h = &flow.AcceptOAuth2ConsentRequest{
		ConsentRequest:  c,
		RememberFor:     rememberFor,
		Remember:        remember,
		ID:              makeID("challenge", network, key),
		RequestedAt:     time.Now().UTC().Add(-time.Minute),
		AuthenticatedAt: authenticatedAt,
		GrantedScope:    []string{"scopea" + key, "scopeb" + key},
		GrantedAudience: []string{"auda" + key, "audb" + key},
		Error:           err,
		HandledAt:       sqlxx.NullTime(time.Now().UTC()),
		// WasUsed:         true,
	}

	return c, h, f
}

func MockLogoutRequest(key string, withClient bool, network string) (c *flow.LogoutRequest) {
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

func MockAuthRequest(key string, authAt bool, network string) (c *flow.LoginRequest, h *flow.HandledLoginRequest, f *flow.Flow) {
	c = &flow.LoginRequest{
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		RequestedAt:    time.Now().UTC().Add(-time.Minute),
		Client:         &client.Client{ID: "fk-client-" + key},
		Subject:        "subject" + key,
		RequestURL:     "https://request-url/path" + key,
		Skip:           true,
		ID:             makeID("challenge", network, key),
		Verifier:       makeID("verifier", network, key),
		RequestedScope: []string{"scopea" + key, "scopeb" + key},
		CSRF:           "csrf" + key,
		SessionID:      sqlxx.NullString(makeID("fk-login-session", network, key)),
	}

	f = flow.NewFlow(c)

	var err = &flow.RequestDeniedError{
		Name:        "error_name" + key,
		Description: "error_description" + key,
		Hint:        "error_hint,omitempty" + key,
		Code:        100,
		Debug:       "error_debug,omitempty" + key,
		Valid:       true,
	}

	var authenticatedAt time.Time
	if authAt {
		authenticatedAt = time.Now().UTC().Add(-time.Minute)
	}

	h = &flow.HandledLoginRequest{
		LoginRequest:           c,
		RememberFor:            120,
		Remember:               true,
		ID:                     makeID("challenge", network, key),
		RequestedAt:            time.Now().UTC().Add(-time.Minute),
		AuthenticatedAt:        sqlxx.NullTime(authenticatedAt),
		Error:                  err,
		Subject:                c.Subject,
		ACR:                    "acr",
		ForceSubjectIdentifier: "forced-subject",
		WasHandled:             false,
	}

	return c, h, f
}

func SaneMockHandleConsentRequest(t *testing.T, m consent.Manager, f *flow.Flow, c *flow.OAuth2ConsentRequest, authAt time.Time, rememberFor int, remember bool, hasError bool) *flow.AcceptOAuth2ConsentRequest {
	var rde *flow.RequestDeniedError
	if hasError {
		rde = &flow.RequestDeniedError{
			Name:        "error_name",
			Description: "error_description",
			Hint:        "error_hint",
			Code:        100,
			Debug:       "error_debug",
			Valid:       true,
		}
	}

	h := &flow.AcceptOAuth2ConsentRequest{
		ConsentRequest:  c,
		RememberFor:     rememberFor,
		Remember:        remember,
		ID:              c.ID,
		RequestedAt:     time.Now().UTC().Add(-time.Minute),
		AuthenticatedAt: sqlxx.NullTime(authAt),
		GrantedScope:    []string{"scopea", "scopeb"},
		GrantedAudience: []string{"auda", "audb"},
		Error:           rde,
		WasHandled:      false,
		HandledAt:       sqlxx.NullTime(time.Now().UTC().Add(-time.Minute)),
	}

	_, err := m.HandleConsentRequest(context.Background(), f, h)
	require.NoError(t, err)

	return h
}

// SaneMockConsentRequest does the same thing as MockConsentRequest but uses less insanity and implicit dependencies.
func SaneMockConsentRequest(t *testing.T, m consent.Manager, f *flow.Flow, skip bool) (c *flow.OAuth2ConsentRequest) {
	c = &flow.OAuth2ConsentRequest{
		RequestedScope:    []string{"scopea", "scopeb"},
		RequestedAudience: []string{"auda", "audb"},
		Skip:              skip,
		Subject:           f.Subject,
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1", "2"},
			UILocales: []string{"fr", "de"},
			Display:   "popup",
		},
		Client:                 f.Client,
		RequestURL:             "https://request-url/path",
		LoginChallenge:         sqlxx.NullString(f.ID),
		LoginSessionID:         f.SessionID,
		ForceSubjectIdentifier: "forced-subject",
		ACR:                    "1",
		AuthenticatedAt:        sqlxx.NullTime(time.Now().UTC().Add(-time.Hour)),
		RequestedAt:            time.Now().UTC().Add(-time.Hour),
		Context:                sqlxx.JSONRawMessage(`{"foo": "bar"}`),

		ID:       uuid.New().String(),
		Verifier: uuid.New().String(),
		CSRF:     uuid.New().String(),
	}

	require.NoError(t, m.CreateConsentRequest(context.Background(), f, c))

	return c
}

// SaneMockAuthRequest does the same thing as MockAuthRequest but uses less insanity and implicit dependencies.
func SaneMockAuthRequest(t *testing.T, m consent.Manager, ls *flow.LoginSession, cl *client.Client) (c *flow.LoginRequest) {
	c = &flow.LoginRequest{
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1", "2"},
			UILocales: []string{"fr", "de"},
			Display:   "popup",
		},
		RequestedAt:    time.Now().UTC().Add(-time.Hour),
		Client:         cl,
		Subject:        ls.Subject,
		RequestURL:     "https://request-url/path",
		Skip:           true,
		RequestedScope: []string{"scopea", "scopeb"},
		SessionID:      sqlxx.NullString(ls.ID),

		CSRF:     uuid.New().String(),
		ID:       uuid.New().String(),
		Verifier: uuid.New().String(),
	}
	_, err := m.CreateLoginRequest(context.Background(), c)
	require.NoError(t, err)
	return c
}

func makeID(base string, network string, key string) string {
	return fmt.Sprintf("%s-%s-%s", base, network, key)
}

func TestHelperNID(r interface {
	client.ManagerProvider
	FlowCipher() *aead.XChaCha20Poly1305
}, t1ValidNID consent.Manager, t2InvalidNID consent.Manager) func(t *testing.T) {
	testClient := client.Client{ID: "2022-03-11-client-nid-test-1"}
	testLS := flow.LoginSession{
		ID:      "2022-03-11-ls-nid-test-1",
		Subject: "2022-03-11-test-1-sub",
	}
	testLR := flow.LoginRequest{
		ID:          "2022-03-11-lr-nid-test-1",
		Subject:     "2022-03-11-test-1-sub",
		Verifier:    "2022-03-11-test-1-ver",
		RequestedAt: time.Now(),
		Client:      &client.Client{ID: "2022-03-11-client-nid-test-1"},
	}
	testHLR := flow.HandledLoginRequest{
		LoginRequest:           &testLR,
		RememberFor:            120,
		Remember:               true,
		ID:                     testLR.ID,
		RequestedAt:            testLR.RequestedAt,
		AuthenticatedAt:        sqlxx.NullTime(time.Now()),
		Error:                  nil,
		Subject:                testLR.Subject,
		ACR:                    "acr",
		ForceSubjectIdentifier: "2022-03-11-test-1-forced-sub",
		WasHandled:             false,
	}

	return func(t *testing.T) {
		ctx := context.Background()
		require.NoError(t, r.ClientManager().CreateClient(ctx, &testClient))
		require.Error(t, t2InvalidNID.CreateLoginSession(ctx, &testLS))
		require.NoError(t, t1ValidNID.CreateLoginSession(ctx, &testLS))

		_, err := t2InvalidNID.CreateLoginRequest(ctx, &testLR)
		require.Error(t, err)
		f, err := t1ValidNID.CreateLoginRequest(ctx, &testLR)
		require.NoError(t, err)

		testLR.ID = x.Must(f.ToLoginChallenge(ctx, r))
		_, err = t2InvalidNID.GetLoginRequest(ctx, testLR.ID)
		require.Error(t, err)
		_, err = t1ValidNID.GetLoginRequest(ctx, testLR.ID)
		require.NoError(t, err)
		_, err = t2InvalidNID.HandleLoginRequest(ctx, f, testLR.ID, &testHLR)
		require.Error(t, err)
		_, err = t1ValidNID.HandleLoginRequest(ctx, f, testLR.ID, &testHLR)
		require.NoError(t, err)
		require.Error(t, t2InvalidNID.ConfirmLoginSession(ctx, &testLS))
		require.NoError(t, t1ValidNID.ConfirmLoginSession(ctx, &testLS))
		ls, err := t2InvalidNID.DeleteLoginSession(ctx, testLS.ID)
		require.Error(t, err)
		assert.Nil(t, ls)
		ls, err = t1ValidNID.DeleteLoginSession(ctx, testLS.ID)
		require.NoError(t, err)
		assert.Equal(t, testLS.ID, ls.ID)
	}
}

type Deps interface {
	FlowCipher() *aead.XChaCha20Poly1305
	contextx.Provider
}

func ManagerTests(deps Deps, m consent.Manager, clientManager client.Manager, fositeManager x.FositeStorer, network string, parallel bool) func(t *testing.T) {
	lr := make(map[string]*flow.LoginRequest)

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
					AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
					Subject:         fmt.Sprintf("subject-%s", k),
				}
				require.NoError(t, m.CreateLoginSession(ctx, loginSession))
				require.NoError(t, m.ConfirmLoginSession(ctx, loginSession))

				lr[k] = &flow.LoginRequest{
					ID:              makeID("fk-login-challenge", network, k),
					Subject:         fmt.Sprintf("subject%s", k),
					SessionID:       sqlxx.NullString(makeID("fk-login-session", network, k)),
					Verifier:        makeID("fk-login-verifier", network, k),
					Client:          &client.Client{ID: fmt.Sprintf("fk-client-%s", k)},
					AuthenticatedAt: sqlxx.NullTime(time.Now()),
					RequestedAt:     time.Now(),
				}

				_, err := m.CreateLoginRequest(ctx, lr[k])
				require.NoError(t, err)
			}
		})

		t.Run("case=auth-session", func(t *testing.T) {
			for _, tc := range []struct {
				s flow.LoginSession
			}{
				{
					s: flow.LoginSession{
						ID:              makeID("session", network, "1"),
						AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-time.Minute).UTC()),
						Subject:         "subject1",
					},
				},
				{
					s: flow.LoginSession{
						ID:              makeID("session", network, "2"),
						AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
						Subject:         "subject2",
					},
				},
			} {
				tc := tc
				t.Run("case=create-get-"+tc.s.ID, func(t *testing.T) {
					_, err := m.GetRememberedLoginSession(ctx, &tc.s, tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error(), "%#v", err)

					err = m.CreateLoginSession(ctx, &tc.s)
					require.NoError(t, err)

					_, err = m.GetRememberedLoginSession(ctx, &tc.s, tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error())

					updatedAuth := time.Time(tc.s.AuthenticatedAt).Add(time.Second)
					tc.s.AuthenticatedAt = sqlxx.NullTime(updatedAuth)
					require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
						ID:              tc.s.ID,
						AuthenticatedAt: sqlxx.NullTime(updatedAuth),
						Subject:         tc.s.Subject,
						Remember:        true,
					}))

					got, err := m.GetRememberedLoginSession(ctx, nil, tc.s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, tc.s.ID, got.ID)
					assert.Equal(t, tc.s.AuthenticatedAt, got.AuthenticatedAt) // this was updated from confirm...
					assert.EqualValues(t, tc.s.Subject, got.Subject)

					// Make sure AuthAt does not equal...
					updatedAuth2 := updatedAuth.Add(1 * time.Second).UTC()
					require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
						ID:              tc.s.ID,
						AuthenticatedAt: sqlxx.NullTime(updatedAuth2),
						Subject:         "some-other-subject",
						Remember:        true,
					}))

					got2, err := m.GetRememberedLoginSession(ctx, nil, tc.s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, tc.s.ID, got2.ID)
					assert.Equal(t, updatedAuth2.Unix(), time.Time(got2.AuthenticatedAt).Unix()) // this was updated from confirm...
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

					_, err = m.GetRememberedLoginSession(ctx, nil, tc.id)
					require.Error(t, err)
				})
			}
		})

		t.Run("case=auth-request", func(t *testing.T) {
			for _, tc := range []struct {
				key    string
				authAt bool
			}{
				{"1", true},
				{"2", true},
				{"3", true},
				{"4", true},
				{"5", true},
				{"6", false},
				{"7", true},
			} {
				t.Run("key="+tc.key, func(t *testing.T) {
					c, h, f := MockAuthRequest(tc.key, tc.authAt, network)
					_ = clientManager.CreateClient(ctx, c.Client) // Ignore errors that are caused by duplication
					loginChallenge := x.Must(f.ToLoginChallenge(ctx, deps))

					_, err := m.GetLoginRequest(ctx, loginChallenge)
					require.Error(t, err)

					f, err = m.CreateLoginRequest(ctx, c)
					require.NoError(t, err)

					loginChallenge = x.Must(f.ToLoginChallenge(ctx, deps))

					got1, err := m.GetLoginRequest(ctx, loginChallenge)
					require.NoError(t, err)
					assert.False(t, got1.WasHandled)
					compareAuthenticationRequest(t, c, got1)

					got1, err = m.HandleLoginRequest(ctx, f, loginChallenge, h)
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got1)

					loginVerifier := x.Must(f.ToLoginVerifier(ctx, deps))

					got2, err := m.VerifyAndInvalidateLoginRequest(ctx, loginVerifier)
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got2.LoginRequest)

					loginChallenge = x.Must(f.ToLoginChallenge(ctx, deps))
					got1, err = m.GetLoginRequest(ctx, loginChallenge)
					require.NoError(t, err)
				})
			}
		})

		t.Run("case=consent-request", func(t *testing.T) {
			for _, tc := range []struct {
				key         string
				remember    bool
				rememberFor int
				hasError    bool
				skip        bool
				authAt      bool
			}{
				{"1", true, 0, false, false, true},
				{"2", true, 0, true, false, true},
				{"3", true, 1, false, false, true},
				{"4", false, 0, false, false, true},
				{"5", true, 120, false, false, true},
				{"6", true, 120, false, true, true},
				{"7", false, 0, false, false, false},
			} {
				t.Run("key="+tc.key, func(t *testing.T) {
					consentRequest, h, f := MockConsentRequest(tc.key, tc.remember, tc.rememberFor, tc.hasError, tc.skip, tc.authAt, "challenge", network)
					_ = clientManager.CreateClient(ctx, consentRequest.Client) // Ignore errors that are caused by duplication
					f.NID = deps.Contextualizer().Network(context.Background(), gofrsuuid.Nil)

					consentChallenge := makeID("challenge", network, tc.key)

					_, err := m.GetConsentRequest(ctx, consentChallenge)
					require.Error(t, err)

					consentChallenge = x.Must(f.ToConsentChallenge(ctx, deps))
					consentRequest.ID = consentChallenge

					err = m.CreateConsentRequest(ctx, f, consentRequest)
					require.NoError(t, err)

					got1, err := m.GetConsentRequest(ctx, consentChallenge)
					require.NoError(t, err)
					compareConsentRequest(t, consentRequest, got1)
					assert.False(t, got1.WasHandled)

					got1, err = m.HandleConsentRequest(ctx, f, h)
					require.NoError(t, err)
					assertx.TimeDifferenceLess(t, time.Now(), time.Time(h.HandledAt), 5)
					compareConsentRequest(t, consentRequest, got1)

					h.GrantedAudience = sqlxx.StringSliceJSONFormat{"new-audience"}
					_, err = m.HandleConsentRequest(ctx, f, h)
					require.NoError(t, err)

					consentVerifier := x.Must(f.ToConsentVerifier(ctx, deps))

					got2, err := m.VerifyAndInvalidateConsentRequest(ctx, consentVerifier)
					require.NoError(t, err)
					consentRequest.ID = got2.ID
					compareConsentRequest(t, consentRequest, got2.ConsentRequest)
					assert.Equal(t, consentRequest.ID, got2.ID)
					assert.Equal(t, h.GrantedAudience, got2.GrantedAudience)

					t.Run("sub=detect double-submit for consent verifier", func(t *testing.T) {
						_, err := m.VerifyAndInvalidateConsentRequest(ctx, consentVerifier)
						require.Error(t, err)
					})

					if tc.hasError {
						assert.True(t, got2.HasError())
					}
					assert.Equal(t, tc.remember, got2.Remember)
					assert.Equal(t, tc.rememberFor, got2.RememberFor)
				})
			}

			for _, tc := range []struct {
				keyC           string
				keyS           string
				expectedLength int
			}{
				{"1", "1", 1},
				{"2", "2", 0},
				// {"3", "3", 0},  // Some consent is given in some other test case. Yay global fixtues :)
				{"4", "4", 0},
				{"1", "2", 0},
				{"2", "1", 0},
				{"5", "5", 1},
				{"6", "6", 0},
			} {
				t.Run("key="+tc.keyC+"-"+tc.keyS, func(t *testing.T) {
					rs, err := m.FindGrantedAndRememberedConsentRequests(ctx, "fk-client-"+tc.keyC, "subject"+tc.keyS)
					if tc.expectedLength == 0 {
						assert.Nil(t, rs)
						assert.EqualError(t, err, consent.ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						assert.Len(t, rs, tc.expectedLength)
					}
				})
			}
		})

		t.Run("case=revoke-auth-request", func(t *testing.T) {
			require.NoError(t, m.CreateLoginSession(ctx, &flow.LoginSession{
				ID:              makeID("rev-session", network, "-1"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-1",
			}))

			require.NoError(t, m.CreateLoginSession(ctx, &flow.LoginSession{
				ID:              makeID("rev-session", network, "-2"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-2",
			}))

			require.NoError(t, m.CreateLoginSession(ctx, &flow.LoginSession{
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
							_, err := m.GetRememberedLoginSession(ctx, nil, id)
							assert.EqualError(t, err, x.ErrNotFound.Error())
						})
					}
				})
			}
		})

		challengerv1 := makeID("challenge", network, "rv1")
		challengerv2 := makeID("challenge", network, "rv2")
		t.Run("case=revoke-used-consent-request", func(t *testing.T) {

			cr1, hcr1, f1 := MockConsentRequest("rv1", false, 0, false, false, false, "fk-login-challenge", network)
			cr2, hcr2, f2 := MockConsentRequest("rv2", false, 0, false, false, false, "fk-login-challenge", network)
			f1.NID = deps.Contextualizer().Network(context.Background(), gofrsuuid.Nil)
			f2.NID = deps.Contextualizer().Network(context.Background(), gofrsuuid.Nil)

			// Ignore duplication errors
			_ = clientManager.CreateClient(ctx, cr1.Client)
			_ = clientManager.CreateClient(ctx, cr2.Client)

			err := m.CreateConsentRequest(ctx, f1, cr1)
			require.NoError(t, err)
			err = m.CreateConsentRequest(ctx, f2, cr2)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(ctx, f1, hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(ctx, f2, hcr2)
			require.NoError(t, err)

			crr1, err := m.VerifyAndInvalidateConsentRequest(ctx, x.Must(f1.ToConsentVerifier(ctx, deps)))
			require.NoError(t, err)
			crr2, err := m.VerifyAndInvalidateConsentRequest(ctx, x.Must(f2.ToConsentVerifier(ctx, deps)))
			require.NoError(t, err)

			require.NoError(t, fositeManager.CreateAccessTokenSession(
				ctx,
				makeID("", network, "trva1"),
				&fosite.Request{Client: cr1.Client, ID: crr1.ID, RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(
				ctx,
				makeID("", network, "rrva1"),
				"",
				&fosite.Request{Client: cr1.Client, ID: crr1.ID, RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateAccessTokenSession(
				ctx,
				makeID("", network, "trva2"),
				&fosite.Request{Client: cr2.Client, ID: crr2.ID, RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(
				ctx,
				makeID("", network, "rrva2"),
				"",
				&fosite.Request{Client: cr2.Client, ID: crr2.ID, RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
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

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetConsentRequest(ctx, id)
							assert.True(t, errors.Is(err, x.ErrNotFound))
						})
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
			f1, err := m.CreateLoginRequest(ctx, lr["rv1"])
			require.NoError(t, err)
			f2, err := m.CreateLoginRequest(ctx, lr["rv2"])
			require.NoError(t, err)

			cr1, hcr1, _ := MockConsentRequest("rv1", true, 0, false, false, false, "fk-login-challenge", network)
			cr2, hcr2, _ := MockConsentRequest("rv2", false, 0, false, false, false, "fk-login-challenge", network)

			// Ignore duplicate errors
			_ = clientManager.CreateClient(ctx, cr1.Client)
			_ = clientManager.CreateClient(ctx, cr2.Client)

			err = m.CreateConsentRequest(ctx, f1, cr1)
			require.NoError(t, err)
			err = m.CreateConsentRequest(ctx, f2, cr2)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(ctx, f1, hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(ctx, f2, hcr2)
			require.NoError(t, err)
			handledConsentRequest1, err := m.VerifyAndInvalidateConsentRequest(ctx, x.Must(f1.ToConsentVerifier(ctx, deps)))
			require.NoError(t, err)
			handledConsentRequest2, err := m.VerifyAndInvalidateConsentRequest(ctx, x.Must(f2.ToConsentVerifier(ctx, deps)))
			require.NoError(t, err)

			for i, tc := range []struct {
				subject    string
				sid        string
				challenges []string
				clients    []string
			}{
				{
					subject:    cr1.Subject,
					sid:        makeID("fk-login-session", network, "rv1"),
					challenges: []string{handledConsentRequest1.ID},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    cr2.Subject,
					sid:        makeID("fk-login-session", network, "rv2"),
					challenges: []string{handledConsentRequest2.ID},
					clients:    []string{"fk-client-rv2"},
				},
				{
					subject:    "subjectrv3",
					sid:        makeID("fk-login-session", network, "rv2"),
					challenges: []string{},
					clients:    []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s/session=%s", i, tc.subject, tc.sid), func(t *testing.T) {
					consents, err := m.FindSubjectsSessionGrantedConsentRequests(ctx, tc.subject, tc.sid, 100, 0)
					assert.Equal(t, len(tc.challenges), len(consents))

					if len(tc.challenges) == 0 {
						assert.EqualError(t, err, consent.ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						for _, consent := range consents {
							assert.Contains(t, tc.challenges, consent.ID)
							assert.Contains(t, tc.clients, consent.ConsentRequest.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(ctx, tc.subject)
					require.NoError(t, err)
					assert.Equal(t, n, len(tc.challenges))

				})
			}

			for i, tc := range []struct {
				subject    string
				challenges []string
				clients    []string
			}{
				{
					subject:    "subjectrv1",
					challenges: []string{handledConsentRequest1.ID},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    "subjectrv2",
					challenges: []string{handledConsentRequest2.ID},
					clients:    []string{"fk-client-rv2"},
				},
				{
					subject:    "subjectrv3",
					challenges: []string{},
					clients:    []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					consents, err := m.FindSubjectsGrantedConsentRequests(ctx, tc.subject, 100, 0)
					assert.Equal(t, len(tc.challenges), len(consents))

					if len(tc.challenges) == 0 {
						assert.EqualError(t, err, consent.ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						for _, consent := range consents {
							assert.Contains(t, tc.challenges, consent.ID)
							assert.Contains(t, tc.clients, consent.ConsentRequest.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(ctx, tc.subject)
					require.NoError(t, err)
					assert.Equal(t, n, len(tc.challenges))

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
				require.NotEqual(t, got.NID, gofrsuuid.Nil)
				got.NID = gofrsuuid.Nil
				assert.EqualValues(t, expect, got)

				expect = &consent.ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-2",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(ctx, expect))

				got, err = m.GetForcedObfuscatedLoginSession(ctx, "fk-client-1", "obfuscated-2")
				require.NotEqual(t, got.NID, gofrsuuid.Nil)
				got.NID = gofrsuuid.Nil
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
					id := uuid.New().String()
					subject := subjects[k%len(subjects)]
					t.Run(fmt.Sprintf("create/session=%s/subject=%s", id, subject), func(t *testing.T) {
						ls := &flow.LoginSession{
							ID:              id,
							AuthenticatedAt: sqlxx.NullTime(time.Now()),
							Subject:         subject,
						}
						require.NoError(t, m.CreateLoginSession(ctx, ls))
						ls.Remember = true
						require.NoError(t, m.ConfirmLoginSession(ctx, ls))

						cl := &client.Client{ID: uuid.New().String()}
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

						ar := SaneMockAuthRequest(t, m, ls, cl)
						f := flow.NewFlow(ar)
						f.NID = deps.Contextualizer().Network(ctx, gofrsuuid.Nil)
						cr := SaneMockConsentRequest(t, m, f, false)
						_ = SaneMockHandleConsentRequest(t, m, f, cr, time.Time{}, 0, false, false)
						_, err = m.VerifyAndInvalidateConsentRequest(ctx, x.Must(f.ToConsentVerifier(ctx, deps)))
						require.NoError(t, err)

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
								assert.Equal(t, e.GetID(), a.GetID())
								assert.Equal(t, e.FrontChannelLogoutURI, a.FrontChannelLogoutURI)
								assert.Equal(t, e.BackChannelLogoutURI, a.BackChannelLogoutURI)
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
						c := MockLogoutRequest(tc.key, tc.withClient, network)
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
			cl := &client.Client{ID: uuid.New().String()}
			require.NoError(t, clientManager.CreateClient(ctx, cl))

			subject := uuid.New().String()
			s := flow.LoginSession{
				ID:              uuid.New().String(),
				AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
				Subject:         subject,
			}

			require.NoError(t, m.CreateLoginSession(ctx, &s))
			require.NoError(t, m.ConfirmLoginSession(ctx, &s))

			lr := &flow.LoginRequest{
				ID:              uuid.New().String(),
				Subject:         uuid.New().String(),
				Verifier:        uuid.New().String(),
				Client:          cl,
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				RequestedAt:     time.Now(),
				SessionID:       sqlxx.NullString(s.ID),
			}

			f, err := m.CreateLoginRequest(ctx, lr)
			require.NoError(t, err)
			expected := &flow.OAuth2ConsentRequest{
				ID:                   x.Must(f.ToConsentChallenge(ctx, deps)),
				Skip:                 true,
				Subject:              subject,
				OpenIDConnectContext: nil,
				Client:               cl,
				ClientID:             cl.ID,
				RequestURL:           "",
				LoginChallenge:       sqlxx.NullString(lr.ID),
				LoginSessionID:       sqlxx.NullString(s.ID),
				Verifier:             uuid.New().String(),
				CSRF:                 uuid.New().String(),
			}
			err = m.CreateConsentRequest(ctx, f, expected)
			require.NoError(t, err)

			result, err := m.GetConsentRequest(ctx, expected.ID)
			require.NoError(t, err)
			assert.EqualValues(t, expected.ID, result.ID)

			_, err = m.DeleteLoginSession(ctx, s.ID)
			require.NoError(t, err)

			result, err = m.GetConsentRequest(ctx, expected.ID)
			require.NoError(t, err)
			assert.EqualValues(t, expected.ID, result.ID)
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

func compareConsentRequest(t *testing.T, a, b *flow.OAuth2ConsentRequest) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.ID, b.ID)
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.CSRF, b.CSRF)
	assert.EqualValues(t, a.Skip, b.Skip)
	assert.EqualValues(t, a.LoginChallenge, b.LoginChallenge)
	assert.EqualValues(t, a.LoginSessionID, b.LoginSessionID)
}

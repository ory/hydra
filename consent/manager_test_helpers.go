// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ory/x/assertx"

	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlxx"

	"github.com/ory/fosite"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/x"
)

func MockConsentRequest(key string, remember bool, rememberFor int, hasError bool, skip bool, authAt bool, loginChallengeBase string, network string) (c *OAuth2ConsentRequest, h *AcceptOAuth2ConsentRequest) {
	c = &OAuth2ConsentRequest{
		ID:                makeID("challenge", network, key),
		RequestedScope:    []string{"scopea" + key, "scopeb" + key},
		RequestedAudience: []string{"auda" + key, "audb" + key},
		Skip:              skip,
		Subject:           "subject" + key,
		OpenIDConnectContext: &OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		Client:                 &client.Client{LegacyClientID: "fk-client-" + key},
		RequestURL:             "https://request-url/path" + key,
		LoginChallenge:         sqlxx.NullString(makeID(loginChallengeBase, network, key)),
		LoginSessionID:         sqlxx.NullString(makeID("fk-login-session", network, key)),
		ForceSubjectIdentifier: "forced-subject",
		Verifier:               makeID("verifier", network, key),
		CSRF:                   "csrf" + key,
		ACR:                    "1",
		AuthenticatedAt:        sqlxx.NullTime(time.Now().UTC().Add(-time.Hour)),
		RequestedAt:            time.Now().UTC().Add(-time.Hour),
		Context:                sqlxx.JSONRawMessage(`{"foo": "bar` + key + `"}`),
	}

	var err *RequestDeniedError
	if hasError {
		err = &RequestDeniedError{
			Name:        "error_name" + key,
			Description: "error_description" + key,
			Hint:        "error_hint,omitempty" + key,
			Code:        100,
			Debug:       "error_debug,omitempty" + key,
			valid:       true,
		}
	}

	var authenticatedAt sqlxx.NullTime
	if authAt {
		authenticatedAt = sqlxx.NullTime(time.Now().UTC().Add(-time.Minute))
	}

	h = &AcceptOAuth2ConsentRequest{
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

	return c, h
}

func MockLogoutRequest(key string, withClient bool, network string) (c *LogoutRequest) {
	var cl *client.Client
	if withClient {
		cl = &client.Client{
			LegacyClientID: "fk-client-" + key,
		}
	}
	return &LogoutRequest{
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

func MockAuthRequest(key string, authAt bool, network string) (c *LoginRequest, h *HandledLoginRequest) {
	c = &LoginRequest{
		OpenIDConnectContext: &OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		RequestedAt:    time.Now().UTC().Add(-time.Minute),
		Client:         &client.Client{LegacyClientID: "fk-client-" + key},
		Subject:        "subject" + key,
		RequestURL:     "https://request-url/path" + key,
		Skip:           true,
		ID:             makeID("challenge", network, key),
		Verifier:       makeID("verifier", network, key),
		RequestedScope: []string{"scopea" + key, "scopeb" + key},
		CSRF:           "csrf" + key,
		SessionID:      sqlxx.NullString(makeID("fk-login-session", network, key)),
	}

	var err = &RequestDeniedError{
		Name:        "error_name" + key,
		Description: "error_description" + key,
		Hint:        "error_hint,omitempty" + key,
		Code:        100,
		Debug:       "error_debug,omitempty" + key,
		valid:       true,
	}

	var authenticatedAt time.Time
	if authAt {
		authenticatedAt = time.Now().UTC().Add(-time.Minute)
	}

	h = &HandledLoginRequest{
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

	return c, h
}

func SaneMockHandleConsentRequest(t *testing.T, m Manager, c *OAuth2ConsentRequest, authAt time.Time, rememberFor int, remember bool, hasError bool) *AcceptOAuth2ConsentRequest {
	var rde *RequestDeniedError
	if hasError {
		rde = &RequestDeniedError{
			Name:        "error_name",
			Description: "error_description",
			Hint:        "error_hint",
			Code:        100,
			Debug:       "error_debug",
			valid:       true,
		}
	}

	h := &AcceptOAuth2ConsentRequest{
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

	_, err := m.HandleConsentRequest(context.Background(), h)
	require.NoError(t, err)
	return h
}

// SaneMockConsentRequest does the same thing as MockConsentRequest but uses less insanity and implicit dependencies.
func SaneMockConsentRequest(t *testing.T, m Manager, ar *LoginRequest, skip bool) (c *OAuth2ConsentRequest) {
	c = &OAuth2ConsentRequest{
		RequestedScope:    []string{"scopea", "scopeb"},
		RequestedAudience: []string{"auda", "audb"},
		Skip:              skip,
		Subject:           ar.Subject,
		OpenIDConnectContext: &OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1", "2"},
			UILocales: []string{"fr", "de"},
			Display:   "popup",
		},
		Client:                 ar.Client,
		RequestURL:             "https://request-url/path",
		LoginChallenge:         sqlxx.NullString(ar.ID),
		LoginSessionID:         ar.SessionID,
		ForceSubjectIdentifier: "forced-subject",
		ACR:                    "1",
		AuthenticatedAt:        sqlxx.NullTime(time.Now().UTC().Add(-time.Hour)),
		RequestedAt:            time.Now().UTC().Add(-time.Hour),
		Context:                sqlxx.JSONRawMessage(`{"foo": "bar"}`),

		ID:       uuid.New().String(),
		Verifier: uuid.New().String(),
		CSRF:     uuid.New().String(),
	}

	require.NoError(t, m.CreateConsentRequest(context.Background(), c))
	return c
}

// SaneMockAuthRequest does the same thing as MockAuthRequest but uses less insanity and implicit dependencies.
func SaneMockAuthRequest(t *testing.T, m Manager, ls *LoginSession, cl *client.Client) (c *LoginRequest) {
	c = &LoginRequest{
		OpenIDConnectContext: &OAuth2ConsentRequestOpenIDConnectContext{
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
	require.NoError(t, m.CreateLoginRequest(context.Background(), c))
	return c
}

func makeID(base string, network string, key string) string {
	return fmt.Sprintf("%s-%s-%s", base, network, key)
}

func TestHelperNID(t1ClientManager client.Manager, t1ValidNID Manager, t2InvalidNID Manager) func(t *testing.T) {
	testClient := client.Client{LegacyClientID: "2022-03-11-client-nid-test-1"}
	testLS := LoginSession{
		ID:      "2022-03-11-ls-nid-test-1",
		Subject: "2022-03-11-test-1-sub",
	}
	testLR := LoginRequest{
		ID:          "2022-03-11-lr-nid-test-1",
		Subject:     "2022-03-11-test-1-sub",
		Verifier:    "2022-03-11-test-1-ver",
		RequestedAt: time.Now(),
		Client:      &client.Client{LegacyClientID: "2022-03-11-client-nid-test-1"},
	}
	testHLR := HandledLoginRequest{
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
		require.NoError(t, t1ClientManager.CreateClient(context.Background(), &testClient))
		require.Error(t, t2InvalidNID.CreateLoginSession(context.Background(), &testLS))
		require.NoError(t, t1ValidNID.CreateLoginSession(context.Background(), &testLS))
		require.Error(t, t2InvalidNID.CreateLoginRequest(context.Background(), &testLR))
		require.NoError(t, t1ValidNID.CreateLoginRequest(context.Background(), &testLR))
		_, err := t2InvalidNID.GetLoginRequest(context.Background(), testLR.ID)
		require.Error(t, err)
		_, err = t1ValidNID.GetLoginRequest(context.Background(), testLR.ID)
		require.NoError(t, err)
		_, err = t2InvalidNID.HandleLoginRequest(context.Background(), testLR.ID, &testHLR)
		require.Error(t, err)
		_, err = t1ValidNID.HandleLoginRequest(context.Background(), testLR.ID, &testHLR)
		require.NoError(t, err)
		require.NoError(t, t2InvalidNID.ConfirmLoginSession(context.Background(), testLS.ID, time.Now(), testLS.Subject, true))
		require.NoError(t, t1ValidNID.ConfirmLoginSession(context.Background(), testLS.ID, time.Now(), testLS.Subject, true))
		require.Error(t, t2InvalidNID.DeleteLoginSession(context.Background(), testLS.ID))
		require.NoError(t, t1ValidNID.DeleteLoginSession(context.Background(), testLS.ID))
	}
}

func ManagerTests(m Manager, clientManager client.Manager, fositeManager x.FositeStorer, network string, parallel bool) func(t *testing.T) {
	lr := make(map[string]*LoginRequest)

	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		t.Run("case=init-fks", func(t *testing.T) {
			for _, k := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "rv1", "rv2"} {
				require.NoError(t, clientManager.CreateClient(context.Background(), &client.Client{LegacyClientID: fmt.Sprintf("fk-client-%s", k)}))

				require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
					ID:              makeID("fk-login-session", network, k),
					AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
					Subject:         fmt.Sprintf("subject-%s", k),
				}))

				lr[k] = &LoginRequest{
					ID:              makeID("fk-login-challenge", network, k),
					Subject:         fmt.Sprintf("subject%s", k),
					SessionID:       sqlxx.NullString(makeID("fk-login-session", network, k)),
					Verifier:        makeID("fk-login-verifier", network, k),
					Client:          &client.Client{LegacyClientID: fmt.Sprintf("fk-client-%s", k)},
					AuthenticatedAt: sqlxx.NullTime(time.Now()),
					RequestedAt:     time.Now(),
				}

				require.NoError(t, m.CreateLoginRequest(context.Background(), lr[k]))
			}
		})

		t.Run("case=auth-session", func(t *testing.T) {
			for _, tc := range []struct {
				s LoginSession
			}{
				{
					s: LoginSession{
						ID:              makeID("session", network, "1"),
						AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-time.Minute).UTC()),
						Subject:         "subject1",
					},
				},
				{
					s: LoginSession{
						ID:              makeID("session", network, "2"),
						AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
						Subject:         "subject2",
					},
				},
			} {
				t.Run("case=create-get-"+tc.s.ID, func(t *testing.T) {
					_, err := m.GetRememberedLoginSession(context.Background(), tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error(), "%#v", err)

					err = m.CreateLoginSession(context.Background(), &tc.s)
					require.NoError(t, err)

					_, err = m.GetRememberedLoginSession(context.Background(), tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error())

					updatedAuth := time.Time(tc.s.AuthenticatedAt).Add(time.Second)
					require.NoError(t, m.ConfirmLoginSession(context.Background(), tc.s.ID, updatedAuth, tc.s.Subject, true))

					got, err := m.GetRememberedLoginSession(context.Background(), tc.s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, tc.s.ID, got.ID)
					assert.Equal(t, updatedAuth.Unix(), time.Time(got.AuthenticatedAt).Unix()) // this was updated from confirm...
					assert.EqualValues(t, tc.s.Subject, got.Subject)

					time.Sleep(time.Second) // Make sure AuthAt does not equal...
					updatedAuth2 := time.Now().Truncate(time.Second).UTC()
					require.NoError(t, m.ConfirmLoginSession(context.Background(), tc.s.ID, updatedAuth2, "some-other-subject", true))

					got2, err := m.GetRememberedLoginSession(context.Background(), tc.s.ID)
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
					err := m.DeleteLoginSession(context.Background(), tc.id)
					require.NoError(t, err)

					_, err = m.GetRememberedLoginSession(context.Background(), tc.id)
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
					c, h := MockAuthRequest(tc.key, tc.authAt, network)
					_ = clientManager.CreateClient(context.Background(), c.Client) // Ignore errors that are caused by duplication

					_, err := m.GetLoginRequest(context.Background(), makeID("challenge", network, tc.key))
					require.Error(t, err)

					require.NoError(t, m.CreateLoginRequest(context.Background(), c))

					got1, err := m.GetLoginRequest(context.Background(), makeID("challenge", network, tc.key))
					require.NoError(t, err)
					assert.False(t, got1.WasHandled)
					compareAuthenticationRequest(t, c, got1)

					got1, err = m.HandleLoginRequest(context.Background(), makeID("challenge", network, tc.key), h)
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got1)

					got2, err := m.VerifyAndInvalidateLoginRequest(context.Background(), makeID("verifier", network, tc.key))
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got2.LoginRequest)
					assert.Equal(t, c.ID, got2.ID)

					_, err = m.VerifyAndInvalidateLoginRequest(context.Background(), makeID("verifier", network, tc.key))
					require.Error(t, err)

					got1, err = m.GetLoginRequest(context.Background(), makeID("challenge", network, tc.key))
					require.NoError(t, err)
					assert.True(t, got1.WasHandled)
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
					c, h := MockConsentRequest(tc.key, tc.remember, tc.rememberFor, tc.hasError, tc.skip, tc.authAt, "challenge", network)
					_ = clientManager.CreateClient(context.Background(), c.Client) // Ignore errors that are caused by duplication

					consentChallenge := makeID("challenge", network, tc.key)

					_, err := m.GetConsentRequest(context.Background(), consentChallenge)
					require.Error(t, err)

					require.NoError(t, m.CreateConsentRequest(context.Background(), c))

					got1, err := m.GetConsentRequest(context.Background(), consentChallenge)
					require.NoError(t, err)
					compareConsentRequest(t, c, got1)
					assert.False(t, got1.WasHandled)

					got1, err = m.HandleConsentRequest(context.Background(), h)
					require.NoError(t, err)
					assertx.TimeDifferenceLess(t, time.Now(), time.Time(h.HandledAt), 5)
					compareConsentRequest(t, c, got1)

					h.GrantedAudience = sqlxx.StringSliceJSONFormat{"new-audience"}
					_, err = m.HandleConsentRequest(context.Background(), h)
					require.NoError(t, err)

					got2, err := m.VerifyAndInvalidateConsentRequest(context.Background(), makeID("verifier", network, tc.key))
					require.NoError(t, err)
					compareConsentRequest(t, c, got2.ConsentRequest)
					assert.Equal(t, c.ID, got2.ID)
					assert.Equal(t, h.GrantedAudience, got2.GrantedAudience)

					// Trying to update this again should return an error because the consent request was used.
					h.GrantedAudience = sqlxx.StringSliceJSONFormat{"new-audience", "new-audience-2"}
					_, err = m.HandleConsentRequest(context.Background(), h)
					require.Error(t, err)

					if tc.hasError {
						assert.True(t, got2.HasError())
					}
					assert.Equal(t, tc.remember, got2.Remember)
					assert.Equal(t, tc.rememberFor, got2.RememberFor)

					_, err = m.VerifyAndInvalidateConsentRequest(context.Background(), makeID("verifier", network, tc.key))
					require.Error(t, err)

					got1, err = m.GetConsentRequest(context.Background(), consentChallenge)
					require.NoError(t, err)
					assert.True(t, got1.WasHandled)
				})
			}

			for _, tc := range []struct {
				keyC           string
				keyS           string
				expectedLength int
			}{
				{"1", "1", 1},
				{"2", "2", 0},
				{"3", "3", 0},
				{"4", "4", 0},
				{"1", "2", 0},
				{"2", "1", 0},
				{"5", "5", 1},
				{"6", "6", 0},
			} {
				t.Run("key="+tc.keyC+"-"+tc.keyS, func(t *testing.T) {
					rs, err := m.FindGrantedAndRememberedConsentRequests(context.Background(), "fk-client-"+tc.keyC, "subject"+tc.keyS)
					if tc.expectedLength == 0 {
						assert.EqualError(t, err, ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						assert.Len(t, rs, tc.expectedLength)
					}
				})
			}
		})

		t.Run("case=revoke-auth-request", func(t *testing.T) {
			require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
				ID:              makeID("rev-session", network, "-1"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-1",
			}))

			require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
				ID:              makeID("rev-session", network, "-2"),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         "subject-2",
			}))

			require.NoError(t, m.CreateLoginSession(context.Background(), &LoginSession{
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
					require.NoError(t, m.RevokeSubjectLoginSession(context.Background(), tc.subject))

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetRememberedLoginSession(context.Background(), id)
							assert.EqualError(t, err, x.ErrNotFound.Error())
						})
					}
				})
			}
		})

		challengerv1 := makeID("challenge", network, "rv1")
		challengerv2 := makeID("challenge", network, "rv2")
		t.Run("case=revoke-used-consent-request", func(t *testing.T) {
			cr1, hcr1 := MockConsentRequest("rv1", false, 0, false, false, false, "fk-login-challenge", network)
			cr2, hcr2 := MockConsentRequest("rv2", false, 0, false, false, false, "fk-login-challenge", network)

			// Ignore duplication errors
			_ = clientManager.CreateClient(context.Background(), cr1.Client)
			_ = clientManager.CreateClient(context.Background(), cr2.Client)

			require.NoError(t, m.CreateConsentRequest(context.Background(), cr1))
			require.NoError(t, m.CreateConsentRequest(context.Background(), cr2))
			_, err := m.HandleConsentRequest(context.Background(), hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(context.Background(), hcr2)
			require.NoError(t, err)

			require.NoError(t, fositeManager.CreateAccessTokenSession(context.Background(), makeID("", network, "trva1"), &fosite.Request{Client: cr1.Client, ID: challengerv1, RequestedAt: time.Now()}))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(context.Background(), makeID("", network, "rrva1"), &fosite.Request{Client: cr1.Client, ID: challengerv1, RequestedAt: time.Now()}))
			require.NoError(t, fositeManager.CreateAccessTokenSession(context.Background(), makeID("", network, "trva2"), &fosite.Request{Client: cr2.Client, ID: challengerv2, RequestedAt: time.Now()}))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(context.Background(), makeID("", network, "rrva2"), &fosite.Request{Client: cr2.Client, ID: challengerv2, RequestedAt: time.Now()}))

			for i, tc := range []struct {
				subject string
				client  string
				at      string
				rt      string
				ids     []string
			}{
				{
					at: makeID("", network, "trva1"), rt: makeID("", network, "rrva1"),
					subject: "subjectrv1",
					client:  "",
					ids:     []string{challengerv1},
				},
				{
					at: makeID("", network, "trva2"), rt: makeID("", network, "rrva2"),
					subject: "subjectrv2",
					client:  "fk-client-rv2",
					ids:     []string{challengerv2},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					_, err := fositeManager.GetAccessTokenSession(context.Background(), tc.at, nil)
					assert.NoError(t, err)
					_, err = fositeManager.GetRefreshTokenSession(context.Background(), tc.rt, nil)
					assert.NoError(t, err)

					if tc.client == "" {
						require.NoError(t, m.RevokeSubjectConsentSession(context.Background(), tc.subject))
					} else {
						require.NoError(t, m.RevokeSubjectClientConsentSession(context.Background(), tc.subject, tc.client))
					}

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetConsentRequest(context.Background(), id)
							assert.True(t, errors.Is(err, x.ErrNotFound))
						})
					}

					r, err := fositeManager.GetAccessTokenSession(context.Background(), tc.at, nil)
					assert.Error(t, err, "%+v", r)
					r, err = fositeManager.GetRefreshTokenSession(context.Background(), tc.rt, nil)
					assert.Error(t, err, "%+v", r)
				})
			}

			require.EqualError(t, m.RevokeSubjectConsentSession(context.Background(), "i-do-not-exist"), x.ErrNotFound.Error())
			require.EqualError(t, m.RevokeSubjectClientConsentSession(context.Background(), "i-do-not-exist", "i-do-not-exist"), x.ErrNotFound.Error())
		})

		t.Run("case=list-used-consent-requests", func(t *testing.T) {
			require.NoError(t, m.CreateLoginRequest(context.Background(), lr["rv1"]))
			require.NoError(t, m.CreateLoginRequest(context.Background(), lr["rv2"]))

			cr1, hcr1 := MockConsentRequest("rv1", true, 0, false, false, false, "fk-login-challenge", network)
			cr2, hcr2 := MockConsentRequest("rv2", false, 0, false, false, false, "fk-login-challenge", network)

			// Ignore duplicate errors
			_ = clientManager.CreateClient(context.Background(), cr1.Client)
			_ = clientManager.CreateClient(context.Background(), cr2.Client)

			require.NoError(t, m.CreateConsentRequest(context.Background(), cr1))
			require.NoError(t, m.CreateConsentRequest(context.Background(), cr2))
			_, err := m.HandleConsentRequest(context.Background(), hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(context.Background(), hcr2)
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
					challenges: []string{challengerv1},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    cr2.Subject,
					sid:        makeID("fk-login-session", network, "rv2"),
					challenges: []string{challengerv2},
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
					consents, err := m.FindSubjectsSessionGrantedConsentRequests(context.Background(), tc.subject, tc.sid, 100, 0)
					assert.Equal(t, len(tc.challenges), len(consents))

					if len(tc.challenges) == 0 {
						assert.EqualError(t, err, ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						for _, consent := range consents {
							assert.Contains(t, tc.challenges, consent.ID)
							assert.Contains(t, tc.clients, consent.ConsentRequest.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(context.Background(), tc.subject)
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
					challenges: []string{challengerv1},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    "subjectrv2",
					challenges: []string{challengerv2},
					clients:    []string{"fk-client-rv2"},
				},
				{
					subject:    "subjectrv3",
					challenges: []string{},
					clients:    []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					consents, err := m.FindSubjectsGrantedConsentRequests(context.Background(), tc.subject, 100, 0)
					assert.Equal(t, len(tc.challenges), len(consents))

					if len(tc.challenges) == 0 {
						assert.EqualError(t, err, ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						for _, consent := range consents {
							assert.Contains(t, tc.challenges, consent.ID)
							assert.Contains(t, tc.clients, consent.ConsentRequest.Client.GetID())
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(context.Background(), tc.subject)
					require.NoError(t, err)
					assert.Equal(t, n, len(tc.challenges))

				})
			}

			t.Run("case=obfuscated", func(t *testing.T) {
				_, err := m.GetForcedObfuscatedLoginSession(context.Background(), "fk-client-1", "obfuscated-1")
				require.True(t, errors.Is(err, x.ErrNotFound))

				expect := &ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-1",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(context.Background(), expect))

				got, err := m.GetForcedObfuscatedLoginSession(context.Background(), "fk-client-1", "obfuscated-1")
				require.NoError(t, err)
				require.NotEqual(t, got.NID, gofrsuuid.Nil)
				got.NID = gofrsuuid.Nil
				assert.EqualValues(t, expect, got)

				expect = &ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-2",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(context.Background(), expect))

				got, err = m.GetForcedObfuscatedLoginSession(context.Background(), "fk-client-1", "obfuscated-2")
				require.NotEqual(t, got.NID, gofrsuuid.Nil)
				got.NID = gofrsuuid.Nil
				require.NoError(t, err)
				assert.EqualValues(t, expect, got)

				_, err = m.GetForcedObfuscatedLoginSession(context.Background(), "fk-client-1", "obfuscated-1")
				require.True(t, errors.Is(err, x.ErrNotFound))
			})

			t.Run("case=ListUserAuthenticatedClientsWithFrontAndBackChannelLogout", func(t *testing.T) {
				// The idea of this test is to create two identities (subjects) with 4 sessions each, where
				// only some sessions have been associated with a client that has a front channel logout url

				subjects := make([]string, 1)
				for k := range subjects {
					subjects[k] = fmt.Sprintf("subject-ListUserAuthenticatedClientsWithFrontAndBackChannelLogout-%d", k)
				}

				sessions := make([]LoginSession, len(subjects)*1)
				frontChannels := map[string][]client.Client{}
				backChannels := map[string][]client.Client{}
				for k := range sessions {
					id := uuid.New().String()
					subject := subjects[k%len(subjects)]
					t.Run(fmt.Sprintf("create/session=%s/subject=%s", id, subject), func(t *testing.T) {
						ls := &LoginSession{
							ID:              id,
							AuthenticatedAt: sqlxx.NullTime(time.Now()),
							Subject:         subject,
						}
						require.NoError(t, m.CreateLoginSession(context.Background(), ls))

						cl := &client.Client{LegacyClientID: uuid.New().String()}
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
						require.NoError(t, clientManager.CreateClient(context.Background(), cl))

						ar := SaneMockAuthRequest(t, m, ls, cl)
						cr := SaneMockConsentRequest(t, m, ar, false)
						_ = SaneMockHandleConsentRequest(t, m, cr, time.Time{}, 0, false, false)

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
						actual, err := m.ListUserAuthenticatedClientsWithFrontChannelLogout(context.Background(), ls.Subject, ls.ID)
						require.NoError(t, err)
						check(t, frontChannels, actual)
					})

					t.Run(fmt.Sprintf("method=ListUserAuthenticatedClientsWithBackChannelLogout/session=%s", ls.ID), func(t *testing.T) {
						actual, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(context.Background(), ls.Subject, ls.ID)
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
							require.NoError(t, clientManager.CreateClient(context.Background(), c.Client)) // Ignore errors that are caused by duplication
						}

						_, err := m.GetLogoutRequest(context.Background(), challenge)
						require.Error(t, err)

						require.NoError(t, m.CreateLogoutRequest(context.Background(), c))

						got2, err := m.GetLogoutRequest(context.Background(), challenge)
						require.NoError(t, err)
						assert.False(t, got2.WasHandled)
						assert.False(t, got2.Accepted)
						compareLogoutRequest(t, c, got2)

						if k%2 == 0 {
							got2, err = m.AcceptLogoutRequest(context.Background(), challenge)
							require.NoError(t, err)
							assert.True(t, got2.Accepted)
							compareLogoutRequest(t, c, got2)

							got3, err := m.VerifyAndInvalidateLogoutRequest(context.Background(), verifier)
							require.NoError(t, err)
							assert.True(t, got3.Accepted)
							assert.True(t, got3.WasHandled)
							compareLogoutRequest(t, c, got3)

							_, err = m.VerifyAndInvalidateLogoutRequest(context.Background(), verifier)
							require.Error(t, err)

							got2, err = m.GetLogoutRequest(context.Background(), challenge)
							require.NoError(t, err)
							compareLogoutRequest(t, got3, got2)
							assert.True(t, got2.WasHandled)
						} else {
							require.NoError(t, m.RejectLogoutRequest(context.Background(), challenge))
							_, err = m.GetLogoutRequest(context.Background(), challenge)
							require.Error(t, err)
						}
					})
				}
			})
		})

		t.Run("case=foreign key regression", func(t *testing.T) {
			cl := &client.Client{LegacyClientID: uuid.New().String()}
			require.NoError(t, clientManager.CreateClient(context.Background(), cl))

			subject := uuid.New().String()
			s := LoginSession{
				ID:              uuid.New().String(),
				AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Minute).Add(-time.Minute).UTC()),
				Subject:         subject,
			}

			err := m.CreateLoginSession(context.Background(), &s)
			require.NoError(t, err)

			lr := &LoginRequest{
				ID:              uuid.New().String(),
				Subject:         uuid.New().String(),
				Verifier:        uuid.New().String(),
				Client:          cl,
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				RequestedAt:     time.Now(),
				SessionID:       sqlxx.NullString(s.ID),
			}

			require.NoError(t, m.CreateLoginRequest(context.Background(), lr))
			expected := &OAuth2ConsentRequest{
				ID:                   uuid.New().String(),
				Skip:                 true,
				Subject:              subject,
				OpenIDConnectContext: nil,
				Client:               cl,
				ClientID:             cl.LegacyClientID,
				RequestURL:           "",
				LoginChallenge:       sqlxx.NullString(lr.ID),
				LoginSessionID:       sqlxx.NullString(s.ID),
				Verifier:             uuid.New().String(),
				CSRF:                 uuid.New().String(),
			}
			require.NoError(t, m.CreateConsentRequest(context.Background(), expected))

			result, err := m.GetConsentRequest(context.Background(), expected.ID)
			require.NoError(t, err)
			assert.EqualValues(t, expected.ID, result.ID)

			require.NoError(t, m.DeleteLoginSession(context.Background(), s.ID))

			result, err = m.GetConsentRequest(context.Background(), expected.ID)
			require.NoError(t, err)
			assert.EqualValues(t, expected.ID, result.ID)
		})
	}
}

func compareLogoutRequest(t *testing.T, a, b *LogoutRequest) {
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

func compareAuthenticationRequest(t *testing.T, a, b *LoginRequest) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.ID, b.ID)
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.CSRF, b.CSRF)
	assert.EqualValues(t, a.Skip, b.Skip)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareConsentRequest(t *testing.T, a, b *OAuth2ConsentRequest) {
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

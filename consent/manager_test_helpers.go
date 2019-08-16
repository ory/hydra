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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
)

func MockConsentRequest(key string, remember bool, rememberFor int, hasError bool, skip bool, authAt bool) (c *ConsentRequest, h *HandledConsentRequest) {
	c = &ConsentRequest{
		Challenge:         "challenge" + key,
		RequestedScope:    []string{"scopea" + key, "scopeb" + key},
		RequestedAudience: []string{"auda" + key, "audb" + key},
		Skip:              skip,
		Subject:           "subject" + key,
		OpenIDConnectContext: &OpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		Client:                 &client.Client{ClientID: "fk-client-" + key},
		RequestURL:             "https://request-url/path" + key,
		LoginChallenge:         "fk-login-challenge-" + key,
		LoginSessionID:         "fk-login-session-" + key,
		ForceSubjectIdentifier: "forced-subject",
		SubjectIdentifier:      "forced-subject",
		Verifier:               "verifier" + key,
		CSRF:                   "csrf" + key,
		ACR:                    "1",
		AuthenticatedAt:        time.Now().UTC().Add(-time.Hour),
		RequestedAt:            time.Now().UTC().Add(-time.Hour),
		Context:                map[string]interface{}{"foo": "bar" + key},
	}

	var err *RequestDeniedError
	if hasError {
		err = &RequestDeniedError{
			Name:        "error_name" + key,
			Description: "error_description" + key,
			Hint:        "error_hint,omitempty" + key,
			Code:        100,
			Debug:       "error_debug,omitempty" + key,
		}
	}

	var authenticatedAt time.Time
	if authAt {
		time.Now().UTC().Add(-time.Minute)
	}

	h = &HandledConsentRequest{
		ConsentRequest:  c,
		RememberFor:     rememberFor,
		Remember:        remember,
		Challenge:       "challenge" + key,
		RequestedAt:     time.Now().UTC().Add(-time.Minute),
		AuthenticatedAt: authenticatedAt,
		GrantedScope:    []string{"scopea" + key, "scopeb" + key},
		GrantedAudience: []string{"auda" + key, "audb" + key},
		Error:           err,
		//WasUsed:         true,
	}

	return c, h
}

func MockLogoutRequest(key string, withClient bool) (c *LogoutRequest) {
	var cl *client.Client
	if withClient {
		cl = &client.Client{
			ClientID: "fk-client-" + key,
		}
	}
	return &LogoutRequest{
		Subject:               "subject" + key,
		Challenge:             "challenge" + key,
		Verifier:              "verifier" + key,
		SessionID:             "session" + key,
		RPInitiated:           true,
		RequestURL:            "http://request-me/",
		PostLogoutRedirectURI: "http://redirect-me/",
		WasUsed:               false,
		Accepted:              false,
		Client:                cl,
	}
}

func MockAuthRequest(key string, authAt bool) (c *LoginRequest, h *HandledLoginRequest) {
	c = &LoginRequest{
		OpenIDConnectContext: &OpenIDConnectContext{
			ACRValues: []string{"1" + key, "2" + key},
			UILocales: []string{"fr" + key, "de" + key},
			Display:   "popup" + key,
		},
		RequestedAt:    time.Now().UTC().Add(-time.Hour),
		Client:         &client.Client{ClientID: "fk-client-" + key},
		Subject:        "subject" + key,
		RequestURL:     "https://request-url/path" + key,
		Skip:           true,
		Challenge:      "challenge" + key,
		Verifier:       "verifier" + key,
		RequestedScope: []string{"scopea" + key, "scopeb" + key},
		CSRF:           "csrf" + key,
		SessionID:      "fk-login-session-" + key,
	}

	var err = &RequestDeniedError{
		Name:        "error_name" + key,
		Description: "error_description" + key,
		Hint:        "error_hint,omitempty" + key,
		Code:        100,
		Debug:       "error_debug,omitempty" + key,
	}

	var authenticatedAt time.Time
	if authAt {
		authenticatedAt = time.Now().UTC().Add(-time.Minute)
	}

	h = &HandledLoginRequest{
		LoginRequest:           c,
		RememberFor:            120,
		Remember:               true,
		Challenge:              "challenge" + key,
		RequestedAt:            time.Now().UTC().Add(-time.Minute),
		AuthenticatedAt:        authenticatedAt,
		Error:                  err,
		Subject:                c.Subject,
		ACR:                    "acr",
		ForceSubjectIdentifier: "forced-subject",
		//WasUsed:                false,
	}

	return c, h
}

func ManagerTests(m Manager, clientManager client.Manager, fositeManager x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("case=init-fks", func(t *testing.T) {
			for _, k := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "rv1", "rv2"} {
				require.NoError(t, clientManager.CreateClient(context.TODO(), &client.Client{ClientID: fmt.Sprintf("fk-client-%s", k)}))

				require.NoError(t, m.CreateLoginSession(context.TODO(), &LoginSession{
					ID:              fmt.Sprintf("fk-login-session-%s", k),
					AuthenticatedAt: time.Now().Round(time.Second).UTC(),
					Subject:         fmt.Sprintf("subject-%s", k),
				}))

				require.NoError(t, m.CreateLoginRequest(context.TODO(), &LoginRequest{
					Challenge:       fmt.Sprintf("fk-login-challenge-%s", k),
					Verifier:        fmt.Sprintf("fk-login-verifier-%s", k),
					Client:          &client.Client{ClientID: fmt.Sprintf("fk-client-%s", k)},
					AuthenticatedAt: time.Now(),
					RequestedAt:     time.Now(),
				}))
			}
		})

		t.Run("case=auth-session", func(t *testing.T) {
			for _, tc := range []struct {
				s LoginSession
			}{
				{
					s: LoginSession{
						ID:              "session1",
						AuthenticatedAt: time.Now().Round(time.Second).Add(-time.Minute).UTC(),
						Subject:         "subject1",
					},
				},
				{
					s: LoginSession{
						ID:              "session2",
						AuthenticatedAt: time.Now().Round(time.Minute).Add(-time.Minute).UTC(),
						Subject:         "subject2",
					},
				},
			} {
				t.Run("case=create-get-"+tc.s.ID, func(t *testing.T) {
					_, err := m.GetRememberedLoginSession(context.TODO(), tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error())

					err = m.CreateLoginSession(context.TODO(), &tc.s)
					require.NoError(t, err)

					_, err = m.GetRememberedLoginSession(context.TODO(), tc.s.ID)
					require.EqualError(t, err, x.ErrNotFound.Error())

					require.NoError(t, m.ConfirmLoginSession(context.TODO(), tc.s.ID, tc.s.Subject, true))

					got, err := m.GetRememberedLoginSession(context.TODO(), tc.s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, tc.s.ID, got.ID)
					assert.NotEqual(t, tc.s.AuthenticatedAt.Unix(), got.AuthenticatedAt.Unix()) // this was updated from confirm...
					assert.EqualValues(t, tc.s.Subject, got.Subject)

					time.Sleep(time.Second) // Make sure AuthAt does not equal...
					require.NoError(t, m.ConfirmLoginSession(context.TODO(), tc.s.ID, "some-other-subject", true))

					got2, err := m.GetRememberedLoginSession(context.TODO(), tc.s.ID)
					require.NoError(t, err)
					assert.EqualValues(t, tc.s.ID, got2.ID)
					assert.NotEqual(t, tc.s.AuthenticatedAt.Unix(), got2.AuthenticatedAt.Unix()) // this was updated from confirm...
					assert.NotEqual(t, got.AuthenticatedAt.Unix(), got2.AuthenticatedAt.Unix())  // this was updated from confirm...
					assert.EqualValues(t, "some-other-subject", got2.Subject)
				})
			}
			for _, tc := range []struct {
				id string
			}{
				{
					id: "session1",
				},
				{
					id: "session2",
				},
			} {
				t.Run("case=delete-get-"+tc.id, func(t *testing.T) {
					err := m.DeleteLoginSession(context.TODO(), tc.id)
					require.NoError(t, err)

					_, err = m.GetRememberedLoginSession(context.TODO(), tc.id)
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
			} {
				t.Run("key="+tc.key, func(t *testing.T) {
					c, h := MockAuthRequest(tc.key, tc.authAt)
					clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

					_, err := m.GetLoginRequest(context.TODO(), "challenge"+tc.key)
					require.Error(t, err)

					require.NoError(t, m.CreateLoginRequest(context.TODO(), c))

					got1, err := m.GetLoginRequest(context.TODO(), "challenge"+tc.key)
					require.NoError(t, err)
					assert.False(t, got1.WasHandled)
					compareAuthenticationRequest(t, c, got1)

					got1, err = m.HandleLoginRequest(context.TODO(), "challenge"+tc.key, h)
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got1)

					got2, err := m.VerifyAndInvalidateLoginRequest(context.TODO(), "verifier"+tc.key)
					require.NoError(t, err)
					compareAuthenticationRequest(t, c, got2.LoginRequest)
					assert.Equal(t, c.Challenge, got2.Challenge)

					_, err = m.VerifyAndInvalidateLoginRequest(context.TODO(), "verifier"+tc.key)
					require.Error(t, err)

					got1, err = m.GetLoginRequest(context.TODO(), "challenge"+tc.key)
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
					c, h := MockConsentRequest(tc.key, tc.remember, tc.rememberFor, tc.hasError, tc.skip, tc.authAt)
					clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

					_, err := m.GetConsentRequest(context.TODO(), "challenge"+tc.key)
					require.Error(t, err)

					require.NoError(t, m.CreateConsentRequest(context.TODO(), c))

					got1, err := m.GetConsentRequest(context.TODO(), "challenge"+tc.key)
					require.NoError(t, err)
					compareConsentRequest(t, c, got1)
					assert.False(t, got1.WasHandled)

					got1, err = m.HandleConsentRequest(context.TODO(), "challenge"+tc.key, h)
					require.NoError(t, err)
					compareConsentRequest(t, c, got1)

					_, err = m.HandleConsentRequest(context.TODO(), "challenge"+tc.key, h)
					require.NoError(t, err)

					got2, err := m.VerifyAndInvalidateConsentRequest(context.TODO(), "verifier"+tc.key)
					require.NoError(t, err)
					compareConsentRequest(t, c, got2.ConsentRequest)
					assert.Equal(t, c.Challenge, got2.Challenge)

					_, err = m.VerifyAndInvalidateConsentRequest(context.TODO(), "verifier"+tc.key)
					require.Error(t, err)

					got1, err = m.GetConsentRequest(context.TODO(), "challenge"+tc.key)
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
					rs, err := m.FindGrantedAndRememberedConsentRequests(context.TODO(), "fk-client-"+tc.keyC, "subject"+tc.keyS)
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
			require.NoError(t, m.CreateLoginSession(context.TODO(), &LoginSession{
				ID:              "rev-session-1",
				AuthenticatedAt: time.Now(),
				Subject:         "subject-1",
			}))

			require.NoError(t, m.CreateLoginSession(context.TODO(), &LoginSession{
				ID:              "rev-session-2",
				AuthenticatedAt: time.Now(),
				Subject:         "subject-2",
			}))

			require.NoError(t, m.CreateLoginSession(context.TODO(), &LoginSession{
				ID:              "rev-session-3",
				AuthenticatedAt: time.Now(),
				Subject:         "subject-1",
			}))

			for i, tc := range []struct {
				subject string
				ids     []string
			}{
				{
					subject: "subject-1",
					ids:     []string{"rev-session-1", "rev-session-3"},
				},
				{
					subject: "subject-2",
					ids:     []string{"rev-session-1", "rev-session-3", "rev-session-2"},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					require.NoError(t, m.RevokeSubjectLoginSession(context.TODO(), tc.subject))

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetRememberedLoginSession(context.TODO(), id)
							assert.EqualError(t, err, x.ErrNotFound.Error())
						})
					}
				})
			}
		})

		t.Run("case=revoke-used-consent-request", func(t *testing.T) {
			cr1, hcr1 := MockConsentRequest("rv1", false, 0, false, false, false)
			cr2, hcr2 := MockConsentRequest("rv2", false, 0, false, false, false)
			clientManager.CreateClient(context.TODO(), cr1.Client)
			clientManager.CreateClient(context.TODO(), cr2.Client)

			require.NoError(t, m.CreateConsentRequest(context.TODO(), cr1))
			require.NoError(t, m.CreateConsentRequest(context.TODO(), cr2))
			_, err := m.HandleConsentRequest(context.TODO(), "challengerv1", hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(context.TODO(), "challengerv2", hcr2)
			require.NoError(t, err)

			fositeManager.CreateAccessTokenSession(context.TODO(), "trva1", &fosite.Request{Client: cr1.Client, ID: "challengerv1", RequestedAt: time.Now()})
			fositeManager.CreateRefreshTokenSession(context.TODO(), "rrva1", &fosite.Request{Client: cr1.Client, ID: "challengerv1", RequestedAt: time.Now()})
			fositeManager.CreateAccessTokenSession(context.TODO(), "trva2", &fosite.Request{Client: cr2.Client, ID: "challengerv2", RequestedAt: time.Now()})
			fositeManager.CreateRefreshTokenSession(context.TODO(), "rrva2", &fosite.Request{Client: cr2.Client, ID: "challengerv2", RequestedAt: time.Now()})

			for i, tc := range []struct {
				subject string
				client  string
				at      string
				rt      string
				ids     []string
			}{
				{
					at: "trva1", rt: "rrva1",
					subject: "subjectrv1",
					client:  "",
					ids:     []string{"challengerv1"},
				},
				{
					at: "trva2", rt: "rrva2",
					subject: "subjectrv2",
					client:  "fk-client-rv2",
					ids:     []string{"challengerv2"},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					_, err := fositeManager.GetAccessTokenSession(context.TODO(), tc.at, nil)
					assert.NoError(t, err)
					_, err = fositeManager.GetRefreshTokenSession(context.TODO(), tc.rt, nil)
					assert.NoError(t, err)

					if tc.client == "" {
						require.NoError(t, m.RevokeSubjectConsentSession(context.TODO(), tc.subject))
					} else {
						require.NoError(t, m.RevokeSubjectClientConsentSession(context.TODO(), tc.subject, tc.client))
					}

					for _, id := range tc.ids {
						t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
							_, err := m.GetConsentRequest(context.TODO(), id)
							assert.EqualError(t, err, x.ErrNotFound.Error())
						})
					}

					r, err := fositeManager.GetAccessTokenSession(context.TODO(), tc.at, nil)
					assert.Error(t, err, "%+v", r)
					r, err = fositeManager.GetRefreshTokenSession(context.TODO(), tc.rt, nil)
					assert.Error(t, err, "%+v", r)
				})
			}
		})

		t.Run("case=list-used-consent-requests", func(t *testing.T) {
			cr1, hcr1 := MockConsentRequest("rv1", true, 0, false, false, false)
			cr2, hcr2 := MockConsentRequest("rv2", false, 0, false, false, false)
			clientManager.CreateClient(context.TODO(), cr1.Client)
			clientManager.CreateClient(context.TODO(), cr2.Client)

			require.NoError(t, m.CreateConsentRequest(context.TODO(), cr1))
			require.NoError(t, m.CreateConsentRequest(context.TODO(), cr2))
			_, err := m.HandleConsentRequest(context.TODO(), "challengerv1", hcr1)
			require.NoError(t, err)
			_, err = m.HandleConsentRequest(context.TODO(), "challengerv2", hcr2)
			require.NoError(t, err)

			for i, tc := range []struct {
				subject    string
				challenges []string
				clients    []string
			}{
				{
					subject:    "subjectrv1",
					challenges: []string{"challengerv1"},
					clients:    []string{"fk-client-rv1"},
				},
				{
					subject:    "subjectrv2",
					challenges: []string{"challengerv2"},
					clients:    []string{"fk-client-rv2"},
				},
				{
					subject:    "subjectrv3",
					challenges: []string{},
					clients:    []string{},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/subject=%s", i, tc.subject), func(t *testing.T) {
					consents, err := m.FindSubjectsGrantedConsentRequests(context.TODO(), tc.subject, 100, 0)
					assert.Equal(t, len(tc.challenges), len(consents))

					if len(tc.challenges) == 0 {
						assert.EqualError(t, err, ErrNoPreviousConsentFound.Error())
					} else {
						require.NoError(t, err)
						for _, consent := range consents {
							assert.Contains(t, tc.challenges, consent.Challenge)
							assert.Contains(t, tc.clients, consent.ConsentRequest.Client.ClientID)
						}
					}

					n, err := m.CountSubjectsGrantedConsentRequests(context.TODO(), tc.subject)
					require.NoError(t, err)
					assert.Equal(t, n, len(tc.challenges))

				})
			}

			t.Run("case=obfuscated", func(t *testing.T) {
				got, err := m.GetForcedObfuscatedLoginSession(context.TODO(), "fk-client-1", "obfuscated-1")
				require.EqualError(t, err, x.ErrNotFound.Error())

				expect := &ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-1",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(context.TODO(), expect))

				got, err = m.GetForcedObfuscatedLoginSession(context.TODO(), "fk-client-1", "obfuscated-1")
				require.NoError(t, err)
				assert.EqualValues(t, expect, got)

				expect = &ForcedObfuscatedLoginSession{
					ClientID:          "fk-client-1",
					Subject:           "subject-1",
					SubjectObfuscated: "obfuscated-2",
				}
				require.NoError(t, m.CreateForcedObfuscatedLoginSession(context.TODO(), expect))

				got, err = m.GetForcedObfuscatedLoginSession(context.TODO(), "fk-client-1", "obfuscated-2")
				require.NoError(t, err)
				assert.EqualValues(t, expect, got)

				got, err = m.GetForcedObfuscatedLoginSession(context.TODO(), "fk-client-1", "obfuscated-1")
				require.EqualError(t, err, x.ErrNotFound.Error())
			})

			t.Run("case=ListUserAuthenticatedClientsWithFrontChannelLogout", func(t *testing.T) {
				for i := 0; i <= 10; i++ {
					c, h := MockConsentRequest(fmt.Sprintf("LUACWFCL-%d", i), false, 0, false, false, false)
					if i == 5 {
						c.Client.FrontChannelLogoutURI = "http://some-url.com/"
					}
					c.LoginSessionID = "" // otherwise we had to create the login session as well..
					c.Subject = "subjectLUACWFCL"
					clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

					lc, _ := MockAuthRequest(fmt.Sprintf("LUACWFCL-%d", i), true)
					lc.Challenge = fmt.Sprintf("fk-login-challenge-LUACWFCL-%d", i)
					lc.SessionID = ""
					require.NoError(t, m.CreateLoginRequest(context.TODO(), lc))

					require.NoError(t, m.CreateConsentRequest(context.TODO(), c))
					_, err = m.HandleConsentRequest(context.TODO(), c.Challenge, h)
					require.NoError(t, err)
				}

				clients, err := m.ListUserAuthenticatedClientsWithFrontChannelLogout(context.TODO(), "subjectLUACWFCL")
				require.NoError(t, err)

				require.Len(t, clients, 1)
				assert.EqualValues(t, "http://some-url.com/", clients[0].FrontChannelLogoutURI)
				assert.EqualValues(t, "fk-client-LUACWFCL-5", clients[0].ClientID)
			})

			t.Run("case=ListUserAuthenticatedClientsWithFrontChannelLogout-duplicateLogin", func(t *testing.T) {
				c, h := MockConsentRequest("LUACWFCL-DUPLICCATE", false, 0, false, false, false)
				c.Client.BackChannelLogoutURI = "http://some-url.com/"

				c.LoginSessionID = ""                                // otherwise we had to create the login session as well..
				c.Subject = "subjectLUACWFCL-DUPLICCATE"             // otherwise we had to create the login session as well..
				clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

				for i := 0; i < 2; i++ {
					lc, _ := MockAuthRequest("LUACWFCL-DUPLICCATE", true)
					lc.Challenge = fmt.Sprintf("fk-login-challenge-LUACWFCL-DUPLICCATE-%d", i)
					lc.Verifier = fmt.Sprintf("fk-login-verifier-LUACWFCL-DUPLICCATE-%d", i)
					lc.SessionID = ""
					require.NoError(t, m.CreateLoginRequest(context.TODO(), lc))

					c.Challenge = fmt.Sprintf("challenge-LUACWFCL-DUPLICCATE-%d", i)
					c.LoginChallenge = fmt.Sprintf("fk-login-challenge-LUACWFCL-DUPLICCATE-%d", i)
					c.Verifier = fmt.Sprintf("verifiers-LUACWFCL-DUPLICCATE-%d", i)
					require.NoError(t, m.CreateConsentRequest(context.TODO(), c))

					h.Challenge = fmt.Sprintf("challenge-LUACWFCL-DUPLICCATE-%d", i)
					_, err = m.HandleConsentRequest(context.TODO(), c.Challenge, h)
					require.NoError(t, err)
				}

				clients, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(context.TODO(), "subjectLUACWFCL-DUPLICCATE")
				require.NoError(t, err)

				require.Len(t, clients, 1)
				assert.EqualValues(t, "http://some-url.com/", clients[0].BackChannelLogoutURI)
				assert.EqualValues(t, "fk-client-LUACWFCL-DUPLICCATE", clients[0].ClientID)
			})

			t.Run("case=ListUserAuthenticatedClientsWithBackChannelLogout", func(t *testing.T) {
				for i := 0; i <= 10; i++ {
					c, h := MockConsentRequest(fmt.Sprintf("LUACWFBL-%d", i), false, 0, false, false, false)
					if i == 5 {
						c.Client.BackChannelLogoutURI = "http://some-url.com/"
					}
					c.LoginSessionID = ""                                // otherwise we had to create the login session as well..
					c.Subject = "subjectLUACWFBL"                        // otherwise we had to create the login session as well..
					clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

					lc, _ := MockAuthRequest(fmt.Sprintf("LUACWFBL-%d", i), true)
					lc.Challenge = fmt.Sprintf("fk-login-challenge-LUACWFBL-%d", i)
					lc.SessionID = ""
					require.NoError(t, m.CreateLoginRequest(context.TODO(), lc))

					require.NoError(t, m.CreateConsentRequest(context.TODO(), c))
					_, err = m.HandleConsentRequest(context.TODO(), c.Challenge, h)
					require.NoError(t, err)
				}

				clients, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(context.TODO(), "subjectLUACWFBL")
				require.NoError(t, err)

				require.Len(t, clients, 1)
				assert.EqualValues(t, "http://some-url.com/", clients[0].BackChannelLogoutURI)
				assert.EqualValues(t, "fk-client-LUACWFBL-5", clients[0].ClientID)
			})

			t.Run("case=ListUserAuthenticatedClientsWithBackChannelLogout-duplicateLogin", func(t *testing.T) {
				c, h := MockConsentRequest("LUACWFBL-DUPLICCATE", false, 0, false, false, false)
				c.Client.BackChannelLogoutURI = "http://some-url.com/"

				c.LoginSessionID = ""                                // otherwise we had to create the login session as well..
				c.Subject = "subjectLUACWFBL-DUPLICCATE"             // otherwise we had to create the login session as well..
				clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication

				for i := 0; i < 2; i++ {
					lc, _ := MockAuthRequest("LUACWFBL-DUPLICCATE", true)
					lc.Challenge = fmt.Sprintf("fk-login-challenge-LUACWFBL-DUPLICCATE-%d", i)
					lc.Verifier = fmt.Sprintf("fk-login-verifier-LUACWFBL-DUPLICCATE-%d", i)
					lc.SessionID = ""
					require.NoError(t, m.CreateLoginRequest(context.TODO(), lc))

					c.Challenge = fmt.Sprintf("challenge-LUACWFBL-DUPLICCATE-%d", i)
					c.LoginChallenge = fmt.Sprintf("fk-login-challenge-LUACWFBL-DUPLICCATE-%d", i)
					c.Verifier = fmt.Sprintf("verifiers-LUACWFBL-DUPLICCATE-%d", i)
					require.NoError(t, m.CreateConsentRequest(context.TODO(), c))

					h.Challenge = fmt.Sprintf("challenge-LUACWFBL-DUPLICCATE-%d", i)
					_, err = m.HandleConsentRequest(context.TODO(), c.Challenge, h)
					require.NoError(t, err)
				}

				clients, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(context.TODO(), "subjectLUACWFBL-DUPLICCATE")
				require.NoError(t, err)

				require.Len(t, clients, 1)
				assert.EqualValues(t, "http://some-url.com/", clients[0].BackChannelLogoutURI)
				assert.EqualValues(t, "fk-client-LUACWFBL-DUPLICCATE", clients[0].ClientID)
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
						c := MockLogoutRequest(tc.key, tc.withClient)
						if tc.withClient {
							clientManager.CreateClient(context.TODO(), c.Client) // Ignore errors that are caused by duplication
						}

						_, err := m.GetLogoutRequest(context.TODO(), "challenge"+tc.key)
						require.Error(t, err)

						require.NoError(t, m.CreateLogoutRequest(context.TODO(), c))

						got2, err := m.GetLogoutRequest(context.TODO(), "challenge"+tc.key)
						require.NoError(t, err)
						assert.False(t, got2.WasUsed)
						assert.False(t, got2.Accepted)
						compareLogoutRequest(t, c, got2)

						if k%2 == 0 {
							got2, err = m.AcceptLogoutRequest(context.TODO(), "challenge"+tc.key)
							require.NoError(t, err)
							assert.True(t, got2.Accepted)
							compareLogoutRequest(t, c, got2)

							got3, err := m.VerifyAndInvalidateLogoutRequest(context.TODO(), "verifier"+tc.key)
							require.NoError(t, err)
							assert.True(t, got3.Accepted)
							assert.True(t, got3.WasUsed)
							compareLogoutRequest(t, c, got3)

							_, err = m.VerifyAndInvalidateLogoutRequest(context.TODO(), "verifier"+tc.key)
							require.Error(t, err)

							got2, err = m.GetLogoutRequest(context.TODO(), "challenge"+tc.key)
							require.NoError(t, err)
							compareLogoutRequest(t, got3, got2)
							assert.True(t, got2.WasUsed)
						} else {
							require.NoError(t, m.RejectLogoutRequest(context.TODO(), "challenge"+tc.key))
							_, err = m.GetLogoutRequest(context.TODO(), "challenge"+tc.key)
							require.Error(t, err)
						}
					})
				}
			})
		})
	}
}

func compareLogoutRequest(t *testing.T, a, b *LogoutRequest) {
	require.True(t, (a.Client != nil && b.Client != nil) || (a.Client == nil && b.Client == nil))
	if a.Client != nil {
		assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	}

	assert.EqualValues(t, a.Challenge, b.Challenge)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.PostLogoutRedirectURI, b.PostLogoutRedirectURI)
	assert.EqualValues(t, a.RPInitiated, b.RPInitiated)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareAuthenticationRequest(t *testing.T, a, b *LoginRequest) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.Challenge, b.Challenge)
	assert.EqualValues(t, *a.OpenIDConnectContext, *b.OpenIDConnectContext)
	assert.EqualValues(t, a.Subject, b.Subject)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.Verifier, b.Verifier)
	assert.EqualValues(t, a.RequestURL, b.RequestURL)
	assert.EqualValues(t, a.CSRF, b.CSRF)
	assert.EqualValues(t, a.Skip, b.Skip)
	assert.EqualValues(t, a.SessionID, b.SessionID)
}

func compareConsentRequest(t *testing.T, a, b *ConsentRequest) {
	assert.EqualValues(t, a.Client.GetID(), b.Client.GetID())
	assert.EqualValues(t, a.Challenge, b.Challenge)
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

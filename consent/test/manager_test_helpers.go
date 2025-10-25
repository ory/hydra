// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/uuidx"
)

func mockConsentRequest(remember bool, rememberFor int, skip bool) *flow.Flow {
	return &flow.Flow{
		ID:               uuidx.NewV4().String(),
		Client:           &client.Client{ID: uuidx.NewV4().String()},
		State:            flow.FlowStateConsentUsed,
		ConsentRequestID: sqlxx.NullString(uuidx.NewV4().String()),
		ConsentSkip:      skip,
		ConsentCSRF:      sqlxx.NullString(uuidx.NewV4().String()),
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues: []string{"1", "2"},
			UILocales: []string{"fr", "de"},
			Display:   "popup",
		},
		Subject:            uuidx.NewV4().String(),
		RequestedScope:     []string{"scope_a", "scope_b"},
		RequestedAudience:  []string{"aud_a", "aud_b"},
		RequestURL:         "https://request-url/path",
		RequestedAt:        time.Now().UTC(),
		ConsentRemember:    remember,
		ConsentRememberFor: pointerx.Ptr(rememberFor),
		GrantedScope:       []string{"scope_a", "scope_b"},
		GrantedAudience:    []string{"aud_a", "aud_b"},
		ConsentHandledAt:   sqlxx.NullTime(time.Now().UTC()),
	}
}

func mockLogoutRequest(withClient bool) (c *flow.LogoutRequest) {
	req := &flow.LogoutRequest{
		Subject:               uuidx.NewV4().String(),
		ID:                    uuidx.NewV4().String(),
		Verifier:              uuidx.NewV4().String(),
		SessionID:             uuidx.NewV4().String(),
		RPInitiated:           true,
		RequestURL:            "http://request-me/",
		PostLogoutRedirectURI: "http://redirect-me/",
		WasHandled:            false,
		Accepted:              false,
	}
	if withClient {
		req.Client = &client.Client{ID: uuidx.NewV4().String()}
	}
	return req
}

func LoginNIDTest(t1ValidNID, t2InvalidNID consent.LoginManager) func(t *testing.T) {
	testLS := flow.LoginSession{
		ID:      "2022-03-11-ls-nid-test-1",
		Subject: "2022-03-11-test-1-sub",
	}
	return func(t *testing.T) {
		ctx := t.Context()

		require.ErrorContains(t, t2InvalidNID.ConfirmLoginSession(ctx, &testLS), "foreign key constraint")
		require.NoError(t, t1ValidNID.ConfirmLoginSession(ctx, &testLS))
		ls, err := t2InvalidNID.DeleteLoginSession(ctx, testLS.ID)
		require.ErrorIs(t, err, sqlcon.ErrNoRows)
		assert.Nil(t, ls)
		ls, err = t1ValidNID.DeleteLoginSession(ctx, testLS.ID)
		require.NoError(t, err)
		assert.EqualValues(t, testLS.ID, ls.ID)
	}
}

type Deps interface {
	contextx.Provider
	x.TracingProvider
	x.NetworkProvider
	config.Provider
}

func LoginManagerTest(t *testing.T, deps Deps, m consent.LoginManager) {
	t.Run("get with random id", func(t *testing.T) {
		_, err := m.GetRememberedLoginSession(t.Context(), uuidx.NewV4().String())
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("create get update", func(t *testing.T) {
		sess := &flow.LoginSession{
			ID:              uuidx.NewV4().String(),
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
			Subject:         uuidx.NewV4().String(),
			Remember:        true,
		}
		require.NoError(t, m.ConfirmLoginSession(t.Context(), sess))

		actual, err := m.GetRememberedLoginSession(t.Context(), sess.ID)
		require.NoError(t, err)
		assert.Equal(t, deps.Networker().NetworkID(t.Context()), sess.NID)
		assert.Equal(t, sess, actual)

		sess.AuthenticatedAt = sqlxx.NullTime(time.Now().Add(10 * time.Minute).Round(time.Second).UTC())
		sess.Subject = uuidx.NewV4().String() // not sure why we should be able to update the subject, but ok...
		require.NoError(t, m.ConfirmLoginSession(t.Context(), sess))

		actual, err = m.GetRememberedLoginSession(t.Context(), sess.ID)
		require.NoError(t, err)
		assert.Equal(t, sess, actual)
	})

	t.Run("get non-remembered session", func(t *testing.T) {
		id := uuidx.NewV4().String()
		require.NoError(t, m.ConfirmLoginSession(t.Context(), &flow.LoginSession{
			ID:              id,
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
			Subject:         uuidx.NewV4().String(),
			Remember:        false,
		}))

		_, err := m.GetRememberedLoginSession(t.Context(), id)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("delete", func(t *testing.T) {
		expected := &flow.LoginSession{
			ID:              uuidx.NewV4().String(),
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
			Subject:         uuidx.NewV4().String(),
			Remember:        true,
		}
		require.NoError(t, m.ConfirmLoginSession(t.Context(), expected))

		deleted, err := m.DeleteLoginSession(t.Context(), expected.ID)
		require.NoError(t, err)
		assert.EqualValues(t, expected, deleted)

		_, err = m.GetRememberedLoginSession(t.Context(), expected.ID)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("revoke by subject", func(t *testing.T) {
		subs := make([]uuid.UUID, 2)
		sessions := make([]*flow.LoginSession, 0)
		for i := range subs {
			subs[i] = uuid.Must(uuid.NewV4())
			sessions = append(sessions, &flow.LoginSession{
				ID:              uuidx.NewV4().String(),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         subs[i].String(),
				Remember:        true,
			}, &flow.LoginSession{
				ID:              uuidx.NewV4().String(),
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         subs[i].String(),
				Remember:        true,
			})
		}
		sessions = append(sessions, &flow.LoginSession{
			ID:              uuidx.NewV4().String(),
			AuthenticatedAt: sqlxx.NullTime(time.Now()),
			Subject:         uuid.Must(uuid.NewV4()).String(),
			Remember:        true,
		})

		for _, s := range sessions {
			require.NoError(t, m.ConfirmLoginSession(t.Context(), s))
		}

		for _, sub := range subs {
			require.NoError(t, m.RevokeSubjectLoginSession(t.Context(), sub.String()))
		}

		for _, s := range sessions[:len(sessions)-1] {
			_, err := m.GetRememberedLoginSession(t.Context(), s.ID)
			assert.ErrorIs(t, err, x.ErrNotFound)
		}

		// ensure the unrelated session still exists
		_, err := m.GetRememberedLoginSession(t.Context(), sessions[len(sessions)-1].ID)
		assert.NoError(t, err)
	})

	t.Run("revoke with random subject", func(t *testing.T) {
		assert.NoError(t, m.RevokeSubjectLoginSession(t.Context(), uuidx.NewV4().String()))
	})
}

func ConsentManagerTests(t *testing.T, deps Deps, m consent.Manager, loginManager consent.LoginManager, clientManager client.Manager, fositeManager x.FositeStorer) {
	t.Run("case=consent-request", func(t *testing.T) {
		for _, tc := range []struct {
			key              string
			remember         bool
			rememberFor      int
			skip             bool
			expectRemembered bool
		}{
			{"1", true, 0, false, true},
			{"3", true, 1, false, false},
			{"4", false, 0, false, false},
			{"5", true, 120, false, true},
			{"6", true, 120, true, false},
		} {
			t.Run("key="+tc.key, func(t *testing.T) {
				f := mockConsentRequest(tc.remember, tc.rememberFor, tc.skip)
				_ = clientManager.CreateClient(t.Context(), f.Client) // Ignore errors that are caused by duplication
				f.NID = deps.Networker().NetworkID(t.Context())

				require.NoError(t, m.CreateConsentSession(t.Context(), f))

				t.Run("sub=detect double-submit for consent verifier", func(t *testing.T) {
					require.ErrorIs(t, m.CreateConsentSession(t.Context(), f), sqlcon.ErrUniqueViolation)
				})

				t.Run("sub=find granted and remembered consent", func(t *testing.T) {
					if tc.rememberFor == 1 {
						// unfortunately the interface does not allow us to set the absolute time, so we have to wait
						time.Sleep(2 * time.Second)
					}
					actual, err := m.FindGrantedAndRememberedConsentRequest(t.Context(), f.ClientID, f.Subject)
					if !tc.expectRemembered {
						assert.Nil(t, actual)
						assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
					} else {
						require.NoError(t, err)
						assert.NotNil(t, actual)
					}
				})
			})
		}

		for _, tc := range []struct{ keyC, keyS string }{
			{"1", "5"},
			{"5", "1"},
		} {
			t.Run(fmt.Sprintf("missmatched client %q and subject %q", tc.keyC, tc.keyS), func(t *testing.T) {
				rs, err := m.FindGrantedAndRememberedConsentRequest(t.Context(), "fk-client-"+tc.keyC, "subject"+tc.keyS)
				assert.Nil(t, rs)
				assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
			})
		}
	})

	t.Run("case=revoke consent request", func(t *testing.T) {
		type tc struct {
			at, rt, subject, client string
			revoke                  func(*testing.T)
		}
		tcs := make([]tc, 2)
		for i := range tcs {
			f := mockConsentRequest(true, 0, false)
			f.NID = deps.Networker().NetworkID(t.Context())

			tcs[i] = tc{
				subject: f.Subject,
				client:  f.Client.ID,
				at:      uuidx.NewV4().String(),
				rt:      uuidx.NewV4().String(),
			}

			require.NoError(t, clientManager.CreateClient(t.Context(), f.Client))
			require.NoError(t, m.CreateConsentSession(t.Context(), f))
			require.NoError(t, fositeManager.CreateAccessTokenSession(t.Context(), tcs[i].at,
				&fosite.Request{Client: f.Client, ID: f.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
			require.NoError(t, fositeManager.CreateRefreshTokenSession(t.Context(), tcs[i].rt, tcs[i].at,
				&fosite.Request{Client: f.Client, ID: f.ConsentRequestID.String(), RequestedAt: time.Now(), Session: &oauth2.Session{DefaultSession: openid.NewDefaultSession()}},
			))
		}
		tcs[0].revoke = func(t *testing.T) {
			require.NoError(t, m.RevokeSubjectConsentSession(t.Context(), tcs[0].subject))
		}
		tcs[1].revoke = func(t *testing.T) {
			require.NoError(t, m.RevokeSubjectClientConsentSession(t.Context(), tcs[1].subject, tcs[1].client))
		}

		for i, tc := range tcs {
			t.Run(fmt.Sprintf("run=%d", i), func(t *testing.T) {
				_, err := fositeManager.GetAccessTokenSession(t.Context(), tc.at, nil)
				require.NoError(t, err)
				_, err = fositeManager.GetRefreshTokenSession(t.Context(), tc.rt, nil)
				require.NoError(t, err)

				tc.revoke(t)

				r, err := fositeManager.GetAccessTokenSession(t.Context(), tc.at, nil)
				assert.ErrorIsf(t, err, fosite.ErrNotFound, "%+v", r)
				r, err = fositeManager.GetRefreshTokenSession(t.Context(), tc.rt, nil)
				assert.ErrorIsf(t, err, fosite.ErrNotFound, "%+v", r)
			})
		}

		t.Run("unknown subject/client return no error", func(t *testing.T) {
			require.NoError(t, m.RevokeSubjectConsentSession(t.Context(), "i-do-not-exist"))
			require.NoError(t, m.RevokeSubjectClientConsentSession(t.Context(), "i-do-not-exist", "i-do-not-exist"))
		})
	})

	t.Run("case=list consents", func(t *testing.T) {
		flows := make([]*flow.Flow, 2)
		for i := range flows {
			f := mockConsentRequest(true, 0, false)
			f.NID = deps.Networker().NetworkID(t.Context())
			f.SessionID = sqlxx.NullString(uuidx.NewV4().String())
			flows[i] = f

			require.NoError(t, clientManager.CreateClient(t.Context(), f.Client))
			require.NoError(t, loginManager.ConfirmLoginSession(t.Context(), &flow.LoginSession{
				ID:              string(f.SessionID),
				NID:             deps.Networker().NetworkID(t.Context()),
				AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).UTC()),
				Subject:         f.Subject,
			}))
			require.NoError(t, m.CreateConsentSession(t.Context(), f))
		}

		t.Run("by subject and session", func(t *testing.T) {
			for i, f := range flows {
				t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
					consents, nextPage, err := m.FindSubjectsSessionGrantedConsentRequests(t.Context(), f.Subject, f.SessionID.String())
					require.NoError(t, err)
					require.Len(t, consents, 1)

					assert.True(t, nextPage.IsLast())
					assert.Equal(t, f.ConsentRequestID, consents[0].ConsentRequestID)
					assert.Equal(t, f.Client.ID, consents[0].Client.GetID())
				})
			}

			t.Run("random subject", func(t *testing.T) {
				_, _, err := m.FindSubjectsSessionGrantedConsentRequests(t.Context(), uuidx.NewV4().String(), flows[0].SessionID.String())
				assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
			})
		})

		for i, f := range flows {
			t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
				consents, nextPage, err := m.FindSubjectsGrantedConsentRequests(t.Context(), f.Subject)
				require.NoError(t, err)
				require.Len(t, consents, 1)
				assert.True(t, nextPage.IsLast())

				assert.Equal(t, f.ConsentRequestID, consents[0].ConsentRequestID)
				assert.Equal(t, f.Client.ID, consents[0].Client.GetID())
			})

			t.Run("random subject", func(t *testing.T) {
				_, _, err := m.FindSubjectsGrantedConsentRequests(t.Context(), uuidx.NewV4().String())
				assert.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
			})
		}

		t.Run("case=ListUserAuthenticatedClientsWithFrontAndBackChannelLogout", func(t *testing.T) {
			// The idea of this test is to create two identities (subjects) with 4 sessions each, where
			// only some sessions have been associated with a client that has a front channel logout url

			subjects := make([]string, 2)
			for k := range subjects {
				subjects[k] = fmt.Sprintf("subject-ListUserAuthenticatedClientsWithFrontAndBackChannelLogout-%d", k)
			}

			sessions := make([]flow.LoginSession, len(subjects)*4)
			frontChannels := map[string][]client.Client{}
			backChannels := map[string][]client.Client{}
			for k := range sessions {
				id := uuidx.NewV4().String()
				subject := subjects[k%len(subjects)]
				t.Run(fmt.Sprintf("create/session=%s/subject=%s", id, subject), func(t *testing.T) {
					ls := &flow.LoginSession{
						ID:              id,
						NID:             deps.Networker().NetworkID(t.Context()),
						AuthenticatedAt: sqlxx.NullTime(time.Now()),
						Subject:         subject,
						Remember:        true,
					}
					require.NoError(t, loginManager.ConfirmLoginSession(t.Context(), ls))

					cl := &client.Client{ID: uuidx.NewV4().String()}
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
					require.NoError(t, clientManager.CreateClient(t.Context(), cl))

					f := &flow.Flow{
						NID: deps.Networker().NetworkID(t.Context()),
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
						ConsentCSRF:      sqlxx.NullString(uuid.Must(uuid.NewV4()).String()),

						LoginCSRF: uuid.Must(uuid.NewV4()).String(),
						ID:        uuid.Must(uuid.NewV4()).String(),
						State:     flow.FlowStateLoginUsed,
					}

					require.NoError(t, m.CreateConsentSession(t.Context(), f))

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
					actual, err := m.ListUserAuthenticatedClientsWithFrontChannelLogout(t.Context(), ls.Subject, ls.ID)
					require.NoError(t, err)
					check(t, frontChannels, actual)
				})

				t.Run(fmt.Sprintf("method=ListUserAuthenticatedClientsWithBackChannelLogout/session=%s", ls.ID), func(t *testing.T) {
					actual, err := m.ListUserAuthenticatedClientsWithBackChannelLogout(t.Context(), ls.Subject, ls.ID)
					require.NoError(t, err)
					check(t, backChannels, actual)
				})
			}
		})
	})
}

func ObfuscatedSubjectManagerTest(t *testing.T, deps Deps, m consent.ObfuscatedSubjectManager, clientManager client.Manager) {
	t.Run("get with random keys", func(t *testing.T) {
		_, err := m.GetForcedObfuscatedLoginSession(t.Context(), uuidx.NewV4().String(), uuidx.NewV4().String())
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("create and retrieve", func(t *testing.T) {
		cl := &client.Client{ID: uuidx.NewV4().String()}
		require.NoError(t, clientManager.CreateClient(t.Context(), cl))
		obfuscatedSession := &consent.ForcedObfuscatedLoginSession{
			ClientID:          cl.ID,
			Subject:           uuidx.NewV4().String(),
			SubjectObfuscated: uuidx.NewV4().String(),
			NID:               deps.Networker().NetworkID(t.Context()),
		}
		require.NoError(t, m.CreateForcedObfuscatedLoginSession(t.Context(), obfuscatedSession))

		actual, err := m.GetForcedObfuscatedLoginSession(t.Context(), cl.ID, obfuscatedSession.SubjectObfuscated)
		require.NoError(t, err)
		assert.EqualValues(t, obfuscatedSession, actual)

		t.Run("with random client fails", func(t *testing.T) {
			_, err = m.GetForcedObfuscatedLoginSession(t.Context(), uuidx.NewV4().String(), obfuscatedSession.SubjectObfuscated)
			assert.ErrorIs(t, err, x.ErrNotFound)
		})

		t.Run("with random obfuscated subject fails", func(t *testing.T) {
			_, err = m.GetForcedObfuscatedLoginSession(t.Context(), cl.ID, uuidx.NewV4().String())
			assert.ErrorIs(t, err, x.ErrNotFound)
		})
	})
}

func LogoutManagerTest(t *testing.T, m consent.LogoutManager, clientManager client.Manager) {
	for _, withClient := range []bool{true, false} {
		t.Run("get with random challenge", func(t *testing.T) {
			_, err := m.GetLogoutRequest(t.Context(), uuidx.NewV4().String())
			assert.ErrorIs(t, err, sqlcon.ErrNoRows)
		})

		t.Run(fmt.Sprintf("with client=%v", withClient), func(t *testing.T) {
			setup := func(t *testing.T) *flow.LogoutRequest {
				req := mockLogoutRequest(withClient)
				if withClient {
					require.NoError(t, clientManager.CreateClient(t.Context(), req.Client))
				}
				require.NoError(t, m.CreateLogoutRequest(t.Context(), req))
				return req
			}

			t.Run("get unhandled", func(t *testing.T) {
				expected := setup(t)

				actual, err := m.GetLogoutRequest(t.Context(), expected.ID)
				require.NoError(t, err)
				assert.False(t, actual.WasHandled)
				assert.False(t, actual.Accepted)
				compareLogoutRequest(t, expected, actual)
			})

			t.Run("accept and verify", func(t *testing.T) {
				expected := setup(t)

				actual, err := m.AcceptLogoutRequest(t.Context(), expected.ID)
				require.NoError(t, err)
				assert.True(t, actual.Accepted)
				assert.False(t, actual.WasHandled)
				compareLogoutRequest(t, expected, actual)

				actual, err = m.VerifyAndInvalidateLogoutRequest(t.Context(), expected.Verifier)
				require.NoError(t, err)
				assert.True(t, actual.Accepted)
				assert.True(t, actual.WasHandled)
				compareLogoutRequest(t, expected, actual)

				t.Run("double verify fails", func(t *testing.T) {
					_, err = m.VerifyAndInvalidateLogoutRequest(t.Context(), expected.Verifier)
					require.NotErrorIs(t, err, x.ErrNotFound)
				})

				t.Run("get verified", func(t *testing.T) {
					actual, err = m.GetLogoutRequest(t.Context(), expected.ID)
					require.NoError(t, err)
					assert.True(t, actual.WasHandled)
					assert.True(t, actual.Accepted)
					compareLogoutRequest(t, expected, actual)
				})
			})

			t.Run("reject", func(t *testing.T) {
				expected := setup(t)

				require.NoError(t, m.RejectLogoutRequest(t.Context(), expected.ID))
				_, err := m.GetLogoutRequest(t.Context(), expected.ID)
				assert.ErrorIs(t, err, sqlcon.ErrNoRows)
			})
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

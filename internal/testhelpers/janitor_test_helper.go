package testhelpers

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/hydra/x"
	"github.com/ory/x/logrusx"

	"github.com/ory/x/sqlxx"
)

type JanitorConsentTestHelper struct {
	uniqueName           string
	flushLoginRequests   []*consent.LoginRequest
	flushConsentRequests []*consent.ConsentRequest
	flushAccessRequests  []*fosite.Request
	flushRefreshRequests []*fosite.AccessRequest
	flushGrants          []*createGrantRequest
	conf                 *config.Provider
	Lifespan             time.Duration
}

type createGrantRequest struct {
	grant trust.Grant
	pk    jose.JSONWebKey
}

const lifespan = time.Hour

func NewConsentJanitorTestHelper(uniqueName string) *JanitorConsentTestHelper {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyAccessTokenLifespan, lifespan)
	conf.MustSet(config.KeyRefreshTokenLifespan, lifespan)
	conf.MustSet(config.KeyConsentRequestMaxAge, lifespan)
	conf.MustSet(config.KeyLogLevel, "trace")

	return &JanitorConsentTestHelper{
		uniqueName:           uniqueName,
		conf:                 conf,
		flushLoginRequests:   genLoginRequests(uniqueName, lifespan),
		flushConsentRequests: genConsentRequests(uniqueName, lifespan),
		flushAccessRequests:  getAccessRequests(uniqueName, lifespan),
		flushRefreshRequests: getRefreshRequests(uniqueName, lifespan),
		flushGrants:          getGrantRequests(uniqueName, lifespan),
		Lifespan:             lifespan,
	}
}

func (j *JanitorConsentTestHelper) GetDSN() string {
	return j.conf.DSN()
}

func (j *JanitorConsentTestHelper) GetConfig() *config.Provider {
	return j.conf
}

func (j *JanitorConsentTestHelper) GetNotAfterTestCycles() map[string]time.Duration {
	return map[string]time.Duration{
		"notAfter24h":   j.Lifespan * 24,
		"notAfter1h30m": j.Lifespan + time.Hour/2,
		"notAfterNow":   0,
	}
}

func (j *JanitorConsentTestHelper) GetRegistry(ctx context.Context, dbname string) (driver.Registry, error) {
	j.conf.MustSet(config.KeyDSN, fmt.Sprintf("sqlite://file:%s?mode=memory&_fk=true&cache=shared", dbname))
	return driver.NewRegistryFromDSN(ctx, j.conf, logrusx.New("test_hydra", "master"))
}

func (j *JanitorConsentTestHelper) AccessTokenNotAfterSetup(ctx context.Context, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		// Create access token clients and session
		for _, r := range j.flushAccessRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client.(*client.Client)))
			require.NoError(t, store.CreateAccessTokenSession(ctx, r.ID, r))
		}

	}
}

func (j *JanitorConsentTestHelper) AccessTokenNotAfterValidate(ctx context.Context, notAfter time.Time, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		ds := new(oauth2.Session)

		accessTokenLifespan := time.Now().Round(time.Second).Add(-j.conf.AccessTokenLifespan())

		for _, r := range j.flushAccessRequests {
			t.Logf("access flush check: %s", r.ID)
			_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
			if j.notAfterCheck(notAfter, accessTokenLifespan, r.RequestedAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) RefreshTokenNotAfterSetup(ctx context.Context, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		// Create refresh token clients and session
		for _, fr := range j.flushRefreshRequests {
			require.NoError(t, cl.CreateClient(ctx, fr.Client.(*client.Client)))
			require.NoError(t, store.CreateRefreshTokenSession(ctx, fr.ID, fr))
		}
	}
}

func (j *JanitorConsentTestHelper) RefreshTokenNotAfterValidate(ctx context.Context, notAfter time.Time, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		ds := new(oauth2.Session)

		refreshTokenLifespan := time.Now().Round(time.Second).Add(-j.conf.RefreshTokenLifespan())

		for _, r := range j.flushRefreshRequests {
			t.Logf("refresh flush check: %s", r.ID)
			_, err = store.GetRefreshTokenSession(ctx, r.ID, ds)
			if j.notAfterCheck(notAfter, refreshTokenLifespan, r.RequestedAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) GrantNotAfterSetup(ctx context.Context, cl client.Manager, gr trust.GrantManager) func(t *testing.T) {
	return func(t *testing.T) {
		for _, fg := range j.flushGrants {
			require.NoError(t, gr.CreateGrant(ctx, fg.grant, fg.pk))
		}
	}
}

func (j *JanitorConsentTestHelper) GrantNotAfterValidate(ctx context.Context, notAfter time.Time, gr trust.GrantManager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// flush won't delete grants that have not yet expired, so use now to check that
		deleteUntil := time.Now().Round(time.Second)
		if deleteUntil.After(notAfter) {
			deleteUntil = notAfter
		}

		for _, r := range j.flushGrants {
			t.Logf("grant flush check: %s", r.grant.Issuer)
			_, err = gr.GetConcreteGrant(ctx, r.grant.ID)

			if deleteUntil.After(r.grant.ExpiresAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) LoginRejectionSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// Create login requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		// Explicit rejection
		for _, r := range j.flushLoginRequests {
			if r.ID == j.flushLoginRequests[0].ID {
				// accept this one
				_, err = cm.HandleLoginRequest(ctx, r.ID, consent.NewHandledLoginRequest(
					r.ID, false, r.RequestedAt, r.AuthenticatedAt))

				require.NoError(t, err)
				continue
			}

			// reject flush-login-2 and 3
			_, err = cm.HandleLoginRequest(ctx, r.ID, consent.NewHandledLoginRequest(
				r.ID, true, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
		}
	}
}

func (j *JanitorConsentTestHelper) LoginRejectionValidate(ctx context.Context, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		// flush-login-2 and 3 should be cleared now
		for _, r := range j.flushLoginRequests {
			t.Logf("check login: %s", r.ID)
			_, err := cm.GetLoginRequest(ctx, r.ID)
			if r.ID == j.flushLoginRequests[0].ID {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) LimitSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// Create login requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		// Reject each request
		for _, r := range j.flushLoginRequests {
			_, err = cm.HandleLoginRequest(ctx, r.ID, consent.NewHandledLoginRequest(
				r.ID, true, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
		}
	}
}

func (j *JanitorConsentTestHelper) LimitValidate(ctx context.Context, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		// flush-login-2 and 3 should be cleared now
		for _, r := range j.flushLoginRequests {
			t.Logf("check login: %s", r.ID)
			_, err := cm.GetLoginRequest(ctx, r.ID)
			if r.ID == j.flushLoginRequests[0].ID {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) ConsentRejectionSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// Create login requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		// Create consent requests
		for _, r := range j.flushConsentRequests {
			require.NoError(t, cm.CreateConsentRequest(ctx, r))
		}

		//Reject the consents
		for _, r := range j.flushConsentRequests {
			if r.ID == j.flushConsentRequests[0].ID {
				// accept this one
				_, err = cm.HandleConsentRequest(ctx, r.ID, consent.NewHandledConsentRequest(
					r.ID, false, r.RequestedAt, r.AuthenticatedAt))
				require.NoError(t, err)
				continue
			}
			_, err = cm.HandleConsentRequest(ctx, r.ID, consent.NewHandledConsentRequest(
				r.ID, true, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
		}
	}
}

func (j *JanitorConsentTestHelper) ConsentRejectionValidate(ctx context.Context, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		for _, r := range j.flushConsentRequests {
			t.Logf("check consent: %s", r.ID)
			_, err = cm.GetConsentRequest(ctx, r.ID)
			if r.ID == j.flushConsentRequests[0].ID {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) LoginTimeoutSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// Create login requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		// Creating at least 1 that has not timed out
		_, err = cm.HandleLoginRequest(ctx, j.flushLoginRequests[0].ID, &consent.HandledLoginRequest{
			ID:              j.flushLoginRequests[0].ID,
			RequestedAt:     j.flushLoginRequests[0].RequestedAt,
			AuthenticatedAt: j.flushLoginRequests[0].AuthenticatedAt,
			WasHandled:      true,
		})

		require.NoError(t, err)
	}
}

func (j *JanitorConsentTestHelper) LoginTimeoutValidate(ctx context.Context, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		for _, r := range j.flushLoginRequests {
			_, err = cm.GetLoginRequest(ctx, r.ID)
			if r.ID == j.flushLoginRequests[0].ID {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}

	}
}

func (j *JanitorConsentTestHelper) ConsentTimeoutSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// Let's reset and accept all login requests to test the consent requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
			_, err = cm.HandleLoginRequest(ctx, r.ID, &consent.HandledLoginRequest{
				ID:              r.ID,
				AuthenticatedAt: r.AuthenticatedAt,
				RequestedAt:     r.RequestedAt,
				WasHandled:      true,
			})
			require.NoError(t, err)
		}

		// Create consent requests
		for _, r := range j.flushConsentRequests {
			require.NoError(t, cm.CreateConsentRequest(ctx, r))
		}

		// Create at least 1 consent request that has been accepted
		_, err = cm.HandleConsentRequest(ctx, j.flushConsentRequests[0].ID, &consent.HandledConsentRequest{
			ID:              j.flushConsentRequests[0].ID,
			WasHandled:      true,
			RequestedAt:     j.flushConsentRequests[0].RequestedAt,
			AuthenticatedAt: j.flushConsentRequests[0].AuthenticatedAt,
		})
		require.NoError(t, err)
	}
}

func (j *JanitorConsentTestHelper) ConsentTimeoutValidate(ctx context.Context, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		for _, r := range j.flushConsentRequests {
			_, err = cm.GetConsentRequest(ctx, r.ID)
			if r.ID == j.flushConsentRequests[0].ID {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}

	}
}

func (j *JanitorConsentTestHelper) LoginConsentNotAfterSetup(ctx context.Context, cm consent.Manager, cl client.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		for _, r := range j.flushConsentRequests {
			require.NoError(t, cm.CreateConsentRequest(ctx, r))
		}
	}
}

func (j *JanitorConsentTestHelper) LoginConsentNotAfterValidate(ctx context.Context, notAfter time.Time, consentRequestLifespan time.Time, cm consent.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		for _, r := range j.flushLoginRequests {
			t.Logf("login flush check:\nNotAfter: %s\nConsentRequest: %s\n%+v\n",
				notAfter.String(), consentRequestLifespan.String(), r)
			_, err = cm.GetLoginRequest(ctx, r.ID)
			// if the lowest between notAfter and consent-request-lifespan is greater than requested_at
			// then the it should expect the value to be deleted.
			if j.notAfterCheck(notAfter, consentRequestLifespan, r.RequestedAt) {
				// value has been deleted here
				require.Error(t, err)
			} else {
				// value has not been deleted here
				require.NoError(t, err)
			}
		}

		for _, r := range j.flushConsentRequests {
			t.Logf("consent flush check:\nNotAfter: %s\nConsentRequest: %s\n%+v\n",
				notAfter.String(), consentRequestLifespan.String(), r)
			_, err = cm.GetConsentRequest(ctx, r.ID)
			// if the lowest between notAfter and consent-request-lifespan is greater than requested_at
			// then the it should expect the value to be deleted.
			if j.notAfterCheck(notAfter, consentRequestLifespan, r.RequestedAt) {
				// value has been deleted here
				require.Error(t, err)
			} else {
				// value has not been deleted here
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) GetConsentRequestLifespan() time.Duration {
	return j.conf.ConsentRequestMaxAge()
}

func (j *JanitorConsentTestHelper) GetAccessTokenLifespan() time.Duration {
	return j.conf.AccessTokenLifespan()
}

func (j *JanitorConsentTestHelper) GetRefreshTokenLifespan() time.Duration {
	return j.conf.RefreshTokenLifespan()
}

func (j *JanitorConsentTestHelper) notAfterCheck(notAfter time.Time, lifespan time.Time, requestedAt time.Time) bool {
	// The database deletes where requested_at time is smaller than the lowest between notAfter and consent-request-lifespan
	// thus we get the lowest value here first to compare later to requested_at
	var lesser time.Time
	// if the lowest between notAfter and consent-request-lifespan is greater than requested_at
	// then the it should expect the value to be deleted.
	if notAfter.Unix() < lifespan.Unix() {
		lesser = notAfter
	} else {
		lesser = lifespan
	}

	// true: value has been deleted
	// false: value still exists
	return lesser.Unix() > requestedAt.Unix()
}

func JanitorTests(conf *config.Provider, consentManager consent.Manager, clientManager client.Manager, fositeManager x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		jt := NewConsentJanitorTestHelper(t.Name())

		conf.MustSet(config.KeyConsentRequestMaxAge, jt.GetConsentRequestLifespan())

		t.Run("case=flush-consent-request-not-after", func(t *testing.T) {

			notAfterTests := jt.GetNotAfterTestCycles()

			for k, v := range notAfterTests {
				jt := NewConsentJanitorTestHelper(k)
				t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
					notAfter := time.Now().Round(time.Second).Add(-v)
					consentRequestLifespan := time.Now().Round(time.Second).Add(-jt.GetConsentRequestLifespan())

					// setup test
					t.Run("step=setup", jt.LoginConsentNotAfterSetup(ctx, consentManager, clientManager))

					// run the cleanup routine
					t.Run("step=cleanup", func(t *testing.T) {
						require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, notAfter, 1000, 100))
					})

					// validate test
					t.Run("step=validate", jt.LoginConsentNotAfterValidate(ctx, notAfter, consentRequestLifespan, consentManager))
				})

			}
		})

		t.Run("case=flush-consent-request-limit", func(t *testing.T) {
			jt := NewConsentJanitorTestHelper("limit")

			t.Run("case=limit", func(t *testing.T) {
				// setup
				t.Run("step=setup", jt.LimitSetup(ctx, consentManager, clientManager))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, time.Now().Round(time.Second), 2, 1))
				})

				// validate
				t.Run("step=validate", jt.LimitValidate(ctx, consentManager))
			})
		})

		t.Run("case=flush-consent-request-rejection", func(t *testing.T) {
			jt := NewConsentJanitorTestHelper("loginRejection")

			t.Run(fmt.Sprintf("case=%s", "loginRejection"), func(t *testing.T) {
				// setup
				t.Run("step=setup", jt.LoginRejectionSetup(ctx, consentManager, clientManager))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, time.Now().Round(time.Second), 1000, 100))
				})

				// validate
				t.Run("step=validate", jt.LoginRejectionValidate(ctx, consentManager))
			})

			jt = NewConsentJanitorTestHelper("consentRejection")

			t.Run(fmt.Sprintf("case=%s", "consentRejection"), func(t *testing.T) {
				// setup
				t.Run("step=setup", jt.ConsentRejectionSetup(ctx, consentManager, clientManager))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, time.Now().Round(time.Second), 1000, 100))
				})

				// validate
				t.Run("step=validate", jt.ConsentRejectionValidate(ctx, consentManager))
			})

		})

		t.Run("case=flush-consent-request-timeout", func(t *testing.T) {
			jt := NewConsentJanitorTestHelper("loginTimeout")

			t.Run(fmt.Sprintf("case=%s", "login-timeout"), func(t *testing.T) {

				// setup
				t.Run("step=setup", jt.LoginTimeoutSetup(ctx, consentManager, clientManager))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, time.Now().Round(time.Second), 1000, 100))
				})

				// validate
				t.Run("step=validate", jt.LoginTimeoutValidate(ctx, consentManager))

			})

			jt = NewConsentJanitorTestHelper("consentTimeout")

			t.Run(fmt.Sprintf("case=%s", "consent-timeout"), func(t *testing.T) {

				// setup
				t.Run("step=setup", jt.ConsentTimeoutSetup(ctx, consentManager, clientManager))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					require.NoError(t, fositeManager.FlushInactiveLoginConsentRequests(ctx, time.Now().Round(time.Second), 1000, 100))
				})

				// validate
				t.Run("step=validate", jt.ConsentTimeoutValidate(ctx, consentManager))

			})
		})
	}
}

func getAccessRequests(uniqueName string, lifespan time.Duration) []*fosite.Request {
	return []*fosite.Request{
		{
			ID:             fmt.Sprintf("%s_flush-access-1", uniqueName),
			RequestedAt:    time.Now().Round(time.Second),
			Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-1", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
		{
			ID:             fmt.Sprintf("%s_flush-access-2", uniqueName),
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-2", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
		{
			ID:             fmt.Sprintf("%s_flush-access-3", uniqueName),
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-3", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
	}
}

func getRefreshRequests(uniqueName string, lifespan time.Duration) []*fosite.AccessRequest {
	var tokenSignature = "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"
	return []*fosite.AccessRequest{
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second),
				ID:             fmt.Sprintf("%s_flush-refresh-1", uniqueName),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-1", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-1", uniqueName), tokenSignature)},
				},
			},
		},
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
				ID:             fmt.Sprintf("%s_flush-refresh-2", uniqueName),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-2", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-2", uniqueName), tokenSignature)},
				},
			},
		},
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
				ID:             fmt.Sprintf("%s_flush-refresh-3", uniqueName),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-3", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-3", uniqueName), tokenSignature)},
				},
			},
		},
	}
}

func genLoginRequests(uniqueName string, lifespan time.Duration) []*consent.LoginRequest {
	return []*consent.LoginRequest{
		{
			ID:             fmt.Sprintf("%s_flush-login-1", uniqueName),
			RequestedScope: []string{"foo", "bar"},
			Subject:        fmt.Sprintf("%s_flush-login-1", uniqueName),
			Client: &client.Client{
				OutfacingID:  fmt.Sprintf("%s_flush-login-consent-1", uniqueName),
				RedirectURIs: []string{"http://redirect"},
			},
			RequestURL:      "http://redirect",
			RequestedAt:     time.Now().Round(time.Second),
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second)),
			Verifier:        fmt.Sprintf("%s_flush-login-1", uniqueName),
		},
		{
			ID:             fmt.Sprintf("%s_flush-login-2", uniqueName),
			RequestedScope: []string{"foo", "bar"},
			Subject:        fmt.Sprintf("%s_flush-login-2", uniqueName),
			Client: &client.Client{
				OutfacingID:  fmt.Sprintf("%s_flush-login-consent-2", uniqueName),
				RedirectURIs: []string{"http://redirect"},
			},
			RequestURL:      "http://redirect",
			RequestedAt:     time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-(lifespan + time.Minute))),
			Verifier:        fmt.Sprintf("%s_flush-login-2", uniqueName),
		},
		{
			ID:             fmt.Sprintf("%s_flush-login-3", uniqueName),
			RequestedScope: []string{"foo", "bar"},
			Subject:        fmt.Sprintf("%s_flush-login-3", uniqueName),
			Client: &client.Client{
				OutfacingID:  fmt.Sprintf("%s_flush-login-consent-3", uniqueName),
				RedirectURIs: []string{"http://redirect"},
			},
			RequestURL:      "http://redirect",
			RequestedAt:     time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-(lifespan + time.Hour))),
			Verifier:        fmt.Sprintf("%s_flush-login-3", uniqueName),
		},
	}
}

func genConsentRequests(uniqueName string, lifespan time.Duration) []*consent.ConsentRequest {
	return []*consent.ConsentRequest{
		{
			ID:                   fmt.Sprintf("%s_flush-consent-1", uniqueName),
			RequestedScope:       []string{"foo", "bar"},
			Subject:              fmt.Sprintf("%s_flush-consent-1", uniqueName),
			OpenIDConnectContext: nil,
			ClientID:             fmt.Sprintf("%s_flush-login-consent-1", uniqueName),
			RequestURL:           "http://redirect",
			LoginChallenge:       sqlxx.NullString(fmt.Sprintf("%s_flush-login-1", uniqueName)),
			RequestedAt:          time.Now().Round(time.Second),
			Verifier:             fmt.Sprintf("%s_flush-consent-1", uniqueName),
		},
		{
			ID:                   fmt.Sprintf("%s_flush-consent-2", uniqueName),
			RequestedScope:       []string{"foo", "bar"},
			Subject:              fmt.Sprintf("%s_flush-consent-2", uniqueName),
			OpenIDConnectContext: nil,
			ClientID:             fmt.Sprintf("%s_flush-login-consent-2", uniqueName),
			RequestURL:           "http://redirect",
			LoginChallenge:       sqlxx.NullString(fmt.Sprintf("%s_flush-login-2", uniqueName)),
			RequestedAt:          time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			Verifier:             fmt.Sprintf("%s_flush-consent-2", uniqueName),
		},
		{
			ID:                   fmt.Sprintf("%s_flush-consent-3", uniqueName),
			RequestedScope:       []string{"foo", "bar"},
			Subject:              fmt.Sprintf("%s_flush-consent-3", uniqueName),
			OpenIDConnectContext: nil,
			ClientID:             fmt.Sprintf("%s_flush-login-consent-3", uniqueName),
			RequestURL:           "http://redirect",
			LoginChallenge:       sqlxx.NullString(fmt.Sprintf("%s_flush-login-3", uniqueName)),
			RequestedAt:          time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			Verifier:             fmt.Sprintf("%s_flush-consent-3", uniqueName),
		},
	}
}

func getGrantRequests(uniqueName string, lifespan time.Duration) []*createGrantRequest {
	return []*createGrantRequest{
		{
			grant: trust.Grant{
				ID:      uuid.New().String(),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-1", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-1", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-1", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-1", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second),
				ExpiresAt: time.Now().Round(time.Second).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-1", uniqueName),
			},
		},
		{
			grant: trust.Grant{
				ID:      uuid.New().String(),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-2", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-2", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-2", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-2", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
				ExpiresAt: time.Now().Round(time.Second).Add(-(lifespan + time.Minute)).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-2", uniqueName),
			},
		},
		{
			grant: trust.Grant{
				ID:      uuid.New().String(),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-3", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-3", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-3", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-3", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
				ExpiresAt: time.Now().Round(time.Second).Add(-(lifespan + time.Hour)).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-3", uniqueName),
			},
		},
	}
}

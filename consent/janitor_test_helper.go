package consent

import (
	"context"
	"fmt"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/x"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/ory/x/sqlxx"
)

type cleanupRoutine = func() error

type JanitorConsentTestHelper struct {
	flushLoginRequests   []*LoginRequest
	flushConsentRequests []*ConsentRequest
	Lifespan             time.Duration
}

func NewConsentJanitorTestHelper(uniqueName string) *JanitorConsentTestHelper {
	var lifespan = time.Hour
	return &JanitorConsentTestHelper{
		flushLoginRequests: []*LoginRequest{
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
		},
		flushConsentRequests: []*ConsentRequest{
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
		},
		Lifespan: lifespan,
	}
}

func (j *JanitorConsentTestHelper) NewHandledLoginRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *HandledLoginRequest {
	var deniedErr *RequestDeniedError
	if hasError {
		deniedErr = &RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			valid:       true,
		}
	}

	return &HandledLoginRequest{
		ID:              challenge,
		Error:           deniedErr,
		WasUsed:         true,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
	}
}

func (j *JanitorConsentTestHelper) NewHandledConsentRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *HandledConsentRequest {
	var deniedErr *RequestDeniedError
	if hasError {
		deniedErr = &RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			valid:       true,
		}
	}

	return &HandledConsentRequest{
		ID:              challenge,
		HandledAt:       sqlxx.NullTime(time.Now().Round(time.Second)),
		Error:           deniedErr,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
		WasUsed:         true,
	}
}

func (j *JanitorConsentTestHelper) LoginRejection(cr cleanupRoutine, cm Manager, cl client.Manager) func(t *testing.T) {
	ctx := context.Background()

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
				_, err = cm.HandleLoginRequest(ctx, r.ID, j.NewHandledLoginRequest(
					r.ID, false, r.RequestedAt, r.AuthenticatedAt))

				require.NoError(t, err)
				continue
			}

			// reject flush-login-2 and 3
			_, err = cm.HandleLoginRequest(ctx, r.ID, j.NewHandledLoginRequest(
				r.ID, true, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
		}

		// run the cleanup routine
		require.NoError(t, cr())

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

func (j *JanitorConsentTestHelper) ConsentRejection(cr cleanupRoutine, cm Manager, cl client.Manager) func(t *testing.T) {
	ctx := context.Background()

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
				_, err = cm.HandleConsentRequest(ctx, r.ID, j.NewHandledConsentRequest(
					r.ID, false, r.RequestedAt, r.AuthenticatedAt))
				require.NoError(t, err)
				continue
			}
			_, err = cm.HandleConsentRequest(ctx, r.ID, j.NewHandledConsentRequest(
				r.ID, true, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
		}

		// run the cleanup routine
		require.NoError(t, cr())

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

func (j *JanitorConsentTestHelper) LoginTimeout(cr cleanupRoutine, cm Manager, cl client.Manager) func(t *testing.T) {
	ctx := context.Background()

	return func(t *testing.T) {
		var err error

		// Create login requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		// Creating at least 1 that has not timed out
		_, err = cm.HandleLoginRequest(ctx, j.flushLoginRequests[0].ID, &HandledLoginRequest{
			ID:              j.flushLoginRequests[0].ID,
			RequestedAt:     j.flushLoginRequests[0].RequestedAt,
			AuthenticatedAt: j.flushLoginRequests[0].AuthenticatedAt,
			WasUsed:         true,
		})

		require.NoError(t, err)

		//run the cleanup routine
		require.NoError(t, cr())

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

func (j *JanitorConsentTestHelper) ConsentTimeout(cr cleanupRoutine, cm Manager, cl client.Manager) func(t *testing.T) {
	ctx := context.Background()

	return func(t *testing.T) {
		var err error

		// Let's reset and accept all login requests to test the consent requests
		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
			_, err = cm.HandleLoginRequest(ctx, r.ID, &HandledLoginRequest{
				ID:              r.ID,
				AuthenticatedAt: r.AuthenticatedAt,
				RequestedAt:     r.RequestedAt,
				WasUsed:         true,
			})
			require.NoError(t, err)
		}

		// Create consent requests
		for _, r := range j.flushConsentRequests {
			require.NoError(t, cm.CreateConsentRequest(ctx, r))
		}

		// Create at least 1 consent request that has been accepted
		_, err = cm.HandleConsentRequest(ctx, j.flushConsentRequests[0].ID, &HandledConsentRequest{
			ID:              j.flushConsentRequests[0].ID,
			WasUsed:         true,
			RequestedAt:     j.flushConsentRequests[0].RequestedAt,
			AuthenticatedAt: j.flushConsentRequests[0].AuthenticatedAt,
		})

		require.NoError(t, err)

		// Run the cleanup routine
		require.NoError(t, cr())

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

func (j *JanitorConsentTestHelper) LoginConsentNotAfter(notAfter time.Time, consentRequestLifespan time.Duration, cr cleanupRoutine, cm Manager, cl client.Manager) func(t *testing.T) {
	ctx := context.Background()

	return func(t *testing.T) {
		var err error

		var lesser time.Time
		if notAfter.Unix() < time.Now().Add(-consentRequestLifespan).Unix() {
			lesser = notAfter
		} else {
			lesser = time.Now().Add(-consentRequestLifespan)
		}

		for _, r := range j.flushLoginRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client))
			require.NoError(t, cm.CreateLoginRequest(ctx, r))
		}

		for _, r := range j.flushConsentRequests {
			require.NoError(t, cm.CreateConsentRequest(ctx, r))
		}

		// run the cleanup routine
		require.NoError(t, cr())

		for _, r := range j.flushLoginRequests {
			t.Logf("login flush check:\nNotAfter: %s\nGreater: %s\nConsentRequest: %s\n%+v\n",
				notAfter.String(), lesser.String(), time.Now().Add(-consentRequestLifespan).String(), r)
			_, err = cm.GetLoginRequest(ctx, r.ID)
			if lesser.Unix() >= r.RequestedAt.Unix() {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}

		for _, r := range j.flushConsentRequests {
			t.Logf("consent flush check:\nNotAfter: %s\nGreater: %s\nConsentRequest: %s\n%+v\n",
				notAfter.String(), lesser.String(), time.Now().Add(-consentRequestLifespan).String(), r)
			_, err = cm.GetConsentRequest(ctx, r.ID)
			if lesser.Unix() >= r.RequestedAt.Unix() {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}

	}
}

func JanitorTests(conf *config.Provider, consentManager Manager, clientManager client.Manager, fositeManager x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		jt := NewConsentJanitorTestHelper("")
		conf.MustSet(config.KeyConsentRequestMaxAge, jt.Lifespan)

		t.Run("case=flush-consent-request-not-after", func(t *testing.T) {
			ctx := context.Background()

			var flushHandler = func(notAfter time.Time) func() error {
				return func() error {

					if err := fositeManager.FlushInactiveLoginConsentRequests(ctx, notAfter); err != nil {
						return err
					}

					return nil

				}
			}

			notAfterTests := map[string]time.Time{
				"notAfter24h":   time.Now().Round(time.Second).Add(-(jt.Lifespan * 24)),
				"notAfter1h30m": time.Now().Round(time.Second).Add(-(jt.Lifespan + time.Hour/2)),
				"notAfterNow":   time.Now().Round(time.Second),
			}

			for k, v := range notAfterTests {
				t.Run(fmt.Sprintf("case=%s", k),
					NewConsentJanitorTestHelper(k).LoginConsentNotAfter(v, conf.ConsentRequestMaxAge(), flushHandler(v), consentManager, clientManager))
			}
		})

		t.Run("case=flush-consent-request", func(t *testing.T) {
			ctx := context.Background()

			var flushHandler = func(notAfter time.Time) func() error {
				return func() error {

					if err := fositeManager.FlushInactiveLoginConsentRequests(ctx, notAfter); err != nil {
						return err
					}
					return nil

				}
			}

			type loginConsentTest = func(func() error, Manager, client.Manager) func(t *testing.T)

			loginConsentTests := map[string]loginConsentTest{
				"loginRejection":   NewConsentJanitorTestHelper("loginRejection").LoginRejection,
				"loginTimeout":     NewConsentJanitorTestHelper("loginTimeout").LoginTimeout,
				"consentRejection": NewConsentJanitorTestHelper("consentRejection").ConsentRejection,
				"consentTimeout":   NewConsentJanitorTestHelper("consentTimeout").ConsentTimeout,
			}

			for k, v := range loginConsentTests {
				t.Run(fmt.Sprintf("case=%s", k),
					v(flushHandler(time.Now().Round(time.Second)), consentManager, clientManager))
			}
		})
	}
}

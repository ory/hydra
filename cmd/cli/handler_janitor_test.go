package cli

import (
	"context"
	"fmt"
	"github.com/ory/x/configx"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlxx"
)

var (
	janitor    = newJanitorHandler()
	janitorCmd = &cobra.Command{
		Use:  "janitor",
		RunE: janitor.Purge,
	}
)

var lifespan = time.Hour
var flushAccessRequests = []*fosite.Request{
	{
		ID:             "flush-access-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{OutfacingID: "flush-access-1"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-access-2",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:         &client.Client{OutfacingID: "flush-access-2"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-access-3",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:         &client.Client{OutfacingID: "flush-access-3"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
}
var tokenSignature = "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"
var flushRefreshRequests = []*fosite.AccessRequest{
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second),
			ID:             "flush-refresh-1",
			Client:         &client.Client{OutfacingID: "flush-refresh-1"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-1", tokenSignature)},
			},
		},
	},
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			ID:             "flush-refresh-2",
			Client:         &client.Client{OutfacingID: "flush-refresh-2"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-2", tokenSignature)},
			},
		},
	},
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			ID:             "flush-refresh-3",
			Client:         &client.Client{OutfacingID: "flush-refresh-3"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-3", tokenSignature)},
			},
		},
	},
}
var flushLoginRequests = []*consent.LoginRequest{
	{
		ID:             "flush-login-1",
		RequestedScope: []string{"foo", "bar"},
		Subject:        "flush-login-1",
		Client: &client.Client{
			OutfacingID:  "flush-login-consent-1",
			RedirectURIs: []string{"http://redirect"},
		},
		RequestURL:      "http://redirect",
		RequestedAt:     time.Now().Round(time.Second),
		AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second)),
		//SessionID:       "flush-login-1",
		Verifier: "flush-login-1",
	},
	{
		ID:             "flush-login-2",
		RequestedScope: []string{"foo", "bar"},
		Subject:        "flush-login-2",
		Client: &client.Client{
			OutfacingID:  "flush-login-consent-2",
			RedirectURIs: []string{"http://redirect"},
		},
		RequestURL:      "http://redirect",
		RequestedAt:     time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-(lifespan + time.Minute))),
		//SessionID:       "flush-login-2",
		Verifier: "flush-login-2",
	},
	{
		ID:             "flush-login-3",
		RequestedScope: []string{"foo", "bar"},
		Subject:        "flush-login-3",
		Client: &client.Client{
			OutfacingID:  "flush-login-consent-3",
			RedirectURIs: []string{"http://redirect"},
		},
		RequestURL:      "http://redirect",
		RequestedAt:     time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second).Add(-(lifespan + time.Hour))),
		//SessionID:       "flush-login-3",
		Verifier: "flush-login-3",
	},
}
var flushConsentRequests = []*consent.ConsentRequest{
	{
		ID:                   "flush-consent-1",
		RequestedScope:       []string{"foo", "bar"},
		Subject:              "flush-consent-1",
		OpenIDConnectContext: nil,
		ClientID:             "flush-login-consent-1",
		RequestURL:           "http://redirect",
		LoginChallenge:       "flush-login-1",
		RequestedAt:          time.Now().Round(time.Second),
		Verifier:             "flush-consent-1",
	},
	{
		ID:                   "flush-consent-2",
		RequestedScope:       []string{"foo", "bar"},
		Subject:              "flush-consent-2",
		OpenIDConnectContext: nil,
		ClientID:             "flush-login-consent-2",
		RequestURL:           "http://redirect",
		LoginChallenge:       "flush-login-2",
		RequestedAt:          time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Verifier:             "flush-consent-2",
	},
	{
		ID:                   "flush-consent-3",
		RequestedScope:       []string{"foo", "bar"},
		Subject:              "flush-consent-3",
		OpenIDConnectContext: nil,
		ClientID:             "flush-login-consent-3",
		RequestURL:           "http://redirect",
		LoginChallenge:       "flush-login-3",
		RequestedAt:          time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Verifier:             "flush-consent-3",
	},
}

func init() {
	janitorCmd.Flags().StringP("keep-if-younger", "k", "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("access-lifespan", "a", "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("refresh-lifespan", "r", "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("consent-request-lifespan", "l", "", "Set the login-consent request lifespan e.g. 1s, 1m, 1h")
	janitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(janitorCmd.PersistentFlags())
}

func connectPostgres(t *testing.T) (*pop.Connection, string) {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	return c, c.URL()
}

func connectMySQL(t *testing.T) (*pop.Connection, string) {
	c := dockertest.ConnectToTestMySQLPop(t)
	return c, fmt.Sprintf("mysql://%s", c.URL())
}

func connectCockroach(t *testing.T) (*pop.Connection, string) {
	c := dockertest.ConnectToTestCockroachDBPop(t)
	return c, c.URL()
}

var dbConnections = map[string]func(*testing.T) (*pop.Connection, string){
	"postgres": connectPostgres,
	"mysql":    connectMySQL,
	//"cockroach": connectCockroach,
}

func loginConsentRejectionRun(t *testing.T, ctx context.Context, dsn string) {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyConsentRequestMaxAge, time.Hour*1)

	conf.MustSet(config.KeyConsentURL, "http://redirect")
	conf.MustSet(config.KeyLoginURL, "http://redirect")

	conf.MustSet(config.KeyLogLevel, "trace")
	conf.MustSet(config.KeyDSN, dsn)

	// Connect to sqlite db
	reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
	require.NoError(t, err)

	require.NoError(t, reg.Persister().MigrateUp(ctx))

	cm := reg.ConsentManager()
	cl := reg.ClientManager()

	jt := consent.JanitorLoginConsentTest{}

	// Create login requests
	for _, r := range flushLoginRequests {
		require.NoError(t, cl.CreateClient(ctx, r.Client))
		require.NoError(t, cm.CreateLoginRequest(ctx, r))
	}

	// Explicit rejection
	for _, r := range flushLoginRequests {
		if r.ID == "flush-login-1" {
			// accept this one
			_, err = cm.HandleLoginRequest(ctx, r.ID, jt.NewHandledLoginRequest(
				r.ID, false, r.RequestedAt, r.AuthenticatedAt))

			require.NoError(t, err)
			continue
		}

		// reject flush-login-2 and 3
		_, err = cm.HandleLoginRequest(ctx, r.ID, jt.NewHandledLoginRequest(
			r.ID, true, r.RequestedAt, r.AuthenticatedAt))
		require.NoError(t, err)
	}

	janitorCmd.SetArgs([]string{reg.Config().DSN()})
	require.NoError(t, janitorCmd.Execute())

	// flush-login-2 and 3 should be cleared now
	for _, r := range flushLoginRequests {
		t.Logf("check login: %s", r.ID)
		_, err := cm.GetLoginRequest(ctx, r.ID)
		if r.ID == "flush-login-1" {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	// Cleanup rejected logins and recreate them to test the consent
	// since consent is dependant on login's acceptance.
	x.CleanSQLPop(t, reg.Persister().Connection(ctx))
	require.NoError(t, reg.Persister().MigrateUp(ctx))

	for _, r := range flushLoginRequests {
		/*if r.ID == "flush-login-1" {
			// this one is already accepted so no need to recreate it
			continue
		}
		require.NoError(t, cl.DeleteClient(ctx, r.Client.OutfacingID))*/
		require.NoError(t, cl.CreateClient(ctx, r.Client))
		require.NoError(t, cm.CreateLoginRequest(ctx, r))

		// accept flush-login-2 and 3
		_, err = cm.HandleLoginRequest(ctx, r.ID, jt.NewHandledLoginRequest(
			r.ID, false, r.RequestedAt, r.AuthenticatedAt))
	}

	// Create consent requests
	for _, r := range flushConsentRequests {
		require.NoError(t, cm.CreateConsentRequest(ctx, r))
	}

	//Reject the consents
	for _, r := range flushConsentRequests {
		if r.ID == "flush-consent-1" {
			// accept this one
			_, err = cm.HandleConsentRequest(ctx, r.ID, jt.NewHandledConsentRequest(
				r.ID, false, r.RequestedAt, r.AuthenticatedAt))
			require.NoError(t, err)
			continue
		}
		_, err = cm.HandleConsentRequest(ctx, r.ID, jt.NewHandledConsentRequest(
			r.ID, true, r.RequestedAt, r.AuthenticatedAt))
		require.NoError(t, err)
	}

	janitorCmd.SetArgs([]string{reg.Config().DSN()})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushConsentRequests {
		t.Logf("check consent: %s", r.ID)
		_, err = cm.GetConsentRequest(ctx, r.ID)
		if r.ID == "flush-consent-1" {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	require.NoError(t, reg.Persister().Connection(ctx).Close())
}

func loginConsentTimeout(t *testing.T, ctx context.Context, dsn string) {
	var err error

	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyConsentRequestMaxAge, time.Hour*1)

	conf.MustSet(config.KeyConsentURL, "http://redirect")
	conf.MustSet(config.KeyLoginURL, "http://redirect")

	conf.MustSet(config.KeyLogLevel, "trace")
	conf.MustSet(config.KeyDSN, dsn)

	reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
	require.NoError(t, err)
	require.NoError(t, reg.Persister().MigrateUp(ctx))

	cm := reg.ConsentManager()
	cl := reg.ClientManager()

	// Create login requests
	for _, r := range flushLoginRequests {
		require.NoError(t, cl.CreateClient(ctx, r.Client))
		require.NoError(t, cm.CreateLoginRequest(ctx, r))
	}

	// Creating at least 1 that has not timed out
	_, err = cm.HandleLoginRequest(ctx, flushLoginRequests[0].ID, &consent.HandledLoginRequest{
		ID:              flushLoginRequests[0].ID,
		RequestedAt:     flushLoginRequests[0].RequestedAt,
		AuthenticatedAt: flushLoginRequests[0].AuthenticatedAt,
		WasUsed:         true,
	})

	require.NoError(t, err)

	// First check if the login's can be purged when they have timed out
	janitorCmd.SetArgs([]string{reg.Config().DSN()})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushLoginRequests {
		_, err = cm.GetLoginRequest(ctx, r.ID)
		if r.ID == flushLoginRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	// Let's reset and accept all login requests to test the consent requests
	for _, r := range flushLoginRequests {
		require.NoError(t, cl.DeleteClient(ctx, r.Client.OutfacingID))
		require.NoError(t, cl.CreateClient(ctx, r.Client))
		require.NoError(t, cm.CreateLoginRequest(ctx, r))
		_, err = cm.HandleLoginRequest(ctx, r.ID, &consent.HandledLoginRequest{
			ID:              r.ID,
			AuthenticatedAt: r.AuthenticatedAt,
			RequestedAt:     r.RequestedAt,
			WasUsed:         true,
		})
		require.NoError(t, err)
	}

	// Create consent requests
	for _, r := range flushConsentRequests {
		require.NoError(t, cm.CreateConsentRequest(ctx, r))
	}

	// Create at least 1 consent request that has been accepted
	_, err = cm.HandleConsentRequest(ctx, flushConsentRequests[0].ID, &consent.HandledConsentRequest{
		ID:              flushConsentRequests[0].ID,
		WasUsed:         true,
		RequestedAt:     flushConsentRequests[0].RequestedAt,
		AuthenticatedAt: flushConsentRequests[0].AuthenticatedAt,
	})

	// Explicit timeout test
	janitorCmd.SetArgs([]string{reg.Config().DSN()})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushConsentRequests {
		_, err = cm.GetConsentRequest(ctx, r.ID)
		if r.ID == flushConsentRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	require.NoError(t, reg.Persister().Connection(ctx).Close())
}

func TestJanitorHandler_PurgeNotAfter(t *testing.T) {
	ctx := context.Background()
	var err error
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyRefreshTokenLifespan, lifespan)
	conf.MustSet(config.KeyConsentRequestMaxAge, lifespan)

	conf.MustSet(config.KeyLogLevel, "trace")
	conf.MustSet(config.KeyDSN, "sqlite://file:purgenotafter?mode=memory&_fk=true&cache=shared")

	reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
	require.NoError(t, err)

	cm := reg.ConsentManager()
	cl := reg.ClientManager()
	store := reg.OAuth2Storage()

	// Create login clients and requests
	for _, r := range flushLoginRequests {
		require.NoError(t, cl.CreateClient(ctx, r.Client))
		require.NoError(t, cm.CreateLoginRequest(ctx, r))
	}

	// Create consent requests, the clients have already been created in login
	for _, r := range flushConsentRequests {
		require.NoError(t, cm.CreateConsentRequest(ctx, r))
	}

	// Create access token clients and session
	for _, r := range flushAccessRequests {
		require.NoError(t, cl.CreateClient(ctx, r.Client.(*client.Client)))
		require.NoError(t, store.CreateAccessTokenSession(ctx, r.ID, r))
	}

	// Create refresh token clients and session
	for _, fr := range flushRefreshRequests {
		require.NoError(t, cl.CreateClient(ctx, fr.Client.(*client.Client)))
		require.NoError(t, store.CreateRefreshTokenSession(ctx, fr.ID, fr))
	}

	ds := new(oauth2.Session)

	// == Test Cycle 1: do not remove anything that is not older than 24 hours ==
	janitorCmd.SetArgs([]string{
		fmt.Sprintf("-k=%s", (time.Hour * 24).String()),
		fmt.Sprintf("-r=%s", conf.RefreshTokenLifespan().String()),
		fmt.Sprintf("-l=%s", conf.ConsentRequestMaxAge().String()),
		reg.Config().DSN(),
	})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushLoginRequests {
		t.Logf("login flush check: %s", r.ID)
		_, err = cm.GetLoginRequest(ctx, r.ID)
		require.NoError(t, err)
	}

	for _, r := range flushConsentRequests {
		t.Logf("consent flush check: %s", r.ID)
		_, err = cm.GetConsentRequest(ctx, r.ID)
		require.NoError(t, err)
	}

	for _, r := range flushAccessRequests {
		t.Logf("access flush check: %s", r.ID)
		_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
		require.NoError(t, err)
	}

	for _, r := range flushRefreshRequests {
		t.Logf("refresh flush check: %s", r.ID)
		_, err = store.GetRefreshTokenSession(ctx, r.ID, ds)
		require.NoError(t, err)
	}

	// == Test Cycle 2: do not remove anything that is older than 1h30min ==
	janitorCmd.SetArgs([]string{
		fmt.Sprintf("-k=%s", (lifespan + time.Hour/2).String()),
		fmt.Sprintf("-r=%s", conf.RefreshTokenLifespan().String()),
		fmt.Sprintf("-l=%s", conf.ConsentRequestMaxAge().String()),
		reg.Config().DSN(),
	})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushLoginRequests {
		t.Logf("login flush check: %s", r.ID)
		_, err = cm.GetLoginRequest(ctx, r.ID)
		if r.ID == flushLoginRequests[2].ID {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	for _, r := range flushConsentRequests {
		t.Logf("consent flush check: %s", r.ID)
		_, err = cm.GetConsentRequest(ctx, r.ID)
		if r.ID == flushConsentRequests[2].ID {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	for _, r := range flushAccessRequests {
		t.Logf("access flush check: %s", r.ID)
		_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
		if r.ID == flushAccessRequests[2].ID {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	for _, r := range flushRefreshRequests {
		t.Logf("refresh flush check: %s", r.Request.ID)
		_, err = store.GetRefreshTokenSession(ctx, r.Request.ID, ds)
		if r.ID == flushRefreshRequests[2].Request.ID {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	// == Test Cycle 3: remove anything that is older than now ==
	janitorCmd.SetArgs([]string{
		fmt.Sprintf("-k=%s", ""), // just keep this here to clear the previous keep-if-younger value...
		fmt.Sprintf("-r=%s", conf.RefreshTokenLifespan().String()),
		fmt.Sprintf("-l=%s", conf.ConsentRequestMaxAge().String()),
		reg.Config().DSN(),
	})
	require.NoError(t, janitorCmd.Execute())

	for _, r := range flushLoginRequests {
		t.Logf("login flush check: %s", r.ID)
		_, err = cm.GetLoginRequest(ctx, r.ID)
		if r.ID == flushLoginRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	for _, r := range flushConsentRequests {
		t.Logf("consent flush check: %s", r.ID)
		_, err = cm.GetConsentRequest(ctx, r.ID)
		if r.ID == flushConsentRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	for _, r := range flushAccessRequests {
		t.Logf("access flush check: %s", r.ID)
		_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
		if r.ID == flushAccessRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	for _, r := range flushRefreshRequests {
		t.Logf("refresh flush check: %s", r.ID)
		_, err = store.GetRefreshTokenSession(ctx, r.ID, ds)
		if r.ID == flushRefreshRequests[0].ID {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	require.NoError(t, reg.Persister().Connection(ctx).Close())
}

func TestJanitorHandler_PurgeLoginConsentRejection(t *testing.T) {
	/*
		Login and Consent also needs to be purged on two conditions besides the KeyConsentRequestMaxAge and notAfter time
		- when a login/consent request was never completed (timed out)
		- when a login/consent request was rejected
	*/
	ctx := context.Background()

	wg := &sync.WaitGroup{}

	for db, connect := range dbConnections {
		wg.Add(1)
		go func(d string, y func(t *testing.T) (*pop.Connection, string)) {
			t.Run("db="+d, func(t *testing.T) {
				defer wg.Done()
				c, dsn := y(t)
				t.Logf("Using dsn: %s\n", dsn)
				x.CleanSQLPop(t, c)
				loginConsentRejectionRun(t, ctx, dsn)
			})
		}(db, connect)
	}

	// sqlite tests
	wg.Add(1)
	go func() {
		t.Run("db=sqlite", func(t *testing.T) {
			defer wg.Done()
			// Setup config
			loginConsentRejectionRun(t, ctx, "sqlite://file:purgerejection?mode=memory&_fk=true&cache=shared")
		})
	}()

	wg.Wait()
}

func TestJanitorHandler_PurgeLoginConsentTimeout(t *testing.T) {
	/*
			Login and Consent also needs to be purged on two conditions besides the KeyConsentRequestMaxAge and notAfter time
			- when a login/consent request was never completed (timed out)
			- when a login/consent request was rejected

		The request is timeout when there are no entries inside the *_handled tables.

	*/
	ctx := context.Background()

	wg := &sync.WaitGroup{}

	for db, connect := range dbConnections {
		wg.Add(1)
		go func(d string, y func(t *testing.T) (*pop.Connection, string)) {
			t.Run("db="+d, func(t *testing.T) {
				defer wg.Done()
				c, dsn := y(t)
				t.Logf("Using dsn: %s\n", dsn)
				x.CleanSQLPop(t, c)
				loginConsentTimeout(t, ctx, dsn)
			})
		}(db, connect)
	}

	// sqlite tests
	wg.Add(1)
	go func() {
		t.Run("db=sqlite", func(t *testing.T) {
			defer wg.Done()
			// Setup config
			loginConsentTimeout(t, ctx, "sqlite://file:purgetimeout?mode=memory&_fk=true&cache=shared")
		})
	}()

	wg.Wait()
}

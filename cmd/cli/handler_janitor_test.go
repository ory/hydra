package cli

import (
	"context"
	"fmt"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/configx"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/x/logrusx"
)

var (
	janitor    = newJanitorHandler()
	janitorCmd = &cobra.Command{
		Use:  "janitor",
		RunE: janitor.Purge,
	}

	keepIfYounger          = "keep-if-younger"
	accessLifespan         = "access-lifespan"
	refreshLifespan        = "refresh-lifespan"
	consentRequestLifespan = "consent-request-lifespan"
	onlyTokens             = "only-tokens"
	onlyRequests           = "only-requests"
)

func init() {
	janitorCmd.Flags().String(keepIfYounger, "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().String(accessLifespan, "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().String(refreshLifespan, "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().String(consentRequestLifespan, "", "Set the login/consent request lifespan e.g. 1s, 1m, 1h")
	janitorCmd.Flags().Bool(onlyRequests, false, "This will only run the cleanup on requests and will skip token cleanup.")
	janitorCmd.Flags().Bool(onlyTokens, false, "This will only run the cleanup on tokens and will skip requests cleanup.")

	janitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(janitorCmd.PersistentFlags())
}

func TestJanitorHandler_PurgeTokenNotAfter(t *testing.T) {
	ctx := context.Background()
	jt := oauth2.NewOauthJanitorTestHelper("")

	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyRefreshTokenLifespan, jt.Lifespan)
	conf.MustSet(config.KeyAccessTokenLifespan, jt.Lifespan)

	conf.MustSet(config.KeyLogLevel, "trace")

	testCycles := map[string]time.Time{
		"notAfter24h":   time.Now().Round(time.Second).Add(-(jt.Lifespan * 24)),
		"notAfter1h30m": time.Now().Round(time.Second).Add(-(jt.Lifespan + time.Hour/2)),
		"notAfterNow":   time.Now().Round(time.Second),
	}

	for k, v := range testCycles {
		conf.MustSet(config.KeyDSN, fmt.Sprintf("sqlite://file:access_%s?mode=memory&_fk=true&cache=shared", k))
		reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
		require.NoError(t, err)
		janitorCmd.SetArgs([]string{
			fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
			fmt.Sprintf("--%s=%s", accessLifespan, conf.AccessTokenLifespan().String()),
			fmt.Sprintf("--%s", onlyTokens),
			conf.DSN(),
		})

		t.Run(fmt.Sprintf("case=access_%s", k),
			oauth2.NewOauthJanitorTestHelper(k).AccessTokenNotAfter(v, conf.AccessTokenLifespan(), janitorCmd.Execute, reg.ClientManager(), reg.OAuth2Storage()))
	}

	for k, v := range testCycles {
		conf.MustSet(config.KeyDSN, fmt.Sprintf("sqlite://file:refresh_%s?mode=memory&_fk=true&cache=shared", k))
		reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
		require.NoError(t, err)
		janitorCmd.SetArgs([]string{
			fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
			fmt.Sprintf("--%s=%s", refreshLifespan, conf.RefreshTokenLifespan().String()),
			fmt.Sprintf("--%s", onlyTokens),
			conf.DSN(),
		})

		t.Run(fmt.Sprintf("case=refresh_%s", k),
			oauth2.NewOauthJanitorTestHelper(k).RefreshTokenNotAfter(v, conf.RefreshTokenLifespan(), janitorCmd.Execute, reg.ClientManager(), reg.OAuth2Storage()))
	}
}

func TestJanitorHandler_PurgeLoginConsentNotAfter(t *testing.T) {
	ctx := context.Background()
	jt := consent.NewConsentJanitorTestHelper("")

	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyConsentRequestMaxAge, jt.Lifespan)
	conf.MustSet(config.KeyLogLevel, "trace")

	testCycles := map[string]time.Duration{
		"notAfter24h":   jt.Lifespan * 24,
		"notAfter1h30m": jt.Lifespan + time.Hour/2,
		"notAfterNow":   0,
	}

	for k, v := range testCycles {
		conf.MustSet(config.KeyDSN, fmt.Sprintf("sqlite://file:%s?mode=memory&_fk=true&cache=shared", k))
		reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
		require.NoError(t, err)
		janitorCmd.SetArgs([]string{
			fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
			fmt.Sprintf("--%s=%s", consentRequestLifespan, conf.ConsentRequestMaxAge().String()),
			fmt.Sprintf("--%s", onlyRequests),
			conf.DSN(),
		})

		t.Run(fmt.Sprintf("case=%s", k),
			consent.NewConsentJanitorTestHelper(k).LoginConsentNotAfter(time.Now().Round(time.Second).Add(-v), conf.ConsentRequestMaxAge(), janitorCmd.Execute, reg.ConsentManager(), reg.ClientManager()))
	}

}

func TestJanitorHandler_PurgeLoginConsent(t *testing.T) {
	/*
		Login and Consent also needs to be purged on two conditions besides the KeyConsentRequestMaxAge and notAfter time
		- when a login/consent request was never completed (timed out)
		- when a login/consent request was rejected
	*/
	ctx := context.Background()
	jt := consent.NewConsentJanitorTestHelper("")

	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyRefreshTokenLifespan, jt.Lifespan)
	conf.MustSet(config.KeyConsentRequestMaxAge, jt.Lifespan)
	conf.MustSet(config.KeyLogLevel, "trace")

	type loginConsentTest = func(func() error, consent.Manager, client.Manager) func(t *testing.T)

	testCycles := map[string]loginConsentTest{
		"loginRejection":   consent.NewConsentJanitorTestHelper("loginRejection").LoginRejection,
		"loginTimeout":     consent.NewConsentJanitorTestHelper("loginTimeout").LoginTimeout,
		"consentRejection": consent.NewConsentJanitorTestHelper("consentRejection").ConsentRejection,
		"consentTimeout":   consent.NewConsentJanitorTestHelper("consentTimeout").ConsentTimeout,
	}

	for k, v := range testCycles {
		conf.MustSet(config.KeyDSN, fmt.Sprintf("sqlite://file:%s?mode=memory&_fk=true&cache=shared", k))
		janitorCmd.SetArgs([]string{
			fmt.Sprintf("--%s", onlyRequests),
			conf.DSN(),
		})
		reg, err := driver.NewRegistryFromDSN(ctx, conf, logrusx.New("test_hydra", "master"))
		require.NoError(t, err)
		t.Run(fmt.Sprintf("case=%s", k), v(janitorCmd.Execute, reg.ConsentManager(), reg.ClientManager()))
	}
}

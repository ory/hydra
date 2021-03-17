package cli

import (
	"context"
	"fmt"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/configx"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
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
	jt := testhelpers.NewConsentJanitorTestHelper("token_not_after")
	testCycles := jt.GetNotAfterTestCycles()

	for k, v := range testCycles {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {

			reg, err := jt.GetRegistry(ctx, k)
			require.NoError(t, err)

			// setup test
			t.Run("step=setup_access", jt.AccessTokenNotAfterSetup(ctx, reg.ClientManager(), reg.OAuth2Storage()))
			t.Run("step=setup_refresh", jt.RefreshTokenNotAfterSetup(ctx, reg.ClientManager(), reg.OAuth2Storage()))

			// run the cleanup routine
			t.Run("step=cleanup", func(t *testing.T) {
				janitorCmd.SetArgs([]string{
					fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", accessLifespan, jt.GetAccessTokenLifespan().String()),
					fmt.Sprintf("--%s=%s", refreshLifespan, jt.GetRefreshTokenLifespan().String()),
					fmt.Sprintf("--%s", onlyTokens),
					jt.GetDSN(),
				})
				require.NoError(t, janitorCmd.Execute())
			})

			// validate test
			notAfter := time.Now().Round(time.Second).Add(-v)
			t.Run("step=validate_access", jt.AccessTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
			t.Run("step=validate_refresh", jt.RefreshTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
		})
	}
}

func TestJanitorHandler_PurgeLoginConsentNotAfter(t *testing.T) {
	ctx := context.Background()

	jt := testhelpers.NewConsentJanitorTestHelper("login_consent_not_after")
	testCycles := jt.GetNotAfterTestCycles()

	for k, v := range testCycles {
		reg, err := jt.GetRegistry(ctx, k)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			// Setup the test
			t.Run("step=setup", jt.LoginConsentNotAfterSetup(ctx, reg.ConsentManager(), reg.ClientManager()))
			// Run the cleanup routine
			t.Run("step=cleanup", func(t *testing.T) {
				janitorCmd.SetArgs([]string{
					fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", consentRequestLifespan, jt.GetConsentRequestLifespan().String()),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				})
				require.NoError(t, janitorCmd.Execute())
			})

			notAfter := time.Now().Round(time.Second).Add(-v)
			consentLifespan := time.Now().Round(time.Second).Add(-jt.GetConsentRequestLifespan())
			t.Run("step=validate", jt.LoginConsentNotAfterValidate(ctx, notAfter, consentLifespan, reg.ConsentManager()))
		})
	}

}

func TestJanitorHandler_PurgeLoginConsent(t *testing.T) {
	/*
		Login and Consent also needs to be purged on two conditions besides the KeyConsentRequestMaxAge and notAfter time
		- when a login/consent request was never completed (timed out)
		- when a login/consent request was rejected
	*/

	t.Run("case=login_consent_timeout", func(t *testing.T) {
		ctx := context.Background()
		jt := testhelpers.NewConsentJanitorTestHelper("login_consent_timeout")
		loginConsentTimeoutSetup := jt.GetLoginConsentTimeoutSetup()
		loginConsentTimeoutValidate := jt.GetLoginConsentTimeoutValidate()

		for k, _ := range loginConsentTimeoutSetup {
			t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
				jt := testhelpers.NewConsentJanitorTestHelper(k)
				reg, err := jt.GetRegistry(ctx, k)
				require.NoError(t, err)

				// setup
				t.Run("step=setup", loginConsentTimeoutSetup[k](ctx, reg.ConsentManager(), reg.ClientManager()))

				// run cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					janitorCmd.SetArgs([]string{
						fmt.Sprintf("--%s", onlyRequests),
						jt.GetDSN(),
					})
					require.NoError(t, janitorCmd.Execute())
				})

				t.Run("step=validate", loginConsentTimeoutValidate[k](ctx, reg.ConsentManager()))
			})
		}
	})

	t.Run("case=login_consent_rejection", func(t *testing.T) {
		ctx := context.Background()
		jt := testhelpers.NewConsentJanitorTestHelper("login_consent_rejection")
		loginConsentRejectionSetup := jt.GetLoginConsentRejectionSetup()
		loginConsentRejectionValidate := jt.GetLoginConsentRejectionValidate()

		for k, _ := range loginConsentRejectionSetup {
			t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
				jt := testhelpers.NewConsentJanitorTestHelper(k)
				reg, err := jt.GetRegistry(ctx, k)
				require.NoError(t, err)

				// setup
				t.Run("step=setup", loginConsentRejectionSetup[k](ctx, reg.ConsentManager(), reg.ClientManager()))

				// cleanup
				t.Run("step=cleanup", func(t *testing.T) {
					janitorCmd.SetArgs([]string{
						fmt.Sprintf("--%s", onlyRequests),
						jt.GetDSN(),
					})
					require.NoError(t, janitorCmd.Execute())
				})

				// validate
				t.Run("step=validate", loginConsentRejectionValidate[k](ctx, reg.ConsentManager()))
			})
		}

	})

}

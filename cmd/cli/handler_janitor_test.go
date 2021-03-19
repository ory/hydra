package cli

import (
	"context"
	"fmt"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/configx"
	"github.com/ory/x/flagx"
	"github.com/spf13/cobra"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	keepIfYounger          = "keep-if-younger"
	accessLifespan         = "access-lifespan"
	refreshLifespan        = "refresh-lifespan"
	consentRequestLifespan = "consent-request-lifespan"
	onlyTokens             = "tokens"
	onlyRequests           = "requests"
)

func newJanitorCmd() *cobra.Command {
	janitor := newJanitorHandler()
	JanitorCmd := &cobra.Command{
		Use:  "janitor",
		RunE: janitor.Purge,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 &&
				!flagx.MustGetBool(cmd, "read-from-env") &&
				len(flagx.MustGetStringSlice(cmd, "config")) == 0 {

				fmt.Printf("%s\n", cmd.UsageString())
				return fmt.Errorf("%s\n%s\n%s\n",
					"A DSN is required as a positional argument when not passing any of the following flags:",
					"- Using the environment variable with flag -e, --read-from-env",
					"- Using the config file with flag -c, --config")
			}

			if (!flagx.MustGetBool(cmd, onlyTokens) && !flagx.MustGetBool(cmd, onlyRequests)) || (flagx.MustGetBool(cmd, onlyTokens) && flagx.MustGetBool(cmd, onlyRequests)) {
				return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
					"Janitor requires either --tokens or --requests to be set")
			}

			return nil
		},
	}
	JanitorCmd.Flags().String(keepIfYounger, "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	JanitorCmd.Flags().String(accessLifespan, "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	JanitorCmd.Flags().String(refreshLifespan, "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	JanitorCmd.Flags().String(consentRequestLifespan, "", "Set the login/consent request lifespan e.g. 1s, 1m, 1h")
	JanitorCmd.Flags().Bool(onlyRequests, false, "This will only run the cleanup on requests and will skip token cleanup.")
	JanitorCmd.Flags().Bool(onlyTokens, false, "This will only run the cleanup on tokens and will skip requests cleanup.")

	JanitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(JanitorCmd.PersistentFlags())
	return JanitorCmd
}

func TestJanitorHandler_PurgeTokenNotAfter(t *testing.T) {
	ctx := context.Background()
	testCycles := testhelpers.NewConsentJanitorTestHelper("").GetNotAfterTestCycles()

	for k, v := range testCycles {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, k)
			require.NoError(t, err)

			// setup test
			t.Run("step=setup-access", jt.AccessTokenNotAfterSetup(ctx, reg.ClientManager(), reg.OAuth2Storage()))
			t.Run("step=setup-refresh", jt.RefreshTokenNotAfterSetup(ctx, reg.ClientManager(), reg.OAuth2Storage()))

			// run the cleanup routine
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", accessLifespan, jt.GetAccessTokenLifespan().String()),
					fmt.Sprintf("--%s=%s", refreshLifespan, jt.GetRefreshTokenLifespan().String()),
					fmt.Sprintf("--%s", onlyTokens),
					jt.GetDSN(),
				)
			})

			// validate test
			notAfter := time.Now().Round(time.Second).Add(-v)
			t.Run("step=validate-access", jt.AccessTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
			t.Run("step=validate-refresh", jt.RefreshTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
		})
	}
}

func TestJanitorHandler_PurgeLoginConsentNotAfter(t *testing.T) {
	ctx := context.Background()

	testCycles := testhelpers.NewConsentJanitorTestHelper("").GetNotAfterTestCycles()

	for k, v := range testCycles {
		jt := testhelpers.NewConsentJanitorTestHelper(k)
		reg, err := jt.GetRegistry(ctx, k)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			// Setup the test
			t.Run("step=setup", jt.LoginConsentNotAfterSetup(ctx, reg.ConsentManager(), reg.ClientManager()))
			// Run the cleanup routine
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s=%s", keepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", consentRequestLifespan, jt.GetConsentRequestLifespan().String()),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				)
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

	t.Run("case=login-consent-timeout", func(t *testing.T) {
		t.Run("case=login-timeout", func(t *testing.T) {
			ctx := context.Background()
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			// setup
			t.Run("step=setup", jt.LoginTimeoutSetup(ctx, reg.ConsentManager(), reg.ClientManager()))

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				)
			})

			t.Run("step=validate", jt.LoginTimeoutValidate(ctx, reg.ConsentManager()))

		})

		t.Run("case=consent-timeout", func(t *testing.T) {
			ctx := context.Background()
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			// setup
			t.Run("step=setup", jt.ConsentTimeoutSetup(ctx, reg.ConsentManager(), reg.ClientManager()))

			// run cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				)
			})

			// validate
			t.Run("step=validate", jt.ConsentTimeoutValidate(ctx, reg.ConsentManager()))
		})

	})

	t.Run("case=login-consent-rejection", func(t *testing.T) {
		ctx := context.Background()

		t.Run("case=login-rejection", func(t *testing.T) {
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			// setup
			t.Run("step=setup", jt.LoginRejectionSetup(ctx, reg.ConsentManager(), reg.ClientManager()))

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				)
			})

			// validate
			t.Run("step=validate", jt.LoginRejectionValidate(ctx, reg.ConsentManager()))
		})

		t.Run("case=consent-rejection", func(t *testing.T) {
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			// setup
			t.Run("step=setup", jt.ConsentRejectionSetup(ctx, reg.ConsentManager(), reg.ClientManager()))

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					fmt.Sprintf("--%s", onlyRequests),
					jt.GetDSN(),
				)
			})

			// validate
			t.Run("step=validate", jt.ConsentRejectionValidate(ctx, reg.ConsentManager()))
		})

	})

}

/*
// TODO: this throws a panic like error instead of a pass on an expected error
func TestJanitorHandler_Arguments(t *testing.T) {
	cmdx.ExecNoErr(t, newJanitorCmd(),
		fmt.Sprintf("--%s", onlyRequests),
		"memory",
	)
	cmdx.ExecNoErr(t, newJanitorCmd(),
		fmt.Sprintf("--%s", onlyTokens),
		"memory",
	)
	cmdx.ExecExpectedErr(t, newJanitorCmd(),
		fmt.Sprintf("--%s", onlyRequests),
		fmt.Sprintf("--%s", onlyTokens),
		"memory",
	)
	cmdx.ExecExpectedErr(t, newJanitorCmd(),
		"memory",
	)
}*/

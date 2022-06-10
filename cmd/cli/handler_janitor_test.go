package cli_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/hydra/cmd"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/cmdx"
)

func newJanitorCmd() *cobra.Command {
	return cmd.NewRootCmd()
}

func TestJanitorHandler_PurgeTokenNotAfter(t *testing.T) {
	ctx := context.Background()
	testCycles := testhelpers.NewConsentJanitorTestHelper("").GetNotAfterTestCycles()

	require.True(t, len(testCycles) > 0)

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
					"janitor",
					fmt.Sprintf("--%s=%s", cli.KeepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", cli.AccessLifespan, jt.GetAccessTokenLifespan().String()),
					fmt.Sprintf("--%s=%s", cli.RefreshLifespan, jt.GetRefreshTokenLifespan().String()),
					fmt.Sprintf("--%s", cli.OnlyTokens),
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
					"janitor",
					fmt.Sprintf("--%s=%s", cli.KeepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", cli.ConsentRequestLifespan, jt.GetConsentRequestLifespan().String()),
					fmt.Sprintf("--%s", cli.OnlyRequests),
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
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyRequests),
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
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyRequests),
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
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyRequests),
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
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyRequests),
					jt.GetDSN(),
				)
			})

			// validate
			t.Run("step=validate", jt.ConsentRejectionValidate(ctx, reg.ConsentManager()))
		})

	})

}

func TestJanitorHandler_Arguments(t *testing.T) {
	cmdx.ExecNoErr(t, cmd.NewRootCmd(),
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		"memory",
	)
	cmdx.ExecNoErr(t, cmd.NewRootCmd(),
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyTokens),
		"memory",
	)

	cmdx.ExecNoErr(t, cmd.NewRootCmd(),
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyLoginSessions),
		"memory")

	_, _, err := cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "hydra janitor [<database-url>] [flags]")

	cmdx.ExecNoErr(t, cmd.NewRootCmd(),
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.Limit, "1000"),
		fmt.Sprintf("--%s=%s", cli.BatchSize, "100"),
		"memory",
	)

	_, _, err = cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.Limit, "0"),
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Values for --limit and --batch-size should both be greater than 0")

	_, _, err = cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.Limit, "-100"),
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Values for --limit and --batch-size should both be greater than 0")

	_, _, err = cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.BatchSize, "0"),
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Values for --limit and --batch-size should both be greater than 0")

	_, _, err = cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.BatchSize, "-100"),
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Values for --limit and --batch-size should both be greater than 0")

	_, _, err = cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		fmt.Sprintf("--%s", cli.OnlyRequests),
		fmt.Sprintf("--%s=%s", cli.Limit, "100"),
		fmt.Sprintf("--%s=%s", cli.BatchSize, "1000"),
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Value for --batch-size must not be greater than value for --limit")
}

func TestJanitorHandler_PurgeLoginSessions(t *testing.T) {
	t.Run("case=cleanup-login-session", func(t *testing.T) {
		t.Run("case=clean session without consent-request", func(t *testing.T) {
			ctx := context.Background()
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			const sessionID = "session_id"

			// setup
			sessionHelper := testhelpers.NewJanitorSessionTestHelper(reg, t.Name())
			sessionHelper.CreateLoginSession(t, ctx, sessionID)

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyLoginSessions),
					jt.GetDSN(),
				)
			})

			// validate
			sessionHelper.ValidateSessionNotExist(t, ctx, sessionID)
		})

		t.Run("case=clean sessions with consent-request", func(t *testing.T) {
			ctx := context.Background()
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			const (
				sessionID = "session_id_1"
			)

			//setup
			sessionHelper := testhelpers.NewJanitorSessionTestHelper(reg, t.Name())
			sessionHelper.CreateLoginSession(t, ctx, sessionID)
			sessionHelper.CreateEnvironmentForSession(t, ctx, sessionID)

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyLoginSessions),
					jt.GetDSN(),
				)
			})

			// validate
			sessionHelper.ValidateSessionExist(t, ctx, sessionID)
		})

		t.Run("case=alive session and flush session", func(t *testing.T) {
			ctx := context.Background()
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, t.Name())
			require.NoError(t, err)

			const (
				aliveSessionID = "session_id_1"
				flushSessionID = "session_id_flush"
			)

			//setup
			sessionHelper := testhelpers.NewJanitorSessionTestHelper(reg, t.Name())
			sessionHelper.CreateLoginSession(t, ctx, aliveSessionID)
			sessionHelper.CreateEnvironmentForSession(t, ctx, aliveSessionID)
			sessionHelper.CreateLoginSession(t, ctx, flushSessionID)

			// cleanup
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					"janitor",
					fmt.Sprintf("--%s", cli.OnlyLoginSessions),
					jt.GetDSN(),
				)
			})

			// validate
			sessionHelper.ValidateSessionExist(t, ctx, aliveSessionID)
			sessionHelper.ValidateSessionNotExist(t, ctx, flushSessionID)
		})

	})
}

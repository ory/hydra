package cli

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/cmdx"
)

func newJanitorCmd() *cobra.Command {
	return newJanitorHandler().Command()
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
					fmt.Sprintf("--%s=%s", KeepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", AccessLifespan, jt.GetAccessTokenLifespan().String()),
					fmt.Sprintf("--%s=%s", RefreshLifespan, jt.GetRefreshTokenLifespan().String()),
					fmt.Sprintf("--%s", OnlyTokens),
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
					fmt.Sprintf("--%s=%s", KeepIfYounger, v.String()),
					fmt.Sprintf("--%s=%s", ConsentRequestLifespan, jt.GetConsentRequestLifespan().String()),
					fmt.Sprintf("--%s", OnlyRequests),
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
					fmt.Sprintf("--%s", OnlyRequests),
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
					fmt.Sprintf("--%s", OnlyRequests),
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
					fmt.Sprintf("--%s", OnlyRequests),
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
					fmt.Sprintf("--%s", OnlyRequests),
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
		fmt.Sprintf("--%s", OnlyRequests),
		"memory",
	)
	cmdx.ExecNoErr(t, newJanitorCmd(),
		fmt.Sprintf("--%s", OnlyTokens),
		"memory",
	)
	cmdx.ExecExpectedErr(t, newJanitorCmd(),
		fmt.Sprintf("--%s", OnlyRequests),
		fmt.Sprintf("--%s", OnlyTokens),
		"memory",
	)
	cmdx.ExecExpectedErr(t, newJanitorCmd(),
		"memory",
	)
}*/

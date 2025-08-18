// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/hydra/v2/cmd"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/cmdx"
)

func newJanitorCmd() *cobra.Command {
	return cmd.NewRootCmd()
}

func TestJanitorHandler_PurgeTokenNotAfter(t *testing.T) {
	ctx := t.Context()

	for k, v := range testhelpers.NotAfterTestCycles {
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
					reg.Config().DSN(),
				)
			})

			// validate test
			notAfter := time.Now().Round(time.Second).Add(-v)
			t.Run("step=validate-access", jt.AccessTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
			t.Run("step=validate-refresh", jt.RefreshTokenNotAfterValidate(ctx, notAfter, reg.OAuth2Storage()))
		})
	}
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
		fmt.Sprintf("--%s", cli.OnlyGrants),
		"memory",
	)

	_, _, err := cmdx.ExecCtx(context.Background(), cmd.NewRootCmd(), nil,
		"janitor",
		"memory")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Janitor requires at least one of --tokens, --requests or --grants to be set")

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

func TestJanitorHandler_PurgeGrantNotAfter(t *testing.T) {
	ctx := t.Context()

	for k, v := range testhelpers.NotAfterTestCycles {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			jt := testhelpers.NewConsentJanitorTestHelper(t.Name())
			reg, err := jt.GetRegistry(ctx, k)
			require.NoError(t, err)

			// setup test
			t.Run("step=setup", jt.GrantNotAfterSetup(ctx, reg.GrantManager()))

			// run the cleanup routine
			t.Run("step=cleanup", func(t *testing.T) {
				cmdx.ExecNoErr(t, newJanitorCmd(),
					"janitor",
					fmt.Sprintf("--%s=%s", cli.KeepIfYounger, v.String()),
					fmt.Sprintf("--%s", cli.OnlyGrants),
					reg.Config().DSN(),
				)
			})

			// validate test
			notAfter := time.Now().Round(time.Second).Add(-v)
			t.Run("step=validate-access", jt.GrantNotAfterValidate(ctx, notAfter, reg.GrantManager()))
		})
	}
}

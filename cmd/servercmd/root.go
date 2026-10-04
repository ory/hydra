// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package servercmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/hydra/v2/driver"
)

func NewRootCmd(opts ...driver.OptionsModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage Ory Hydra",
	}
	cmdx.EnableUsageTemplating(cmd)
	RegisterCommandRecursive(cmd, opts...)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command, opts ...driver.OptionsModifier) {
	createCmd := cmd.NewCreateCmd()
	createCmd.AddCommand(
		cmd.NewCreateClientsCommand(),
		cmd.NewCreateJWKSCmd(),
	)

	getCmd := cmd.NewGetCmd()
	getCmd.AddCommand(
		cmd.NewGetClientsCmd(),
		cmd.NewGetJWKSCmd(),
	)

	deleteCmd := cmd.NewDeleteCmd()
	deleteCmd.AddCommand(
		cmd.NewDeleteClientCmd(),
		cmd.NewDeleteJWKSCommand(),
		cmd.NewDeleteAccessTokensCmd(),
	)

	listCmd := cmd.NewListCmd()
	listCmd.AddCommand(cmd.NewListClientsCmd())

	updateCmd := cmd.NewUpdateCmd()
	updateCmd.AddCommand(cmd.NewUpdateClientCmd())

	importCmd := cmd.NewImportCmd()
	importCmd.AddCommand(
		cmd.NewImportClientCmd(),
		cmd.NewKeysImportCmd(),
	)

	performCmd := cmd.NewPerformCmd()
	performCmd.AddCommand(
		cmd.NewPerformClientCredentialsCmd(),
		cmd.NewPerformAuthorizationCodeCmd(),
		cmd.NewPerformDeviceCodeCmd(),
	)

	revokeCmd := cmd.NewRevokeCmd()
	revokeCmd.AddCommand(cmd.NewRevokeTokenCmd())

	introspectCmd := cmd.NewIntrospectCmd()
	introspectCmd.AddCommand(cmd.NewIntrospectTokenCmd())

	migrateCmd := NewMigrateCmd()
	migrateCmd.AddCommand(NewMigrateSQLCmd(opts))
	migrateCmd.AddCommand(NewMigrateStatusCmd(opts))

	serveCmd := NewServeCmd()
	serveCmd.AddCommand(NewServeAdminCmd(opts))
	serveCmd.AddCommand(NewServePublicCmd(opts))
	serveCmd.AddCommand(NewServeAllCmd(opts))

	parent.AddCommand(
		createCmd,
		getCmd,
		deleteCmd,
		listCmd,
		updateCmd,
		importCmd,
		performCmd,
		introspectCmd,
		revokeCmd,
		migrateCmd,
		serveCmd,
		NewJanitorCmd(opts),
		NewVersionCmd(),
	)
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	c := NewRootCmd()
	if err := c.Execute(); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(c.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}

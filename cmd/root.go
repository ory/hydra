// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"

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
	createCmd := NewCreateCmd()
	createCmd.AddCommand(
		NewCreateClientsCommand(),
		NewCreateJWKSCmd(),
	)

	getCmd := NewGetCmd()
	getCmd.AddCommand(
		NewGetClientsCmd(),
		NewGetJWKSCmd(),
	)

	deleteCmd := NewDeleteCmd()
	deleteCmd.AddCommand(
		NewDeleteClientCmd(),
		NewDeleteJWKSCommand(),
		NewDeleteAccessTokensCmd(),
	)

	listCmd := NewListCmd()
	listCmd.AddCommand(NewListClientsCmd())

	updateCmd := NewUpdateCmd()
	updateCmd.AddCommand(NewUpdateClientCmd())

	importCmd := NewImportCmd()
	importCmd.AddCommand(
		NewImportClientCmd(),
		NewKeysImportCmd(),
	)

	performCmd := NewPerformCmd()
	performCmd.AddCommand(
		NewPerformClientCredentialsCmd(),
		NewPerformAuthorizationCodeCmd(),
		NewPerformDeviceCodeCmd(),
	)

	revokeCmd := NewRevokeCmd()
	revokeCmd.AddCommand(NewRevokeTokenCmd())

	introspectCmd := NewIntrospectCmd()
	introspectCmd.AddCommand(NewIntrospectTokenCmd())

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

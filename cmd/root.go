// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/ory/x/cmdx"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/servicelocatorx"

	"github.com/spf13/cobra"
)

func NewRootCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage Ory Hydra",
	}
	cmdx.EnableUsageTemplating(cmd)
	RegisterCommandRecursive(cmd, slOpts, dOpts, cOpts)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command, slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) {
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
	)

	revokeCmd := NewRevokeCmd()
	revokeCmd.AddCommand(NewRevokeTokenCmd())

	introspectCmd := NewIntrospectCmd()
	introspectCmd.AddCommand(NewIntrospectTokenCmd())

	migrateCmd := NewMigrateCmd()
	migrateCmd.AddCommand(NewMigrateGenCmd())
	migrateCmd.AddCommand(NewMigrateSqlCmd(slOpts, dOpts, cOpts))
	migrateCmd.AddCommand(NewMigrateStatusCmd(slOpts, dOpts, cOpts))

	serveCmd := NewServeCmd()
	serveCmd.AddCommand(NewServeAdminCmd(slOpts, dOpts, cOpts))
	serveCmd.AddCommand(NewServePublicCmd(slOpts, dOpts, cOpts))
	serveCmd.AddCommand(NewServeAllCmd(slOpts, dOpts, cOpts))

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
		NewJanitorCmd(slOpts, dOpts, cOpts),
		NewVersionCmd(),
	)
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := NewRootCmd(nil, nil, nil).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

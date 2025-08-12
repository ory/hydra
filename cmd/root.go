// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/servicelocatorx"
)

func NewRootCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage Ory Hydra",
	}
	cmdx.EnableUsageTemplating(cmd)
	RegisterCommandRecursive(cmd, slOpts, dOpts)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command, slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier) {
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
	migrateCmd.AddCommand(NewMigrateSQLCmd(slOpts, dOpts))
	migrateCmd.AddCommand(NewMigrateStatusCmd(slOpts, dOpts))

	serveCmd := NewServeCmd()
	serveCmd.AddCommand(NewServeAdminCmd(slOpts, dOpts))
	serveCmd.AddCommand(NewServePublicCmd(slOpts, dOpts))
	serveCmd.AddCommand(NewServeAllCmd(slOpts, dOpts))

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
		NewJanitorCmd(slOpts, dOpts),
		NewVersionCmd(),
	)
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	c := NewRootCmd(nil, nil)
	if err := c.Execute(); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(c.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}

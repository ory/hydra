/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/ory/x/cmdx"

	"github.com/ory/hydra/driver"
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

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

	"github.com/ory/hydra/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/servicelocatorx"

	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
func NewRootCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage Ory Hydra",
	}
	RegisterCommandRecursive(cmd, slOpts, dOpts, cOpts)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command, slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) {
	createCmd := NewCreateCmd(parent)
	parent.AddCommand(createCmd)
	createCmd.AddCommand(NewCreateClientsCommand(parent))
	createCmd.AddCommand(NewCreateJWKSCmd(parent))

	getCmd := NewGetCmd(parent)
	parent.AddCommand(getCmd)
	getCmd.AddCommand(NewGetClientsCmd(parent))
	getCmd.AddCommand(NewGetJWKSCmd(parent))

	deleteCmd := NewDeleteCmd(parent)
	parent.AddCommand(deleteCmd)
	deleteCmd.AddCommand(NewDeleteClientCmd(parent))
	deleteCmd.AddCommand(NewDeleteJWKSCommand(parent))
	deleteCmd.AddCommand(NewDeleteAccessTokensCmd(parent))

	listCmd := NewListCmd(parent)
	parent.AddCommand(listCmd)
	listCmd.AddCommand(NewListClientsCmd(parent))

	updateCmd := NewUpdateCmd(parent)
	parent.AddCommand(updateCmd)
	updateCmd.AddCommand(NewUpdateClientCmd(parent))

	importCmd := NewImportCmd(parent)
	parent.AddCommand(importCmd)
	importCmd.AddCommand(NewImportClientCmd(parent))
	importCmd.AddCommand(NewKeysImportCmd(parent))

	performCmd := NewPerformCmd(parent)
	parent.AddCommand(performCmd)
	performCmd.AddCommand(NewPerformClientCredentialsCmd(parent))
	performCmd.AddCommand(NewPerformAuthorizationCodeCmd(parent))

	revokeCmd := NewRevokeCmd(parent)
	parent.AddCommand(revokeCmd)
	revokeCmd.AddCommand(NewRevokeTokenCmd(parent))

	introspectCmd := NewIntrospectCmd(parent)
	parent.AddCommand(introspectCmd)
	introspectCmd.AddCommand(NewIntrospectTokenCmd(parent))

	parent.AddCommand(NewJanitorCmd(slOpts, dOpts, cOpts))

	migrateCmd := NewMigrateCmd()
	parent.AddCommand(migrateCmd)
	migrateCmd.AddCommand(NewMigrateGenCmd())
	migrateCmd.AddCommand(NewMigrateSqlCmd(slOpts, dOpts, cOpts))

	serveCmd := NewServeCmd()
	parent.AddCommand(serveCmd)
	serveCmd.AddCommand(NewServeAdminCmd(slOpts, dOpts, cOpts))
	serveCmd.AddCommand(NewServePublicCmd(slOpts, dOpts, cOpts))
	serveCmd.AddCommand(NewServeAllCmd(slOpts, dOpts, cOpts))

	parent.AddCommand(NewVersionCmd())
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := NewRootCmd(nil, nil, nil).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

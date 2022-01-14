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

	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage Ory Hydra",
	}
	RegisterCommandRecursive(cmd)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command) {
	// Clients
	clientCmd := NewClientsCmd()
	parent.AddCommand(clientCmd)
	clientCmd.AddCommand(NewClientsCreateCmd())
	clientCmd.AddCommand(NewClientsDeleteCmd())
	clientCmd.AddCommand(NewClientsGetCmd())
	clientCmd.AddCommand(NewClientsImportCmd())
	clientCmd.AddCommand(NewClientsImportCmd())
	clientCmd.AddCommand(NewClientsImportCmd())
	clientCmd.AddCommand(NewClientsListCmd())
	clientCmd.AddCommand(NewKeysImportCmd())
	clientCmd.AddCommand(NewClientsUpdateCmd())

	parent.AddCommand(NewJanitorCmd())

	keyCmd := NewKeysCmd()
	parent.AddCommand(keyCmd)
	keyCmd.AddCommand(NewKeysCreateCmd())
	keyCmd.AddCommand(NewKeysDeleteCmd())
	keyCmd.AddCommand(NewKeysGetCmd())
	keyCmd.AddCommand(NewKeysImportCmd())

	migrateCmd := NewMigrateCmd()
	parent.AddCommand(migrateCmd)
	migrateCmd.AddCommand(NewMigrateSqlCmd())

	serveCmd := NewServeCmd()
	parent.AddCommand(serveCmd)
	serveCmd.AddCommand(NewServeAdminCmd())
	serveCmd.AddCommand(NewServePublicCmd())
	serveCmd.AddCommand(NewServeAllCmd())

	tokenCmd := NewTokenCmd()
	parent.AddCommand(tokenCmd)
	tokenCmd.AddCommand(NewTokenClientCmd())
	tokenCmd.AddCommand(NewTokenDeleteCmd())
	tokenCmd.AddCommand(NewTokenFlushCmd())
	tokenCmd.AddCommand(NewTokenIntrospectCmd())
	tokenCmd.AddCommand(NewTokenRevokeCmd())
	tokenCmd.AddCommand(NewTokenUserCmd())

	parent.AddCommand(NewVersionCmd())
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

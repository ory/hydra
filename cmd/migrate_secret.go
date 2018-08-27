// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

// secretCmd represents the secret command
var migrateSecretCmd = &cobra.Command{
	Use:   "secret <database-url>",
	Short: "Rotates system secrets",
	Long: `This command rotates the system secret and reconfigures the store.
Example:

	OLD_SYSTEM_SECRET=old-secret... NEW_SYSTEM_SECRET=new-secret... hydra migrate secret postgres://...

You can read in the database URL using the -e flag, for example:
	export DATABASE_URL=...
	hydra migrate secret^^ -e
`,
	Run: cmdHandler.Migration.MigrateSecret,
}

func init() {
	migrateCmd.AddCommand(migrateSecretCmd)

	migrateSecretCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database URL from the environment variable DATABASE_URL.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

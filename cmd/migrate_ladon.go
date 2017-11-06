// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
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

import "github.com/spf13/cobra"

// migrateLadonCmd represents the ladon command
var migrateLadonCmd = &cobra.Command{
	Use:   "ladon 0.6.0 <database-url>",
	Short: "Migrates Ladon SQL schema to version 0.6.0",
	Long: `Hydra version 0.8.0 includes a breaking schema change from Ladon which was introduced
with Ladon version 0.6.0. This script applies the neccessary migrations by copying data from the old tables
to the new ones. This command might take some time, depending on how many policies are in your store.

Do not run this command on a fresh installation.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

### WARNING ###

Before running this command on an existing database, create a back up!
`,
	Run: cmdHandler.Migration.MigrateLadon050To060,
}

func init() {
	migrateCmd.AddCommand(migrateLadonCmd)
}

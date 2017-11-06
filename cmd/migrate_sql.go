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

// migrateSqlCmd represents the sql command
var migrateSqlCmd = &cobra.Command{
	Use:   "sql <database-url>",
	Short: "Create SQL schemas and apply migration plans",
	Long: `Run this command on a fresh SQL installation and when you upgrade Hydra to a new minor version. For example,
upgrading Hydra 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

### WARNING ###

Before running this command on an existing database, create a back up!
`,
	Run: cmdHandler.Migration.MigrateSQL,
}

func init() {
	migrateCmd.AddCommand(migrateSqlCmd)
}

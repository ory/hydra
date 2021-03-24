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
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
)

// migrateSqlCmd represents the sql command
func NewMigrateSqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql <database-url>",
		Short: "Create SQL schemas and apply migration plans",
		Long: `Run this command on a fresh SQL installation and when you upgrade Hydra to a new minor version. For example,
upgrading Hydra 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

You can read in the database URL using the -e flag, for example:
	export DSN=...
	hydra migrate sql -e

### WARNING ###

Before running this command on an existing database, create a back up!
`,
		Run: cli.NewHandler().Migration.MigrateSQL,
	}

	cmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	cmd.Flags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")

	return cmd
}

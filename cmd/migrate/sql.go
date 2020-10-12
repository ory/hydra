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

package migrate

import (
	"context"
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"
	"github.com/spf13/cobra"
)

const (
	flagReadEnv = "read-from-env"
	flagYes     = "yes"
)

// migrateSqlCmd represents the sql command
var migrateSqlCmd = &cobra.Command{
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
	RunE: migrateSQLUp,
	Args: func(cmd *cobra.Command, args []string) error {
		if _, err := cmd.Flags().GetBool(flagReadEnv); len(args) != 1 && err != nil {
			return fmt.Errorf("expected one arg or %s flag to be set", flagReadEnv)
		}
		return nil
	},
}

func init() {
	migrateSqlCmd.Flags().BoolP(flagReadEnv, "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	migrateSqlCmd.Flags().BoolP(flagYes, "y", false, "If set all confirmation requests are accepted without user interaction.")
}

func migrateSQLUp(cmd *cobra.Command, args []string) error {
	var d driver.Driver

	if flagx.MustGetBool(cmd, "read-from-env") {
		d = driver.NewDefaultDriver(logrusx.New("", ""), false, nil, "", "", "", false)
		if len(d.Configuration().DSN()) == 0 {
			return fmt.Errorf("when using flag --%s/-e, environment variable DSN must be set", flagReadEnv)
		}
	} else {
		viper.Set(configuration.ViperKeyDSN, args[0])
		d = driver.NewDefaultDriver(logrusx.New("", ""), false, nil, "", "", "", false)
	}

	p := d.Registry().Persister()
	conn := p.Connection(context.Background())
	if conn == nil {
		return fmt.Errorf("migrations can only be executed against a SQL-compatible driver but DSN is not a SQL source.\n")
	}

	if err := conn.Open(); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not open the database connection:\n%+v\n", err)
		return cmdx.FailSilently(cmd)
	}

	// convert migration tables
	if err := p.PrepareMigration(context.Background()); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not convert the migration table:\n%+v\n", err)
		return cmdx.FailSilently(cmd)
	}

	// print migration status
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "The following migration is planned:")
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	if err := p.MigrationStatus(context.Background(), cmd.OutOrStdout()); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not get the migration status:\n%+v\n", err)
		return cmdx.FailSilently(cmd)
	}

	if !flagx.MustGetBool(cmd, "yes") {
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "To skip the next question use flag --yes (at your own risk).")
		if !cmdx.AskForConfirmation("Do you wish to execute this migration plan?", cmd.InOrStdin(), cmd.OutOrStdout()) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Migration aborted.")
			return nil
		}
	}

	// apply migrations
	if err := p.MigrateUp(context.Background()); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not apply migrations:\n%+v\n", err)
		return cmdx.FailSilently(cmd)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Successfully applied migrations!")
	return nil
}

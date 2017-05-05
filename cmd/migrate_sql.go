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

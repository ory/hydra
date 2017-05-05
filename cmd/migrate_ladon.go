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

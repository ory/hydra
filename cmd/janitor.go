package cmd

import (
	"github.com/ory/x/configx"
	"github.com/spf13/cobra"
	"time"
)

var janitorCmd = &cobra.Command{
	Use:   "janitor",
	Short: "Clean the database of old tokens and login/consent requests",
	Long: `This command will cleanup any expired oauth2 tokens as well as login/consent requests.

	### Warning ###

	This is a destructive command and will purge data directly from the database.
	Please use this command with caution if you need to keep historic data for any reason.`,
	Run: cmdHandler.Janitor.Purge,
}

func init() {
	RootCmd.AddCommand(janitorCmd)
	configx.RegisterFlags(janitorCmd.PersistentFlags())
	janitorCmd.Flags().Duration("keep-if-younger", time.Duration(0), "Keep database records that are younger than a specified duration e.g. 1m, 1h, 1d.")
	janitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")

}

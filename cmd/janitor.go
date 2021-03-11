package cmd

import (
	"github.com/ory/x/configx"
	"github.com/spf13/cobra"
)

var janitorCmd = &cobra.Command{
	Use:   "janitor <database-url>",
	Short: "Clean the database of old tokens and login/consent requests",
	Long: `This command will cleanup any expired oauth2 tokens as well as login/consent requests.

	### Warning ###

	This is a destructive command and will purge data directly from the database.
	Please use this command with caution if you need to keep historic data for any reason.`,
	Run: cmdHandler.Janitor.Purge,
}

func init() {
	RootCmd.AddCommand(janitorCmd)
	janitorCmd.Flags().StringP("keep-if-younger", "k", "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("access-lifespan", "a", "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("refresh-lifespan", "r", "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("consent-request-lifespan", "c", "", "Set the login-consent request lifespan e.g. 1s, 1m, 1h")
	janitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(janitorCmd.PersistentFlags())
}

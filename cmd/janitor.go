package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/x/configx"
	"github.com/ory/x/flagx"
)

var janitorCmd = &cobra.Command{
	Use:   "janitor <database-url>",
	Short: "Clean the database of old tokens and login/consent requests",
	Long: `This command will cleanup any expired oauth2 tokens as well as login/consent requests.

Janitor can be used in several ways.

1. By passing the database connection string (DSN) as an argument
Pass the database url (dsn) as an argument to janitor. E.g. janitor <database-url>

2. By passing the DSN as an environment variable
	export DSN=...
	janitor -e

3. By passing a configuration file containing the DSN
janitor -c /path/to/conf.yml

4. Extra *optional* parameters can also be added such as
janitor <database-url> --keep-if-younger=23h --access-lifespan=1h --refresh-lifespan=2h --consent-request-lifespan=10m

Each flag also has a shorthand flag which can be used. To see all options run janitor -h

### Warning ###

This is a destructive command and will purge data directly from the database.
Please use this command with caution if you need to keep historic data for any reason.
`,
	RunE: cmdHandler.Janitor.Purge,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 &&
			!flagx.MustGetBool(cmd, "read-from-env") &&
			len(flagx.MustGetStringSlice(cmd, "config")) == 0 {

			fmt.Printf("%s\n", cmd.UsageString())
			return fmt.Errorf("%s\n%s\n%s\n",
				"A DSN is required as a positional argument when not passing any of the following flags:",
				"- Using the environment variable with flag -e, --read-from-env",
				"- Using the config file with flag -c, --config")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(janitorCmd)
	janitorCmd.Flags().StringP("keep-if-younger", "k", "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("access-lifespan", "a", "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("refresh-lifespan", "r", "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	janitorCmd.Flags().StringP("consent-request-lifespan", "l", "", "Set the login-consent request lifespan e.g. 1s, 1m, 1h")
	janitorCmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(janitorCmd.PersistentFlags())
}

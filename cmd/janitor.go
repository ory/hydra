package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/x/configx"
)

func NewJanitorCmd() *cobra.Command {
	cmd := &cobra.Command{
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
janitor <database-url> --keep-if-younger 23h --access-lifespan 1h --refresh-lifespan 2d --consent-request-lifespan 10m

5. Running only a certain cleanup
janitor <database-url> --tokens

or

janitor <database-url> --requests

or both

janitor <database-url> --tokens --requests

### Warning ###

This is a destructive command and will purge data directly from the database.
Please use this command with caution if you need to keep historic data for any reason.
`,
		RunE: cli.NewHandler().Janitor.RunE,
		Args: cli.NewHandler().Janitor.Args,
	}
	cmd.Flags().String(cli.KeepIfYounger, "", "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	cmd.Flags().String(cli.AccessLifespan, "", "Set the access token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().String(cli.RefreshLifespan, "", "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().String(cli.ConsentRequestLifespan, "", "Set the login/consent request lifespan e.g. 1s, 1m, 1h")
	cmd.Flags().Bool(cli.OnlyRequests, false, "This will only run the cleanup on requests and will skip token cleanup.")
	cmd.Flags().Bool(cli.OnlyTokens, false, "This will only run the cleanup on tokens and will skip requests cleanup.")
	cmd.Flags().BoolP(cli.ReadFromEnv, "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(cmd.PersistentFlags())
	return cmd

}

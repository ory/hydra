package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/x/configx"
)

func NewJanitorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "janitor [<database-url>]",
		Short: "Clean the database of old tokens, login/consent requests and jwt grant issuers",
		Long: `This command will cleanup any expired oauth2 tokens as well as login/consent requests.
This will select records to delete with a limit and delete records in batch to ensure that no table locking issues arise in big production databases.

### Warning ###

This command is in beta. Proceed with caution!

This is a destructive command and will purge data directly from the database.
Please use this command with caution if you need to keep historic data for any reason.

###############

Janitor can be used in several ways.

1. By passing the database connection string (DSN) as an argument
   Pass the database url (dsn) as an argument to janitor. E.g. janitor <database-url>
2. By passing the DSN as an environment variable

		export DSN=...
		janitor -e

3. By passing a configuration file containing the DSN
   janitor -c /path/to/conf.yml
4. Extra *optional* parameters can also be added such as

		janitor --keep-if-younger 23h --access-lifespan 1h --refresh-lifespan 40h --consent-request-lifespan 10m <database-url>

5. Running only a certain cleanup

		janitor --tokens <database-url>

   or

		janitor --requests <database-url>

    or

		janitor --grants <database-url>

   or any combination of them

		janitor --tokens --requests --grants <database-url>
`,
		RunE: cli.NewHandler().Janitor.RunE,
		Args: cli.NewHandler().Janitor.Args,
	}
	cmd.Flags().Int(cli.Limit, 10000, "Limit the number of records retrieved from database for deletion.")
	cmd.Flags().Int(cli.BatchSize, 100, "Define how many records are deleted with each iteration.")
	cmd.Flags().Duration(cli.KeepIfYounger, 0, "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(cli.AccessLifespan, 0, "Set the access token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(cli.RefreshLifespan, 0, "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(cli.ConsentRequestLifespan, 0, "Set the login/consent request lifespan e.g. 1s, 1m, 1h")
	cmd.Flags().Bool(cli.OnlyRequests, false, "This will only run the cleanup on requests and will skip token and trust relationships cleanup.")
	cmd.Flags().Bool(cli.OnlyTokens, false, "This will only run the cleanup on tokens and will skip requests and trust relationships cleanup.")
	cmd.Flags().Bool(cli.OnlyGrants, false, "This will only run the cleanup on trust relationships and will skip requests and token cleanup.")
	cmd.Flags().BoolP(cli.ReadFromEnv, "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(cmd.PersistentFlags())
	return cmd

}

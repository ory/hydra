package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ory/hydra/persistence"

	"github.com/pkg/errors"

	"github.com/ory/x/flagx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/errorsx"
)

const (
	KeepIfYounger          = "keep-if-younger"
	AccessLifespan         = "access-lifespan"
	RefreshLifespan        = "refresh-lifespan"
	ConsentRequestLifespan = "consent-request-lifespan"
	OnlyTokens             = "tokens"
	OnlyRequests           = "requests"
	ReadFromEnv            = "read-from-env"
	Config                 = "config"
)

type JanitorHandler struct{}

func newJanitorHandler() *JanitorHandler {
	return &JanitorHandler{}
}

func (j *JanitorHandler) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "janitor <database-url>",
		Short: "BETA - Clean the database of old tokens and login/consent requests",
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
janitor <database-url> --keep-if-younger 23h --access-lifespan 1h --refresh-lifespan 40h --consent-request-lifespan 10m

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
		RunE: purge,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 &&
				!flagx.MustGetBool(cmd, ReadFromEnv) &&
				len(flagx.MustGetStringSlice(cmd, Config)) == 0 {

				fmt.Printf("%s\n", cmd.UsageString())
				return fmt.Errorf("%s\n%s\n%s\n",
					"A DSN is required as a positional argument when not passing any of the following flags:",
					"- Using the environment variable with flag -e, --read-from-env",
					"- Using the config file with flag -c, --config")
			}

			if !flagx.MustGetBool(cmd, OnlyTokens) && !flagx.MustGetBool(cmd, OnlyRequests) {
				return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
					"Janitor requires either --tokens or --requests or both to be set")
			}

			return nil
		},
	}
	cmd.Flags().Duration(KeepIfYounger, 0, "Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(AccessLifespan, 0, "Set the access token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(RefreshLifespan, 0, "Set the refresh token lifespan e.g. 1s, 1m, 1h.")
	cmd.Flags().Duration(ConsentRequestLifespan, 0, "Set the login/consent request lifespan e.g. 1s, 1m, 1h")
	cmd.Flags().Bool(OnlyRequests, false, "This will only run the cleanup on requests and will skip token cleanup.")
	cmd.Flags().Bool(OnlyTokens, false, "This will only run the cleanup on tokens and will skip requests cleanup.")
	cmd.Flags().BoolP(ReadFromEnv, "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	configx.RegisterFlags(cmd.PersistentFlags())
	return cmd
}

func purge(cmd *cobra.Command, args []string) error {
	var d driver.Registry

	co := []configx.OptionModifier{
		configx.WithFlags(cmd.Flags()),
		configx.SkipValidation(),
	}

	keys := map[string]string{
		AccessLifespan:         config.KeyAccessTokenLifespan,
		RefreshLifespan:        config.KeyRefreshTokenLifespan,
		ConsentRequestLifespan: config.KeyConsentRequestMaxAge,
	}

	for k, v := range keys {
		if x := flagx.MustGetDuration(cmd, k); x > 0 {
			co = append(co, configx.WithValue(v, x))
		}
	}

	notAfter := time.Now()

	if keepYounger := flagx.MustGetDuration(cmd, KeepIfYounger); keepYounger > 0 {
		notAfter = notAfter.Add(-keepYounger)
	}

	if !flagx.MustGetBool(cmd, ReadFromEnv) && len(flagx.MustGetStringSlice(cmd, Config)) == 0 {
		co = append(co, configx.WithValue(config.KeyDSN, args[0]))
	}

	do := []driver.OptionsModifier{
		driver.DisableValidation(),
		driver.DisablePreloading(),
		driver.WithOptions(co...),
	}

	d = driver.New(cmd.Context(), do...)

	if len(d.Config().DSN()) == 0 {
		return fmt.Errorf("%s\n%s\n%s\n", cmd.UsageString(),
			"When using flag -e, environment variable DSN must be set.",
			"When using flag -c, the dsn property should be set.")
	}

	if err := d.Init(cmd.Context()); err != nil {
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Janitor can only be executed against a SQL-compatible driver but DSN is not a SQL source.")
	}

	p := d.Persister()

	var routineFlags []string

	if flagx.MustGetBool(cmd, OnlyTokens) {
		routineFlags = append(routineFlags, OnlyTokens)
	}

	if flagx.MustGetBool(cmd, OnlyRequests) {
		routineFlags = append(routineFlags, OnlyRequests)
	}

	return cleanupRun(cmd.Context(), notAfter, addRoutine(p, routineFlags...)...)
}

func addRoutine(p persistence.Persister, names ...string) []cleanupRoutine {
	var routines []cleanupRoutine
	for _, n := range names {
		switch n {
		case OnlyTokens:
			routines = append(routines, cleanup(p.FlushInactiveAccessTokens, "access tokens"))
			routines = append(routines, cleanup(p.FlushInactiveRefreshTokens, "refresh tokens"))
		case OnlyRequests:
			routines = append(routines, cleanup(p.FlushInactiveLoginConsentRequests, "login-consent requests"))
		}
	}
	return routines
}

type cleanupRoutine func(ctx context.Context, notAfter time.Time) error

func cleanup(cr cleanupRoutine, routineName string) cleanupRoutine {
	return func(ctx context.Context, notAfter time.Time) error {
		if err := cr(ctx, notAfter); err != nil {
			return errors.Wrap(errorsx.WithStack(err), fmt.Sprintf("Could not cleanup inactive %s", routineName))
		}
		fmt.Printf("Successfully completed Janitor run on %s\n", routineName)
		return nil
	}
}

func cleanupRun(ctx context.Context, notAfter time.Time, routines ...cleanupRoutine) error {
	if len(routines) == 0 {
		return errors.New("clean up run received 0 routines")
	}

	for _, r := range routines {
		if err := r(ctx, notAfter); err != nil {
			return err
		}
	}
	return nil
}

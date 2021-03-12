package cli

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/flagx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/errorsx"
)

type JanitorHandler struct{}

func newJanitorHandler() *JanitorHandler {
	return &JanitorHandler{}
}

func (j *JanitorHandler) Purge(cmd *cobra.Command, args []string) error {
	var d driver.Registry

	co := []configx.OptionModifier{
		configx.WithFlags(cmd.Flags()),
		configx.SkipValidation(),
	}

	keys := map[string]string{
		"access-lifespan":          config.KeyAccessTokenLifespan,
		"refresh-lifespan":         config.KeyRefreshTokenLifespan,
		"consent-request-lifespan": config.KeyConsentRequestMaxAge,
	}

	for k, v := range keys {
		if x := flagx.MustGetString(cmd, k); x != "" {
			if xp, err := time.ParseDuration(x); err == nil {
				co = append(co, configx.WithValue(v, xp))
			}
		}
	}

	notAfter := time.Now()

	if keepYounger := flagx.MustGetString(cmd, "keep-if-younger"); keepYounger != "" {
		if keepYoungerDuration, err := time.ParseDuration(keepYounger); err == nil {
			notAfter = notAfter.Add(-keepYoungerDuration)
		}
	}

	if !flagx.MustGetBool(cmd, "read-from-env") && len(flagx.MustGetStringSlice(cmd, "config")) == 0 {
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

	p := d.Persister()

	conn := p.Connection(cmd.Context())

	if conn == nil {
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Janitor can only be executed against a SQL-compatible driver but DSN is not a SQL source.")
	}

	if err := conn.Open(); err != nil {
		return errors.Wrap(errorsx.WithStack(err), "Could not open the database connection")
	}

	if err := p.FlushInactiveAccessTokens(cmd.Context(), notAfter); err != nil {
		return errors.Wrap(errorsx.WithStack(err), "Could not flush inactive access tokens")
	}

	if err := p.FlushInactiveRefreshTokens(cmd.Context(), notAfter); err != nil {
		return errors.Wrap(errorsx.WithStack(err), "Could not flush inactive refresh tokens")
	}

	if err := p.FlushInactiveLoginConsentRequests(cmd.Context(), notAfter); err != nil {
		return errors.Wrap(errorsx.WithStack(err), "Could not flush inactive login/consent requests")
	}

	fmt.Print("Successfully completed Janitor!\n")
	return nil
}

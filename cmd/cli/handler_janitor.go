package cli

import (
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/errorsx"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type JanitorHandler struct{}

func newJanitorHandler() *JanitorHandler {
	return &JanitorHandler{}
}

func (j *JanitorHandler) Purge(cmd *cobra.Command, args []string) {
	var d driver.Registry

	co := []configx.OptionModifier{
		configx.WithFlags(cmd.Flags()),
		configx.SkipValidation(),
	}

	keys := map[string]string{
		"access-lifespan":        config.KeyAccessTokenLifespan,
		"refresh-lifespan":       config.KeyRefreshTokenLifespan,
		"consent-request-lifespan": config.KeyConsentRequestMaxAge,
	}

	for k, v := range keys {
		if x, err := cmd.Flags().GetString(k); err == nil && x != "" {
			if xp, err := time.ParseDuration(x); err == nil {
				co = append(co, configx.WithValue(v, xp))
			}
		}
	}

	notAfter := time.Now()

	if keepYounger, err := cmd.Flags().GetString("keep-if-younger"); err == nil && keepYounger != "" {
		if keepYoungerDuration, err := time.ParseDuration(keepYounger); err == nil {
			notAfter = notAfter.Add(-keepYoungerDuration)
		}
	}

	if ok, _ := cmd.Flags().GetBool("read-from-env"); !ok {
		if len(args) == 0 {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
			return
		}
		co = append(co, configx.WithValue(config.KeyDSN, args[0]))
	}

	do := []driver.OptionsModifier{
		driver.DisableValidation(),
		driver.DisablePreloading(),
		driver.WithOptions(co...),
	}

	d = driver.New(cmd.Context(), do...)

	if len(d.Config().DSN()) == 0 {
		fmt.Println(fmt.Sprintf("%s\n%s", cmd.UsageString(),
			"When using flag -e, environment variable DSN must be set"))
		os.Exit(1)
		return
	}

	p := d.Persister()

	conn := p.Connection(cmd.Context())

	if conn == nil {
		fmt.Println(fmt.Sprintf("%s\n%s\n", cmd.UsageString(),
			"Janitor can only be executed against a SQL-compatible driver but DSN is not a SQL source."))
		os.Exit(1)
		return
	}

	if err := conn.Open(); err != nil {
		fmt.Printf("Could not open the database connection:\n%+v\n", err)
		os.Exit(1)
		return
	}

	if err := p.FlushInactiveAccessTokens(cmd.Context(), notAfter); err != nil {
		fmt.Printf("Could not flush inactive access tokens:\n%+v\n", errorsx.WithStack(err))
		os.Exit(1)
		return
	}

	if err := p.FlushInactiveRefreshTokens(cmd.Context(), notAfter); err != nil {
		fmt.Printf("Could not flush inactive refresh tokens:\n%+v\n", errorsx.WithStack(err))
		os.Exit(1)
		return
	}

	if err := p.FlushInactiveLoginConsentRequests(cmd.Context(), notAfter); err != nil {
		fmt.Printf("Could not flush inactive login/consent requests:\n%+v\n", errorsx.WithStack(err))
		os.Exit(1)
		return
	}

	fmt.Println("Successfully completed Janitor!")
}

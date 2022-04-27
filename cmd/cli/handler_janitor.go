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
	Limit                  = "limit"
	BatchSize              = "batch-size"
	KeepIfYounger          = "keep-if-younger"
	AccessLifespan         = "access-lifespan"
	RefreshLifespan        = "refresh-lifespan"
	ConsentRequestLifespan = "consent-request-lifespan"
	OnlyTokens             = "tokens"
	OnlyRequests           = "requests"
	OnlyGrants             = "grants"
	ReadFromEnv            = "read-from-env"
	Config                 = "config"
)

type JanitorHandler struct{}

func NewJanitorHandler() *JanitorHandler {
	return &JanitorHandler{}
}

func (_ *JanitorHandler) Args(cmd *cobra.Command, args []string) error {
	if len(args) == 0 &&
		!flagx.MustGetBool(cmd, ReadFromEnv) &&
		len(flagx.MustGetStringSlice(cmd, Config)) == 0 {

		fmt.Printf("%s\n", cmd.UsageString())
		return fmt.Errorf("%s\n%s\n%s\n",
			"A DSN is required as a positional argument when not passing any of the following flags:",
			"- Using the environment variable with flag -e, --read-from-env",
			"- Using the config file with flag -c, --config")
	}

	if !flagx.MustGetBool(cmd, OnlyTokens) && !flagx.MustGetBool(cmd, OnlyRequests) && !flagx.MustGetBool(cmd, OnlyGrants) {
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Janitor requires at least one of --tokens, --requests or --grants to be set")
	}

	limit := flagx.MustGetInt(cmd, Limit)
	batchSize := flagx.MustGetInt(cmd, BatchSize)
	if limit <= 0 || batchSize <= 0 {
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Values for --limit and --batch-size should both be greater than 0")
	}
	if batchSize > limit {
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Value for --batch-size must not be greater than value for --limit")
	}

	return nil
}

func (_ *JanitorHandler) RunE(cmd *cobra.Command, args []string) error {
	return purge(cmd, args)
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

	p := d.Persister()

	limit := flagx.MustGetInt(cmd, Limit)
	batchSize := flagx.MustGetInt(cmd, BatchSize)

	var routineFlags []string

	if flagx.MustGetBool(cmd, OnlyTokens) {
		routineFlags = append(routineFlags, OnlyTokens)
	}

	if flagx.MustGetBool(cmd, OnlyRequests) {
		routineFlags = append(routineFlags, OnlyRequests)
	}

	if flagx.MustGetBool(cmd, OnlyGrants) {
		routineFlags = append(routineFlags, OnlyGrants)
	}

	return cleanupRun(cmd.Context(), notAfter, limit, batchSize, addRoutine(p, routineFlags...)...)
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
		case OnlyGrants:
			routines = append(routines, cleanup(p.FlushInactiveGrants, "grants"))
		}
	}
	return routines
}

type cleanupRoutine func(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

func cleanup(cr cleanupRoutine, routineName string) cleanupRoutine {
	return func(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
		if err := cr(ctx, notAfter, limit, batchSize); err != nil {
			return errors.Wrap(errorsx.WithStack(err), fmt.Sprintf("Could not cleanup inactive %s", routineName))
		}
		fmt.Printf("Successfully completed Janitor run on %s\n", routineName)
		return nil
	}
}

func cleanupRun(ctx context.Context, notAfter time.Time, limit int, batchSize int, routines ...cleanupRoutine) error {
	if len(routines) == 0 {
		return errors.New("clean up run received 0 routines")
	}

	for _, r := range routines {
		if err := r(ctx, notAfter, limit, batchSize); err != nil {
			return err
		}
	}
	return nil
}

// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ory/x/servicelocatorx"

	"github.com/ory/hydra/v2/persistence"

	"github.com/pkg/errors"

	"github.com/ory/x/flagx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
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

type JanitorHandler struct {
	slOpts []servicelocatorx.Option
	dOpts  []driver.OptionsModifier
	cOpts  []configx.OptionModifier
}

func NewJanitorHandler(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *JanitorHandler {
	return &JanitorHandler{
		slOpts: slOpts,
		dOpts:  dOpts,
		cOpts:  cOpts,
	}
}

func (*JanitorHandler) Args(cmd *cobra.Command, args []string) error {
	if len(args) == 0 &&
		!flagx.MustGetBool(cmd, ReadFromEnv) &&
		len(flagx.MustGetStringSlice(cmd, Config)) == 0 {

		fmt.Printf("%s\n", cmd.UsageString())
		//lint:ignore ST1005 formatted error string used in CLI output
		return fmt.Errorf("%s\n%s\n%s\n",
			"A DSN is required as a positional argument when not passing any of the following flags:",
			"- Using the environment variable with flag -e, --read-from-env",
			"- Using the config file with flag -c, --config")
	}

	if !flagx.MustGetBool(cmd, OnlyTokens) && !flagx.MustGetBool(cmd, OnlyRequests) && !flagx.MustGetBool(cmd, OnlyGrants) {
		//lint:ignore ST1005 formatted error string used in CLI output
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Janitor requires at least one of --tokens, --requests or --grants to be set")
	}

	limit := flagx.MustGetInt(cmd, Limit)
	batchSize := flagx.MustGetInt(cmd, BatchSize)
	if limit <= 0 || batchSize <= 0 {
		//lint:ignore ST1005 formatted error string used in CLI output
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Values for --limit and --batch-size should both be greater than 0")
	}
	if batchSize > limit {
		//lint:ignore ST1005 formatted error string used in CLI output
		return fmt.Errorf("%s\n%s\n", cmd.UsageString(),
			"Value for --batch-size must not be greater than value for --limit")
	}

	return nil
}

func (j *JanitorHandler) RunE(cmd *cobra.Command, args []string) error {
	return purge(cmd, args, servicelocatorx.NewOptions(j.slOpts...), j.dOpts)
}

func purge(cmd *cobra.Command, args []string, sl *servicelocatorx.Options, dOpts []driver.OptionsModifier) error {
	ctx := cmd.Context()
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

	do := append(dOpts,
		driver.DisableValidation(),
		driver.DisablePreloading(),
		driver.WithOptions(co...),
	)

	d, err := driver.New(ctx, sl, do)
	if err != nil {
		return errors.Wrap(err, "Could not create driver")
	}

	if len(d.Config().DSN()) == 0 {
		//lint:ignore ST1005 formatted error string used in CLI output
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

	return cleanupRun(cmd.Context(), notAfter, limit, batchSize, addRoutine(cmd.OutOrStdout(), p, routineFlags...)...)
}

func addRoutine(out io.Writer, p persistence.Persister, names ...string) []cleanupRoutine {
	var routines []cleanupRoutine
	for _, n := range names {
		switch n {
		case OnlyTokens:
			routines = append(routines, cleanup(out, p.FlushInactiveAccessTokens, "access tokens"))
			routines = append(routines, cleanup(out, p.FlushInactiveRefreshTokens, "refresh tokens"))
		case OnlyRequests:
			routines = append(routines, cleanup(out, p.FlushInactiveLoginConsentRequests, "login-consent requests"))
		case OnlyGrants:
			routines = append(routines, cleanup(out, p.FlushInactiveGrants, "grants"))
		}
	}
	return routines
}

type cleanupRoutine func(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

func cleanup(out io.Writer, cr cleanupRoutine, routineName string) cleanupRoutine {
	return func(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
		if err := cr(ctx, notAfter, limit, batchSize); err != nil {
			return errors.Wrap(errorsx.WithStack(err), fmt.Sprintf("Could not cleanup inactive %s", routineName))
		}
		fmt.Fprintf(out, "Successfully completed Janitor run on %s\n", routineName)
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

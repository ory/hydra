// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/configx"
	"github.com/ory/x/popx"
)

type MigrateHandler struct {
	dOpts []driver.OptionsModifier
}

func newMigrateHandler(dOpts []driver.OptionsModifier) *MigrateHandler {
	return &MigrateHandler{
		dOpts: dOpts,
	}
}

func (h *MigrateHandler) makeMigrationManager(cmd *cobra.Command, args []string) (*sql.MigrationManager, error) {
	opts := append([]driver.OptionsModifier{
		driver.WithConfigOptions(
			configx.SkipValidation(),
			configx.WithFlags(cmd.Flags())),
		driver.DisableValidation(),
		driver.DisablePreloading(),
		driver.SkipNetworkInit(),
	}, h.dOpts...)
	if len(args) > 0 {
		opts = append(opts, driver.WithConfigOptions(
			configx.WithValue(config.KeyDSN, args[0]),
		))
	}

	d, err := driver.New(
		cmd.Context(),
		opts...)
	if err != nil {
		return nil, err
	}
	if len(d.Config().DSN()) == 0 {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "No DSN provided. Please provide a DSN as the first argument or set the DSN environment variable.")
		return nil, cmdx.FailSilently(cmd)
	}

	return d.Migrator(), nil
}

func (h *MigrateHandler) MigrateSQLUp(cmd *cobra.Command, args []string) (err error) {
	p, err := h.makeMigrationManager(cmd, args)
	if err != nil {
		return err
	}
	return popx.MigrateSQLUp(cmd, p)
}

func (h *MigrateHandler) MigrateSQLDown(cmd *cobra.Command, args []string) (err error) {
	p, err := h.makeMigrationManager(cmd, args)
	if err != nil {
		return err
	}
	return popx.MigrateSQLDown(cmd, p)
}

func (h *MigrateHandler) MigrateStatus(cmd *cobra.Command, args []string) error {
	p, err := h.makeMigrationManager(cmd, args)
	if err != nil {
		return err
	}
	return popx.MigrateStatus(cmd, p)
}

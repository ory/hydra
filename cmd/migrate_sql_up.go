// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/popx"
	"github.com/ory/x/servicelocatorx"
)

func NewMigrateSQLUpCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *cobra.Command {
	return popx.NewMigrateSQLDownCmd("hydra", cli.NewHandler(slOpts, dOpts, cOpts).Migration.MigrateSQLUp)
}

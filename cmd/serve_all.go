// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/server"
	"github.com/ory/hydra/v2/driver"
)

func NewServeAllCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Serves both public and administrative HTTP/2 APIs",
		Long: `Starts a process which listens on two ports for public and administrative HTTP/2 API requests.

If you want more granular control (e.g. different TLS settings) over each API group (administrative, public) you
can run "serve admin" and "serve public" separately.

This command exposes a variety of controls via environment variables. You can
set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux,
you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 hydra"

All possible controls are listed below. This command exposes exposes command line flags, which are listed below
the controls section.

` + serveControls,
		RunE: server.RunServeAll(dOpts),
	}
}

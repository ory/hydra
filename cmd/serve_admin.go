// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/servicelocatorx"

	"github.com/ory/hydra/v2/cmd/server"
)

// adminCmd represents the admin command
func NewServeAdminCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "Serves Administrative HTTP/2 APIs",
		Long: `This command opens one port and listens to HTTP/2 API requests. The exposed API handles administrative
requests like managing OAuth 2.0 Clients, JSON Web Keys, login and consent sessions, and others.

This command is configurable using the same options available to "serve public" and "serve all".

It is generally recommended to use this command only if you require granular control over the administrative and public APIs.
For example, you might want to run different TLS certificates or CORS settings on the public and administrative API.

This command does not work with the "memory" database. Both services (administrative, public) MUST use the same database
connection to be able to synchronize.

` + serveControls,
		RunE: server.RunServeAdmin(slOpts, dOpts, cOpts),
	}
}

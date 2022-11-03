// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/configx"
)

var serveControls = `## Configuration

Ory Hydra can be configured using environment variables as well as a configuration file. For more information
on configuration options, open the configuration documentation:

>> https://www.ory.sh/hydra/docs/reference/configuration <<
`

// serveCmd represents the host command
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Parent command for starting public and administrative HTTP/2 APIs",
		Long: `Ory Hydra exposes two ports, a public and an administrative port. The public port is responsible
for handling requests from the public internet, such as the OAuth 2.0 Authorize and Token URLs. The administrative
port handles administrative requests like creating OAuth 2.0 Clients, managing JSON Web Keys, and managing User Login
and Consent sessions.

It is recommended to run "hydra serve all". If you need granular control over CORS settings or similar, you may
want to run "hydra serve admin" and "hydra serve public" separately.

To learn more about each individual command, run:

- hydra help serve all
- hydra help serve admin
- hydra help serve public

All sub-commands share command line flags and configuration options.

` + serveControls,
	}

	configx.RegisterFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().Bool("dev", false, "Disables critical security checks to improve local development experience. Do not use in production.")
	cmd.PersistentFlags().Bool("sqa-opt-out", false, "Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa")

	return cmd
}

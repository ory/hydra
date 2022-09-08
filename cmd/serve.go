/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

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

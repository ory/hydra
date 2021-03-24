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

	"github.com/ory/hydra/cmd/cli"
)

func NewClientsImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <path/to/file.json> [<path/to/other/file.json>...]",
		Short: "Import OAuth 2.0 Clients from one or more JSON files",
		Long: `This command reads in each listed JSON file and imports their contents as OAuth 2.0 Clients.

The format for the JSON file is:

{
  "client_id": "...",
  "client_secret": "...",
  // ... all other fields of the OAuth 2.0 Client model are allowed here
}

Please be aware that this command does not update existing clients. If the client exists already, this command will fail.

Example:
	hydra clients import client-1.json

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or "--keybase" flag, for example:
	hydra clients import client-1.json --keybase keybase_username
`,
		Run: cli.NewHandler().Clients.ImportClients,
	}

	// encrypt client secret options
	cmd.Flags().String("pgp-key", "", "Base64 encoded PGP encryption key for encrypting client secret")
	cmd.Flags().String("pgp-key-url", "", "PGP encryption key URL for encrypting client secret")
	cmd.Flags().String("keybase", "", "Keybase username for encrypting client secret")

	return cmd
}

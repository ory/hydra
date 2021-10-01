// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
)

// keysImportCmd represents the import command
func NewKeysImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <set> <file-1> [<file-2> [<file-3 [<...>]]]",
		Short: "Imports OAuth 2.0 Clients from one or more JSON files to the JSON Web Key Store",
		Long: `This command allows you to import OAuth 2.0 Clients from one or more JSON files to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If the JSON Web Key Set exists already,
the imported keys will be added to that set. Otherwise, a new set will be created.

Please be aware that importing a private key does not automatically import its public key as well.

Examples:
	hydra keys import my-set ./path/to/jwk.json ./path/to/jwk-2.json
	hydra keys import my-set ./path/to/rsa.key ./path/to/rsa.pub --default-key-id cae6b214-fb1e-4ebc-9019-95286a62eabc
`,
		Run: cli.NewHandler().Keys.ImportKeys,
	}
	cmd.Flags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
	cmd.Flags().String("default-key-id", "", "A fallback value for keys without \"kid\" attribute to be stored with a common \"kid\", e.g. private/public keypairs")
	return cmd
}

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
)

// keysImportCmd represents the import command
var keysImportCmd = &cobra.Command{
	Use:   "import <set> <file-1> [<file-2> [<file-3 [<...>]]]",
	Short: "Imports cryptographic keys of any format to the JSON Web Key Store",
	Long: `This command allows you to import cryptographic keys to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If the JSON Web Key Set exists already,
the imported keys will be added to that set. Otherwise, a new set will be created.

Please be aware that importing a private key does not automatically import its public key as well.

Examples:
	hydra keys import my-set ./path/to/jwk.json ./path/to/jwk-2.json
	hydra keys import my-set ./path/to/rsa.key ./path/to/rsa.pub
`,
	Run: cmdHandler.Keys.ImportKeys,
}

func init() {
	keysCmd.AddCommand(keysImportCmd)
	keysImportCmd.Flags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
}

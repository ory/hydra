// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
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

// createCmd represents the create command
var clientsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new OAuth2 client",
	Long: `This command creates a basic OAuth2 client. Always specify at least one redirect url.

Example:
  hydra clients create -n "my app" -c http://localhost/cb -g authorization_code -r code -a core,foobar
`,
	Run: cmdHandler.Clients.CreateClient,
}

func init() {
	clientsCmd.AddCommand(clientsCreateCmd)
	clientsCreateCmd.Flags().String("id", "", "Give the client this id")
	clientsCreateCmd.Flags().StringSliceP("callbacks", "c", []string{}, "REQUIRED list of allowed callback URLs")
	clientsCreateCmd.Flags().StringSliceP("grant-types", "g", []string{"authorization_code"}, "A list of allowed grant types")
	clientsCreateCmd.Flags().StringSliceP("response-types", "r", []string{"code"}, "A list of allowed response types")
	clientsCreateCmd.Flags().StringSliceP("allowed-scopes", "a", []string{""}, "A list of allowed scopes")
	clientsCreateCmd.Flags().Bool("is-public", false, "Use this flag to create a public client")
	clientsCreateCmd.Flags().String("secret", "", "Provide the client's secret")
	clientsCreateCmd.Flags().StringP("name", "n", "", "The client's name")
}

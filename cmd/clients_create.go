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
)

// createCmd represents the create command
var clientsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new OAuth 2.0 Client",
	Long: `This command creates an OAuth 2.0 Client which can be used to perform various OAuth 2.0 Flows like
the Authorize Code, Implicit, Refresh flow.

ORY Hydra implements the OpenID Connect Dynamic Client registration specification. Most flags are supported by this command
as well.

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
	clientsCreateCmd.Flags().StringSliceP("scope", "a", []string{""}, "The scope the client is allowed to request")
	clientsCreateCmd.Flags().StringSlice("audience", []string{}, "The audience this client is allowed to request")
	clientsCreateCmd.Flags().String("token-endpoint-auth-method", "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are \"client_secret_post\", \"client_secret_basic\", \"private_key_jwt\", and \"none\"")
	clientsCreateCmd.Flags().String("jwks-uri", "", "Define the URL where the JSON Web Key Set should be fetched from when performing the \"private_key_jwt\" client authentication method")
	clientsCreateCmd.Flags().String("policy-uri", "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data")
	clientsCreateCmd.Flags().String("tos-uri", "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client")
	clientsCreateCmd.Flags().String("client-uri", "", "A URL string of a web page providing information about the client")
	clientsCreateCmd.Flags().String("logo-uri", "", "A URL string that references a logo for the client")
	clientsCreateCmd.Flags().String("subject-type", "public", "A URL string that references a logo for the client")
	clientsCreateCmd.Flags().String("secret", "", "Provide the client's secret")
	clientsCreateCmd.Flags().StringP("name", "n", "", "The client's name")
}

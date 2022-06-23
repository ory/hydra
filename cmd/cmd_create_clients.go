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

func NewClientsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Create a new OAuth 2.0 Client",
		Example: `hydra create client -n "my app" -c http://localhost/cb -g authorization_code -r code -a core,foobar
`,
		Long: `This command creates an OAuth 2.0 Client which can be used to perform various OAuth 2.0 Flows like
the Authorize Code, Implicit, Refresh flow.

Ory Hydra implements the OpenID Connect Dynamic Client registration specification. Most flags are supported by this command
as well.

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or "--keybase" flag, for example:

  hydra create client -n "my app" -g client_credentials -r token -a core,foobar --keybase keybase_username
`,
		Run: cli.NewHandler().Clients.CreateClient,
	}
	cmd.Flags().StringSlice("metadata", nil, "Metadata is an arbitrary JSON String of your choosing.")
	cmd.Flags().String("owner", "", "The owner of this client, typically email addresses or a user ID.")
	cmd.Flags().StringSlice("contact", nil, "A list representing ways to contact people responsible for this client, typically email addresses.")
	cmd.Flags().StringSlice("request-uri", nil, "Array of request_uri values that are pre-registered by the RP for use at the OP.")
	cmd.Flags().String("request-object-signing-alg", "RS256", "Algorithm that must be used for signing Request Objects sent to the OP.")
	cmd.Flags().String("sector-identifier-uri", "", "URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.")
	cmd.Flags().StringSliceP("redirect-uri", "c", []string{}, "List of allowed OAuth2 Redirect URIs.")
	cmd.Flags().StringSliceP("grant-type", "g", []string{"authorization_code"}, "A list of allowed grant types.")
	cmd.Flags().StringSliceP("response-type", "r", []string{"code"}, "A list of allowed response types.")
	cmd.Flags().StringSliceP("scope", "a", []string{""}, "The scope the client is allowed to request.")
	cmd.Flags().StringSlice("audience", []string{}, "The audience this client is allowed to request.")
	cmd.Flags().String("token-endpoint-auth-method", "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are `client_secret_post`, `client_secret_basic`, `private_key_jwt`, and `none`.")
	cmd.Flags().String("jwks-uri", "", "Define the URL where the JSON Web Key Set should be fetched from when performing the `private_key_jwt` client authentication method.")
	cmd.Flags().String("policy-uri", "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.")
	cmd.Flags().String("tos-uri", "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.")
	cmd.Flags().String("client-uri", "", "A URL string of a web page providing information about the client")
	cmd.Flags().String("logo-uri", "", "A URL string that references a logo for the client")
	cmd.Flags().StringSlice("allowed-cors-origin", []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	cmd.Flags().String("subject-type", "public", "A identifier algorithm. Valid values are `public` and `pairwise`.")
	cmd.Flags().String("secret", "", "Provide the client's secret.")
	cmd.Flags().StringP("name", "n", "", "The client's name.")
	cmd.Flags().StringSlice("post-logout-callback", []string{}, "List of allowed URLs to be redirected to after a logout.")

	// back-channel logout options
	cmd.Flags().Bool("backchannel-logout-session-required", false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the backchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String("backchannel-logout-callback", "", "Client URL that will cause the client to log itself out when sent a Logout Token by Hydra.")

	// front-channel logout options
	cmd.Flags().Bool("frontchannel-logout-session-required", false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the frontchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String("frontchannel-logout-callback", "", "Client URL that will cause the client to log itself out when rendered in an iframe by Hydra.")

	// encrypt client secret options
	cmd.Flags().String("pgp-key", "", "Base64 encoded PGP encryption key for encrypting client secret.")
	cmd.Flags().String("pgp-key-url", "", "PGP encryption key URL for encrypting client secret.")
	cmd.Flags().String("keybase", "", "Keybase username for encrypting client secret.")

	return cmd
}

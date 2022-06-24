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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go"
	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/cmd/cli"
)

const (
	allowedCORSOrigin                 = "allowed-cors-origin"
	audience                          = "audience"
	backchannelLogoutCallback         = "backchannel-logout-callback"
	name                              = "name"
	clientURI                         = "client-uri"
	contact                           = "contact"
	frontchannelLogoutSessionRequired = "frontchannel-logout-session-required"
	frontchannelLogoutCallback        = "frontchannel-logout-callback"
	grantType                         = "grant-type"
	jwksURI                           = "jwks-uri"
	logoURI                           = "logo-uri"
	metadata                          = "metadata"
	owner                             = "owner"
	policyURI                         = "policy-uri"
	postLogoutCallback                = "post-logout-callback"
	redirectURI                       = "redirect-uri"
	requestObjectSigningAlg           = "request-object-signing-alg"
	requestURI                        = "request-uri"
	responseType                      = "response-type"
	scope                             = "scope"
	sectorIdentifierURI               = "sector-identifier-uri"
	subjectType                       = "subject-type"
	tokenEndpointAuthMethod           = "token-endpoint-auth-method"
	secret                            = "secret"
	tosURI                            = "tos-uri"
	backchannelLogoutSessionRequired  = "backchannel-logout-session-required"
)

func NewCreateClientsCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client",
		Short:   "Create an OAuth 2.0 Client",
		Example: fmt.Sprintf(`%[1]s create client -n "my app" -c http://localhost/cb -g authorization_code -r code -a core,foobar`, root.Use),
		Long: fmt.Sprintf(`This command creates an OAuth 2.0 Client which can be used to perform various OAuth 2.0 Flows like
the Authorize Code, Implicit, Refresh flow. This command allows settings all fields defined in the OpenID Connect Dynamic Client Registration standard.

To encrypt a auto-generated OAuth2 Client Secret, use flags `+"`--pgp-key`"+`, `+"`--pgp-key-url`"+` or `+"`--keybase`"+` flag, for example:

  %[1]s create client -n "my app" -g client_credentials -r token -a core,foobar --keybase keybase_username
`, root.Use),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			ek, encryptSecret, err := cli.NewEncryptionKey(cmd, nil)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to load encryption key: %s", err)
				return err
			}

			secret := flagx.MustGetString(cmd, secret)
			client, _, err := m.AdminApi.CreateOAuth2Client(cmd.Context()).OAuth2Client(hydra.OAuth2Client{
				AllowedCorsOrigins:                flagx.MustGetStringSlice(cmd, allowedCORSOrigin),
				Audience:                          flagx.MustGetStringSlice(cmd, audience),
				BackchannelLogoutSessionRequired:  pointerx.Bool(flagx.MustGetBool(cmd, backchannelLogoutSessionRequired)),
				BackchannelLogoutUri:              pointerx.String(flagx.MustGetString(cmd, backchannelLogoutCallback)),
				ClientName:                        pointerx.String(flagx.MustGetString(cmd, name)),
				ClientSecret:                      pointerx.String(secret),
				ClientUri:                         pointerx.String(flagx.MustGetString(cmd, clientURI)),
				Contacts:                          flagx.MustGetStringSlice(cmd, contact),
				FrontchannelLogoutSessionRequired: pointerx.Bool(flagx.MustGetBool(cmd, frontchannelLogoutSessionRequired)),
				FrontchannelLogoutUri:             pointerx.String(flagx.MustGetString(cmd, frontchannelLogoutCallback)),
				GrantTypes:                        flagx.MustGetStringSlice(cmd, grantType),
				JwksUri:                           pointerx.String(flagx.MustGetString(cmd, jwksURI)),
				LogoUri:                           pointerx.String(flagx.MustGetString(cmd, logoURI)),
				Metadata:                          json.RawMessage(flagx.MustGetString(cmd, metadata)),
				Owner:                             pointerx.String(flagx.MustGetString(cmd, owner)),
				PolicyUri:                         pointerx.String(flagx.MustGetString(cmd, policyURI)),
				PostLogoutRedirectUris:            flagx.MustGetStringSlice(cmd, postLogoutCallback),
				RedirectUris:                      flagx.MustGetStringSlice(cmd, redirectURI),
				RequestObjectSigningAlg:           pointerx.String(flagx.MustGetString(cmd, requestObjectSigningAlg)),
				RequestUris:                       flagx.MustGetStringSlice(cmd, requestURI),
				ResponseTypes:                     flagx.MustGetStringSlice(cmd, responseType),
				Scope:                             pointerx.String(strings.Join(flagx.MustGetStringSlice(cmd, scope), "")),
				SectorIdentifierUri:               pointerx.String(flagx.MustGetString(cmd, sectorIdentifierURI)),
				SubjectType:                       pointerx.String(flagx.MustGetString(cmd, subjectType)),
				TokenEndpointAuthMethod:           pointerx.String(flagx.MustGetString(cmd, tokenEndpointAuthMethod)),
				TosUri:                            pointerx.String(flagx.MustGetString(cmd, tosURI)),
			}).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			if client.ClientSecret == nil && len(secret) > 0 {
				client.ClientSecret = pointerx.String(secret)
			}

			if encryptSecret && client.ClientSecret != nil {
				enc, err := ek.Encrypt([]byte(*client.ClientSecret))
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to encrypt client secret: %s", err)
					return cmdx.FailSilently(cmd)
				}

				client.ClientSecret = pointerx.String(enc.Base64Encode())
			}

			cmdx.PrintRow(cmd, (*outputOAuth2Client)(client))
			return nil
		},
	}

	cmd.Flags().String(metadata, "{}", "Metadata is an arbitrary JSON String of your choosing.")
	cmd.Flags().String(owner, "", "The owner of this client, typically email addresses or a user ID.")
	cmd.Flags().StringSlice(contact, nil, "A list representing ways to contact people responsible for this client, typically email addresses.")
	cmd.Flags().StringSlice(requestURI, nil, "Array of request_uri values that are pre-registered by the RP for use at the OP.")
	cmd.Flags().String(requestObjectSigningAlg, "RS256", "Algorithm that must be used for signing Request Objects sent to the OP.")
	cmd.Flags().String(sectorIdentifierURI, "", "URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.")
	cmd.Flags().StringSliceP(redirectURI, "c", []string{}, "List of allowed OAuth2 Redirect URIs.")
	cmd.Flags().StringSliceP(grantType, "g", []string{"authorization_code"}, "A list of allowed grant types.")
	cmd.Flags().StringSliceP(responseType, "r", []string{"code"}, "A list of allowed response types.")
	cmd.Flags().StringSliceP(scope, "a", []string{}, "The scope the client is allowed to request.")
	cmd.Flags().StringSlice(audience, []string{}, "The audience this client is allowed to request.")
	cmd.Flags().String(tokenEndpointAuthMethod, "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are `client_secret_post`, `client_secret_basic`, `private_key_jwt`, and `none`.")
	cmd.Flags().String(jwksURI, "", "Define the URL where the JSON Web Key Set should be fetched from when performing the `private_key_jwt` client authentication method.")
	cmd.Flags().String(policyURI, "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.")
	cmd.Flags().String(tosURI, "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.")
	cmd.Flags().String(clientURI, "", "A URL string of a web page providing information about the client")
	cmd.Flags().String(logoURI, "", "A URL string that references a logo for the client")
	cmd.Flags().StringSlice(allowedCORSOrigin, []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	cmd.Flags().String(subjectType, "public", "A identifier algorithm. Valid values are `public` and `pairwise`.")
	cmd.Flags().String(secret, "", "Provide the client's secret.")
	cmd.Flags().StringP(name, "n", "", "The client's name.")
	cmd.Flags().StringSlice(postLogoutCallback, []string{}, "List of allowed URLs to be redirected to after a logout.")

	// back-channel logout options
	cmd.Flags().Bool(backchannelLogoutSessionRequired, false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the backchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String(backchannelLogoutCallback, "", "Client URL that will cause the client to log itself out when sent a Logout Token by Hydra.")

	// front-channel logout options
	cmd.Flags().Bool(frontchannelLogoutSessionRequired, false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the frontchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String(frontchannelLogoutCallback, "", "Client URL that will cause the client to log itself out when rendered in an iframe by Hydra.")

	// encrypt client secret options
	cmd.Flags().String(cli.FlagEncryptionPGPKey, "", "Base64 encoded PGP encryption key for encrypting client secret.")
	cmd.Flags().String(cli.FlagEncryptionPGPKeyURL, "", "PGP encryption key URL for encrypting client secret.")
	cmd.Flags().String(cli.FlagEncryptionKeybase, "", "Keybase username for encrypting client secret.")

	return cmd
}

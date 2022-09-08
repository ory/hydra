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
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/urlx"
)

func NewPerformClientCredentialsCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client-credentials",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`%[1]s perform client-credentials --client-id ... --client-secret`, parent.Use),
		Short:   "Perform the OAuth2 Client Credentials Flow",
		Long: `Performs the OAuth 2.0 Client Credentials Flow. Useful to exchange a client_id and client_secret for an access_token.
using the CLI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, target, err := cmdx.NewClient(cmd)
			if err != nil {
				return err
			}

			ctx := context.WithValue(cmd.Context(), oauth2.HTTPClient, hc)

			scopes := flagx.MustGetStringSlice(cmd, "scope")
			audience := flagx.MustGetStringSlice(cmd, "audience")

			clientID := flagx.MustGetString(cmd, "client-id")
			clientSecret := flagx.MustGetString(cmd, "client-secret")
			if clientID == "" || clientSecret == "" {
				fmt.Print(cmd.UsageString())
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET.")
				return cmdx.FailSilently(cmd)
			}

			oauthConfig := clientcredentials.Config{
				ClientID:       clientID,
				ClientSecret:   clientSecret,
				TokenURL:       urlx.AppendPaths(target, "/oauth2/token").String(),
				Scopes:         scopes,
				EndpointParams: url.Values{"audience": {strings.Join(audience, " ")}},
			}

			t, err := oauthConfig.Token(ctx)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not retrieve access token because: %s", err)
				return cmdx.FailSilently(cmd)
			}

			cmdx.PrintRow(cmd, (*outputOAuth2Token)(t))
			return nil
		},
	}

	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID.")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET.")
	cmd.Flags().StringSlice("scope", []string{}, "OAuth2 scope to request.")
	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience.")

	return cmd
}

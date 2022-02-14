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
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ory/hydra/cmd/cli"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/urlx"
)

type transporter struct {
	*http.Transport
	FakeTLSTermination bool
}

func (t *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.FakeTLSTermination {
		req.Header.Set("X-Forwarded-Proto", "https")
	}

	return t.Transport.RoundTrip(req)
}

func NewTokenClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "An exemplary OAuth 2.0 Client performing the OAuth 2.0 Client Credentials Flow",
		Long: `Performs the OAuth 2.0 Client Credentials Flow. This command will help you to see if Ory Hydra has
been configured properly.

This command should not be used for anything else than manual testing or demo purposes. The server will terminate on error
and success.`,
		Run: func(cmd *cobra.Command, args []string) {
			/* #nosec G402 - we want to support dev environments, hence tls trickery */
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
				Transport: &transporter{
					FakeTLSTermination: flagx.MustGetBool(cmd, "fake-tls-termination"),
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: flagx.MustGetBool(cmd, "skip-tls-verify")},
					},
				},
			})

			scopes := flagx.MustGetStringSlice(cmd, "scope")
			audience := flagx.MustGetStringSlice(cmd, "audience")
			cu := cli.RemoteURI(cmd)

			clientID := flagx.MustGetString(cmd, "client-id")
			clientSecret := flagx.MustGetString(cmd, "client-secret")
			if clientID == "" || clientSecret == "" {
				fmt.Print(cmd.UsageString())
				fmt.Println("Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET.")
				return
			}

			oauthConfig := clientcredentials.Config{
				ClientID:       clientID,
				ClientSecret:   clientSecret,
				TokenURL:       urlx.AppendPaths(cu, "/oauth2/token").String(),
				Scopes:         scopes,
				EndpointParams: url.Values{"audience": {strings.Join(audience, " ")}},
			}

			t, err := oauthConfig.Token(ctx)
			cmdx.Must(err, "Could not retrieve access token because: %s", err)

			if flagx.MustGetBool(cmd, "verbose") {
				fmt.Printf("%+v\n", t)
			} else {
				fmt.Printf("%s\n", t.AccessToken)
			}
		},
	}
	cmd.Flags().StringSlice("scope", []string{}, "OAuth2 scope to request")
	cmd.Flags().BoolP("verbose", "v", false, "Toggle verbose output mode")
	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")
	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	cmd.PersistentFlags().String("endpoint", os.Getenv("HYDRA_URL"), "Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_URL")
	return cmd
}

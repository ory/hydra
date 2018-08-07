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

	"github.com/ory/go-convenience/urlx"
	"github.com/ory/hydra/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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

// tokenClientCmd represents the self command
var tokenClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate an OAuth2 token the client grant type",
	Long:  "This command uses the CLI's credentials to create an access token.",
	Run: func(cmd *cobra.Command, args []string) {
		fakeTlsTermination, _ := cmd.Flags().GetBool("fake-tls-termination")
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
			Transport: &transporter{
				FakeTLSTermination: fakeTlsTermination,
				Transport:          &http.Transport{},
			},
		})

		if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
			// fmt.Println("Warning: Skipping TLS Certificate Verification.")
			ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
				Transport: &transporter{
					FakeTLSTermination: fakeTlsTermination,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				},
			})
		}

		scopes, _ := cmd.Flags().GetStringSlice("scope")

		cu, err := url.Parse(c.GetClusterURLWithoutTailingSlashOrFail(cmd))
		pkg.Must(err, `Unable to parse cluster url ("%s"): %s`, c.GetClusterURLWithoutTailingSlashOrFail(cmd), err)

		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		if clientID == "" || clientSecret == "" {
			fmt.Print(cmd.UsageString())
			fmt.Println("Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET.")
			return
		}

		oauthConfig := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     urlx.AppendPaths(cu, "/oauth2/token").String(),
			Scopes:       scopes,
		}

		t, err := oauthConfig.Token(ctx)
		pkg.Must(err, "Could not retrieve access token because: %s", err)

		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Printf("%+v\n", t)
		} else {
			fmt.Printf("%s\n", t.AccessToken)
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenClientCmd)

	tokenClientCmd.Flags().StringSlice("scope", []string{}, "OAuth2 scope to request")
	tokenClientCmd.Flags().BoolP("verbose", "v", false, "Toggle verbose output mode")
	tokenClientCmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	tokenClientCmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")
	tokenClientCmd.PersistentFlags().String("endpoint", os.Getenv("HYDRA_URL"), "Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_URL")

}

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
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

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

		scopes, _ := cmd.Flags().GetStringSlice("scopes")

		oauthConfig := clientcredentials.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
			Scopes:       scopes,
		}

		t, err := oauthConfig.Token(ctx)
		pkg.Must(err, "Could not retrieve access token because: %s", err)
		fmt.Printf("%s\n", t.AccessToken)
	},
}

func init() {
	tokenCmd.AddCommand(tokenClientCmd)

	tokenClientCmd.Flags().StringSlice("scopes", []string{"hydra", "hydra.*"}, "User a specific set of scopes")
}

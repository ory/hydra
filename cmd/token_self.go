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

// tokenSelfCmd represents the self command
var tokenSelfCmd = &cobra.Command{
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

		oauthConfig := clientcredentials.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
			Scopes: []string{
				"hydra",
			},
		}

		t, err := oauthConfig.Token(ctx)
		pkg.Must(err, "Could not authenticate, because: %s\n", err)
		fmt.Printf("%s\n", t.AccessToken)
	},
}

func init() {
	tokenCmd.AddCommand(tokenSelfCmd)
}

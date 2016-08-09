package cmd

import (
	"fmt"

	"crypto/tls"
	"net/http"

	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// tokenSelfCmd represents the self command
var tokenSelfCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate an OAuth2 token the client grant type",
	Long:  "This command uses the CLI's credentials to create an access token.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
			fmt.Println("Warning: Skipping TLS Certificate Verification.")
			ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}})
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

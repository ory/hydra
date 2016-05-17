package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"crypto/tls"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/oauth2"
	"net/http"
	"golang.org/x/net/context"

)

// tokenSelfCmd represents the self command
var tokenSelfCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate an OAuth2 token the client grant type",
	Long: "This command uses the CLI's credentials to create an access token.",
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
				"core",
				"hydra",
			},
		}

		t, err := oauthConfig.Token(ctx)
		pkg.Must(err, "Could not authenticate, because: %s\n", err)
		fmt.Printf("%s", t.AccessToken)
	},
}

func init() {
	tokenCmd.AddCommand(tokenSelfCmd)
}

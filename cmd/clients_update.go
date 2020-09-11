package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var clientsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an OAuth 2.0 Client",
	Long: `This command updates an OAuth 2.0 Client by its ID.

Example:
  hydra clients update client-1 -n "my app" -c http://localhost/cb -g authorization_code -r code -a core,foobar

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or "--keybase" flag, for example:
  hydra clients update client-1 -n "my app" -g client_credentials -r token -a core,foobar --keybase keybase_username
`,
	Run: cmdHandler.Clients.UpdateClient,
}

func init() {
	clientsCmd.AddCommand(clientsUpdateCmd)
	clientsUpdateCmd.Flags().StringSliceP("callbacks", "c", []string{}, "REQUIRED list of allowed callback URLs")
	clientsUpdateCmd.Flags().StringSliceP("grant-types", "g", []string{"authorization_code"}, "A list of allowed grant types")
	clientsUpdateCmd.Flags().StringSliceP("response-types", "r", []string{"code"}, "A list of allowed response types")
	clientsUpdateCmd.Flags().StringSliceP("scope", "a", []string{""}, "The scope the client is allowed to request")
	clientsUpdateCmd.Flags().StringSlice("audience", []string{}, "The audience this client is allowed to request")
	clientsUpdateCmd.Flags().String("token-endpoint-auth-method", "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are \"client_secret_post\", \"client_secret_basic\", \"private_key_jwt\", and \"none\"")
	clientsUpdateCmd.Flags().String("jwks-uri", "", "Define the URL where the JSON Web Key Set should be fetched from when performing the \"private_key_jwt\" client authentication method")
	clientsUpdateCmd.Flags().String("policy-uri", "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data")
	clientsUpdateCmd.Flags().String("tos-uri", "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client")
	clientsUpdateCmd.Flags().String("client-uri", "", "A URL string of a web page providing information about the client")
	clientsUpdateCmd.Flags().String("logo-uri", "", "A URL string that references a logo for the client")
	clientsUpdateCmd.Flags().StringSlice("allowed-cors-origins", []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	clientsUpdateCmd.Flags().String("subject-type", "public", "A identifier algorithm. Valid values are \"public\" and \"pairwise\"")
	clientsUpdateCmd.Flags().String("secret", "", "Provide the client's secret")
	clientsUpdateCmd.Flags().StringP("name", "n", "", "The client's name")
	clientsUpdateCmd.Flags().StringSlice("post-logout-callbacks", []string{}, "List of allowed URLs to be redirected to after a logout")

	// encrypt client secret options
	clientsUpdateCmd.Flags().String("pgp-key", "", "Base64 encoded PGP encryption key for encrypting client secret")
	clientsUpdateCmd.Flags().String("pgp-key-url", "", "PGP encryption key URL for encrypting client secret")
	clientsUpdateCmd.Flags().String("keybase", "", "Keybase username for encrypting client secret")
}

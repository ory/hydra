package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
)

// updateCmd represents the update command
func NewClientsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an entire OAuth 2.0 Client",
		Long: `This command replaces an OAuth 2.0 Client by its ID.

Please be aware that this command replaces the entire client.
To update only the name, a full client should be provided, for example:
  hydra clients update client-1 -n "my updated app" -c http://localhost/cb -g authorization_code -r code -a core,foobar

If only the name flag (-n "my updated app") is provided, the all other fields are updated to their default values.

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or "--keybase" flag, for example:
  hydra clients update client-1 -n "my updated app" -g client_credentials -r token -a core,foobar --keybase keybase_username
`,
		Run: cli.NewHandler().Clients.UpdateClient,
	}

	cmd.Flags().StringSliceP("callbacks", "c", []string{}, "REQUIRED list of allowed callback URLs")
	cmd.Flags().StringSliceP("grant-types", "g", []string{"authorization_code"}, "A list of allowed grant types")
	cmd.Flags().StringSliceP("response-types", "r", []string{"code"}, "A list of allowed response types")
	cmd.Flags().StringSliceP("scope", "a", []string{""}, "The scope the client is allowed to request")
	cmd.Flags().StringSlice("audience", []string{}, "The audience this client is allowed to request")
	cmd.Flags().String("token-endpoint-auth-method", "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are \"client_secret_post\", \"client_secret_basic\", \"private_key_jwt\", and \"none\"")
	cmd.Flags().String("jwks-uri", "", "Define the URL where the JSON Web Key Set should be fetched from when performing the \"private_key_jwt\" client authentication method")
	cmd.Flags().String("policy-uri", "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data")
	cmd.Flags().String("tos-uri", "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client")
	cmd.Flags().String("client-uri", "", "A URL string of a web page providing information about the client")
	cmd.Flags().String("logo-uri", "", "A URL string that references a logo for the client")
	cmd.Flags().StringSlice("allowed-cors-origins", []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	cmd.Flags().String("subject-type", "public", "A identifier algorithm. Valid values are \"public\" and \"pairwise\"")
	cmd.Flags().String("secret", "", "Provide the client's secret")
	cmd.Flags().StringP("name", "n", "", "The client's name")
	cmd.Flags().StringSlice("post-logout-callbacks", []string{}, "List of allowed URLs to be redirected to after a logout")

	// back-channel logout options
	cmd.Flags().Bool("backchannel-logout-session-required", false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the backchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String("backchannel-logout-callback", "", "Client URL that will cause the client to log itself out when sent a Logout Token by Hydra.")

	// front-channel logout options
	cmd.Flags().Bool("frontchannel-logout-session-required", false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the frontchannel-logout-callback is used. If omitted, the default value is false.")
	cmd.Flags().String("frontchannel-logout-callback", "", "Client URL that will cause the client to log itself out when rendered in an iframe by Hydra.")

	// encrypt client secret options
	cmd.Flags().String("pgp-key", "", "Base64 encoded PGP encryption key for encrypting client secret")
	cmd.Flags().String("pgp-key-url", "", "PGP encryption key URL for encrypting client secret")
	cmd.Flags().String("keybase", "", "Keybase username for encrypting client secret")

	return cmd
}

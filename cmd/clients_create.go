package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var clientsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new OAuth2 client",
	Long: `This command creates a basic OAuth2 client. Always specify at least one redirect url.

Example:
  hydra clients create -n "my app" -c=[http://localhost/cb] -g [authorization_code] -r [code] -a [core,foobar]
`,
	Run: cmdHandler.Clients.CreateClient,
}

func init() {
	clientsCmd.AddCommand(clientsCreateCmd)
	clientsCreateCmd.Flags().String("id", "", "Give the client this id")
	clientsCreateCmd.Flags().StringSliceP("callbacks", "c", []string{}, "REQUIRED list of allowed callback URLs")
	clientsCreateCmd.Flags().StringSliceP("grant-types", "g", []string{"authorization_code"}, "A list of allowed grant types")
	clientsCreateCmd.Flags().StringSliceP("response-types", "r", []string{"code"}, "A list of allowed response types")
	clientsCreateCmd.Flags().StringSliceP("allowed-scopes", "a", []string{""}, "A list of allowed scopes")
	clientsCreateCmd.Flags().Bool("is-public", false, "Use this flag to create a public client")
	clientsCreateCmd.Flags().StringP("name", "n", "", "The client's name")
}

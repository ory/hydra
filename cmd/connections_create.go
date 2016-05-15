package cmd

import (
	"github.com/spf13/cobra"
)

// connectionsCreate represents the create command
var connectionsCreate = &cobra.Command{
	Use:   "create <provider> <local> <remote>",
	Short: "Associate local identites with remote ones",
	Long: `Use a user id from your database as the local argument.
The provider is the name of the SSO provider, e.g. "google", "twitter", "facebook".
The remote argument is the user's id from the SSO provider.

Example:
  create google peter@foobar.com googleid:jd92joafj`,
	Run:cmdHandler.Connections.CreateConnection,
}

func init() {
	connectionsCmd.AddCommand(connectionsCreate)
}

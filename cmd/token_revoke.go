package cmd

import (
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke <token>",
	Short: "Revoke an access or refresh token",
	Run:   cmdHandler.Revocation.RevokeToken,
}

func init() {
	tokenCmd.AddCommand(tokenRevokeCmd)

}

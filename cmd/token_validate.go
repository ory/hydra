package cmd

import (
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var tokenValidatorCmd = &cobra.Command{
	Use:   "validate <token>",
	Short: "Check if an access token is valid.",
	Run:   cmdHandler.Warden.IsAuthorized,
}

func init() {
	tokenCmd.AddCommand(tokenValidatorCmd)
	tokenValidatorCmd.Flags().StringSlice("scopes", []string{""}, "Additionally check if scope was granted")
	tokenValidatorCmd.Flags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
}

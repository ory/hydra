package cmd

import (
	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate OAuth2 tokens",
}

func init() {
	RootCmd.AddCommand(tokenCmd)
	tokenCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
	tokenCmd.PersistentFlags().Bool("fake-tls-termination", false, `fake tls termination by adding "X-Forwarded-Proto: https"" to http headers`)
}

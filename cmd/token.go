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
}

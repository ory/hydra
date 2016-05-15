package cmd

import (
	"github.com/spf13/cobra"
)

// policiesSubjectsCmd represents the subjects command
var policiesSubjectsCmd = &cobra.Command{
	Use:   "subjects",
	Short: "Manage a policy's subject matches",
}

func init() {
	policiesCmd.AddCommand(policiesSubjectsCmd)
}

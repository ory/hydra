package cmd

import (
	"github.com/spf13/cobra"
)

// policiesSubjectsAddCmd represents the add command
var policiesSubjectsAddCmd = &cobra.Command{
	Use:   "add <policy> <subject> [<subject>...]",
	Short: "Add subjects to the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies subjects hydra add my-policy john@org.com <[peter|max]>@org.com`,
	Run: cmdHandler.Policies.AddSubjectToPolicy,
}

func init() {
	policiesSubjectsCmd.AddCommand(policiesSubjectsAddCmd)
}

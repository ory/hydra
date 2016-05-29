package cmd

import (
	"github.com/spf13/cobra"
)

// policiesSubjectsDeleteCmd represents the remove command
var policiesSubjectsDeleteCmd = &cobra.Command{
	Use:   "delete <policy> <subject> [<subject>...]",
	Short: "Remove subjects from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies subjects delete my-policy john@org.com <[peter|max]>@org.com`,
	Run: cmdHandler.Policies.RemoveSubjectFromPolicy,
}

func init() {
	policiesSubjectsCmd.AddCommand(policiesSubjectsDeleteCmd)
}

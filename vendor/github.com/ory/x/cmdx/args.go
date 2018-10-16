package cmdx

import "github.com/spf13/cobra"

func MinArgs(cmd *cobra.Command, args []string, min int) {
	if len(args) < min {
		Fatalf(`%s

Expected %d command line arguments but got %d.`, cmd.UsageString(), min, len(args))
	}
}

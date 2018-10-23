package cmdx

import (
	"github.com/spf13/cobra"
)

func MinArgs(cmd *cobra.Command, args []string, min int) {
	if len(args) < min {
		Fatalf(`%s

Expected at least %d command line arguments but only got %d.`, cmd.UsageString(), min, len(args))
	}
}

func ExactArgs(cmd *cobra.Command, args []string, min int) {
	if len(args) < min {
		Fatalf(`%s

Expected exactly %d command line arguments but got %d.`, cmd.UsageString(), min, len(args))
	}
}

func RangeArgs(cmd *cobra.Command, args []string, allowed []int) {
	for _, a := range allowed {
		if len(args) == a {
			return
		}
	}
	Fatalf(`%s

Expected exact %v command line arguments but got %d.`, cmd.UsageString(), allowed, len(args))
}

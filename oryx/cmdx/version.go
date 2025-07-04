// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version returns a *cobra.Command that handles the `version` command.
func Version(gitTag, gitHash, buildTime *string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the build version, build time, and git hash",
		Run: func(cmd *cobra.Command, args []string) {
			if len(*gitTag) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine version because the build process did not properly configure it.")
			} else {
				fmt.Printf("Version:			%s\n", *gitTag)
			}

			if len(*gitHash) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine build commit because the build process did not properly configure it.")
			} else {
				fmt.Printf("Build Commit:	%s\n", *gitHash)
			}

			if len(*buildTime) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine build timestamp because the build process did not properly configure it.")
			} else {
				fmt.Printf("Build Timestamp:	%s\n", *buildTime)
			}
		},
	}
}

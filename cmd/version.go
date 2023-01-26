// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/ory/hydra/v2/driver/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display this binary's version, build time and git hash of this build",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:    %s\n", config.Version)
			fmt.Printf("Git Hash:   %s\n", config.Commit)
			fmt.Printf("Build Time: %s\n", config.Date)
		},
	}
}

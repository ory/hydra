package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time"
)

var (
	Version   = "dev-master"
	BuildTime = time.Now().String()
	GitHash   = "undefined"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this binary's version, build time and git hash of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Hash:   %s\n", GitHash)
		fmt.Printf("Build Time: %s\n", BuildTime)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

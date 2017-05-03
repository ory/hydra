//go:generate swagger generate spec
package main

import (
	"os"

	"github.com/ory/hydra/cmd"
	"github.com/pkg/profile"
)

func main() {
	if os.Getenv("PROFILING") == "cpu" {
		defer profile.Start(profile.CPUProfile).Stop()
	} else if os.Getenv("PROFILING") == "memory" {
		defer profile.Start(profile.MemProfile).Stop()
	}

	cmd.Execute()
}

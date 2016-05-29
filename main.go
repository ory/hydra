package main

import (
	"os"

	"github.com/ory-am/hydra/cmd"
	"github.com/pkg/profile"
)

func main() {
	if os.Getenv("HYDRA_PROFILING") == "1" {
		defer profile.Start().Stop()
	}
	cmd.Execute()
}

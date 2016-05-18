package main

import (
	"github.com/ory-am/hydra/cmd"
	"github.com/pkg/profile"
	"os"
)

func main() {
	if os.Getenv("HYDRA_PROFILING") == "1" {
		defer profile.Start().Stop()
	}
	cmd.Execute()
}

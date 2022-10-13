// Copyright © 2022 Ory Corp

package main

import (
	"github.com/ory/hydra/cmd"
	"github.com/ory/x/profilex"
)

func main() {
	defer profilex.Profile().Stop()

	cmd.Execute()
}

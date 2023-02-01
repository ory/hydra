// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/profilex"
)

func main() {
	defer profilex.Profile().Stop()

	cmd.Execute()
}

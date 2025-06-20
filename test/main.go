// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"

	"github.com/ory/pop/v6"

	"github.com/ory/hydra/v2/flow"
)

func main() {
	var session flow.LoginSession

	fmt.Printf("%+v", pop.NewModel(&session, context.Background()).Columns().Readable().SelectString())
}

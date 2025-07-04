// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/x/jsonnetsecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := jsonnetsecure.NewJsonnetCmd().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

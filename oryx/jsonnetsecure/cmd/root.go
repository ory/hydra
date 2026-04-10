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
	if err := jsonnetsecure.NewJsonnetCmd().ExecuteContext(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

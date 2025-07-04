// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ory/x/cmdx"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		cmdx.Fatalf("Expects exactly one input parameter")
	}
	err := filepath.Walk(args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, "vendor") {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			p, err := filepath.Abs(filepath.Join(args[1], path))
			if err != nil {
				return err
			}
			fmt.Println(p)
		}

		return nil
	})

	cmdx.Must(err, "%s", err)
}

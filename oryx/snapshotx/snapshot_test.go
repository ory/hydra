// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package snapshotx

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteMatches(t *testing.T) {
	files := map[string][]byte{}
	// Iterate over all json files
	require.NoError(t, filepath.Walk("fixtures", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[info.Name()] = f
		return nil
	}))

	for k, f := range files {
		t.Run(fmt.Sprintf("file=%s/fn", k), func(t *testing.T) {
			var tc struct {
				Content      json.RawMessage `json:"content"`
				IgnoreNested []string        `json:"ignore_nested"`
				IgnoreExact  []string        `json:"ignore_exact"`
			}
			require.NoError(t, json.Unmarshal(f, &tc))
			SnapshotT(t, tc.Content, ExceptNestedKeys(tc.IgnoreNested...), ExceptPaths(tc.IgnoreExact...))
		})
	}
}

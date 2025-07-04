// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/snapshotx"
)

func TestEmbedSources(t *testing.T) {
	t.Run("fixtures", func(t *testing.T) {
		require.NoError(t, filepath.Walk("fixture/embed", func(p string, i fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if i.IsDir() {
				return nil
			}

			t.Run("fixture="+i.Name(), func(t *testing.T) {
				t.Parallel()

				input, err := os.ReadFile(p)
				require.NoError(t, err)

				actual, err := EmbedSources(input, WithIgnoreKeys(
					"ignore_this_key",
				))
				require.NoError(t, err)

				snapshotx.SnapshotT(t, actual)
			})

			return nil
		}))
	})

	t.Run("only embeds base64", func(t *testing.T) {
		actual, err := EmbedSources([]byte(`{"key":"https://foobar.com", "bar":"base64://YXNkZg=="}`), WithOnlySchemes(
			"base64",
		))
		require.NoError(t, err)

		snapshotx.SnapshotT(t, actual)
	})

	t.Run("fails on invalid source", func(t *testing.T) {
		expected := []byte(`{"foo":"base64://invalid}`)
		actual, err := EmbedSources(expected)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(actual))
	})
}

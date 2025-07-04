// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPagination(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetErr(io.Discard)
	page, perPage, err := ParsePaginationArgs(cmd, "1", "2")
	require.NoError(t, err)
	assert.EqualValues(t, 1, page)
	assert.EqualValues(t, 2, perPage)

	_, _, err = ParsePaginationArgs(cmd, "abcd", "")
	require.Error(t, err)
}

func TestTokenPagination(t *testing.T) {
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetErr(&stderr)
	RegisterTokenPaginationFlags(cmd)
	require.NoError(t, cmd.Flags().Set(FlagPageToken, "1"))
	require.NoError(t, cmd.Flags().Set(FlagPageSize, "2"))

	page, perPage, err := ParseTokenPaginationArgs(cmd)
	require.NoError(t, err, stderr.String())
	assert.EqualValues(t, "1", page)
	assert.EqualValues(t, 2, perPage)
}

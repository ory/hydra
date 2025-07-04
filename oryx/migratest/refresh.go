// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build refresh
// +build refresh

package migratest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func WriteFixtureOnError(t *testing.T, err error, actual interface{}, location string) {
	content, err := json.MarshalIndent(actual, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(location), 0777))
	require.NoError(t, os.WriteFile(location, content, 0666))
}

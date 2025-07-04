// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migratest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ContainsExpectedIds(t *testing.T, path string, ids []string) {
	files, err := os.ReadDir(path)
	require.NoError(t, err)

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			expected := strings.TrimSuffix(filepath.Base(f.Name()), ".json")
			assert.Contains(t, ids, expected)
		}
	}
}

func CompareWithFixture(t *testing.T, actual interface{}, prefix string, id string) {
	location := filepath.Join("fixtures", prefix, id+".json")
	//#nosec G304 -- false positive
	expected, err := os.ReadFile(location)
	WriteFixtureOnError(t, err, actual, location)

	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)

	if !assert.JSONEq(t, string(expected), string(actualJSON)) {
		WriteFixtureOnError(t, nil, actual, location)
	}
}

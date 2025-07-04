// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package assertx

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/sjson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func PrettifyJSONPayload(t testing.TB, payload interface{}) string {
	t.Helper()
	o, err := json.MarshalIndent(payload, "", "  ")
	require.NoError(t, err)
	return string(o)
}

func EqualAsJSON(t testing.TB, expected, actual interface{}, args ...interface{}) {
	t.Helper()
	var eb, ab bytes.Buffer
	if len(args) == 0 {
		args = []interface{}{PrettifyJSONPayload(t, actual)}
	}

	require.NoError(t, json.NewEncoder(&eb).Encode(expected), args...)
	require.NoError(t, json.NewEncoder(&ab).Encode(actual), args...)
	assert.JSONEq(t, strings.TrimSpace(eb.String()), strings.TrimSpace(ab.String()), args...)
}

func EqualAsJSONExcept(t testing.TB, expected, actual interface{}, except []string, args ...interface{}) {
	t.Helper()
	var eb, ab bytes.Buffer
	if len(args) == 0 {
		args = []interface{}{PrettifyJSONPayload(t, actual)}
	}

	require.NoError(t, json.NewEncoder(&eb).Encode(expected), args...)
	require.NoError(t, json.NewEncoder(&ab).Encode(actual), args...)

	var err error
	ebs, abs := eb.String(), ab.String()
	for _, k := range except {
		ebs, err = sjson.Delete(ebs, k)
		require.NoError(t, err)

		abs, err = sjson.Delete(abs, k)
		require.NoError(t, err)
	}

	assert.JSONEq(t, strings.TrimSpace(ebs), strings.TrimSpace(abs), args...)
}

// Deprecated: use assert.WithinDuration instead
func TimeDifferenceLess(t testing.TB, t1, t2 time.Time, seconds int) {
	t.Helper()
	assert.WithinDuration(t, t1, t2, time.Duration(seconds)*time.Second)
}

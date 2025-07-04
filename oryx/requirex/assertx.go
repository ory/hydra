// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package requirex

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func EqualAsJSON(t *testing.T, expected, actual interface{}, args ...interface{}) {
	var eb, ab bytes.Buffer
	require.NoError(t, json.NewEncoder(&eb).Encode(expected))
	require.NoError(t, json.NewEncoder(&ab).Encode(actual))
	require.JSONEq(t, eb.String(), ab.String(), args...)
}

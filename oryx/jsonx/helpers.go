// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalJSONString(t *testing.T, i interface{}) string {
	out, err := json.Marshal(i)
	require.NoError(t, err)
	return string(out)
}

// Deprecated: this function does nothing helpful
func TestUnmarshalJSON(t *testing.T, in []byte, i interface{}) {
	require.NoError(t, json.Unmarshal(in, i))
}

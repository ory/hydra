// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringInSlice(t *testing.T) {
	for k, c := range []struct {
		needle   string
		haystack []string
		ok       bool
	}{
		{needle: "foo", haystack: []string{"foo", "bar"}, ok: true},
		{needle: "bar", haystack: []string{"foo", "bar"}, ok: true},
		{needle: "baz", haystack: []string{"foo", "bar"}, ok: false},
		{needle: "foo", haystack: []string{"bar"}, ok: false},
		{needle: "bar", haystack: []string{"bar"}, ok: true},
		{needle: "foo", haystack: []string{}, ok: false},
	} {
		assert.Equal(t, c.ok, StringInSlice(c.needle, c.haystack), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}

func TestEscapeJSONString(t *testing.T) {
	for _, str := range []string{"", "foobar", `foo"bar`, `foo\bar`, "foo\n\tbar"} {
		escaped := EscapeJSONString(str)
		var unmarshaled string
		err := json.Unmarshal([]byte(`"`+escaped+`"`), &unmarshaled)
		require.NoError(t, err, str)
		assert.Equal(t, str, unmarshaled, str)
	}
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGetJSONKeys(t *testing.T) {
	type A struct {
		B string
	}

	for _, tc := range []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name: "simple struct",
			input: struct {
				A, B string
			}{},
			expected: []string{"A", "B"},
		},
		{
			name: "struct with json tags",
			input: struct {
				A string `json:"a"`
				B string `json:"b"`
			}{},
			expected: []string{"a", "b"},
		},
		{
			name: "struct with unexported field",
			input: struct {
				A, b string
				C    string `json:"c"`
			}{},
			expected: []string{"A", "c"},
		},
		{
			name: "struct with omitempty",
			input: struct {
				A string `json:"a"`
				B string `json:"b,omitempty"`
			}{
				B: "we have to set this to a non-empty value because gjson keys collection will not work otherwise",
			},
			expected: []string{"a", "b"},
		},
		{
			name: "pointer to struct",
			input: &struct {
				A string
			}{},
			expected: []string{"A"},
		},
		{
			name: "embedded struct",
			input: struct {
				A
			}{},
			expected: []string{"B"},
		},
		{
			name: "nested structs",
			input: struct {
				A struct {
					B string `json:"b"`
				} `json:"a"`
			}{},
			expected: []string{"a.b"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, AllValidJSONKeys(tc.input))

			// collect keys with gjson, which only works reliably for non-omitempty fields
			var collectKeys func(gjson.Result) []string
			collectKeys = func(res gjson.Result) []string {
				var keys []string
				res.ForEach(func(key, value gjson.Result) bool {
					if value.IsObject() {
						childKeys := collectKeys(value)
						for _, k := range childKeys {
							keys = append(keys, key.String()+"."+k)
						}
					} else {
						keys = append(keys, key.String())
					}
					return true
				})
				return keys
			}
			assert.ElementsMatch(t, tc.expected, collectKeys(gjson.Parse(TestMarshalJSONString(t, tc.input))))
		})
	}
}

func TestResultGetValidKey(t *testing.T) {
	t.Run("case=fails on invalid key", func(t *testing.T) {
		r := ParseEnsureKeys(struct{ A string }{}, []byte("{}"))
		assert.Panics(t, func() {
			r.GetRequireValidKey(&panicFail{}, "b")
		})
	})

	t.Run("case=does not fail on valid key", func(t *testing.T) {
		r := ParseEnsureKeys(struct{ A string }{}, []byte(`{"A":"a"}`))
		var v string
		require.NotPanics(t, func() {
			v = r.GetRequireValidKey(&panicFail{}, "A").Str
		})
		assert.Equal(t, "a", v)
	})

	t.Run("case=nested key", func(t *testing.T) {
		r := ParseEnsureKeys(struct{ A struct{ B string } }{}, []byte(`{"A":{"B":"b"}}`))
		var v string
		require.NotPanics(t, func() {
			v = r.GetRequireValidKey(&panicFail{}, "A.B").Str
		})
		assert.Equal(t, "b", v)
	})
}

var _ require.TestingT = (*panicFail)(nil)

type panicFail struct{}

func (*panicFail) Errorf(string, ...interface{}) {}

func (*panicFail) FailNow() {
	panic("failing")
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package snapshotx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type (
	Opt     = func(*options)
	options struct {
		modifiers []func(t *testing.T, raw []byte) []byte
		name      string
	}
)

func ExceptPaths(keys ...string) Opt {
	return func(o *options) {
		o.modifiers = append(o.modifiers, func(t *testing.T, raw []byte) []byte {
			for _, key := range keys {
				var err error
				raw, err = sjson.DeleteBytes(raw, key)
				require.NoError(t, err)
			}
			return raw
		})
	}
}

func ExceptNestedKeys(nestedKeys ...string) Opt {
	return func(o *options) {
		o.modifiers = append(o.modifiers, func(t *testing.T, raw []byte) []byte {
			parsed := gjson.ParseBytes(raw)
			require.True(t, parsed.IsObject() || parsed.IsArray())
			return deleteMatches(t, "", parsed, nestedKeys, []string{}, raw)
		})
	}
}

func WithReplacement(str, replace string) Opt {
	return func(o *options) {
		o.modifiers = append(o.modifiers, func(t *testing.T, raw []byte) []byte {
			return bytes.ReplaceAll(raw, []byte(str), []byte(replace))
		})
	}
}

func WithName(name string) Opt {
	return func(o *options) {
		o.name = name
	}
}

func newOptions(opts ...Opt) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *options) applyModifiers(t *testing.T, compare []byte) []byte {
	for _, modifier := range o.modifiers {
		compare = modifier(t, compare)
	}
	return compare
}

var snapshot = cupaloy.New(cupaloy.SnapshotFileExtension(".json"))

func SnapshotTJSON[C ~string | ~[]byte](t *testing.T, compare C, opts ...Opt) {
	SnapshotT(t, json.RawMessage(compare), opts...)
}

func SnapshotT(t *testing.T, actual any, opts ...Opt) {
	t.Helper()
	compare, err := json.MarshalIndent(actual, "", "  ")
	require.NoErrorf(t, err, "%+v", actual)

	o := newOptions(opts...)
	compare = o.applyModifiers(t, compare)

	if o.name == "" {
		snapshot.SnapshotT(t, compare)
	} else {
		name := strings.ReplaceAll(t.Name()+"_"+o.name, "/", "-")
		require.NoError(t, snapshot.SnapshotWithName(name, compare))
	}
}

// SnapshotTExcept is deprecated in favor of SnapshotT with Opt.
//
// DEPRECATED: please use SnapshotT instead
func SnapshotTExcept(t *testing.T, actual interface{}, except []string) {
	t.Helper()
	compare, err := json.MarshalIndent(actual, "", "  ")
	require.NoError(t, err, "%+v", actual)
	for _, e := range except {
		compare, err = sjson.DeleteBytes(compare, e)
		require.NoError(t, err, "%s", e)
	}

	snapshot.SnapshotT(t, compare)
}

func deleteMatches(t *testing.T, key string, result gjson.Result, matches []string, parents []string, content []byte) []byte {
	path := parents
	if key != "" {
		path = append(parents, key)
	}

	if result.IsObject() {
		result.ForEach(func(key, value gjson.Result) bool {
			content = deleteMatches(t, key.String(), value, matches, path, content)
			return true
		})
	} else if result.IsArray() {
		var i int
		result.ForEach(func(_, value gjson.Result) bool {
			content = deleteMatches(t, fmt.Sprintf("%d", i), value, matches, path, content)
			i++
			return true
		})
	}

	if slices.Contains(matches, key) {
		content, err := sjson.DeleteBytes(content, strings.Join(path, "."))
		require.NoError(t, err)
		return content
	}

	return content
}

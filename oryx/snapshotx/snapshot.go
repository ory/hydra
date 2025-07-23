// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package snapshotx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"

	"github.com/ory/x/stringslice"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"
)

type (
	ExceptOpt interface {
		apply(t *testing.T, raw []byte) []byte
	}
	exceptPaths      []string
	exceptNestedKeys []string
	replacement      struct{ str, replacement string }
)

func (e exceptPaths) apply(t *testing.T, raw []byte) []byte {
	for _, ee := range e {
		var err error
		raw, err = sjson.DeleteBytes(raw, ee)
		require.NoError(t, err)
	}
	return raw
}

func (e exceptNestedKeys) apply(t *testing.T, raw []byte) []byte {
	parsed := gjson.ParseBytes(raw)
	require.True(t, parsed.IsObject() || parsed.IsArray())
	return deleteMatches(t, "", parsed, e, []string{}, raw)
}

func (r *replacement) apply(_ *testing.T, raw []byte) []byte {
	return bytes.ReplaceAll(raw, []byte(r.str), []byte(r.replacement))
}

func ExceptPaths(keys ...string) ExceptOpt {
	return exceptPaths(keys)
}

func ExceptNestedKeys(nestedKeys ...string) ExceptOpt {
	return exceptNestedKeys(nestedKeys)
}

func WithReplacement(str, replace string) ExceptOpt {
	return &replacement{str: str, replacement: replace}
}

func SnapshotTJSON(t *testing.T, compare []byte, except ...ExceptOpt) {
	t.Helper()
	for _, e := range except {
		compare = e.apply(t, compare)
	}

	cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(true),
		cupaloy.SnapshotFileExtension(".json"),
	).SnapshotT(t, pretty.Pretty(compare))
}

func SnapshotTJSONString(t *testing.T, str string, except ...ExceptOpt) {
	t.Helper()
	SnapshotTJSON(t, []byte(str), except...)
}

func SnapshotT(t *testing.T, actual interface{}, except ...ExceptOpt) {
	t.Helper()
	compare, err := json.MarshalIndent(actual, "", "  ")
	require.NoErrorf(t, err, "%+v", actual)
	for _, e := range except {
		compare = e.apply(t, compare)
	}

	cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(true),
		cupaloy.SnapshotFileExtension(".json"),
	).SnapshotT(t, compare)
}

// SnapshotTExcept is deprecated in favor of SnapshotT with ExceptOpt.
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

	cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(true),
		cupaloy.SnapshotFileExtension(".json"),
	).SnapshotT(t, compare)
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

	if stringslice.Has(matches, key) {
		content, err := sjson.DeleteBytes(content, strings.Join(path, "."))
		require.NoError(t, err)
		return content
	}

	return content
}

// SnapshotTExceptMatchingKeys works like SnapshotTExcept but deletes keys that match the given matches recursively.
//
// So instead of having deeply nested keys like `foo.bar.baz.0.key_to_delete` you can have `key_to_delete` and
// all occurences of `key_to_delete` will be removed.
//
// DEPRECATED: please use SnapshotT instead
func SnapshotTExceptMatchingKeys(t *testing.T, actual interface{}, matches []string) {
	t.Helper()
	compare, err := json.MarshalIndent(actual, "", "  ")
	require.NoError(t, err, "%+v", actual)

	parsed := gjson.ParseBytes(compare)
	require.True(t, parsed.IsObject() || parsed.IsArray())
	compare = deleteMatches(t, "", parsed, matches, []string{}, compare)

	cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(true),
		cupaloy.SnapshotFileExtension(".json"),
	).SnapshotT(t, compare)
}

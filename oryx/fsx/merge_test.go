// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fsx

import (
	"testing"
	"testing/fstest"

	"github.com/laher/mergefs"
	"github.com/stretchr/testify/assert"
)

var (
	a = fstest.MapFS{
		"a":     &fstest.MapFile{},
		"dir/c": &fstest.MapFile{},
	}
	b = fstest.MapFS{
		"b":     &fstest.MapFile{},
		"dir/d": &fstest.MapFile{},
	}
	x = fstest.MapFS{
		"x":     &fstest.MapFile{},
		"dir/y": &fstest.MapFile{},
	}
)

func TestMergeFS(t *testing.T) {
	assert.NoError(t, fstest.TestFS(
		Merge(a, b),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
	))

	assert.NoError(t, fstest.TestFS(
		Merge(a, b, x),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
	assert.NoError(t, fstest.TestFS(
		Merge(x, b, a),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
	assert.NoError(t, fstest.TestFS(
		Merge(Merge(a, b), x),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
	assert.NoError(t, fstest.TestFS(
		Merge(Merge(x, b), a),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
}

func TestLaherMergeFS(t *testing.T) {
	assert.Error(t, fstest.TestFS(
		mergefs.Merge(a, b),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
	))

	t.Skip("laher/mergefs does not handle recursive merges correctly")

	assert.NoError(t, fstest.TestFS(
		mergefs.Merge(mergefs.Merge(a, b), x),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
	assert.NoError(t, fstest.TestFS(
		mergefs.Merge(a, mergefs.Merge(b, x)),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
	assert.NoError(t, fstest.TestFS(
		mergefs.Merge(x, mergefs.Merge(b, a)),
		"a",
		"b",
		"dir",
		"dir/c",
		"dir/d",
		"dir/y",
		"x",
	))
}

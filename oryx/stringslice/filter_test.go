// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var filter = func(a string) func(b string) bool {
		return func(b string) bool {
			return a == b
		}
	}

	assert.EqualValues(t, []string{"bar"}, Filter([]string{"foo", "bar"}, filter("foo")))
	assert.EqualValues(t, []string{"foo"}, Filter([]string{"foo", "bar"}, filter("bar")))
	assert.EqualValues(t, []string{"foo", "bar"}, Filter([]string{"foo", "bar"}, filter("baz")))
}

func TestTrimEmptyFilter(t *testing.T) {
	assert.EqualValues(t, []string{}, TrimEmptyFilter([]string{" ", "  ", "    "}, unicode.IsSpace))
	assert.EqualValues(t, []string{"a"}, TrimEmptyFilter([]string{"a", " ", "  ", "    "}, unicode.IsSpace))
}

func TestTrimSpaceEmptyFilter(t *testing.T) {
	assert.EqualValues(t, []string{}, TrimSpaceEmptyFilter([]string{" ", "  ", "    "}))
	assert.EqualValues(t, []string{"a"}, TrimSpaceEmptyFilter([]string{"a", " ", "  ", "    "}))
}

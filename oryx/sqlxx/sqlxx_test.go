// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlxx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type st struct {
	Foo  string `db:"foo"`
	Bar  string `db:"bar,omitempty"`
	Barn string `db:"barn,omitempty"`
	Baz  string `db:"-"`
	Zab  string
}

func TestNamedUpdateArguments(t *testing.T) {
	assert.Equal(t,
		"UPDATE foo SET foo=:foo, bar=:bar",
		fmt.Sprintf("UPDATE foo SET %s", NamedUpdateArguments(new(st), "barn")),
	)
}

func TestExpectNamedInsert(t *testing.T) {
	columns, arguments := NamedInsertArguments(new(st), "barn")
	assert.Equal(t,
		"INSERT INTO foo (foo, bar) VALUES (:foo, :bar)",
		fmt.Sprintf("INSERT INTO foo (%s) VALUES (%s)", columns, arguments),
	)
}

func TestGetDBFieldNames(t *testing.T) {
	t.Run("get all db field names", func(t *testing.T) {
		fieldNames := GetDBFieldNames[st](true, nil)
		assert.ElementsMatch(t, []string{"foo", "bar", "barn"}, fieldNames)
	})

	t.Run("with exclusions", func(t *testing.T) {
		fieldNames := GetDBFieldNames[st](true, []string{"barn"})
		assert.ElementsMatch(t, []string{"foo", "bar"}, fieldNames)

		fieldNames = GetDBFieldNames[st](true, []string{"barn", "foo"})
		assert.ElementsMatch(t, []string{"bar"}, fieldNames)
	})

	t.Run("fields with - tag are excluded", func(t *testing.T) {
		fieldNames := GetDBFieldNames[st](true, nil)
		assert.NotContains(t, fieldNames, "baz")
	})

	t.Run("fields without db tag are excluded", func(t *testing.T) {
		fieldNames := GetDBFieldNames[st](true, nil)
		assert.NotContains(t, fieldNames, "zab")
	})
}

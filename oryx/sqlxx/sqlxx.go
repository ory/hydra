// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlxx

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

// GetDBFieldNames extracts all database field names from a struct based on the `db` tags using sqlx.
// Fields without a `db` tag, with a `db:"-"` tag, or listed in the `exclude` parameter are omitted.
// Returns a slice of field names as strings.
//
//	type Simple struct {
//		Foo  string `db:"foo"`
//		Bar  string `db:"bar"`
//		Baz  string `db:"baz"`
//		Baz  string `db:"-"`    // Excluded due to "-" tag
//		Qux  string             // Excluded due to missing db tag
//	}
//
//	fields := GetDBFieldNames[Simple](true, []string{"baz"})
//	// Returns: ["foo", "bar"]
func GetDBFieldNames[M any](strict bool, excludeColumns []string) []string {
	// Create a mapper that uses the "db" tag
	mapper := reflectx.NewMapper("db")

	// Get field names from the structs
	fields := mapper.TypeMap(reflectx.Deref(reflect.TypeOf((*M)(nil)))).Names

	// Extract just the field names
	fieldNames := make([]string, 0, len(fields))
	for _, f := range fields {
		if (strict && f.Field.Tag == "") || f.Path == "" || f.Name == "" || slices.Contains(excludeColumns, f.Name) {
			continue
		}
		fieldNames = append(fieldNames, f.Name)
	}

	return fieldNames
}

func keys(t any, exclude []string) []string {
	tt := reflect.TypeOf(t)
	if tt.Kind() == reflect.Pointer {
		tt = tt.Elem()
	}
	ks := make([]string, 0, tt.NumField())
	for i := range tt.NumField() {
		f := tt.Field(i)
		key, _, _ := strings.Cut(f.Tag.Get("db"), ",")
		if key != "" && key != "-" && !slices.Contains(exclude, key) {
			ks = append(ks, key)
		}
	}
	return ks
}

// NamedInsertArguments returns columns and arguments for SQL INSERT statements based on a struct's tags. Does
// not work with nested structs or maps!
//
//	type st struct {
//		Foo string `db:"foo"`
//		Bar string `db:"bar,omitempty"`
//		Baz string `db:"-"`
//		Zab string
//	}
//	columns, arguments := NamedInsertArguments(new(st))
//	query := fmt.Sprintf("INSERT INTO foo (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
func NamedInsertArguments(t any, exclude ...string) (columns string, arguments string) {
	keys := keys(t, exclude)
	return strings.Join(keys, ", "),
		":" + strings.Join(keys, ", :")
}

// NamedUpdateArguments returns columns and arguments for SQL UPDATE statements based on a struct's tags. Does
// not work with nested structs or maps!
//
//	type st struct {
//		Foo string `db:"foo"`
//		Bar string `db:"bar,omitempty"`
//		Baz string `db:"-"`
//		Zab string
//	}
//	query := fmt.Sprintf("UPDATE foo SET %s", NamedUpdateArguments(new(st)))
//	// UPDATE foo SET foo=:foo, bar=:bar
func NamedUpdateArguments(t any, exclude ...string) string {
	keys := keys(t, exclude)
	statements := make([]string, len(keys))

	for k, key := range keys {
		statements[k] = fmt.Sprintf("%s=:%s", key, key)
	}

	return strings.Join(statements, ", ")
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlxx

import (
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/pkg/errors"
)

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

func OnConflictDoNothing(dialect string, columnNoop string) string {
	if dialect == "mysql" {
		return fmt.Sprintf(" ON DUPLICATE KEY UPDATE `%s` = `%s` ", columnNoop, columnNoop)
	} else {
		return ` ON CONFLICT DO NOTHING `
	}
}

// ExtractSchemeFromDSN returns the scheme (e.g. `mysql`, `postgres`, etc) component in a DSN string,
// as well as the remaining part of the DSN after the scheme separator.
// It is an error to not have a scheme present.
// This makes sense in the context of a DSN to be able to identify which database is in use.
func ExtractSchemeFromDSN(dsn string) (string, string, error) {
	scheme, afterSchemeSeparator, schemeSeparatorFound := strings.Cut(dsn, "://")
	if !schemeSeparatorFound {
		return "", "", errors.New("invalid DSN: missing scheme separator")
	}
	if scheme == "" {
		return "", "", errors.New("invalid DSN: empty scheme")
	}

	return scheme, afterSchemeSeparator, nil
}

// ExtractDbNameFromDSN returns the database name component in a DSN string.
func ExtractDbNameFromDSN(dsn string) (string, error) {
	_, afterScheme, err := ExtractSchemeFromDSN(dsn)
	if err != nil {
		return "", err
	}

	_, afterSlash, slashFound := strings.Cut(afterScheme, "/")
	if !slashFound {
		return "", nil
	}

	dbName, _, _ := strings.Cut(afterSlash, "?")
	return dbName, nil
}

// ReplaceSchemeInDSN replaces the scheme (e.g. `mysql`, `postgres`, etc) in a DSN string with another one.
// This is necessary for example when using `cockroach` as a scheme, but using the postgres driver to connect to the database,
// and this driver only accepts `postgres` as a scheme.
func ReplaceSchemeInDSN(dsn string, newScheme string) (string, error) {
	_, afterSchemeSeparator, err := ExtractSchemeFromDSN(dsn)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return newScheme + "://" + afterSchemeSeparator, nil
}

// DSNRedacted parses a database DSN and returns a redacted form as a string.
// It replaces any password with "xxxxx" just like `url.Redacted()`.
// Only the password is redacted, not the username.
// This function is necessary because MySQL uses a DSN format not compatible with `url.Parse`.
// Additionally and as a consequence of the point above, the scheme is expected to be present and non-empty.
// This function is less strict that `url.Parse` in the case of MySQL.
// It also does not escape any characters in the username, whereas `url.String()`/`url.Redacted` does.
func DSNRedacted(dsn string) (string, error) {
	scheme, afterSchemeSeparator, err := ExtractSchemeFromDSN(dsn)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// If this is not MySQL, we simply delegate the work to `url.Parse`.
	if scheme != "mysql" {
		u, err := url.Parse(dsn)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return u.Redacted(), nil
	}

	// MySQL has a weird DSN syntax not conforming to a standard URL, of the form:
	// `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`
	// We only need to parse up to `@` in order to redact the password. The rest is left as-is.

	usernamePassword, afterUsernamePassword, usernamePasswordSeparatorFound := strings.Cut(afterSchemeSeparator, "@")
	if !usernamePasswordSeparatorFound {
		afterUsernamePassword = afterSchemeSeparator
	}

	username, password, hasPassword := strings.Cut(usernamePassword, ":")
	// We only insert a redacted password in the final result if a password was provided in the input.
	// This behavior matches the one of `url.Redacted()`.
	if hasPassword {
		password = ":xxxxx"
	}

	res := scheme + "://"
	if usernamePasswordSeparatorFound {
		res += username + password + "@"
	}
	res += afterUsernamePassword
	return res, nil
}

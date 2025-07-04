// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
)

func ParameterizedMigrationContent(params map[string]interface{}) func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error) {
	return func(mf Migration, c *pop.Connection, b []byte, usingTemplate bool) (string, error) {
		content := ""
		if usingTemplate {
			t := template.New("migration")
			t.Funcs(SQLTemplateFuncs)
			t, err := t.Parse(string(b))
			if err != nil {
				return "", errors.Wrapf(err, "could not parse template %s", mf.Path)
			}
			var bb bytes.Buffer
			err = t.Execute(&bb, struct {
				IsSQLite       bool
				IsCockroach    bool
				IsMySQL        bool
				IsMariaDB      bool
				IsPostgreSQL   bool
				DialectDetails *pop.ConnectionDetails
				Parameters     map[string]interface{}
			}{
				IsSQLite:       c.Dialect.Name() == "sqlite3",
				IsCockroach:    c.Dialect.Name() == "cockroach",
				IsMySQL:        c.Dialect.Name() == "mysql",
				IsMariaDB:      c.Dialect.Name() == "mariadb",
				IsPostgreSQL:   c.Dialect.Name() == "postgres",
				DialectDetails: c.Dialect.Details(),
				Parameters:     params,
			})
			if err != nil {
				return "", errors.Wrapf(err, "could not execute migration template %s", mf.Path)
			}
			content = bb.String()
		} else {
			content = string(b)
		}

		return content, nil
	}
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package batch

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/reflectx"

	"github.com/ory/x/dbal"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/pop/v6"

	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"

	"github.com/ory/x/sqlxx"
)

type (
	insertQueryArgs struct {
		TableName    string
		ColumnsDecl  string
		Columns      []string
		Placeholders string
	}
	quoter interface {
		Quote(key string) string
	}
	TracerConnection struct {
		Tracer     *otelx.Tracer
		Connection *pop.Connection
	}
)

func buildInsertQueryArgs[T any](ctx context.Context, dialect string, mapper *reflectx.Mapper, quoter quoter, models []*T) insertQueryArgs {
	var (
		v     T
		model = pop.NewModel(v, ctx)

		columns        []string
		quotedColumns  []string
		placeholders   []string
		placeholderRow []string
	)

	for _, col := range model.Columns().Cols {
		columns = append(columns, col.Name)
		placeholderRow = append(placeholderRow, "?")
	}

	// We sort for the sole reason that the test snapshots are deterministic.
	sort.Strings(columns)

	for _, col := range columns {
		quotedColumns = append(quotedColumns, quoter.Quote(col))
	}

	// We generate a list (for every row one) of VALUE statements here that
	// will be substituted by their column values later:
	//
	//	(?, ?, ?, ?),
	//	(?, ?, ?, ?),
	//	(?, ?, ?, ?)
	for _, m := range models {
		m := reflect.ValueOf(m)

		pl := make([]string, len(placeholderRow))
		copy(pl, placeholderRow)

		// There is a special case - when using CockroachDB we want to generate
		// UUIDs using "gen_random_uuid()" which ends up in a VALUE statement of:
		//
		//	(gen_random_uuid(), ?, ?, ?),
		for k := range placeholderRow {
			if columns[k] != "id" {
				continue
			}

			field := mapper.FieldByName(m, columns[k])
			val, ok := field.Interface().(uuid.UUID)
			if !ok {
				continue
			}

			if val == uuid.Nil && dialect == dbal.DriverCockroachDB {
				pl[k] = "gen_random_uuid()"
				break
			}
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(pl, ", ")))
	}

	return insertQueryArgs{
		TableName:    quoter.Quote(model.TableName()),
		ColumnsDecl:  strings.Join(quotedColumns, ", "),
		Columns:      columns,
		Placeholders: strings.Join(placeholders, ",\n"),
	}
}

func buildInsertQueryValues[T any](dialect string, mapper *reflectx.Mapper, columns []string, models []*T, nowFunc func() time.Time) (values []any, err error) {
	for _, m := range models {
		m := reflect.ValueOf(m)

		now := nowFunc()
		// Append model fields to args
		for _, c := range columns {
			field := mapper.FieldByName(m, c)

			switch c {
			case "created_at":
				if pop.IsZeroOfUnderlyingType(field.Interface()) {
					field.Set(reflect.ValueOf(now))
				}
			case "updated_at":
				field.Set(reflect.ValueOf(now))
			case "id":
				if value, ok := field.Interface().(uuid.UUID); ok && value != uuid.Nil {
					break // breaks switch, not for
				} else if value, ok := field.Interface().(string); ok && len(value) > 0 {
					break // breaks switch, not for
				} else if dialect == dbal.DriverCockroachDB {
					// This is a special case:
					// 1. We're using cockroach
					// 2. It's the primary key field ("ID")
					// 3. A UUID was not yet set.
					//
					// If all these conditions meet, the VALUE statement will look as such:
					//
					//	(gen_random_uuid(), ?, ?, ?, ...)
					//
					// For that reason, we do not add the ID value to the list of arguments,
					// because one of the arguments is using a built-in and thus doesn't need a value.
					continue // break switch, not for
				}

				id, err := uuid.NewV4()
				if err != nil {
					return nil, err
				}
				field.Set(reflect.ValueOf(id))
			}

			values = append(values, field.Interface())

			// Special-handling for *sqlxx.NullTime: mapper.FieldByName sets this to a zero time.Time,
			// but we want a nil pointer instead.
			if i, ok := field.Interface().(*sqlxx.NullTime); ok {
				if time.Time(*i).IsZero() {
					field.Set(reflect.Zero(field.Type()))
				}
			}
		}
	}

	return values, nil
}

type createOptions struct {
	onConflict string
}

type option func(*createOptions)

func OnConflictDoNothing() func(*createOptions) {
	return func(o *createOptions) {
		o.onConflict = "ON CONFLICT DO NOTHING"
	}
}

// CreateFromSlice is a helper around Create that accepts a slice of models
// instead of a slice of model pointers.
func CreateFromSlice[T any](ctx context.Context, p *TracerConnection, models []T, opts ...option) (err error) {
	var ptrs []*T
	for k := range models {
		ptrs = append(ptrs, &models[k])
	}
	return Create(ctx, p, ptrs, opts...)
}

// Create batch-inserts the given models into the database using a single INSERT statement.
// The models are either all created or none.
func Create[T any](ctx context.Context, p *TracerConnection, models []*T, opts ...option) (err error) {
	ctx, span := p.Tracer.Tracer().Start(ctx, "persistence.sql.batch.Create")
	defer otelx.End(span, &err)

	if len(models) == 0 {
		return nil
	}

	options := &createOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var v T
	model := pop.NewModel(v, ctx)

	conn := p.Connection
	quoter, ok := conn.Dialect.(quoter)
	if !ok {
		return errors.Errorf("store is not a quoter: %T", conn.Store)
	}

	queryArgs := buildInsertQueryArgs(ctx, conn.Dialect.Name(), conn.TX.Mapper, quoter, models)
	values, err := buildInsertQueryValues(conn.Dialect.Name(), conn.TX.Mapper, queryArgs.Columns, models, func() time.Time { return time.Now().UTC().Truncate(time.Microsecond) })
	if err != nil {
		return err
	}

	var returningClause string
	if conn.Dialect.Name() != dbal.DriverMySQL {
		// PostgreSQL, CockroachDB, SQLite support RETURNING.
		returningClause = fmt.Sprintf("RETURNING %s", model.IDField())
	}

	query := conn.Dialect.TranslateSQL(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES\n%s\n%s\n%s",
		queryArgs.TableName,
		queryArgs.ColumnsDecl,
		queryArgs.Placeholders,
		options.onConflict,
		returningClause,
	))

	rows, err := conn.TX.QueryContext(ctx, query, values...)
	if err != nil {
		return sqlcon.HandleError(err)
	}
	defer rows.Close()

	// Hydrate the models from the RETURNING clause.
	//
	// Databases not supporting RETURNING will just return 0 rows.
	count := 0
	for rows.Next() {
		if err := setModelID(rows, pop.NewModel(models[count], ctx)); err != nil {
			return err
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return sqlcon.HandleError(err)
	}

	return sqlcon.HandleError(err)
}

// setModelID was copy & pasted from pop. It basically sets
// the primary key to the given value read from the SQL row.
func setModelID(row *sql.Rows, model *pop.Model) error {
	el := reflect.ValueOf(model.Value).Elem()
	fbn := el.FieldByName("ID")
	if !fbn.IsValid() {
		return errors.New("model does not have a field named id")
	}

	pkt, err := model.PrimaryKeyType()
	if err != nil {
		return errors.WithStack(err)
	}

	switch pkt {
	case "UUID":
		var id uuid.UUID
		if err := row.Scan(&id); err != nil {
			return errors.WithStack(err)
		}
		fbn.Set(reflect.ValueOf(id))
	default:
		var id interface{}
		if err := row.Scan(&id); err != nil {
			return errors.WithStack(err)
		}
		v := reflect.ValueOf(id)
		switch fbn.Kind() {
		case reflect.Int, reflect.Int64:
			fbn.SetInt(v.Int())
		default:
			fbn.Set(reflect.ValueOf(id))
		}
	}

	return nil
}

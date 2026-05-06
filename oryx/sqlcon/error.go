// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

import (
	"database/sql"
	stderrs "errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"modernc.org/sqlite"

	"github.com/ory/herodot"
)

// ErrUniqueViolation is returned when a SQL INSERT / UPDATE command returns a conflict.
func ErrUniqueViolation() *herodot.DefaultError {
	return &herodot.DefaultError{
		CodeField:     http.StatusConflict,
		GRPCCodeField: codes.AlreadyExists,
		StatusField:   http.StatusText(http.StatusConflict),
		ErrorField:    "Unable to insert or update resource because a resource with that value exists already",
	}
}

// ErrNoRows is returned when a SQL SELECT statement returns no rows.
func ErrNoRows() *herodot.DefaultError {
	return &herodot.DefaultError{
		CodeField:     http.StatusNotFound,
		GRPCCodeField: codes.NotFound,
		StatusField:   http.StatusText(http.StatusNotFound),
		ErrorField:    "Unable to locate the resource",
	}
}

// ErrConcurrentUpdate is returned when the database is unable to serialize access due to a concurrent update.
func ErrConcurrentUpdate() *herodot.DefaultError {
	return &herodot.DefaultError{
		CodeField:     http.StatusBadRequest,
		GRPCCodeField: codes.Aborted,
		StatusField:   http.StatusText(http.StatusBadRequest),
		ErrorField:    "Unable to serialize access due to a concurrent update in another session",
	}
}

func ErrNoSuchTable() *herodot.DefaultError {
	return &herodot.DefaultError{
		CodeField:     http.StatusInternalServerError,
		GRPCCodeField: codes.Internal,
		StatusField:   http.StatusText(http.StatusInternalServerError),
		ErrorField:    "Unable to locate the table",
	}
}

func ErrNoSuchColumn() *herodot.DefaultError {
	return &herodot.DefaultError{
		CodeField:     http.StatusInternalServerError,
		GRPCCodeField: codes.Internal,
		StatusField:   http.StatusText(http.StatusInternalServerError),
		ErrorField:    "Unable to locate the column",
	}
}

func handlePostgres(err error, sqlState string) error {
	switch sqlState {
	case "23505": // "unique_violation"
		return errors.WithStack(ErrUniqueViolation().WithWrap(err))
	case "40001", // "serialization_failure" in CRDB
		"CR000": // "serialization_failure"
		return errors.WithStack(ErrConcurrentUpdate().WithWrap(err))
	case "42P01": // "no such table"
		return errors.WithStack(ErrNoSuchTable().WithWrap(err))
	case "42703": // "no such column"
		return errors.WithStack(ErrNoSuchColumn().WithWrap(err))
	}
	return errors.WithStack(err)
}

// HandleError returns the right sqlcon.Err* depending on the input error.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	type stater = interface {
		error
		SQLState() string
	}

	if stderrs.Is(err, sql.ErrNoRows) {
		return errors.WithStack(ErrNoRows())
	}
	if e, ok := stderrs.AsType[stater](err); ok {
		return errors.WithStack(handlePostgres(err, e.SQLState()))
	}
	if e, ok := stderrs.AsType[*pq.Error](err); ok {
		return errors.WithStack(handlePostgres(err, string(e.Code)))
	}
	if e, ok := stderrs.AsType[*pgconn.PgError](err); ok {
		return errors.WithStack(handlePostgres(err, e.Code))
	}
	if e, ok := stderrs.AsType[*mysql.MySQLError](err); ok {
		switch e.Number {
		case 1062: // https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry
			return errors.WithStack(ErrUniqueViolation().WithWrap(err))
		case 1146: // https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_no_such_table
			return errors.WithStack(ErrNoSuchTable().WithWrap(e))
		case 1054: // https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_bad_field_error
			return errors.WithStack(ErrNoSuchColumn().WithWrap(e))
		case 1213: // https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_lock_deadlock
			return errors.WithStack(ErrConcurrentUpdate().WithWrap(e))
		}
	}
	if e, ok := stderrs.AsType[*sqlite.Error](err); ok {
		return handleSqlite(e, err)
	}

	return errors.WithStack(err)
}

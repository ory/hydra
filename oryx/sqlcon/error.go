// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

import (
	"database/sql"
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
)

var (
	// ErrUniqueViolation is returned when^a SQL INSERT / UPDATE command returns a conflict.
	ErrUniqueViolation = &herodot.DefaultError{
		CodeField:     http.StatusConflict,
		GRPCCodeField: codes.AlreadyExists,
		StatusField:   http.StatusText(http.StatusConflict),
		ErrorField:    "Unable to insert or update resource because a resource with that value exists already",
	}
	// ErrNoRows is returned when a SQL SELECT statement returns no rows.
	ErrNoRows = &herodot.DefaultError{
		CodeField:     http.StatusNotFound,
		GRPCCodeField: codes.NotFound,
		StatusField:   http.StatusText(http.StatusNotFound),
		ErrorField:    "Unable to locate the resource",
	}
	// ErrConcurrentUpdate is returned when the database is unable to serialize access due to a concurrent update.
	ErrConcurrentUpdate = &herodot.DefaultError{
		CodeField:     http.StatusBadRequest,
		GRPCCodeField: codes.Aborted,
		StatusField:   http.StatusText(http.StatusBadRequest),
		ErrorField:    "Unable to serialize access due to a concurrent update in another session",
	}
	ErrNoSuchTable = &herodot.DefaultError{
		CodeField:     http.StatusInternalServerError,
		GRPCCodeField: codes.Internal,
		StatusField:   http.StatusText(http.StatusInternalServerError),
		ErrorField:    "Unable to locate the table",
	}
)

func handlePostgres(err error, sqlState string) error {
	switch sqlState {
	case "23505": // "unique_violation"
		return errors.WithStack(ErrUniqueViolation.WithWrap(err))
	case "40001", // "serialization_failure" in CRDB
		"CR000": // "serialization_failure"
		return errors.WithStack(ErrConcurrentUpdate.WithWrap(err))
	case "42P01": // "no such table"
		return errors.WithStack(ErrNoSuchTable.WithWrap(err))
	}
	return errors.WithStack(err)
}

type stater interface {
	SQLState() string
}

// HandleError returns the right sqlcon.Err* depending on the input error.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	var st stater
	if errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(ErrNoRows)
	} else if errors.As(err, &st) {
		return handlePostgres(err, st.SQLState())
	} else if e := new(pq.Error); errors.As(err, &e) {
		return handlePostgres(err, string(e.Code))
	} else if e := new(pgconn.PgError); errors.As(err, &e) {
		return handlePostgres(err, e.Code)
	} else if e := new(mysql.MySQLError); errors.As(err, &e) {
		switch e.Number {
		case 1062:
			return errors.WithStack(ErrUniqueViolation.WithWrap(err))
		case 1146:
			return errors.WithStack(ErrNoSuchTable.WithWrap(e))
		}
	}

	if err := handleSqlite(err); err != nil {
		return err
	}

	return errors.WithStack(err)
}

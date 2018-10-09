package sqlcon

import (
	"database/sql"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
)

var (
	ErrUniqueViolation = &herodot.DefaultError{
		CodeField:   http.StatusConflict,
		StatusField: http.StatusText(http.StatusConflict),
		ErrorField:  "Unable to insert or update resource because a resource with that value exists already",
	}
	ErrNoRows = &herodot.DefaultError{
		CodeField:   http.StatusNotFound,
		StatusField: http.StatusText(http.StatusNotFound),
		ErrorField:  "Unable to locate the resource",
	}
)

func HandleError(err error) error {
	if err == sql.ErrNoRows {
		return errors.WithStack(ErrNoRows)
	}

	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		case "unique_violation":
			return errors.Wrap(ErrUniqueViolation, err.Error())
		}
		return errors.WithStack(err)
	}

	if err, ok := err.(*mysql.MySQLError); ok {
		switch err.Number {
		case 1062:
			return errors.Wrap(ErrUniqueViolation, err.Error())
		}
		return errors.WithStack(err)
	}

	return errors.WithStack(err)
}

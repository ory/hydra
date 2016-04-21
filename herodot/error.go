package herodot

import (
	"net/http"

	"github.com/go-errors/errors"
)

type Error struct {
	*errors.Error
	Code int
}

var (
	ErrNotFound = &Error{
		Error: errors.New("Not found"),
		Code:  http.StatusNotFound,
	}
	ErrUnauthorized = &Error{
		Error: errors.New("Unauthorized"),
		Code:  http.StatusUnauthorized,
	}
	ErrBadRequest = &Error{
		Error: errors.New("Bad request"),
		Code:  http.StatusBadRequest,
	}
	ErrForbidden = &Error{
		Error: errors.New("Forbidden"),
		Code:  http.StatusForbidden,
	}
)

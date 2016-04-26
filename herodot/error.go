package herodot

import (
	"net/http"

	"github.com/go-errors/errors"
)

type Error struct {
	Err  *errors.Error
	Code int
}

func (e Error) Error() string {
	return e.Err.Error()
}

var (
	ErrNotFound = &Error{
		Err:  errors.New("Not found"),
		Code: http.StatusNotFound,
	}
	ErrUnauthorized = &Error{
		Err:  errors.New("Unauthorized"),
		Code: http.StatusUnauthorized,
	}
	ErrBadRequest = &Error{
		Err:  errors.New("Bad request"),
		Code: http.StatusBadRequest,
	}
	ErrForbidden = &Error{
		Err:  errors.New("Forbidden"),
		Code: http.StatusForbidden,
	}
)

func ToError(err error) *Error {
	if e, ok := err.(*errors.Error); ok {
		return ToError(e.Err)
	} else if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{
		Err: errors.New(err),
	}
}

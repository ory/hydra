package herodot

import "net/http"

var ErrorNotFound = DefaultError{
	StatusField: http.StatusText(http.StatusNotFound),
	ErrorField:  "The requested resource could not be found",
	CodeField:   http.StatusNotFound,
}

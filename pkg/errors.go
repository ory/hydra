package pkg

import (
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
)

var (
	ErrNotFound = errors.New("Not found")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrForbidden = errors.New("Forbidden")
)

func LogError(err *errors.Error) {

}

func ForwardToErrorHandler(w http.ResponseWriter, r *http.Request, err error, errorHandlerURL url.URL) {
	q := errorHandlerURL.Query()
	q.Set("error", err.Error())
	errorHandlerURL.RawQuery = q.Encode()

	http.Redirect(w, r, errorHandlerURL.String(), http.StatusFound)
}

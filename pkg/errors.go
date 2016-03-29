package pkg

import (
	"github.com/go-errors/errors"
	"net/http"
	"net/url"
)

func LogError(err *errors.Error) {

}

func ForwardToErrorHandler(w http.ResponseWriter, r *http.Request, err error, errorHandlerURL url.URL) {
	q := errorHandlerURL.Query()
	q.Set("error", err.Error())
	errorHandlerURL.RawQuery = q.Encode()

	http.Redirect(w,r,errorHandlerURL.String(), http.StatusFound)
}

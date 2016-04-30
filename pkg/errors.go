package pkg

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/herodot"
)

var (
	ErrNotFound     = errors.New("Not found")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrForbidden    = errors.New("Forbidden")
)

func LogError(err error) {
	if e, ok := err.(*herodot.Error); ok {
		log.WithError(e).WithField("stack", e.Err.ErrorStack()).Printf("Got error.")
	} else if e, ok := err.(*errors.Error); ok {
		log.WithError(e).WithField("stack", e.ErrorStack()).Printf("Got error.")
	} else {
		log.WithError(err).Printf("Got error.")
	}
}

func ForwardToErrorHandler(w http.ResponseWriter, r *http.Request, err error, errorHandlerURL url.URL) {
	q := errorHandlerURL.Query()
	q.Set("error", err.Error())
	errorHandlerURL.RawQuery = q.Encode()

	http.Redirect(w, r, errorHandlerURL.String(), http.StatusFound)
}

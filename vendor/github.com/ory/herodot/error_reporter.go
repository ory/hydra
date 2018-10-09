package herodot

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func defaultReporter(logger logrus.FieldLogger, args ...interface{}) func(w http.ResponseWriter, r *http.Request, code int, err error) {
	return func(w http.ResponseWriter, r *http.Request, code int, err error) {
		if logger == nil {
			logger = logrus.StandardLogger()
			logger.Warning("No logger was set in json, defaulting to standard logger.")
		}

		trace := fmt.Sprintf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
		if e, ok := err.(stackTracer); ok {
			trace = fmt.Sprintf("Stack trace: %+v", e.StackTrace())
		}

		richError := toDefaultError(err, r.Header.Get("X-Request-ID"))
		logger.
			WithField("request-id", richError.RequestID()).
			WithField("writer", "JSON").
			WithField("trace", trace).
			WithField("code", code).
			WithField("reason", richError.Reason()).
			WithField("details", richError.Details()).
			WithField("status", richError.Status()).
			WithError(err).
			Error(args...)
	}
}

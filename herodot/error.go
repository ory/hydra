package herodot

import (
	"net/http"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"reflect"
)

type Error struct {
	OriginalError error  `json:"-"`
	StatusCode    int    `json:"code"`
	Description   string `json:"description,omitempty"`
	Name          string `json:"name"`
}

func (e Error) Error() string {
	return e.OriginalError.Error()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func ToError(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	} else if e, ok := err.(causer); ok {
		// causer must be at end of logical loop or fosite error detection wont work
		return ToError(e.Cause())
	} else if rfcErr := fosite.ErrorToRFC6749Error(err); rfcErr.Name != fosite.UnknownErrorName {
		return &Error{
			OriginalError: err,
			StatusCode:    rfcErr.StatusCode,
			Description:   rfcErr.Description,
			Name:          rfcErr.Name,
		}
	}

	return &Error{
		OriginalError: err,
		Description:   fmt.Sprintf("Could not unwrap error of type %s", reflect.TypeOf(err)),
		Name:          "internal-error",
		StatusCode:    http.StatusInternalServerError,
	}
}

func LogError(err error, id string, code int) {
	logrus.WithError(err).WithField("request_id", id).WithField("status", code).Errorln("An error occurred")
	if e, ok := err.(stackTracer); ok {
		logrus.Debugf("Stack trace: %+v", e.StackTrace())
	} else if e, ok := errors.Cause(err).(stackTracer); ok {
		logrus.Debugf("Stack trace: %+v", e.StackTrace())
	} else if e, ok := err.(*Error); ok {
		LogError(e.OriginalError, id, code)
	} else {
		logrus.Debugf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
	}
}

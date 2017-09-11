package pkg

import (
	"reflect"

	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotFound = &RichError{
		Status: http.StatusNotFound,
		error:  errors.New("Not found"),
	}
)

type RichError struct {
	Status int
	error
}

func (e *RichError) StatusCode() int {
	return e.Status
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogError(err error, logger log.FieldLogger) {
	if e, ok := errors.Cause(err).(stackTracer); ok {
		log.WithError(err).Errorln("An error occurred")
		log.Debugf("Stack trace: %+v", e.StackTrace())
	} else {
		log.WithError(err).Errorln("An error occurred")
		log.Debugf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
	}
}

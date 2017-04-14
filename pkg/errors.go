package pkg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"reflect"
)

var (
	ErrNotFound = errors.New("Not found")
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogError(err error) {
	if e, ok := errors.Cause(err).(stackTracer); ok {
		log.WithError(err).Errorln("An error occurred")
		log.Debugf("Stack trace: %+v", e.StackTrace())
	} else {
		log.WithError(err).Errorln("An error occurred")
		log.Debugf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
	}
}

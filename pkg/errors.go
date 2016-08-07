package pkg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/herodot"
	perr "github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("Not found")
)

type stackTracer interface {
	StackTrace() perr.StackTrace
}

func LogError(err error) {
	if e, ok := err.(*herodot.Error); ok {
		log.WithError(err).WithField("stack", e.Err.ErrorStack()).Infoln("An error occured")
	} else if e, ok := err.(*errors.Error); ok {
		log.WithError(err).WithField("stack", e.ErrorStack()).Infoln("An error occurred")
	} else if e, ok := err.(stackTracer); ok {
		log.WithError(err).WithField("stack", e.StackTrace()).Infoln("An error occured")
	} else {
		log.WithError(err).Infoln("An error occured")
	}
}

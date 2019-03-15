package x

import (
	"github.com/ory/herodot"
	"github.com/sirupsen/logrus"
)

type RegistryLogger interface {
	Logger() logrus.FieldLogger
}

type RegistryWriter interface {
	Writer() herodot.Writer
}

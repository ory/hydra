package x

import (
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/ory/herodot"
)

type RegistryLogger interface {
	Logger() logrus.FieldLogger
}

type RegistryWriter interface {
	Writer() herodot.Writer
}

type RegistryCookieStore interface {
	CookieStore() sessions.Store
}

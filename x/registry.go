package x

import (
	"context"

	"github.com/gorilla/sessions"

	"github.com/ory/x/tracing"

	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"
)

type RegistryLogger interface {
	Logger() *logrusx.Logger
	AuditLogger() *logrusx.Logger
}

type RegistryWriter interface {
	Writer() herodot.Writer
}

type RegistryCookieStore interface {
	CookieStore() sessions.Store
}

type TracingProvider interface {
	Tracer(ctx context.Context) *tracing.Tracer
}

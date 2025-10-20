// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/herodot"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
)

type RegistryLogger interface {
	Logger() *logrusx.Logger
	AuditLogger() *logrusx.Logger
}

type RegistryWriter interface {
	Writer() herodot.Writer
}

type RegistryCookieStore interface {
	CookieStore(ctx context.Context) (sessions.Store, error)
}

type TracingProvider interface {
	Tracer(ctx context.Context) *otelx.Tracer
}

type HTTPClientProvider interface {
	HTTPClient(ctx context.Context, opts ...httpx.ResilientOptions) *retryablehttp.Client
}

type Networker interface {
	NetworkID(ctx context.Context) uuid.UUID
}

type NetworkProvider interface {
	Networker() Networker
}

type Transactor interface {
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

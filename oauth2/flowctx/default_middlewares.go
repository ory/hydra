// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

const (
	FlowCookie         = "ory_hydra_flow"
	LoginSessionCookie = "ory_hydra_loginsession"
)

type Handler func(next httprouter.Handle) httprouter.Handle

func Chain(middlewares ...*Middleware) Handler {
	return func(next httprouter.Handle) httprouter.Handle {
		for _, mw := range middlewares {
			next = mw.Handle(next)
		}
		return next
	}
}

func DefaultHandler(d Dependencies) Handler {
	return Chain(
		NewMiddleware(FlowCookie, d),
		NewMiddleware(LoginSessionCookie, d),
	)
}

// WithDefaultValues returns a context with default values for the flow and login session cookies.
func WithDefaultValues(ctx context.Context) context.Context {
	return context.WithValue(
		context.WithValue(
			ctx,
			contextKey(FlowCookie), &Value{},
		),
		contextKey(LoginSessionCookie), &Value{})
}

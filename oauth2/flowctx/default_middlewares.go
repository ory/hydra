// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx

import "github.com/julienschmidt/httprouter"

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

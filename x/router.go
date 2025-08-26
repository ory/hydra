// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"github.com/urfave/negroni"

	"github.com/ory/x/httprouterx"
	"github.com/ory/x/serverx"
)

func NewRouterPublic(metricsHandler negroni.Handler) *httprouterx.RouterPublic {
	router := httprouterx.NewRouterPublic(metricsHandler)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

func NewRouterAdmin(metricsHandler negroni.Handler) *httprouterx.RouterAdmin {
	router := httprouterx.NewRouterAdminWithPrefix(metricsHandler)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

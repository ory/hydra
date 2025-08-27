// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/serverx"
)

func NewRouterPublic(metricsHandler *prometheusx.MetricsManager) *httprouterx.RouterPublic {
	router := httprouterx.NewRouterPublic(metricsHandler)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

func NewRouterAdmin(metricsHandler *prometheusx.MetricsManager) *httprouterx.RouterAdmin {
	router := httprouterx.NewRouterAdminWithPrefix(metricsHandler)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

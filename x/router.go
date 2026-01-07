// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/serverx"
)

func NewRouterPublic(metricsManager *prometheusx.MetricsManager) *httprouterx.RouterPublic {
	router := httprouterx.NewRouterPublic(metricsManager)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

func NewRouterAdmin(metricsManager *prometheusx.MetricsManager) *httprouterx.RouterAdmin {
	router := httprouterx.NewRouterAdminWithPrefix(metricsManager)
	router.Handler("", "/", serverx.DefaultNotFoundHandler)
	return router
}

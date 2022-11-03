// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"net/url"

	"github.com/julienschmidt/httprouter"

	"github.com/ory/x/httprouterx"

	"github.com/ory/x/serverx"
)

func NewRouterPublic() *httprouterx.RouterPublic {
	router := httprouter.New()
	router.NotFound = serverx.DefaultNotFoundHandler
	return httprouterx.NewRouterPublic()
}

func NewRouterAdmin(f func(context.Context) *url.URL) *httprouterx.RouterAdmin {
	router := httprouterx.NewRouterAdminWithPrefix("/admin", f)
	router.NotFound = serverx.DefaultNotFoundHandler
	return router
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/ory/x/prometheusx"
)

const AdminPrefix = "/admin"

type (
	router struct {
		Mux            *http.ServeMux
		prefix         string
		metricsManager *prometheusx.MetricsManager
	}
	RouterAdmin  struct{ router }
	RouterPublic struct{ router }
)

var DefaultTestMetricsManager = prometheusx.NewMetricsManager("", "", "", "")

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin(metricsManager *prometheusx.MetricsManager) *RouterAdmin {
	return &RouterAdmin{router: router{
		Mux:            http.NewServeMux(),
		metricsManager: metricsManager,
	}}
}

func NewTestRouterAdmin(_ *testing.T) *RouterAdmin {
	return NewRouterAdmin(DefaultTestMetricsManager)
}

func NewTestRouterAdminWithPrefix(_ *testing.T) *RouterAdmin {
	return NewRouterAdminWithPrefix(DefaultTestMetricsManager)
}

func (r *RouterAdmin) ToPublic() *RouterPublic {
	return &RouterPublic{router: router{
		Mux:            r.Mux,
		metricsManager: r.metricsManager,
	}}
}

// NewRouterPublic returns a public router.
func NewRouterPublic(metricsManager *prometheusx.MetricsManager) *RouterPublic {
	return &RouterPublic{router: router{
		Mux:            http.NewServeMux(),
		metricsManager: metricsManager,
	}}
}

func NewTestRouterPublic(_ *testing.T) *RouterPublic {
	return NewRouterPublic(DefaultTestMetricsManager)
}

// NewRouterAdminWithPrefix creates a new router with the admin prefix.
func NewRouterAdminWithPrefix(metricsManager *prometheusx.MetricsManager) *RouterAdmin {
	r := NewRouterAdmin(metricsManager)
	r.prefix = AdminPrefix
	return r
}

func (r *router) GET(route string, handle http.HandlerFunc) {
	r.handle(http.MethodGet, route, handle)
}

func (r *router) HEAD(route string, handle http.HandlerFunc) {
	r.handle(http.MethodHead, route, handle)
}

func (r *router) POST(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPost, route, handle)
}

func (r *router) PUT(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPut, route, handle)
}

func (r *router) PATCH(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPatch, route, handle)
}

func (r *router) DELETE(route string, handle http.HandlerFunc) {
	r.handle(http.MethodDelete, route, handle)
}

func (r *router) Handler(method, route string, handler http.Handler) {
	r.handle(method, route, handler)
}

func (r *router) handle(method string, route string, handler http.Handler) {
	r.Mux.HandleFunc(method+" "+path.Join(r.prefix, route), func(w http.ResponseWriter, req *http.Request) {
		// In order the get the right metrics for the right path, `req.Pattern` must have been filled by the http router.
		// This is the case at this point, but not before e.g. when the prometheus middleware runs as a negroni middleware:
		// the http router has not run yet and `req.Pattern` is empty.
		r.metricsManager.ServeHTTP(w, req, handler.ServeHTTP)
	})
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.Mux.ServeHTTP(w, req) }

func TrimTrailingSlashNegroni(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

	next(rw, r)
}

func NoCacheNegroni(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		rw.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	}

	next(rw, r)
}

func AddAdminPrefixIfNotPresentNegroni(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !strings.HasPrefix(r.URL.Path, AdminPrefix) {
		r.URL.Path = path.Join(AdminPrefix, r.URL.Path)
	}

	next(rw, r)
}

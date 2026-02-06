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
		mux     *http.ServeMux
		prefix  string
		metrics *prometheusx.HTTPMetrics
	}
	RouterAdmin  struct{ router }
	RouterPublic struct{ router }

	Router interface {
		http.Handler
		GET(route string, handle http.HandlerFunc)
		HEAD(route string, handle http.HandlerFunc)
		POST(route string, handle http.HandlerFunc)
		PUT(route string, handle http.HandlerFunc)
		PATCH(route string, handle http.HandlerFunc)
		DELETE(route string, handle http.HandlerFunc)
		Handler(method, route string, handler http.Handler)
	}
)

func newRouter(metrics *prometheusx.HTTPMetrics) *router {
	return &router{
		mux:     http.NewServeMux(),
		metrics: metrics,
	}
}

// NewRouter creates a new general purpose router. It should only be used when neither the admin nor the public router is applicable.
func NewRouter(metrics *prometheusx.HTTPMetrics) Router { return newRouter(metrics) }

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin(metrics *prometheusx.HTTPMetrics) *RouterAdmin {
	return &RouterAdmin{router: *newRouter(metrics)}
}

func NewTestRouterAdmin(_ testing.TB) *RouterAdmin           { return NewRouterAdmin(nil) }
func NewTestRouterAdminWithPrefix(_ testing.TB) *RouterAdmin { return NewRouterAdminWithPrefix(nil) }
func NewTestRouterPublic(_ testing.TB) *RouterPublic         { return NewRouterPublic(nil) }

func (r *RouterAdmin) ToPublic() *RouterPublic {
	return &RouterPublic{router: router{
		mux:     r.mux,
		metrics: r.metrics,
		prefix:  "", // do not copy the admin prefix
	}}
}

// NewRouterPublic returns a public router.
func NewRouterPublic(metrics *prometheusx.HTTPMetrics) *RouterPublic {
	return &RouterPublic{router: *newRouter(metrics)}
}

// NewRouterAdminWithPrefix creates a new router with the admin prefix.
func NewRouterAdminWithPrefix(metricsManager *prometheusx.HTTPMetrics) *RouterAdmin {
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
	route = path.Join(r.prefix, route)
	if r.metrics != nil {
		handler = r.metrics.Instrument(handler, prometheusx.GetLabelForPattern(route))
	}
	r.mux.Handle(method+" "+route, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.mux.ServeHTTP(w, req) }

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

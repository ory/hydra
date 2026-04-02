// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/urfave/negroni"
)

const AdminPrefix = "/admin"

type (
	router struct {
		mux    *http.ServeMux
		prefix string
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
		Handle(pattern string, handler http.Handler)
		HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	}
)

func newRouter() *router {
	return &router{
		mux: http.NewServeMux(),
	}
}

// NewRouter creates a new general purpose router. It should only be used when neither the admin nor the public router is applicable.
func NewRouter() Router { return newRouter() }

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin() *RouterAdmin {
	return &RouterAdmin{router: *newRouter()}
}

func NewTestRouterAdmin(_ testing.TB) *RouterAdmin           { return NewRouterAdmin() }
func NewTestRouterAdminWithPrefix(_ testing.TB) *RouterAdmin { return NewRouterAdminWithPrefix() }
func NewTestRouterPublic(_ testing.TB) *RouterPublic         { return NewRouterPublic() }

func (r *RouterAdmin) ToPublic() *RouterPublic {
	return &RouterPublic{router: router{
		mux:    r.mux,
		prefix: "", // do not copy the admin prefix
	}}
}

// NewRouterPublic returns a public router.
func NewRouterPublic() *RouterPublic {
	return &RouterPublic{router: *newRouter()}
}

// NewRouterAdminWithPrefix creates a new router with the admin prefix.
func NewRouterAdminWithPrefix() *RouterAdmin {
	r := NewRouterAdmin()
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
	r.mux.Handle(method+" "+route, handler)
}

// Handle registers the handler for the given pattern. It does not prepend the admin prefix, so the caller is responsible for ensuring that the pattern is correct.
func (r *router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler for the given pattern. It does not prepend the admin prefix, so the caller is responsible for ensuring that the pattern is correct.
func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(pattern, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.mux.ServeHTTP(w, req) }

func TrimTrailingSlashNegroni(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Don't trim the root path to avoid redirect loops
	if r.URL.Path != "/" {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}

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

func PopulatePatternNegroni[R http.Handler](r R) negroni.Handler {
	var mux *http.ServeMux
	switch v := any(r).(type) {
	case *RouterPublic:
		mux = v.mux
	case *RouterAdmin:
		mux = v.mux
	case *router:
		mux = v.mux
	case *http.ServeMux:
		mux = v
	default:
		panic(fmt.Sprintf("unsupported router type %T", r))
	}
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		_, pattern := mux.Handler(req)
		req.Pattern = pattern
		next(rw, req)
	})
}

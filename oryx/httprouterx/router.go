// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"context"
	"net/http"
	"path"
	"strings"
)

const (
	AdminPrefix = "/admin"

	afterMatchHooks contextKey = 1
)

type (
	contextKey         int
	AfterMatchHookFunc func(*http.Request)

	router struct {
		mux    *http.ServeMux
		prefix string
	}
	RouterAdmin  struct{ router }
	RouterPublic struct{ router }

	Router interface {
		http.Handler

		Handle(pattern string, handler http.Handler)
		HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

		GET(route string, handle http.HandlerFunc)
		OPTIONS(route string, handle http.HandlerFunc)
		HEAD(route string, handle http.HandlerFunc)
		POST(route string, handle http.HandlerFunc)
		PUT(route string, handle http.HandlerFunc)
		PATCH(route string, handle http.HandlerFunc)
		DELETE(route string, handle http.HandlerFunc)
	}
)

func newRouter() *router { return &router{mux: http.NewServeMux()} }

// NewRouter creates a new general purpose router. It should only be used when neither the admin nor the public router is applicable.
func NewRouter() Router { return newRouter() }

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin() *RouterAdmin {
	return &RouterAdmin{router: *newRouter()}
}

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

func (r *router) OPTIONS(route string, handle http.HandlerFunc) {
	r.handle(http.MethodOptions, route, handle)
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

func (r *router) handle(method string, route string, handler http.Handler) {
	route = path.Join(r.prefix, route)
	r.mux.Handle(method+" "+route, afterMatchHookExecutor(handler))
}

func afterMatchHookExecutor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hooks, ok := r.Context().Value(afterMatchHooks).([]AfterMatchHookFunc); ok {
			for _, hook := range hooks {
				hook(r)
			}
		}
		h.ServeHTTP(w, r)
	})
}

// Handle registers the handler for the given pattern. It does not prepend the admin prefix, so the caller is responsible for ensuring that the pattern is correct.
func (r *router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, afterMatchHookExecutor(handler))
}

// HandleFunc registers the handler for the given pattern. It does not prepend the admin prefix, so the caller is responsible for ensuring that the pattern is correct.
func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Handle(pattern, http.HandlerFunc(handler))
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
		if r.URL.RawPath != "" {
			r.URL.RawPath = path.Join(AdminPrefix, r.URL.RawPath)
		}
	}

	next(rw, r)
}

// WithAfterMatchHook adds the passed AfterMatchHookFunc to the request's context. The hook is called after the
// router matched the request, but before it is passed to the handler. This allows middlewares to access the matched
// pattern and other information that is only available after the router has matched the request.
// Note that this hook is called multiple times per request if routers are nested, the hook should anticipate that.
func WithAfterMatchHook(req *http.Request, newHooks ...AfterMatchHookFunc) *http.Request {
	hooks, _ := req.Context().Value(afterMatchHooks).([]AfterMatchHookFunc)
	ctx := context.WithValue(req.Context(), afterMatchHooks, append(hooks, newHooks...))
	return req.WithContext(ctx)
}

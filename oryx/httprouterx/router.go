// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/urfave/negroni"
)

const AdminPrefix = "/admin"

// RouterPublic wraps httprouter.Mux
type RouterPublic struct {
	Mux *http.ServeMux
}

// NewRouterPublic returns a public router.
func NewRouterPublic() *RouterPublic {
	return &RouterPublic{Mux: http.NewServeMux()}
}

func (r *RouterPublic) GET(path string, handle http.HandlerFunc) {
	r.Handle("GET", path, handle)
}

func (r *RouterPublic) HEAD(path string, handle http.HandlerFunc) {
	r.Handle("HEAD", path, handle)
}

func (r *RouterPublic) POST(path string, handle http.HandlerFunc) {
	r.Handle("POST", path, handle)
}

func (r *RouterPublic) PUT(path string, handle http.HandlerFunc) {
	r.Handle("PUT", path, handle)
}

func (r *RouterPublic) PATCH(path string, handle http.HandlerFunc) {
	r.Handle("PATCH", path, handle)
}

func (r *RouterPublic) DELETE(path string, handle http.HandlerFunc) {
	r.Handle("DELETE", path, handle)
}

func (r *RouterPublic) Handle(method, path string, handle http.Handler) {
	r.Mux.Handle((method + " " + path), handle)
}

func (r *RouterPublic) HandleFunc(method, path string, handler http.HandlerFunc) {
	r.Mux.HandleFunc(method+" "+path, handler)
}

func (r *RouterPublic) Handler(method, path string, handler http.Handler) {
	r.Mux.Handle(method+" "+path, handler)
}

type baseURLProvider func(ctx context.Context) *url.URL

// RouterAdmin is a router able to prefix routes
type RouterAdmin struct {
	Mux            *http.ServeMux
	prefix         string
	metricsHandler negroni.Handler
}

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin(metricsHandler negroni.Handler) *RouterAdmin {
	return &RouterAdmin{
		Mux:            http.NewServeMux(),
		prefix:         AdminPrefix,
		metricsHandler: metricsHandler,
	}
}

func RouterAdminToPublic(r *RouterAdmin) *RouterPublic {
	return &RouterPublic{
		Mux: r.Mux,
	}
}

// NewRouterAdminWithPrefix creates a new router with is prefixed.
//
//	NewRouterAdminWithPrefix("/admin", func(context.Context) *url.URL { return &url.URL{/*...*/} })
func NewRouterAdminWithPrefix(prefix string) *RouterAdmin {
	if prefix != "" {
		prefix = "/" + strings.TrimPrefix(strings.TrimSuffix(prefix, "/"), "/")
	}

	return &RouterAdmin{
		Mux:    http.NewServeMux(),
		prefix: prefix,
	}
}

func (r *RouterAdmin) GET(route string, handle http.HandlerFunc) {
	r.handle(http.MethodGet, route, handle)
}

func (r *RouterAdmin) HEAD(route string, handle http.HandlerFunc) {
	r.handle(http.MethodHead, route, handle)
}

func (r *RouterAdmin) POST(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPost, route, handle)
}

func (r *RouterAdmin) PUT(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPut, route, handle)
}

func (r *RouterAdmin) PATCH(route string, handle http.HandlerFunc) {
	r.handle(http.MethodPatch, route, handle)
}

func (r *RouterAdmin) DELETE(route string, handle http.HandlerFunc) {
	r.handle(http.MethodDelete, route, handle)
}

func (r *RouterAdmin) Handle(method, route string, handle http.HandlerFunc) {
	r.handle(method, route, handle)
}

func (r *RouterAdmin) HandleFunc(method, route string, handler http.HandlerFunc) {
	r.handle(method, route, handler)
}

func (r *RouterAdmin) Handler(method, route string, handler http.Handler) {
	r.handle(method, route, handler)
}

func (router *RouterAdmin) handle(method string, route string, handler http.Handler) {
	router.Mux.HandleFunc(method+" "+path.Join(router.prefix, route), func(w http.ResponseWriter, r *http.Request) {
		// In order the get the right metrics for the right path, `r.Pattern` must have been filled by the http router.
		// This is the case at this point, but not before e.g. when the prometheus middleware runs as a negroni middleware:
		// the http router has not run yet and `r.Pattern` is empty.
		router.metricsHandler.ServeHTTP(w, r, handler.ServeHTTP)
	})
}

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

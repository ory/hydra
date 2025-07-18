// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// RouterPublic wraps httprouter.Router
type RouterPublic struct {
	*httprouter.Router
}

// NewRouterPublic returns a public router.
func NewRouterPublic() *RouterPublic {
	return &RouterPublic{
		Router: httprouter.New(),
	}
}

func (r *RouterPublic) GET(path string, handle httprouter.Handle) {
	r.Handle("GET", path, NoCacheHandle(handle))
}

func (r *RouterPublic) HEAD(path string, handle httprouter.Handle) {
	r.Handle("HEAD", path, NoCacheHandle(handle))
}

func (r *RouterPublic) POST(path string, handle httprouter.Handle) {
	r.Handle("POST", path, NoCacheHandle(handle))
}

func (r *RouterPublic) PUT(path string, handle httprouter.Handle) {
	r.Handle("PUT", path, NoCacheHandle(handle))
}

func (r *RouterPublic) PATCH(path string, handle httprouter.Handle) {
	r.Handle("PATCH", path, NoCacheHandle(handle))
}

func (r *RouterPublic) DELETE(path string, handle httprouter.Handle) {
	r.Handle("DELETE", path, NoCacheHandle(handle))
}

func (r *RouterPublic) Handle(method, path string, handle httprouter.Handle) {
	r.Router.Handle(method, path, NoCacheHandle(handle))
}

func (r *RouterPublic) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.Handler(method, path, NoCacheHandler(handler))
}

func (r *RouterPublic) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, path, NoCacheHandler(handler))
}

type baseURLProvider func(ctx context.Context) *url.URL

// RouterAdmin is a router able to prefix routes
type RouterAdmin struct {
	*httprouter.Router
	prefix          string
	baseURLProvider baseURLProvider
}

// NewRouterAdmin creates a new admin router.
func NewRouterAdmin() *RouterAdmin {
	return &RouterAdmin{
		Router: httprouter.New(),
	}
}

// NewRouterAdminWithPrefixAndRouter wraps NewRouterAdminWithPrefix and additionally sets the base router.
func NewRouterAdminWithPrefixAndRouter(root *httprouter.Router, prefix string, baseURLProvider baseURLProvider) *RouterAdmin {
	router := NewRouterAdminWithPrefix(prefix, baseURLProvider)
	router.Router = root
	return router
}

// NewRouterAdminWithPrefix creates a new router with is prefixed.
//
//	NewRouterAdminWithPrefix("/admin", func(context.Context) *url.URL { return &url.URL{/*...*/} })
func NewRouterAdminWithPrefix(prefix string, baseURLProvider baseURLProvider) *RouterAdmin {
	if prefix != "" {
		prefix = "/" + strings.TrimPrefix(strings.TrimSuffix(prefix, "/"), "/")
	}

	return &RouterAdmin{
		Router:          httprouter.New(),
		prefix:          prefix,
		baseURLProvider: baseURLProvider,
	}
}

func (r *RouterAdmin) GET(route string, handle httprouter.Handle) {
	r.handle(http.MethodGet, route, handle)
}

func (r *RouterAdmin) HEAD(route string, handle httprouter.Handle) {
	r.handle(http.MethodHead, route, handle)
}

func (r *RouterAdmin) POST(route string, handle httprouter.Handle) {
	r.handle(http.MethodPost, route, handle)
}

func (r *RouterAdmin) PUT(route string, handle httprouter.Handle) {
	r.handle(http.MethodPut, route, handle)
}

func (r *RouterAdmin) PATCH(route string, handle httprouter.Handle) {
	r.handle(http.MethodPatch, route, handle)
}

func (r *RouterAdmin) DELETE(route string, handle httprouter.Handle) {
	r.handle(http.MethodDelete, route, handle)
}

func (r *RouterAdmin) Handle(method, route string, handle httprouter.Handle) {
	r.handle(method, route, handle)
}

func (r *RouterAdmin) HandlerFunc(method, route string, handler http.HandlerFunc) {
	r.handleNative(method, route, handler)
}

func (r *RouterAdmin) Handler(method, route string, handler http.Handler) {
	r.Router.Handler(method, path.Join(r.prefix, route), NoCacheHandler(handler))
}

func (r *RouterAdmin) Lookup(method, route string) {
	r.Router.Lookup(method, path.Join(r.prefix, route))
}

func (r *RouterAdmin) handle(method string, route string, handle httprouter.Handle) {
	if len(r.prefix) == 0 {
		r.Router.Handle(method, route, NoCacheHandle(handle))
		return
	}

	r.Router.Handler(method, route, NoCacheHandler(r.handleRedirect()))
	r.Router.Handle(method, path.Join(r.prefix, route), NoCacheHandle(handle))
}

func (r *RouterAdmin) handleNative(method string, route string, handle http.Handler) {
	if len(r.prefix) == 0 {
		r.Router.Handler(method, route, NoCacheHandler(handle))
		return
	}

	r.Router.Handler(method, route, NoCacheHandler(r.handleRedirect()))
	r.Router.Handler(method, path.Join(r.prefix, route), NoCacheHandler(handle))
}

func (r *RouterAdmin) handleRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, rr *http.Request) {
		baseURL := r.baseURLProvider(rr.Context())

		dest := *rr.URL
		dest.Host = baseURL.Host
		dest.Scheme = baseURL.Scheme
		dest.Path = strings.TrimPrefix(dest.Path, r.prefix)
		dest.Path = path.Join(baseURL.Path, r.prefix, dest.Path)

		http.Redirect(w, rr, dest.String(), http.StatusTemporaryRedirect)
	}
}

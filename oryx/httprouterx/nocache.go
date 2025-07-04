// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NoCache adds `Cache-Control: private, no-cache, no-store, must-revalidate` to the response header.
func NoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
}

// NoCacheHandle wraps httprouter.Handle with `Cache-Control: private, no-cache, no-store, must-revalidate` headers.
func NoCacheHandle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		NoCache(w)
		handle(w, r, ps)
	}
}

// NoCacheHandlerFunc wraps http.HandlerFunc with `Cache-Control: private, no-cache, no-store, must-revalidate` headers.
func NoCacheHandlerFunc(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NoCache(w)
		handle(w, r)
	}
}

// NoCacheHandler wraps http.HandlerFunc with `Cache-Control: private, no-cache, no-store, must-revalidate` headers.
func NoCacheHandler(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NoCache(w)
		handle.ServeHTTP(w, r)
	})
}

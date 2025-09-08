// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package serverx

import (
	_ "embed"
	"net/http"

	"github.com/ory/herodot/httputil"
)

//go:embed 404.html
var page404HTML string

//go:embed 404.json
var page404JSON string

// DefaultNotFoundHandler is a default handler for handling 404 errors.
var DefaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var contentType, body string
	switch httputil.NegotiateContentType(r, []string{
		"text/html",
		"text/plain",
		"application/json",
	}, "text/html") {
	case "text/plain":
		contentType = "text/plain"
		body = "Error 404 - The requested route does not exist. Make sure you are using the right path, domain, and port."
	case "application/json":
		contentType = "application/json"
		body = page404JSON
	default:
		fallthrough
	case "text/html":
		contentType = "text/html"
		body = page404HTML
	}

	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(body))
})

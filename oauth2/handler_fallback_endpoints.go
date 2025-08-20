// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"html/template"
	"net/http"

	"github.com/ory/hydra/v2/driver/config"
)

func (h *Handler) fallbackHandler(title, heading string, sc int, configKey string) func(w http.ResponseWriter, r *http.Request) {
	if title == "" {
		title = "The request could not be executed because a mandatory configuration key is missing or malformed"
	}

	if heading == "" {
		heading = "The request could not be executed because a mandatory configuration key is missing or malformed"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.r.Logger().Errorf(`A request failed because configuration key "%s" is missing or malformed.`, configKey)

		t, err := template.New(configKey).Parse(`<html>
<head>
	<title>{{ .Title }}</title>
</head>
<body>
<h1>
	{{ .Heading }}
</h1>
<p>
	You are seeing this page because configuration key <code>{{ .Key }}</code> is not set.
</p>
<p>
	If you are an administrator, please read <a href="https://www.ory.sh/docs">the guide</a> to understand what you
	need to do. If you are a user, please contact the administrator.
</p>
</body>
</html>`)
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}

		w.WriteHeader(sc)
		if err := t.Execute(w, struct {
			Title   string
			Heading string
			Key     string
		}{Title: title, Heading: heading, Key: configKey}); err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	}
}

func (h *Handler) DefaultErrorHandler(w http.ResponseWriter, r *http.Request) {
	h.r.Logger().WithRequest(r).Error("A client requested the default error URL, environment variable URLS_ERROR is probably not set.")

	t, err := template.New("consent").Parse(`
<html>
<head>
	<title>An OAuth 2.0 Error Occurred</title>
</head>
<body>
<h1>
	The OAuth2 request resulted in an error.
</h1>
<ul>
	<li>Error: {{ .Name }}</li>
	<li>Description: {{ .Description }}</li>
	<li>Hint: {{ .Hint }}</li>
	<li>Debug: {{ .Debug }}</li>
</ul>
<p>
	You are seeing this page because configuration key <code>{{ .Key }}</code> is not set.
</p>
<p>
	If you are an administrator, please read <a href="https://www.ory.sh/docs">the guide</a> to understand what you
	need to do. If you are a user, please contact the administrator.
</p>
</body>
</html>
`)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	if err := t.Execute(w, struct {
		Name        string
		Description string
		Hint        string
		Debug       string
		Key         string
	}{
		Name:        r.URL.Query().Get("error"),
		Description: r.URL.Query().Get("error_description"),
		Hint:        r.URL.Query().Get("error_hint"),
		Debug:       r.URL.Query().Get("error_debug"),
		Key:         config.KeyErrorURL,
	}); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
}

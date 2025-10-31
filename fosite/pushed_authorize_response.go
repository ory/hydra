// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "net/http"

// PushedAuthorizeResponse is the response object for PAR
type PushedAuthorizeResponse struct {
	RequestURI string `json:"request_uri"`
	ExpiresIn  int    `json:"expires_in"`
	Header     http.Header
	Extra      map[string]interface{}
}

// GetRequestURI gets
func (a *PushedAuthorizeResponse) GetRequestURI() string {
	return a.RequestURI
}

// SetRequestURI sets
func (a *PushedAuthorizeResponse) SetRequestURI(requestURI string) {
	a.RequestURI = requestURI
}

// GetExpiresIn gets
func (a *PushedAuthorizeResponse) GetExpiresIn() int {
	return a.ExpiresIn
}

// SetExpiresIn sets
func (a *PushedAuthorizeResponse) SetExpiresIn(seconds int) {
	a.ExpiresIn = seconds
}

// GetHeader gets
func (a *PushedAuthorizeResponse) GetHeader() http.Header {
	return a.Header
}

// AddHeader adds
func (a *PushedAuthorizeResponse) AddHeader(key, value string) {
	a.Header.Add(key, value)
}

// SetExtra sets
func (a *PushedAuthorizeResponse) SetExtra(key string, value interface{}) {
	a.Extra[key] = value
}

// GetExtra gets
func (a *PushedAuthorizeResponse) GetExtra(key string) interface{} {
	return a.Extra[key]
}

// ToMap converts to a map
func (a *PushedAuthorizeResponse) ToMap() map[string]interface{} {
	a.Extra["request_uri"] = a.RequestURI
	a.Extra["expires_in"] = a.ExpiresIn
	return a.Extra
}

// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"net/http"
	"net/url"
)

// AuthorizeResponse is an implementation of AuthorizeResponder
type AuthorizeResponse struct {
	Header     http.Header
	Parameters url.Values
	code       string
}

func NewAuthorizeResponse() *AuthorizeResponse {
	return &AuthorizeResponse{
		Header:     http.Header{},
		Parameters: url.Values{},
	}
}

func (a *AuthorizeResponse) GetCode() string {
	return a.code
}

func (a *AuthorizeResponse) GetHeader() http.Header {
	return a.Header
}

func (a *AuthorizeResponse) AddHeader(key, value string) {
	a.Header.Add(key, value)
}

func (a *AuthorizeResponse) GetParameters() url.Values {
	return a.Parameters
}

func (a *AuthorizeResponse) AddParameter(key, value string) {
	if key == "code" {
		a.code = value
	}
	a.Parameters.Add(key, value)
}

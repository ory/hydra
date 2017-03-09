package fosite

import (
	"net/http"
	"net/url"
)

// AuthorizeResponse is an implementation of AuthorizeResponder
type AuthorizeResponse struct {
	Header   http.Header
	Query    url.Values
	Fragment url.Values
	code     string
}

func NewAuthorizeResponse() *AuthorizeResponse {
	return &AuthorizeResponse{
		Header:   http.Header{},
		Query:    url.Values{},
		Fragment: url.Values{},
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

func (a *AuthorizeResponse) GetQuery() url.Values {
	return a.Query
}

func (a *AuthorizeResponse) GetFragment() url.Values {
	return a.Fragment
}

func (a *AuthorizeResponse) AddQuery(key, value string) {
	if key == "code" {
		a.code = value
	}
	a.Query.Add(key, value)
}

func (a *AuthorizeResponse) AddFragment(key, value string) {
	if key == "code" {
		a.code = value
	}
	a.Fragment.Add(key, value)
}

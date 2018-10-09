/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

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

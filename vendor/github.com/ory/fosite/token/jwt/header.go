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

package jwt

import "github.com/dgrijalva/jwt-go"

// Headers is the jwt headers
type Headers struct {
	Extra map[string]interface{}
}

func NewHeaders() *Headers {
	return &Headers{Extra: map[string]interface{}{}}
}

// ToMap will transform the headers to a map structure
func (h *Headers) ToMap() map[string]interface{} {
	var filter = map[string]bool{"alg": true, "typ": true}
	var extra = map[string]interface{}{}

	// filter known values from extra.
	for k, v := range h.Extra {
		if _, ok := filter[k]; !ok {
			extra[k] = v
		}
	}

	return extra
}

// Add will add a key-value pair to the extra field
func (h *Headers) Add(key string, value interface{}) {
	if h.Extra == nil {
		h.Extra = make(map[string]interface{})
	}
	h.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (h *Headers) Get(key string) interface{} {
	return h.Extra[key]
}

// ToMapClaims will return a jwt-go MapClaims representaion
func (h Headers) ToMapClaims() jwt.MapClaims {
	return h.ToMap()
}

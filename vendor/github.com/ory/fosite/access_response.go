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
	"strings"
	"time"
)

func NewAccessResponse() *AccessResponse {
	return &AccessResponse{
		Extra: map[string]interface{}{},
	}
}

type AccessResponse struct {
	Extra       map[string]interface{}
	AccessToken string
	TokenType   string
}

func (a *AccessResponse) SetScopes(scopes Arguments) {
	a.SetExtra("scope", strings.Join(scopes, " "))
}

func (a *AccessResponse) SetExpiresIn(expiresIn time.Duration) {
	a.SetExtra("expires_in", int64(expiresIn/time.Second))
}

func (a *AccessResponse) SetExtra(key string, value interface{}) {
	a.Extra[key] = value
}

func (a *AccessResponse) GetExtra(key string) interface{} {
	return a.Extra[key]
}

func (a *AccessResponse) SetAccessToken(token string) {
	a.AccessToken = token
}

func (a *AccessResponse) SetTokenType(name string) {
	a.TokenType = name
}

func (a *AccessResponse) GetAccessToken() string {
	return a.AccessToken
}

func (a *AccessResponse) GetTokenType() string {
	return a.TokenType
}

func (a *AccessResponse) ToMap() map[string]interface{} {
	a.Extra["access_token"] = a.GetAccessToken()
	a.Extra["token_type"] = a.GetTokenType()
	return a.Extra
}

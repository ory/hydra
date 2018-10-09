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

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

// IDTokenClaims represent the claims used in open id connect requests
type IDTokenClaims struct {
	JTI                                 string
	Issuer                              string
	Subject                             string
	Audience                            []string
	Nonce                               string
	ExpiresAt                           time.Time
	IssuedAt                            time.Time
	RequestedAt                         time.Time
	AuthTime                            time.Time
	AccessTokenHash                     string
	AuthenticationContextClassReference string
	CodeHash                            string
	Extra                               map[string]interface{}
}

// ToMap will transform the headers to a map structure
func (c *IDTokenClaims) ToMap() map[string]interface{} {
	var ret = Copy(c.Extra)
	ret["sub"] = c.Subject
	ret["iss"] = c.Issuer
	ret["jti"] = c.JTI

	if len(c.Audience) > 0 {
		ret["aud"] = c.Audience
	} else {
		ret["aud"] = []string{}
	}

	if len(c.Nonce) >= 0 {
		ret["nonce"] = c.Nonce
	}

	if len(c.AccessTokenHash) > 0 {
		ret["at_hash"] = c.AccessTokenHash
	}

	if len(c.JTI) == 0 {
		ret["jti"] = uuid.New()
	}

	if len(c.CodeHash) > 0 {
		ret["c_hash"] = c.CodeHash
	}

	if !c.AuthTime.IsZero() {
		ret["auth_time"] = c.AuthTime.Unix()
	}

	if len(c.AuthenticationContextClassReference) > 0 {
		ret["acr"] = c.AuthenticationContextClassReference
	}

	ret["iat"] = float64(c.IssuedAt.Unix())
	ret["exp"] = float64(c.ExpiresAt.Unix())
	ret["rat"] = float64(c.RequestedAt.Unix())
	return ret

}

// Add will add a key-value pair to the extra field
func (c *IDTokenClaims) Add(key string, value interface{}) {
	if c.Extra == nil {
		c.Extra = make(map[string]interface{})
	}
	c.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (c *IDTokenClaims) Get(key string) interface{} {
	return c.ToMap()[key]
}

// ToMapClaims will return a jwt-go MapClaims representation
func (c IDTokenClaims) ToMapClaims() jwt.MapClaims {
	return c.ToMap()
}

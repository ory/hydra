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

// JWTClaims represent a token's claims.
type JWTClaims struct {
	Subject   string
	Issuer    string
	Audience  []string
	JTI       string
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	Scope     []string
	Extra     map[string]interface{}
}

// ToMap will transform the headers to a map structure
func (c *JWTClaims) ToMap() map[string]interface{} {
	var ret = Copy(c.Extra)

	ret["jti"] = c.JTI
	if c.JTI == "" {
		ret["jti"] = uuid.New()
	}

	ret["sub"] = c.Subject
	ret["iss"] = c.Issuer
	ret["aud"] = c.Audience

	if !c.IssuedAt.IsZero() {
		ret["iat"] = float64(c.IssuedAt.Unix()) // jwt-go does not support int64 as datatype
	}

	if !c.NotBefore.IsZero() {
		ret["nbf"] = float64(c.NotBefore.Unix()) // jwt-go does not support int64 as datatype
	}

	ret["exp"] = float64(c.ExpiresAt.Unix()) // jwt-go does not support int64 as datatype

	if c.Scope != nil {
		ret["scp"] = c.Scope
	}

	return ret
}

// FromMap will set the claims based on a mapping
func (c *JWTClaims) FromMap(m map[string]interface{}) {
	c.Extra = make(map[string]interface{})
	for k, v := range m {
		switch k {
		case "jti":
			if s, ok := v.(string); ok {
				c.JTI = s
			}
		case "sub":
			if s, ok := v.(string); ok {
				c.Subject = s
			}
		case "iss":
			if s, ok := v.(string); ok {
				c.Issuer = s
			}
		case "aud":
			if s, ok := v.(string); ok {
				c.Audience = []string{s}
			} else if s, ok := v.([]string); ok {
				c.Audience = s
			}
		case "iat":
			switch v.(type) {
			case float64:
				c.IssuedAt = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.IssuedAt = time.Unix(v.(int64), 0).UTC()
			}
		case "nbf":
			switch v.(type) {
			case float64:
				c.NotBefore = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.NotBefore = time.Unix(v.(int64), 0).UTC()
			}
		case "exp":
			switch v.(type) {
			case float64:
				c.ExpiresAt = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.ExpiresAt = time.Unix(v.(int64), 0).UTC()
			}
		case "scp":
			switch v.(type) {
			case []string:
				c.Scope = v.([]string)
			case []interface{}:
				c.Scope = make([]string, len(v.([]interface{})))
				for i, vi := range v.([]interface{}) {
					if s, ok := vi.(string); ok {
						c.Scope[i] = s
					}
				}
			}
		default:
			c.Extra[k] = v
		}
	}
}

// Add will add a key-value pair to the extra field
func (c *JWTClaims) Add(key string, value interface{}) {
	if c.Extra == nil {
		c.Extra = make(map[string]interface{})
	}
	c.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (c JWTClaims) Get(key string) interface{} {
	return c.ToMap()[key]
}

// ToMapClaims will return a jwt-go MapClaims representation
func (c JWTClaims) ToMapClaims() jwt.MapClaims {
	return c.ToMap()
}

// FromMapClaims will populate claims from a jwt-go MapClaims representation
func (c *JWTClaims) FromMapClaims(mc jwt.MapClaims) {
	c.FromMap(mc)
}

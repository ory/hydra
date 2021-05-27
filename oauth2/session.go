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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2

import (
	"time"

	"github.com/mohae/deepcopy"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

type Session struct {
	*openid.DefaultSession `json:"idToken"`
	Extra                  map[string]interface{} `json:"extra"`
	KID                    string
	ClientID               string
	ConsentChallenge       string
	ExcludeNotBeforeClaim  bool
}

func NewSession(subject string) *Session {
	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims:  new(jwt.IDTokenClaims),
			Headers: new(jwt.Headers),
			Subject: subject,
		},
		Extra: map[string]interface{}{},
	}
}

/**
helper function to check if a string is existent in a slice
param: s: string-slice to be searched
param: e: string to find
returns: bool: true if s contains e; false otherwise
*/
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *Session) GetJWTClaims() jwt.JWTClaimsContainer {
	//results in import cycle
	//conf := internal.NewConfigurationWithDefaults()

	//results in nil pointer exception
	//conf1 := config.Provider{}

	//customTestClaims := conf.GetTestingCustomClaims()

	//here you would need to pass the context, for which you would need to change strategy_jwt.go in fosite, so no option
	//reg := driver.New(ctx)
	//customTestClaims := reg.Config(ctx).GetTestingCustomClaims()

	var customTestClaims []string

	//a slice of claims that are reserved and should not be overridden
	var reservedClaims = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti", "client_id", "scp", "ext"}

	var allowedClaims []string

	for _, cc := range customTestClaims {
		//check if custom claim is part of reserved claims
		if !contains(reservedClaims, cc) {
			//if not so, we cann allow it, so add it to allowed claims
			allowedClaims = append(allowedClaims, cc)
		}
	}

	//our new extra map which will be added to the jwt
	var topLevelExtraWithMirrorExt = map[string]interface{}{}

	//a copy of our original extra claims, because otherwise we run into reference errors
	var ext = jwt.Copy(s.Extra)

	//setting every allowed claim top level in jwt with respective value
	for _, allowedClaim := range allowedClaims {
		topLevelExtraWithMirrorExt[allowedClaim] = ext[allowedClaim]
	}

	//for every other claim that was already reserved and for mirroring, add original extra under "ext"
	topLevelExtraWithMirrorExt["ext"] = s.Extra

	claims := &jwt.JWTClaims{
		Subject: s.Subject,
		Issuer:  s.DefaultSession.Claims.Issuer,
		//set our custom extra map as claims.Extra
		Extra:     topLevelExtraWithMirrorExt,
		ExpiresAt: s.GetExpiresAt(fosite.AccessToken),
		IssuedAt:  time.Now(),

		// No need to set the audience because that's being done by fosite automatically.
		// Audience:  s.Audience,

		// The JTI MUST NOT BE FIXED or refreshing tokens will yield the SAME token
		// JTI:       s.JTI,

		// These are set by the DefaultJWTStrategy
		// Scope:     s.Scope,

		// Setting these here will cause the token to have the same iat/nbf values always
		// IssuedAt:  s.DefaultSession.Claims.IssuedAt,
		// NotBefore: s.DefaultSession.Claims.IssuedAt,
	}
	if !s.ExcludeNotBeforeClaim {
		claims.NotBefore = claims.IssuedAt
	}

	if claims.Extra == nil {
		claims.Extra = map[string]interface{}{}
	}

	claims.Extra["client_id"] = s.ClientID
	return claims
}

func (s *Session) GetJWTHeader() *jwt.Headers {
	return &jwt.Headers{
		Extra: map[string]interface{}{"kid": s.KID},
	}
}

func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(fosite.Session)
}

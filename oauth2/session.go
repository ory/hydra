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
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/mohae/deepcopy"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"

	"github.com/ory/x/stringslice"
)

type Session struct {
	*openid.DefaultSession `json:"id_token"`
	Extra                  map[string]interface{} `json:"extra"`
	KID                    string                 `json:"kid"`
	ClientID               string                 `json:"client_id"`
	ConsentChallenge       string                 `json:"consent_challenge"`
	ExcludeNotBeforeClaim  bool                   `json:"exclude_not_before_claim"`
	AllowedTopLevelClaims  []string               `json:"allowed_top_level_claims"`
}

func NewSession(subject string) *Session {
	return NewSessionWithCustomClaims(subject, nil)
}

func NewSessionWithCustomClaims(subject string, allowedTopLevelClaims []string) *Session {
	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims:  new(jwt.IDTokenClaims),
			Headers: new(jwt.Headers),
			Subject: subject,
		},
		Extra:                 map[string]interface{}{},
		AllowedTopLevelClaims: allowedTopLevelClaims,
	}
}

func (s *Session) GetJWTClaims() jwt.JWTClaimsContainer {
	//a slice of claims that are reserved and should not be overridden
	var reservedClaims = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti", "client_id", "scp", "ext"}

	//remove any reserved claims from the custom claims
	allowedClaimsFromConfigWithoutReserved := stringslice.Filter(s.AllowedTopLevelClaims, func(s string) bool {
		return stringslice.Has(reservedClaims, s)
	})

	//our new extra map which will be added to the jwt
	var topLevelExtraWithMirrorExt = map[string]interface{}{}

	//setting every allowed claim top level in jwt with respective value
	for _, allowedClaim := range allowedClaimsFromConfigWithoutReserved {
		topLevelExtraWithMirrorExt[allowedClaim] = s.Extra[allowedClaim]
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

var keyRewrites = map[string]string{
	"Extra":                          "extra",
	"KID":                            "kid",
	"ClientID":                       "client_id",
	"ConsentChallenge":               "consent_challenge",
	"ExcludeNotBeforeClaim":          "exclude_not_before_claim",
	"AllowedTopLevelClaims":          "allowed_top_level_claims",
	"idToken.Headers.Extra":          "id_token.headers.extra",
	"idToken.ExpiresAt":              "id_token.expires_at",
	"idToken.Username":               "id_token.username",
	"idToken.Subject":                "id_token.subject",
	"idToken.Claims.JTI":             "id_token.id_token_claims.jti",
	"idToken.Claims.Issuer":          "id_token.id_token_claims.iss",
	"idToken.Claims.Subject":         "id_token.id_token_claims.sub",
	"idToken.Claims.Audience":        "id_token.id_token_claims.aud",
	"idToken.Claims.Nonce":           "id_token.id_token_claims.nonce",
	"idToken.Claims.ExpiresAt":       "id_token.id_token_claims.exp",
	"idToken.Claims.IssuedAt":        "id_token.id_token_claims.iat",
	"idToken.Claims.RequestedAt":     "id_token.id_token_claims.rat",
	"idToken.Claims.AuthTime":        "id_token.id_token_claims.auth_time",
	"idToken.Claims.AccessTokenHash": "id_token.id_token_claims.at_hash",
	"idToken.Claims.AuthenticationContextClassReference": "id_token.id_token_claims.acr",
	"idToken.Claims.AuthenticationMethodsReferences":     "id_token.id_token_claims.amr",
	"idToken.Claims.CodeHash":                            "id_token.id_token_claims.c_hash",
	"idToken.Claims.Extra":                               "id_token.id_token_claims.ext",
}

func (s *Session) UnmarshalJSON(in []byte) (err error) {
	type t Session
	interpret := in
	parsed := gjson.ParseBytes(in)

	for orig, update := range keyRewrites {
		if !parsed.Get(orig).Exists() {
			continue
		}
		interpret, err = sjson.SetRawBytes(interpret, update, []byte(parsed.Get(orig).Raw))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	for orig := range keyRewrites {
		interpret, err = sjson.DeleteBytes(interpret, orig)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if parsed.Get("idToken").Exists() {
		interpret, err = sjson.DeleteBytes(interpret, "idToken")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	var tt t
	if err := json.Unmarshal(interpret, &tt); err != nil {
		return errors.WithStack(err)
	}

	*s = Session(tt)
	return nil
}

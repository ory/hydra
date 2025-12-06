// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"bytes"
	"context"
	"slices"
	"testing"
	"time"

	jjson "github.com/go-jose/go-jose/v3/json"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/x/logrusx"
)

// swagger:ignore
type Session struct {
	*openid.DefaultSession `json:"id_token"`
	Extra                  map[string]interface{} `json:"extra"`
	KID                    string                 `json:"kid"`
	ClientID               string                 `json:"client_id"`
	ConsentChallenge       string                 `json:"consent_challenge"`
	ExcludeNotBeforeClaim  bool                   `json:"exclude_not_before_claim"`
	AllowedTopLevelClaims  []string               `json:"allowed_top_level_claims"`
	MirrorTopLevelClaims   bool                   `json:"mirror_top_level_claims"`
}

func NewTestSession(t testing.TB, subject string) *Session {
	provider := config.MustNew(t, logrusx.New("", ""))
	return NewSessionWithCustomClaims(t.Context(), provider, subject)
}

func NewSessionWithCustomClaims(ctx context.Context, p *config.DefaultProvider, subject string) *Session {
	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims:    new(jwt.IDTokenClaims),
			Headers:   new(jwt.Headers),
			Subject:   subject,
			ExpiresAt: make(map[fosite.TokenType]time.Time),
		},
		Extra:                 map[string]interface{}{},
		AllowedTopLevelClaims: p.AllowedTopLevelClaims(ctx),
		MirrorTopLevelClaims:  p.MirrorTopLevelClaims(ctx),
		ExcludeNotBeforeClaim: p.ExcludeNotBeforeClaim(ctx),
	}
}

func (s *Session) GetJWTClaims() jwt.JWTClaimsContainer {
	// remove any reserved claims from the custom claims
	allowedClaimsFromConfigWithoutReserved := slices.DeleteFunc(s.AllowedTopLevelClaims, func(s string) bool {
		switch s {
		// these claims are reserved and should not be overridden
		case "iss", "sub", "aud", "exp", "nbf", "iat", "jti", "client_id", "scp", "ext":
			return true
		}
		return false
	})

	// our new extra map which will be added to the jwt
	topLevelExtraWithMirrorExt := make(map[string]interface{}, len(allowedClaimsFromConfigWithoutReserved)+2)
	topLevelExtraWithMirrorExt["client_id"] = s.ClientID

	// setting every allowed claim top level in jwt with respective value
	for _, allowedClaim := range allowedClaimsFromConfigWithoutReserved {
		if cl, ok := s.Extra[allowedClaim]; ok {
			topLevelExtraWithMirrorExt[allowedClaim] = cl
		}
	}

	// for every other claim that was already reserved and for mirroring, add original extra under "ext"
	if s.MirrorTopLevelClaims {
		topLevelExtraWithMirrorExt["ext"] = s.Extra
	}

	claims := &jwt.JWTClaims{
		Subject: s.Subject,
		Issuer:  s.DefaultSession.Claims.Issuer,
		// set our custom extra map as claims.Extra
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

func (s *Session) UnmarshalJSON(original []byte) (err error) {
	transformed := original
	originalParsed := gjson.ParseBytes(original)

	for oldKey, newKey := range keyRewrites {
		if !originalParsed.Get(oldKey).Exists() {
			continue
		}
		transformed, err = sjson.SetRawBytes(transformed, newKey, []byte(originalParsed.Get(oldKey).Raw))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	for orig := range keyRewrites {
		transformed, err = sjson.DeleteBytes(transformed, orig)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if originalParsed.Get("idToken").Exists() {
		transformed, err = sjson.DeleteBytes(transformed, "idToken")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// https://github.com/go-jose/go-jose/issues/144
	dec := jjson.NewDecoder(bytes.NewReader(transformed))
	dec.SetNumberType(jjson.UnmarshalIntOrFloat)
	type t Session
	if err := dec.Decode((*t)(s)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetExtraClaims implements ExtraClaimsSession for Session.
// The returned value can be modified in-place.
func (s *Session) GetExtraClaims() map[string]interface{} {
	if s == nil {
		return nil
	}

	if s.Extra == nil {
		s.Extra = make(map[string]interface{})
	}

	return s.Extra
}

// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

var (
	rsaKey = gen.MustRSAKey()
	j      = &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return rsaKey, nil
			},
		},
		Config: &fosite.Config{},
	}
)

// returns a valid JWT type. The JWTClaims.ExpiresAt time is intentionally
// left empty to ensure it is pulled from the session's ExpiresAt map for
// the given fosite.TokenType.
var jwtValidCase = func(tokenType fosite.TokenType) *fosite.Request {
	r := &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     map[string]interface{}{"foo": "bar"},
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType: time.Now().UTC().Add(time.Hour),
			},
		},
	}
	r.SetRequestedScopes([]string{"email", "offline"})
	r.GrantScope("email")
	r.GrantScope("offline")
	r.SetRequestedAudience([]string{"group0"})
	r.GrantAudience("group0")
	return r
}

var jwtValidCaseWithZeroRefreshExpiry = func(tokenType fosite.TokenType) *fosite.Request {
	r := &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     map[string]interface{}{"foo": "bar"},
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType:           time.Now().UTC().Add(time.Hour),
				fosite.RefreshToken: {},
			},
		},
	}
	r.SetRequestedScopes([]string{"email", "offline"})
	r.GrantScope("email")
	r.GrantScope("offline")
	r.SetRequestedAudience([]string{"group0"})
	r.GrantAudience("group0")
	return r
}

var jwtValidCaseWithRefreshExpiry = func(tokenType fosite.TokenType) *fosite.Request {
	r := &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     map[string]interface{}{"foo": "bar"},
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType:           time.Now().UTC().Add(time.Hour),
				fosite.RefreshToken: time.Now().UTC().Add(time.Hour * 2).Round(time.Hour),
			},
		},
	}
	r.SetRequestedScopes([]string{"email", "offline"})
	r.GrantScope("email")
	r.GrantScope("offline")
	r.SetRequestedAudience([]string{"group0"})
	r.GrantAudience("group0")
	return r
}

// returns an expired JWT type. The JWTClaims.ExpiresAt time is intentionally
// left empty to ensure it is pulled from the session's ExpiresAt map for
// the given fosite.TokenType.
var jwtExpiredCase = func(tokenType fosite.TokenType) *fosite.Request {
	r := &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				ExpiresAt: time.Now().UTC().Add(-time.Minute),
				Extra:     map[string]interface{}{"foo": "bar"},
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType: time.Now().UTC().Add(-time.Hour),
			},
		},
	}
	r.SetRequestedScopes([]string{"email", "offline"})
	r.GrantScope("email")
	r.GrantScope("offline")
	r.SetRequestedAudience([]string{"group0"})
	r.GrantAudience("group0")
	return r
}

func TestAccessToken(t *testing.T) {
	for s, scopeField := range []jwt.JWTScopeFieldEnum{
		jwt.JWTScopeFieldList,
		jwt.JWTScopeFieldString,
		jwt.JWTScopeFieldBoth,
	} {
		for k, c := range []struct {
			r    *fosite.Request
			pass bool
		}{
			{
				r:    jwtValidCase(fosite.AccessToken),
				pass: true,
			},
			{
				r:    jwtExpiredCase(fosite.AccessToken),
				pass: false,
			},
			{
				r:    jwtValidCaseWithZeroRefreshExpiry(fosite.AccessToken),
				pass: true,
			},
			{
				r:    jwtValidCaseWithRefreshExpiry(fosite.AccessToken),
				pass: true,
			},
		} {
			t.Run(fmt.Sprintf("case=%d/%d", s, k), func(t *testing.T) {
				j.Config = &fosite.Config{
					JWTScopeClaimKey: scopeField,
				}
				token, signature, err := j.GenerateAccessToken(context.Background(), c.r)
				assert.NoError(t, err)

				parts := strings.Split(token, ".")
				require.Len(t, parts, 3, "%s - %v", token, parts)
				assert.Equal(t, parts[2], signature)

				rawPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
				require.NoError(t, err)
				var payload map[string]interface{}
				err = json.Unmarshal(rawPayload, &payload)
				require.NoError(t, err)
				if scopeField == jwt.JWTScopeFieldList || scopeField == jwt.JWTScopeFieldBoth {
					scope, ok := payload["scp"]
					require.True(t, ok)
					assert.Equal(t, []interface{}{"email", "offline"}, scope)
				}
				if scopeField == jwt.JWTScopeFieldString || scopeField == jwt.JWTScopeFieldBoth {
					scope, ok := payload["scope"]
					require.True(t, ok)
					assert.Equal(t, "email offline", scope)
				}

				extraClaimsSession, ok := c.r.GetSession().(fosite.ExtraClaimsSession)
				require.True(t, ok)
				claims := extraClaimsSession.GetExtraClaims()
				assert.Equal(t, "bar", claims["foo"])
				// Returned, but will be ignored by the introspect handler.
				assert.Equal(t, "peter", claims["sub"])
				assert.Equal(t, []string{"group0"}, claims["aud"])
				// Scope field is always a string.
				assert.Equal(t, "email offline", claims["scope"])

				validate := oauth2.CallSignature(token, j)
				err = j.ValidateAccessToken(context.Background(), c.r, token)
				if c.pass {
					assert.NoError(t, err)
					assert.Equal(t, signature, validate)
				} else {
					assert.Error(t, err)
				}
			})
		}
	}
}

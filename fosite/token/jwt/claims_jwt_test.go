// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ory/hydra/v2/fosite/token/jwt"
)

var jwtClaims = &JWTClaims{
	Subject:   "peter",
	IssuedAt:  time.Now().UTC().Round(time.Second),
	Issuer:    "fosite",
	NotBefore: time.Now().UTC().Round(time.Second),
	Audience:  []string{"tests"},
	ExpiresAt: time.Now().UTC().Add(time.Hour).Round(time.Second),
	JTI:       "abcdef",
	Scope:     []string{"email", "offline"},
	Extra: map[string]interface{}{
		"foo": "bar",
		"baz": "bar",
	},
	ScopeField: JWTScopeFieldList,
}

var jwtClaimsMap = map[string]interface{}{
	"sub": jwtClaims.Subject,
	"iat": jwtClaims.IssuedAt.Unix(),
	"iss": jwtClaims.Issuer,
	"nbf": jwtClaims.NotBefore.Unix(),
	"aud": jwtClaims.Audience,
	"exp": jwtClaims.ExpiresAt.Unix(),
	"jti": jwtClaims.JTI,
	"scp": []string{"email", "offline"},
	"foo": jwtClaims.Extra["foo"],
	"baz": jwtClaims.Extra["baz"],
}

func TestClaimAddGetString(t *testing.T) {
	jwtClaims.Add("foo", "bar")
	assert.Equal(t, "bar", jwtClaims.Get("foo"))
}

func TestClaimsToMapSetsID(t *testing.T) {
	assert.NotEmpty(t, (&JWTClaims{}).ToMap()["jti"])
}

func TestAssert(t *testing.T) {
	assert.Nil(t, (&JWTClaims{ExpiresAt: time.Now().UTC().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&JWTClaims{ExpiresAt: time.Now().UTC().Add(-2 * time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&JWTClaims{NotBefore: time.Now().UTC().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.Nil(t, (&JWTClaims{NotBefore: time.Now().UTC().Add(-time.Hour)}).
		ToMapClaims().Valid())
	assert.Nil(t, (&JWTClaims{ExpiresAt: time.Now().UTC().Add(time.Hour),
		NotBefore: time.Now().UTC().Add(-time.Hour)}).ToMapClaims().Valid())
}

func TestClaimsToMap(t *testing.T) {
	assert.Equal(t, jwtClaimsMap, jwtClaims.ToMap())
}

func TestClaimsFromMap(t *testing.T) {
	var claims JWTClaims
	claims.FromMap(jwtClaimsMap)
	assert.Equal(t, jwtClaims, &claims)
}

func TestScopeFieldString(t *testing.T) {
	jwtClaimsWithString := jwtClaims.WithScopeField(JWTScopeFieldString)
	// Making a copy of jwtClaimsMap.
	jwtClaimsMapWithString := jwtClaims.ToMap()
	delete(jwtClaimsMapWithString, "scp")
	jwtClaimsMapWithString["scope"] = "email offline"
	assert.Equal(t, jwtClaimsMapWithString, map[string]interface{}(jwtClaimsWithString.ToMapClaims()))
	var claims JWTClaims
	claims.FromMap(jwtClaimsMapWithString)
	assert.Equal(t, jwtClaimsWithString, &claims)
}

func TestScopeFieldBoth(t *testing.T) {
	jwtClaimsWithBoth := jwtClaims.WithScopeField(JWTScopeFieldBoth)
	// Making a copy of jwtClaimsMap
	jwtClaimsMapWithBoth := jwtClaims.ToMap()
	jwtClaimsMapWithBoth["scope"] = "email offline"
	assert.Equal(t, jwtClaimsMapWithBoth, map[string]interface{}(jwtClaimsWithBoth.ToMapClaims()))
	var claims JWTClaims
	claims.FromMap(jwtClaimsMapWithBoth)
	assert.Equal(t, jwtClaimsWithBoth, &claims)
}

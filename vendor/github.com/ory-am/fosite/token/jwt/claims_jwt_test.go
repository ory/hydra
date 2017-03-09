package jwt_test

import (
	"testing"
	"time"

	. "github.com/ory-am/fosite/token/jwt"
	"github.com/stretchr/testify/assert"
)

var jwtClaims = &JWTClaims{
	Subject:   "peter",
	IssuedAt:  time.Now().Round(time.Second),
	Issuer:    "fosite",
	NotBefore: time.Now().Round(time.Second),
	Audience:  "tests",
	ExpiresAt: time.Now().Add(time.Hour).Round(time.Second),
	JTI:       "abcdef",
	Scope:     []string{"email", "offline"},
	Extra: map[string]interface{}{
		"foo": "bar",
		"baz": "bar",
	},
}

var jwtClaimsMap = map[string]interface{}{
	"sub": jwtClaims.Subject,
	"iat": float64(jwtClaims.IssuedAt.Unix()),
	"iss": jwtClaims.Issuer,
	"nbf": float64(jwtClaims.NotBefore.Unix()),
	"aud": jwtClaims.Audience,
	"exp": float64(jwtClaims.ExpiresAt.Unix()),
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
	assert.Nil(t, (&JWTClaims{ExpiresAt: time.Now().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&JWTClaims{ExpiresAt: time.Now().Add(-2 * time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&JWTClaims{NotBefore: time.Now().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&JWTClaims{NotBefore: time.Now().Add(-time.Hour)}).
		ToMapClaims().Valid())
	assert.Nil(t, (&JWTClaims{ExpiresAt: time.Now().Add(time.Hour),
		NotBefore: time.Now().Add(-time.Hour)}).ToMapClaims().Valid())
}

func TestClaimsToMap(t *testing.T) {
	assert.Equal(t, jwtClaimsMap, jwtClaims.ToMap())
}

func TestClaimsFromMap(t *testing.T) {
	var claims JWTClaims
	claims.FromMap(jwtClaimsMap)
	assert.Equal(t, jwtClaims, &claims)
}

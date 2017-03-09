package jwt_test

import (
	"testing"
	"time"

	. "github.com/ory-am/fosite/token/jwt"
	"github.com/stretchr/testify/assert"
)

var idTokenClaims = &IDTokenClaims{
	Subject:         "peter",
	IssuedAt:        time.Now().Round(time.Second),
	Issuer:          "fosite",
	Audience:        "tests",
	ExpiresAt:       time.Now().Add(time.Hour).Round(time.Second),
	AuthTime:        time.Now(),
	AccessTokenHash: "foobar",
	CodeHash:        "barfoo",
	Extra: map[string]interface{}{
		"foo": "bar",
		"baz": "bar",
	},
}

func TestIDTokenAssert(t *testing.T) {
	assert.Nil(t, (&IDTokenClaims{ExpiresAt: time.Now().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.NotNil(t, (&IDTokenClaims{ExpiresAt: time.Now().Add(-time.Hour)}).
		ToMapClaims().Valid())
}

func TestIDTokenClaimsToMap(t *testing.T) {
	assert.Equal(t, map[string]interface{}{
		"sub":       idTokenClaims.Subject,
		"iat":       float64(idTokenClaims.IssuedAt.Unix()),
		"iss":       idTokenClaims.Issuer,
		"aud":       idTokenClaims.Audience,
		"nonce":     idTokenClaims.Nonce,
		"exp":       float64(idTokenClaims.ExpiresAt.Unix()),
		"foo":       idTokenClaims.Extra["foo"],
		"baz":       idTokenClaims.Extra["baz"],
		"at_hash":   idTokenClaims.AccessTokenHash,
		"c_hash":    idTokenClaims.CodeHash,
		"auth_time": idTokenClaims.AuthTime.Unix(),
	}, idTokenClaims.ToMap())
}

// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestIDTokenAssert(t *testing.T) {
	assert.NoError(t, (&IDTokenClaims{ExpiresAt: time.Now().UTC().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.Error(t, (&IDTokenClaims{ExpiresAt: time.Now().UTC().Add(-time.Hour)}).
		ToMapClaims().Valid())

	assert.NotEmpty(t, (new(IDTokenClaims)).ToMapClaims()["jti"])
}

func TestIDTokenClaimsToMap(t *testing.T) {
	idTokenClaims := &IDTokenClaims{
		JTI:                                 "foo-id",
		Subject:                             "peter",
		IssuedAt:                            time.Now().UTC().Round(time.Second),
		Issuer:                              "fosite",
		Audience:                            []string{"tests"},
		ExpiresAt:                           time.Now().UTC().Add(time.Hour).Round(time.Second),
		AuthTime:                            time.Now().UTC(),
		RequestedAt:                         time.Now().UTC(),
		AccessTokenHash:                     "foobar",
		CodeHash:                            "barfoo",
		AuthenticationContextClassReference: "acr",
		AuthenticationMethodsReferences:     []string{"amr"},
		Extra: map[string]interface{}{
			"foo": "bar",
			"baz": "bar",
		},
	}
	assert.Equal(t, map[string]interface{}{
		"jti":       idTokenClaims.JTI,
		"sub":       idTokenClaims.Subject,
		"iat":       idTokenClaims.IssuedAt.Unix(),
		"rat":       idTokenClaims.RequestedAt.Unix(),
		"iss":       idTokenClaims.Issuer,
		"aud":       idTokenClaims.Audience,
		"exp":       idTokenClaims.ExpiresAt.Unix(),
		"foo":       idTokenClaims.Extra["foo"],
		"baz":       idTokenClaims.Extra["baz"],
		"at_hash":   idTokenClaims.AccessTokenHash,
		"c_hash":    idTokenClaims.CodeHash,
		"auth_time": idTokenClaims.AuthTime.Unix(),
		"acr":       idTokenClaims.AuthenticationContextClassReference,
		"amr":       idTokenClaims.AuthenticationMethodsReferences,
	}, idTokenClaims.ToMap())

	idTokenClaims.Nonce = "foobar"
	assert.Equal(t, map[string]interface{}{
		"jti":       idTokenClaims.JTI,
		"sub":       idTokenClaims.Subject,
		"iat":       idTokenClaims.IssuedAt.Unix(),
		"rat":       idTokenClaims.RequestedAt.Unix(),
		"iss":       idTokenClaims.Issuer,
		"aud":       idTokenClaims.Audience,
		"exp":       idTokenClaims.ExpiresAt.Unix(),
		"foo":       idTokenClaims.Extra["foo"],
		"baz":       idTokenClaims.Extra["baz"],
		"at_hash":   idTokenClaims.AccessTokenHash,
		"c_hash":    idTokenClaims.CodeHash,
		"auth_time": idTokenClaims.AuthTime.Unix(),
		"acr":       idTokenClaims.AuthenticationContextClassReference,
		"amr":       idTokenClaims.AuthenticationMethodsReferences,
		"nonce":     idTokenClaims.Nonce,
	}, idTokenClaims.ToMap())
}

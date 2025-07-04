// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwtx

import (
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/mapx"
)

// Claims represents a JSON Web Token's standard claims.
type Claims struct {
	// Audience identifies the recipients that the JWT is intended for.
	Audience []string `json:"aud"`

	// Issuer identifies the principal that issued the JWT.
	Issuer string `json:"iss"`

	// Subject identifies the principal that is the subject of the JWT.
	Subject string `json:"sub"`

	// ExpiresAt identifies the expiration time on or after which the JWT most not be accepted for processing.
	ExpiresAt time.Time `json:"exp"`

	// IssuedAt identifies the time at which the JWT was issued.
	IssuedAt time.Time `json:"iat"`

	// NotBefore identifies the time before which the JWT must not be accepted for processing.
	NotBefore time.Time `json:"nbf"`

	// JTI provides a unique identifier for the JWT.
	JTI string `json:"jti"`
}

// ParseMapStringInterfaceClaims converts map[string]interface{} to *Claims.
func ParseMapStringInterfaceClaims(claims map[string]interface{}) *Claims {
	c := make(map[interface{}]interface{})
	for k, v := range claims {
		c[k] = v
	}
	return ParseMapInterfaceInterfaceClaims(c)
}

// ParseMapInterfaceInterfaceClaims converts map[interface{}]interface{} to *Claims.
func ParseMapInterfaceInterfaceClaims(claims map[interface{}]interface{}) *Claims {
	result := &Claims{
		Issuer:  mapx.GetStringDefault(claims, "iss", ""),
		Subject: mapx.GetStringDefault(claims, "sub", ""),
		JTI:     mapx.GetStringDefault(claims, "jti", ""),
	}

	if aud, err := mapx.GetString(claims, "aud"); err == nil {
		result.Audience = []string{aud}
	} else if errors.Is(err, mapx.ErrKeyCanNotBeTypeAsserted) {
		if aud, err := mapx.GetStringSlice(claims, "aud"); err == nil {
			result.Audience = aud
		} else {
			result.Audience = []string{}
		}
	} else {
		result.Audience = []string{}
	}

	if exp, err := mapx.GetTime(claims, "exp"); err == nil {
		result.ExpiresAt = exp
	}

	if iat, err := mapx.GetTime(claims, "iat"); err == nil {
		result.IssuedAt = iat
	}

	if nbf, err := mapx.GetTime(claims, "nbf"); err == nil {
		result.NotBefore = nbf
	}

	return result
}

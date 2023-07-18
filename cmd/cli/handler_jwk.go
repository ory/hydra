// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	jose "github.com/go-jose/go-jose/v3"
)

func ToSDKFriendlyJSONWebKey(key interface{}, kid, use string) jose.JSONWebKey {
	var alg string

	if jwk, ok := key.(*jose.JSONWebKey); ok {
		key = jwk.Key
		if jwk.KeyID != "" {
			kid = jwk.KeyID
		}
		if jwk.Use != "" {
			use = jwk.Use
		}
		if jwk.Algorithm != "" {
			alg = jwk.Algorithm
		}
	}

	return jose.JSONWebKey{
		KeyID:     kid,
		Use:       use,
		Algorithm: alg,
		Key:       key,
	}
}

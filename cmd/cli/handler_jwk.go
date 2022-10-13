// Copyright © 2022 Ory Corp

package cli

import (
	jose "gopkg.in/square/go-jose.v2"
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

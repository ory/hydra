// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/json"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"

	hydra "github.com/ory/hydra-client-go/v2"
)

// OnlyPublicSDKKeys strips the private parts from a key set so that it is safe
// to print.
func OnlyPublicSDKKeys(in []hydra.JsonWebKey) (out []hydra.JsonWebKey, _ error) {
	var interim []jose.JSONWebKey
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(&in); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	if err := json.NewDecoder(&b).Decode(&interim); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON Web Key Set")
	}

	for i, key := range interim {
		interim[i] = key.Public()
	}

	b.Reset()
	if err := json.NewEncoder(&b).Encode(&interim); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	var keys []hydra.JsonWebKey
	if err := json.NewDecoder(&b).Decode(&keys); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON Web Key Set")
	}

	return keys, nil
}

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

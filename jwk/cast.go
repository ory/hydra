// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"crypto/rsa"

	"github.com/ory/x/josex"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
)

func MustRSAPublic(key *jose.JSONWebKey) *rsa.PublicKey {
	res, err := ToRSAPublic(key)
	if err != nil {
		panic(err.Error())
	}

	return res
}

func ToRSAPublic(key *jose.JSONWebKey) (*rsa.PublicKey, error) {
	pk := josex.ToPublicKey(key)
	res, ok := pk.Key.(*rsa.PublicKey)
	if !ok {
		return res, errors.Errorf("Could not convert key to RSA Public Key, got: %T", pk.Key)
	}

	return res, nil
}

func MustRSAPrivate(key *jose.JSONWebKey) *rsa.PrivateKey {
	res, err := ToRSAPrivate(key)
	if err != nil {
		panic(err.Error())
	}

	return res
}

func ToRSAPrivate(key *jose.JSONWebKey) (*rsa.PrivateKey, error) {
	res, ok := key.Key.(*rsa.PrivateKey)
	if !ok {
		return res, errors.New("Could not convert key to RSA Private Key.")
	}

	return res, nil
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwksx

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

// GenerateSigningKeys generates a JSON Web Key Set for signing.
func GenerateSigningKeys(id, alg string, bits int) (*jose.JSONWebKeySet, error) {
	if id == "" {
		id = uuid.Must(uuid.NewV4()).String()
	}

	key, err := generate(jose.SignatureAlgorithm(alg), bits)
	if err != nil {
		return nil, err
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    alg,
				Use:          "sig",
				Key:          key,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}

// GenerateSigningKeysAvailableAlgorithms lists available algorithms that are supported by GenerateSigningKeys.
func GenerateSigningKeysAvailableAlgorithms() []string {
	return []string{
		string(jose.HS256), string(jose.HS384), string(jose.HS512),
		string(jose.ES256), string(jose.ES384), string(jose.ES512), string(jose.EdDSA),
		string(jose.RS256), string(jose.RS384), string(jose.RS512), string(jose.PS256), string(jose.PS384), string(jose.PS512),
	}
}

// generate generates keypair for corresponding SignatureAlgorithm.
func generate(alg jose.SignatureAlgorithm, bits int) (crypto.PrivateKey, error) {
	switch alg {
	case jose.ES256, jose.ES384, jose.ES512, jose.EdDSA:
		keylen := map[jose.SignatureAlgorithm]int{
			jose.ES256: 256,
			jose.ES384: 384,
			jose.ES512: 521, // sic!
			jose.EdDSA: 256,
		}
		if bits != 0 && bits != keylen[alg] {
			return nil, errors.Errorf(`jwksx: "%s" does not support arbitrary key length`, alg)
		}
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		if bits == 0 {
			bits = 2048
		}
		if bits < 2048 {
			return nil, errors.Errorf(`jwksx: key size must be at least 2048 bit for algorithm "%s"`, alg)
		}
	case jose.HS256:
		if bits == 0 {
			bits = 256
		}
		if bits < 256 {
			return nil, errors.Errorf(`jwksx: key size must be at least 256 bit for algorithm "%s"`, alg)
		}
	case jose.HS384:
		if bits == 0 {
			bits = 384
		}
		if bits < 384 {
			return nil, errors.Errorf(`jwksx: key size must be at least 2038448 bit for algorithm "%s"`, alg)
		}
	case jose.HS512:
		if bits == 0 {
			bits = 1024
		}
		if bits < 512 {
			return nil, errors.Errorf(`jwksx: key size must be at least 512 bit for algorithm "%s"`, alg)
		}
	}

	switch alg {
	case jose.ES256:
		// The cryptographic operations are implemented using constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		return key, errors.Wrapf(err, "jwks: unable to generate key")
	case jose.ES384:
		// NB: The cryptographic operations do not use constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		return key, errors.Wrapf(err, "jwks: unable to generate key")
	case jose.ES512:
		// NB: The cryptographic operations do not use constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		return key, errors.Wrapf(err, "jwks: unable to generate key")
	case jose.EdDSA:
		_, key, err := ed25519.GenerateKey(rand.Reader)
		return key, errors.Wrapf(err, "jwks: unable to generate key")
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		key, err := rsa.GenerateKey(rand.Reader, bits)
		return key, errors.Wrapf(err, "jwks: unable to generate key")
	case jose.HS256, jose.HS384, jose.HS512:
		if bits%8 != 0 {
			return nil, errors.Errorf(`jwksx: key size must be a multiple of 8 for algorithm "%s" but got: %d`, alg, bits)
		}

		key := make([]byte, bits/8)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return nil, errors.Wrapf(err, "jwks: unable to generate key")
		}
		return key, nil
	default:
		return nil, errors.Errorf(`jwksx: available algorithms are "%+v" but unknown algorithm was requested: "%s"`, GenerateSigningKeysAvailableAlgorithms(), alg)
	}
}

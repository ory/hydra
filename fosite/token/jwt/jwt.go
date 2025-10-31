// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package jwt is able to generate and validate json web tokens.
// Follows https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32

package jwt

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"strings"

	"github.com/go-jose/go-jose/v3"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
)

type Signer interface {
	Generate(ctx context.Context, claims MapClaims, header Mapper) (string, string, error)
	Validate(ctx context.Context, token string) (string, error)
	Hash(ctx context.Context, in []byte) ([]byte, error)
	Decode(ctx context.Context, token string) (*Token, error)
	GetSignature(ctx context.Context, token string) (string, error)
	GetSigningMethodLength(ctx context.Context) int
}

var SHA256HashSize = crypto.SHA256.Size()

type GetPrivateKeyFunc func(ctx context.Context) (interface{}, error)

// DefaultSigner is responsible for generating and validating JWT challenges
type DefaultSigner struct {
	GetPrivateKey GetPrivateKeyFunc
}

// Generate generates a new authorize code or returns an error. set secret
func (j *DefaultSigner) Generate(ctx context.Context, claims MapClaims, header Mapper) (string, string, error) {
	key, err := j.GetPrivateKey(ctx)
	if err != nil {
		return "", "", err
	}

	switch t := key.(type) {
	case *jose.JSONWebKey:
		return generateToken(claims, header, jose.SignatureAlgorithm(t.Algorithm), t)
	case jose.JSONWebKey:
		return generateToken(claims, header, jose.SignatureAlgorithm(t.Algorithm), t)
	case *rsa.PrivateKey:
		return generateToken(claims, header, jose.RS256, t)
	case *ecdsa.PrivateKey:
		return generateToken(claims, header, jose.ES256, t)
	case jose.OpaqueSigner:
		switch tt := t.Public().Key.(type) {
		case *rsa.PrivateKey:
			alg := jose.RS256
			if len(t.Algs()) > 0 {
				alg = t.Algs()[0]
			}

			return generateToken(claims, header, alg, t)
		case *ecdsa.PrivateKey:
			alg := jose.ES256
			if len(t.Algs()) > 0 {
				alg = t.Algs()[0]
			}

			return generateToken(claims, header, alg, t)
		default:
			return "", "", errors.Errorf("unsupported private / public key pairs: %T, %T", t, tt)
		}
	default:
		return "", "", errors.Errorf("unsupported private key type: %T", t)
	}
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (j *DefaultSigner) Validate(ctx context.Context, token string) (string, error) {
	key, err := j.GetPrivateKey(ctx)
	if err != nil {
		return "", err
	}

	if t, ok := key.(*jose.JSONWebKey); ok {
		key = t.Key
	}

	switch t := key.(type) {
	case *rsa.PrivateKey:
		return validateToken(token, t.PublicKey)
	case *ecdsa.PrivateKey:
		return validateToken(token, t.PublicKey)
	case jose.OpaqueSigner:
		return validateToken(token, t.Public().Key)
	default:
		return "", errors.New("Unable to validate token. Invalid PrivateKey type")
	}
}

// Decode will decode a JWT token
func (j *DefaultSigner) Decode(ctx context.Context, token string) (*Token, error) {
	key, err := j.GetPrivateKey(ctx)
	if err != nil {
		return nil, err
	}

	if t, ok := key.(*jose.JSONWebKey); ok {
		key = t.Key
	}

	switch t := key.(type) {
	case *rsa.PrivateKey:
		return decodeToken(token, t.PublicKey)
	case *ecdsa.PrivateKey:
		return decodeToken(token, t.PublicKey)
	case jose.OpaqueSigner:
		return decodeToken(token, t.Public().Key)
	default:
		return nil, errors.New("Unable to decode token. Invalid PrivateKey type")
	}
}

// GetSignature will return the signature of a token
func (j *DefaultSigner) GetSignature(ctx context.Context, token string) (string, error) {
	return getTokenSignature(token)
}

// Hash will return a given hash based on the byte input or an error upon fail
func (j *DefaultSigner) Hash(ctx context.Context, in []byte) ([]byte, error) {
	return hashSHA256(in)
}

// GetSigningMethodLength will return the length of the signing method
func (j *DefaultSigner) GetSigningMethodLength(ctx context.Context) int {
	return SHA256HashSize
}

func generateToken(claims MapClaims, header Mapper, signingMethod jose.SignatureAlgorithm, privateKey interface{}) (rawToken string, sig string, err error) {
	if header == nil || claims == nil {
		err = errors.New("either claims or header is nil")
		return
	}

	token := NewWithClaims(signingMethod, claims)
	token.Header = assign(token.Header, header.ToMap())

	rawToken, err = token.SignedString(privateKey)
	if err != nil {
		return
	}

	sig, err = getTokenSignature(rawToken)
	return
}

func decodeToken(token string, verificationKey interface{}) (*Token, error) {
	keyFunc := func(*Token) (interface{}, error) { return verificationKey, nil }
	return ParseWithClaims(token, MapClaims{}, keyFunc)
}

func validateToken(tokenStr string, verificationKey interface{}) (string, error) {
	_, err := decodeToken(tokenStr, verificationKey)
	if err != nil {
		return "", err
	}
	return getTokenSignature(tokenStr)
}

func getTokenSignature(token string) (string, error) {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return "", errors.New("header, body and signature must all be set")
	}
	return split[2], nil
}

func hashSHA256(in []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(in)
	if err != nil {
		return []byte{}, errorsx.WithStack(err)
	}
	return hash.Sum([]byte{}), nil
}

func assign(a, b map[string]interface{}) map[string]interface{} {
	for k, w := range b {
		if _, ok := a[k]; ok {
			continue
		}
		a[k] = w
	}
	return a
}

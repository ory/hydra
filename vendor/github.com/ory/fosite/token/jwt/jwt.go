/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

// Package jwt is able to generate and validate json web tokens.
// Follows https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32

package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

type JWTStrategy interface {
	Generate(ctx context.Context, claims jwt.Claims, header Mapper) (string, string, error)
	Validate(ctx context.Context, token string) (string, error)
	Hash(ctx context.Context, in []byte) ([]byte, error)
	Decode(ctx context.Context, token string) (*jwt.Token, error)
	GetSignature(ctx context.Context, token string) (string, error)
	GetSigningMethodLength() int
}

// RS256JWTStrategy is responsible for generating and validating JWT challenges
type RS256JWTStrategy struct {
	PrivateKey *rsa.PrivateKey
}

// Generate generates a new authorize code or returns an error. set secret
func (j *RS256JWTStrategy) Generate(ctx context.Context, claims jwt.Claims, header Mapper) (string, string, error) {
	if header == nil || claims == nil {
		return "", "", errors.New("Either claims or header is nil.")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header = assign(token.Header, header.ToMap())

	var sig, sstr string
	var err error
	if sstr, err = token.SigningString(); err != nil {
		return "", "", errors.WithStack(err)
	}

	if sig, err = token.Method.Sign(sstr, j.PrivateKey); err != nil {
		return "", "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s.%s", sstr, sig), sig, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (j *RS256JWTStrategy) Validate(ctx context.Context, token string) (string, error) {
	if _, err := j.Decode(ctx, token); err != nil {
		return "", errors.WithStack(err)
	}

	return j.GetSignature(ctx, token)
}

// Decode will decode a JWT token
func (j *RS256JWTStrategy) Decode(ctx context.Context, token string) (*jwt.Token, error) {
	// Parse the token.
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return &j.PrivateKey.PublicKey, nil
	})

	if err != nil {
		return parsedToken, errors.WithStack(err)
	} else if !parsedToken.Valid {
		return parsedToken, errors.WithStack(fosite.ErrInactiveToken)
	}

	return parsedToken, err
}

// GetSignature will return the signature of a token
func (j *RS256JWTStrategy) GetSignature(ctx context.Context, token string) (string, error) {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return "", errors.New("Header, body and signature must all be set")
	}
	return split[2], nil
}

// Hash will return a given hash based on the byte input or an error upon fail
func (j *RS256JWTStrategy) Hash(ctx context.Context, in []byte) ([]byte, error) {
	// SigningMethodRS256
	hash := sha256.New()
	_, err := hash.Write(in)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return hash.Sum([]byte{}), nil
}

// GetSigningMethodLength will return the length of the signing method
func (j *RS256JWTStrategy) GetSigningMethodLength() int {
	return jwt.SigningMethodRS256.Hash.Size()
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

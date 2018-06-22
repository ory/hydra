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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk

import (
	"crypto/rsa"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/ory/fosite/token/jwt"
	"github.com/pkg/errors"
)

func NewRS256JWTStrategy(m Manager, set string) (*RS256JWTStrategy, error) {
	j := &RS256JWTStrategy{
		Manager:          m,
		RS256JWTStrategy: &jwt.RS256JWTStrategy{},
		Set:              set,
	}
	if err := j.refresh(); err != nil {
		return nil, err
	}
	return j, nil
}

type RS256JWTStrategy struct {
	RS256JWTStrategy *jwt.RS256JWTStrategy
	Manager          Manager
	Set              string

	publicKey    *rsa.PublicKey
	privateKey   *rsa.PrivateKey
	publicKeyID  string
	privateKeyID string
}

func (j *RS256JWTStrategy) Hash(in []byte) ([]byte, error) {
	return j.RS256JWTStrategy.Hash(in)
}

// GetSigningMethodLength will return the length of the signing method
func (j *RS256JWTStrategy) GetSigningMethodLength() int {
	return j.RS256JWTStrategy.GetSigningMethodLength()
}

func (j *RS256JWTStrategy) GetSignature(token string) (string, error) {
	return j.RS256JWTStrategy.GetSignature(token)
}

func (j *RS256JWTStrategy) Generate(claims jwt2.Claims, header jwt.Mapper) (string, string, error) {
	if err := j.refresh(); err != nil {
		return "", "", err
	}

	return j.RS256JWTStrategy.Generate(claims, header)
}

func (j *RS256JWTStrategy) Validate(token string) (string, error) {
	if err := j.refresh(); err != nil {
		return "", err
	}

	return j.RS256JWTStrategy.Validate(token)
}

func (j *RS256JWTStrategy) Decode(token string) (*jwt2.Token, error) {
	if err := j.refresh(); err != nil {
		return nil, err
	}

	return j.RS256JWTStrategy.Decode(token)
}

func (j *RS256JWTStrategy) GetPublicKeyID() (string, error) {
	if err := j.refresh(); err != nil {
		return "", err
	}

	return j.publicKeyID, nil
}

func (j *RS256JWTStrategy) refresh() error {
	keys, err := j.Manager.GetKeySet(j.Set)
	if err != nil {
		return err
	}

	public, err := FindKeyByPrefix(keys, "public")
	if err != nil {
		return err
	}

	private, err := FindKeyByPrefix(keys, "private")
	if err != nil {
		return err
	}

	if k, ok := private.Key.(*rsa.PrivateKey); !ok {
		return errors.New("unable to type assert key to *rsa.PublicKey")
	} else {
		j.privateKey = k
		j.RS256JWTStrategy.PrivateKey = k
	}

	if k, ok := public.Key.(*rsa.PublicKey); !ok {
		return errors.New("unable to type assert key to *rsa.PublicKey")
	} else {
		j.publicKey = k
		j.publicKeyID = public.KeyID
	}

	return nil
}

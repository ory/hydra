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
 */

package jwk

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

type RS256Generator struct {
	KeyLength int
}

func (g *RS256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	if g.KeyLength < 4096 {
		g.KeyLength = 4096
	}

	key, err := rsa.GenerateKey(rand.Reader, g.KeyLength)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	if id == "" {
		id = uuid.New()
	}

	// jose does not support this...
	key.Precomputed = rsa.PrecomputedValues{}
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    "RS256",
				Key:          key,
				Use:          use,
				KeyID:        ider("private", id),
				Certificates: []*x509.Certificate{},
			},
			{
				Algorithm:    "RS256",
				Use:          use,
				Key:          &key.PublicKey,
				KeyID:        ider("public", id),
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}

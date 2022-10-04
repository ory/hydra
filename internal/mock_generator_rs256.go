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

package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

type veryInsecureRS256Generator struct{}

func (g *veryInsecureRS256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	/* #nosec G403 - this is ok because this generator is only used in tests. */
	key, err := rsa.GenerateKey(rand.Reader, 512)
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
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}

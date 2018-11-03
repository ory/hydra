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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

func First(keys []jose.JSONWebKey) *jose.JSONWebKey {
	if len(keys) == 0 {
		return nil
	}
	return &keys[0]
}

func FindKeyByPrefix(set *jose.JSONWebKeySet, prefix string) (key *jose.JSONWebKey, err error) {
	keys, err := FindKeysByPrefix(set, prefix)
	if err != nil {
		return nil, err
	}

	return First(keys.Keys), nil
}

func FindKeysByPrefix(set *jose.JSONWebKeySet, prefix string) (*jose.JSONWebKeySet, error) {
	keys := new(jose.JSONWebKeySet)

	for _, k := range set.Keys {
		if len(k.KeyID) >= len(prefix)+1 && k.KeyID[:len(prefix)+1] == prefix+":" {
			keys.Keys = append(keys.Keys, k)
		}
	}

	if len(keys.Keys) == 0 {
		return nil, errors.Errorf("Unable to find key with prefix %s in JSON Web Key Set", prefix)
	}

	return keys, nil
}

func PEMBlockForKey(key interface{}) (*pem.Block, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, errors.New("Invalid key type")
	}
}

func ider(typ, id string) string {
	if id == "" {
		id = uuid.New()
	}
	return fmt.Sprintf("%s:%s", typ, id)
}

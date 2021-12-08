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
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/x"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

func EnsureAsymmetricKeypairExists(ctx context.Context, r InternalRegistry, alg, set string) error {
	_, _, err := GetOrGenerateKeys(ctx, r, r.KeyManager(), set, set, alg)
	return err
}

func GetOrGenerateKeys(ctx context.Context, r InternalRegistry, m Manager, set, kid, alg string) (public, private *jose.JSONWebKey, err error) {
	keys, err := m.GetKeySet(ctx, set)
	if errors.Is(err, x.ErrNotFound) || keys != nil && len(keys.Keys) == 0 {
		r.Logger().Warnf("JSON Web Key Set \"%s\" does not exist yet, generating new key pair...", set)
		keys, err = m.GenerateAndPersistKeySet(ctx, set, kid, alg, "sig")
		if err != nil {
			return nil, nil, err
		}
	} else if err != nil {
		return nil, nil, err
	}

	pubKey, pubKeyErr := FindPublicKey(keys)
	privKey, privKeyErr := FindPrivateKey(keys)
	if pubKeyErr == nil && privKeyErr == nil {
		return pubKey, privKey, nil
	} else {
		if pubKeyErr != nil {
			r.Logger().Warnf("Public JSON Web Key not found in JSON Web Key Set %s, generating new key pair...", set)
		} else {
			r.Logger().Warnf("Private JSON Web Key not found in JSON Web Key Set %s, generating new key pair...", set)
		}
		keys, err = m.GenerateAndPersistKeySet(ctx, set, kid, alg, "sig")
		if err != nil {
			return nil, nil, err
		}
		pubKey, err := FindPublicKey(keys)
		if err != nil {
			return nil, nil, err
		}
		privKey, err := FindPrivateKey(keys)
		if err != nil {
			return nil, nil, err
		}
		return pubKey, privKey, nil
	}
}

func First(keys []jose.JSONWebKey) *jose.JSONWebKey {
	if len(keys) == 0 {
		return nil
	}
	return &keys[0]
}

func FindPublicKey(set *jose.JSONWebKeySet) (key *jose.JSONWebKey, err error) {
	keys := ExcludePrivateKeys(set)
	if len(keys.Keys) == 0 {
		return nil, errors.New("key not found")
	}

	return First(keys.Keys), nil
}

func FindPrivateKey(set *jose.JSONWebKeySet) (key *jose.JSONWebKey, err error) {
	keys := ExcludePublicKeys(set)
	if len(keys.Keys) == 0 {
		return nil, errors.New("key not found")
	}

	return First(keys.Keys), nil
}

func ExcludePublicKeys(set *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	keys := new(jose.JSONWebKeySet)

	for _, k := range set.Keys {
		_, ecdsaOk := k.Key.(*ecdsa.PublicKey)
		_, ed25519OK := k.Key.(ed25519.PublicKey)
		_, rsaOK := k.Key.(*rsa.PublicKey)

		if !ecdsaOk && !ed25519OK && !rsaOK {
			keys.Keys = append(keys.Keys, k)
		}
	}
	return keys
}

func ExcludePrivateKeys(set *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	keys := new(jose.JSONWebKeySet)

	for _, k := range set.Keys {
		_, ecdsaOk := k.Key.(*ecdsa.PublicKey)
		_, ed25519OK := k.Key.(ed25519.PublicKey)
		_, rsaOK := k.Key.(*rsa.PublicKey)

		if ecdsaOk || ed25519OK || rsaOK {
			keys.Keys = append(keys.Keys, k)
		}
	}
	return keys
}

func ExcludeOpaquePrivateKeys(set *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	keys := new(jose.JSONWebKeySet)

	for _, k := range set.Keys {
		if _, opaque := k.Key.(jose.OpaqueSigner); !opaque {
			keys.Keys = append(keys.Keys, k)
		}
	}
	return keys
}

func PEMBlockForKey(key interface{}) (*pem.Block, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, errorsx.WithStack(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	case ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return nil, errorsx.WithStack(err)
		}
		return &pem.Block{Type: "PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, errors.New("Invalid key type")
	}
}

func Ider(typ, id string) string {
	if id == "" {
		id = uuid.New()
	}
	return fmt.Sprintf("%s:%s", typ, id)
}

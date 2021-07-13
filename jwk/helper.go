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

func EnsureAsymmetricKeypairExists(ctx context.Context, r InternalRegistry, g KeyGenerator, set string) error {
	_, _, err := AsymmetricKeypair(ctx, r, g, set)
	return err
}

func AsymmetricKeypair(ctx context.Context, r InternalRegistry, g KeyGenerator, set string) (public, private *jose.JSONWebKey, err error) {
	priv, err := GetOrCreateKey(ctx, r, g, set, "private")
	if err != nil {
		return nil, nil, err
	}

	pub, err := GetOrCreateKey(ctx, r, g, set, "public")
	if err != nil {
		return nil, nil, err
	}

	return pub, priv, nil
}

func GetOrCreateKey(ctx context.Context, r InternalRegistry, g KeyGenerator, set, prefix string) (*jose.JSONWebKey, error) {
	keys, err := r.KeyManager().GetKeySet(ctx, set)
	if errors.Is(err, x.ErrNotFound) || keys != nil && len(keys.Keys) == 0 {
		r.Logger().Warnf("JSON Web Key Set \"%s\" does not exist yet, generating new key pair...", set)
		keys, err = createKey(ctx, r, g, set)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	key, err := FindKeyByPrefix(keys, prefix)
	if err != nil {
		r.Logger().Warnf("JSON Web Key with prefix %s not found in JSON Web Key Set %s, generating new key pair...", prefix, set)

		keys, err = createKey(ctx, r, g, set)
		if err != nil {
			return nil, err
		}

		key, err = FindKeyByPrefix(keys, prefix)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}

func createKey(ctx context.Context, r InternalRegistry, g KeyGenerator, set string) (*jose.JSONWebKeySet, error) {
	keys, err := g.Generate(uuid.New(), "sig")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not generate JSON Web Key Set \"%s\".", set)
	}

	for i, k := range keys.Keys {
		k.Use = "sig"
		keys.Keys[i] = k
	}

	if err = r.KeyManager().AddKeySet(ctx, set, keys); err != nil {
		return nil, errors.Wrapf(err, "Could not persist JSON Web Key Set \"%s\".", set)
	}

	return keys, nil
}

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
			return nil, errorsx.WithStack(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
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

// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"github.com/ory/x/josex"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/x"

	"github.com/ory/x/errorsx"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
)

func EnsureAsymmetricKeypairExists(ctx context.Context, r InternalRegistry, alg, set string) error {
	_, err := GetOrGenerateKeySetPrivateKey(ctx, r.KeyManager(), set, set, alg)
	return err
}

func GetOrGenerateKeySetPrivateKey(ctx context.Context, m Manager, set, kid, alg string) (*jose.JSONWebKey, error) {
	keySet, err := GetOrGenerateKeySet(ctx, m, set, kid, alg)
	if err != nil {
		return nil, err
	}

	privKey, err := FindPrivateKey(keySet)
	if err == nil {
		return privKey, nil
	}

	keySet, err = m.GenerateAndPersistKeySet(ctx, set, kid, alg, "sig")
	if err != nil {
		return nil, err
	}

	return FindPrivateKey(keySet)
}

func GetOrGenerateKeySet(ctx context.Context, m Manager, set, kid, alg string) (*jose.JSONWebKeySet, error) {
	keys, err := m.GetKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return nil, err
	} else if keys != nil && len(keys.Keys) > 0 {
		return keys, nil
	}

	return m.GenerateAndPersistKeySet(ctx, set, kid, alg, "sig")
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
		if !k.IsPublic() {
			keys.Keys = append(keys.Keys, k)
		}
	}

	return keys
}

func ExcludePrivateKeys(set *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	keys := new(jose.JSONWebKeySet)
	for i := range set.Keys {
		keys.Keys = append(keys.Keys, josex.ToPublicKey(&set.Keys[i]))
	}
	return keys
}

func ExcludeOpaquePrivateKeys(set *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	keys := new(jose.JSONWebKeySet)
	for i := range set.Keys {
		if _, opaque := set.Keys[i].Key.(jose.OpaqueSigner); opaque {
			keys.Keys = append(keys.Keys, josex.ToPublicKey(&set.Keys[i]))
		} else {
			keys.Keys = append(keys.Keys, set.Keys[i])
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

func OnlyPublicSDKKeys(in []hydra.JsonWebKey) (out []hydra.JsonWebKey, _ error) {
	var interim []jose.JSONWebKey
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(&in); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	if err := json.NewDecoder(&b).Decode(&interim); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	for i, key := range interim {
		interim[i] = key.Public()
	}

	b.Reset()
	if err := json.NewEncoder(&b).Encode(&interim); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	var keys []hydra.JsonWebKey
	if err := json.NewDecoder(&b).Decode(&keys); err != nil {
		return nil, errors.Wrap(err, "failed to encode JSON Web Key Set")
	}

	return keys, nil
}

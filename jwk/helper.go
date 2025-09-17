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
	"sync"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/josex"
)

var mapLock sync.RWMutex
var locks = map[string]*sync.RWMutex{}

func getLock(set string) *sync.RWMutex {
	mapLock.Lock()
	defer mapLock.Unlock()
	if _, ok := locks[set]; !ok {
		locks[set] = new(sync.RWMutex)
	}
	return locks[set]
}

func EnsureAsymmetricKeypairExists(t testing.TB, r InternalRegistry, alg, set string) {
	_, err := GetOrGenerateKeys(t.Context(), r, r.KeyManager(), set, alg)
	require.NoError(t, err)
}

func GetOrGenerateKeys(ctx context.Context, r InternalRegistry, m Manager, set, alg string) (private *jose.JSONWebKey, err error) {
	getLock(set).Lock()
	defer getLock(set).Unlock()

	keys, err := m.GetKeySet(ctx, set)
	if errors.Is(err, x.ErrNotFound) || err == nil && len(keys.Keys) == 0 {
		r.Logger().Warnf("JSON Web Key Set %q does not exist yet, generating new key pair...", set)
		keys, err = m.GenerateAndPersistKeySet(ctx, set, "", alg, "sig")
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	privKey, privKeyErr := FindPrivateKey(keys)
	if privKeyErr == nil {
		return privKey, nil
	}
	r.Logger().WithField("jwks", set).Warnf("JSON Web Key not found in JSON Web Key Set %s, generating new key pair...", set)

	keys, err = m.GenerateAndPersistKeySet(ctx, set, "", alg, "sig")
	if err != nil {
		return nil, err
	}

	return FindPrivateKey(keys)
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
			return nil, errors.WithStack(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	case ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return nil, errors.WithStack(err)
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

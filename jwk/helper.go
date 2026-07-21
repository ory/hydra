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
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/singleflight"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/josex"
	"github.com/ory/x/popx"
)

// generateFlight deduplicates concurrent key-set generation per network, set,
// and kid, so that a fresh system under parallel load creates exactly one key
// pair per set and replica.
var generateFlight singleflight.Group

// generateTimeout bounds a single key-set generation, including the RSA key
// generation and the database round trips. It applies to the detached flight
// context, which is otherwise not canceled when the triggering request is.
const generateTimeout = 30 * time.Second

func EnsureAsymmetricKeypairExists(t testing.TB, r InternalRegistry, alg, set string) {
	_, err := GetOrGenerateKeys(t.Context(), r, set, alg)
	require.NoError(t, err)
}

func GetOrGenerateKeys(ctx context.Context, r InternalRegistry, set, alg string) (private *jose.JSONWebKey, err error) {
	keys, err := r.KeyManager().GetKeySet(ctx, set)
	if err == nil {
		if privKey, findErr := FindPrivateKey(keys); findErr == nil {
			return privKey, nil
		}
	} else if !errors.Is(err, x.ErrNotFound) {
		return nil, err
	}

	keys, err = getOrGenerateKeySet(ctx, r, set, "", alg, "sig")
	if err != nil {
		return nil, err
	}

	return FindPrivateKey(keys)
}

// getOrGenerateKeySet returns the key set, generating and persisting a new key
// pair unless the set already contains a private key. Concurrent calls for the
// same network, set, and kid share a single generation and receive its result.
func getOrGenerateKeySet(ctx context.Context, r InternalRegistry, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	if popx.InTransaction(ctx) {
		// The caller runs inside a database transaction (e.g. the token
		// endpoint wraps fosite's NewAccessResponse in one). Generate on the
		// caller's own transaction: a flight on a separate connection
		// deadlocks with the locks that transaction holds (SQLite), and keys
		// it generates are only visible to and roll back with that
		// transaction, so they must not be shared with other callers.
		return readOrGenerateKeySet(ctx, r, set, kid, alg, use)
	}

	// The network ID scopes the flight to one tenant. The NUL separators keep
	// distinct (set, kid) tuples from mapping to the same key; results are
	// shared between all callers on the same key, so the key must never be
	// ambiguous.
	key := r.Networker().NetworkID(ctx).String() + "\x00" + set + "\x00" + kid

	// A flight is bounded by generateTimeout, so waiting much longer than that
	// means it died without delivering (e.g. runtime.Goexit in a test); error
	// out rather than blocking forever.
	ctx, cancel := context.WithTimeout(ctx, 2*generateTimeout)
	defer cancel()

	ch := generateFlight.DoChan(key, func() (_ any, err error) {
		// A panic must surface as an error: singleflight would otherwise
		// crash the whole process when channel waiters are present.
		defer func() {
			if e := recover(); e != nil {
				err = errors.Errorf("panic during JSON Web Key generation for set %q: %v", set, e)
			}
		}()

		// The flight must not be canceled by the request that happened to
		// start it: other requests may be waiting on the result. Context
		// values (network ID, tracing) are preserved.
		fctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), generateTimeout)
		defer cancel()

		return readOrGenerateKeySet(fctx, r, set, kid, alg, use)
	})

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result := <-ch:
		if result.Err != nil {
			return nil, result.Err
		}
		keys, _ := result.Val.(*jose.JSONWebKeySet)
		return keys, nil
	}
}

// readOrGenerateKeySet reads the key set and generates and persists a new key
// pair if it does not contain a private key. It runs on the caller's context,
// including any transaction the context carries.
func readOrGenerateKeySet(ctx context.Context, r InternalRegistry, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	keys, err := r.KeyManager().GetKeySet(ctx, set)
	switch {
	case err == nil && len(keys.Keys) > 0:
		// Another caller may have generated the keys since the caller's read.
		if _, findErr := FindPrivateKey(keys); findErr == nil {
			return keys, nil
		}
		r.Logger().WithField("jwks", set).Warnf("JSON Web Key not found in JSON Web Key Set %s, generating new key pair...", set)
	case err == nil || errors.Is(err, x.ErrNotFound):
		r.Logger().Warnf("JSON Web Key Set %q does not exist yet, generating new key pair...", set)
	default:
		return nil, err
	}

	return r.KeyManager().GenerateAndPersistKeySet(ctx, set, kid, alg, use)
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

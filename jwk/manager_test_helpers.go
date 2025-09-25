// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/assertx"
)

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func canonicalizeThumbprints(js []jose.JSONWebKey) []jose.JSONWebKey {
	for k, v := range js {
		js[k] = canonicalizeKeyThumbprints(&v)
	}
	return js
}

func canonicalizeKeyThumbprints(v *jose.JSONWebKey) jose.JSONWebKey {
	if len(v.CertificateThumbprintSHA1) == 0 {
		v.CertificateThumbprintSHA1 = nil
	}
	if len(v.CertificateThumbprintSHA256) == 0 {
		v.CertificateThumbprintSHA256 = nil
	}
	return *v
}

func TestHelperManagerKey(m Manager, algo string, keys *jose.JSONWebKeySet, suffix string) func(t *testing.T) {
	priv := canonicalizeThumbprints(keys.Key(suffix))
	var pub []jose.JSONWebKey
	for _, k := range priv {
		pub = append(pub, canonicalizeThumbprints([]jose.JSONWebKey{k.Public()})...)
	}

	return func(t *testing.T) {
		ctx := t.Context()

		set := algo + uuid.Must(uuid.NewV4()).String()

		_, err := m.GetKey(ctx, set, suffix)
		assert.NotNil(t, err)

		err = m.AddKey(ctx, set, First(priv))
		require.NoError(t, err)

		got, err := m.GetKey(ctx, set, suffix)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, priv, canonicalizeThumbprints(got.Keys))

		addKey := First(pub)
		addKey.KeyID = uuid.Must(uuid.NewV4()).String()
		err = m.AddKey(ctx, set, addKey)
		require.NoError(t, err)

		got, err = m.GetKey(ctx, set, suffix)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, priv, canonicalizeThumbprints(got.Keys))

		// Because MySQL
		time.Sleep(time.Second * 2)

		newKID := "new-key-id:" + suffix
		pub[0].KeyID = newKID
		pub[0].Use = "sig"
		err = m.AddKey(ctx, set, First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(ctx, set, newKID)
		require.NoError(t, err)
		newKey := First(got.Keys)
		assert.EqualValues(t, "sig", newKey.Use)

		newKey.Use = "enc"
		err = m.UpdateKey(ctx, set, newKey)
		require.NoError(t, err)
		updated, err := m.GetKey(ctx, set, newKID)
		require.NoError(t, err)
		updatedKey := First(updated.Keys)
		assert.EqualValues(t, "enc", updatedKey.Use)

		keys, err = m.GetKeySet(ctx, set)
		require.NoError(t, err)
		var found bool
		for _, k := range keys.Keys {
			if k.KeyID == newKID {
				found = true
				break
			}
		}
		assert.True(t, found, "Key not found in key set: %s / %s\n%+v", keys, newKID)

		beforeDeleteKeysCount := len(keys.Keys)
		err = m.DeleteKey(ctx, set, suffix)
		require.NoError(t, err)

		_, err = m.GetKey(ctx, set, suffix)
		require.Error(t, err)

		keys, err = m.GetKeySet(ctx, set)
		require.NoError(t, err)
		assert.EqualValues(t, beforeDeleteKeysCount-1, len(keys.Keys))
	}
}

func TestHelperManagerKeySet(m Manager, algo string, keys *jose.JSONWebKeySet, suffix string, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()

		if parallel {
			t.Parallel()
		}
		set := uuid.Must(uuid.NewV4()).String()
		_, err := m.GetKeySet(ctx, algo+set)
		require.Error(t, err)

		err = m.AddKeySet(ctx, algo+set, keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(ctx, algo+set)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, canonicalizeThumbprints(keys.Key(suffix)), canonicalizeThumbprints(got.Key(suffix)))
		assertx.EqualAsJSON(t, canonicalizeThumbprints(keys.Key(suffix)), canonicalizeThumbprints(got.Key(suffix)))

		for i := range got.Keys {
			got.Keys[i].Use = "enc"
		}
		err = m.UpdateKeySet(ctx, algo+set, got)
		require.NoError(t, err)

		updated, err := m.GetKeySet(ctx, algo+set)
		require.NoError(t, err)
		assert.EqualValues(t, "enc", updated.Key(suffix)[0].Public().Use)
		assert.EqualValues(t, "enc", updated.Key(suffix)[0].Use)

		err = m.DeleteKeySet(ctx, algo+set)
		require.NoError(t, err)

		_, err = m.GetKeySet(ctx, algo+set)
		require.Error(t, err)
	}
}

func TestHelperManagerGenerateAndPersistKeySet(m Manager, alg string, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()

		if parallel {
			t.Parallel()
		}
		_, err := m.GetKeySet(ctx, "foo")
		require.Error(t, err)

		keys, err := m.GenerateAndPersistKeySet(ctx, "foo", "bar", alg, "sig")
		require.NoError(t, err)
		genPub, err := FindPublicKey(keys)
		require.NoError(t, err)
		require.NotEmpty(t, genPub)
		genPriv, err := FindPrivateKey(keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(ctx, "foo")
		require.NoError(t, err)
		gotPub, err := FindPublicKey(got)
		require.NoError(t, err)
		require.NotEmpty(t, gotPub)
		gotPriv, err := FindPrivateKey(got)
		require.NoError(t, err)

		assertx.EqualAsJSON(t, canonicalizeKeyThumbprints(genPub), canonicalizeKeyThumbprints(gotPub))

		assert.EqualValues(t, genPriv.KeyID, gotPriv.KeyID)

		err = m.DeleteKeySet(ctx, "foo")
		require.NoError(t, err)

		_, err = m.GetKeySet(ctx, "foo")
		require.Error(t, err)
	}
}

func TestHelperNID(t1ValidNID, t2InvalidNID Manager) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		jwks, err := GenerateJWK(jose.RS256, "2022-03-11-ks-1-kid", "test")
		require.NoError(t, err)
		require.Error(t, t2InvalidNID.AddKey(ctx, "2022-03-11-k-1", &jwks.Keys[0]))
		require.NoError(t, t1ValidNID.AddKey(ctx, "2022-03-11-k-1", &jwks.Keys[0]))
		require.Error(t, t2InvalidNID.AddKeySet(ctx, "2022-03-11-ks-1", jwks))
		require.NoError(t, t1ValidNID.AddKeySet(ctx, "2022-03-11-ks-1", jwks))
		require.NoError(t, t2InvalidNID.DeleteKey(ctx, "2022-03-11-ks-1", jwks.Keys[0].KeyID)) // Delete doesn't report error if key doesn't exist
		require.NoError(t, t1ValidNID.DeleteKey(ctx, "2022-03-11-ks-1", jwks.Keys[0].KeyID))
		_, err = t2InvalidNID.GenerateAndPersistKeySet(ctx, "2022-03-11-ks-2", "2022-03-11-ks-2-kid", "RS256", "sig")
		require.Error(t, err)
		gks2, err := t1ValidNID.GenerateAndPersistKeySet(ctx, "2022-03-11-ks-2", "2022-03-11-ks-2-kid", "RS256", "sig")
		require.NoError(t, err)

		_, err = t1ValidNID.GetKey(ctx, "2022-03-11-ks-2", gks2.Keys[0].KeyID)
		require.NoError(t, err)
		_, err = t2InvalidNID.GetKey(ctx, "2022-03-11-ks-2", gks2.Keys[0].KeyID)
		require.Error(t, err)

		_, err = t1ValidNID.GetKeySet(ctx, "2022-03-11-ks-2")
		require.NoError(t, err)
		_, err = t2InvalidNID.GetKeySet(ctx, "2022-03-11-ks-2")
		require.Error(t, err)
		updatedKey := &gks2.Keys[0]
		updatedKey.Use = "enc"
		require.Error(t, t2InvalidNID.UpdateKey(ctx, "2022-03-11-ks-2", updatedKey))
		require.NoError(t, t1ValidNID.UpdateKey(ctx, "2022-03-11-ks-2", updatedKey))
		gks2.Keys[0].Use = "enc"
		require.Error(t, t2InvalidNID.UpdateKeySet(ctx, "2022-03-11-ks-2", gks2))
		require.NoError(t, t1ValidNID.UpdateKeySet(ctx, "2022-03-11-ks-2", gks2))
		require.NoError(t, t2InvalidNID.DeleteKeySet(ctx, "2022-03-11-ks-2"))
		require.NoError(t, t1ValidNID.DeleteKeySet(ctx, "2022-03-11-ks-2"))
	}
}

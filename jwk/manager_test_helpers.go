// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ory/x/assertx"

	"github.com/ory/x/errorsx"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errorsx.WithStack(err)
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
		set := algo + uuid.NewString()

		_, err := m.GetKey(context.TODO(), set, suffix)
		assert.NotNil(t, err)

		err = m.AddKey(context.TODO(), set, First(priv))
		require.NoError(t, err)

		got, err := m.GetKey(context.TODO(), set, suffix)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, priv, canonicalizeThumbprints(got.Keys))

		addKey := First(pub)
		addKey.KeyID = uuid.NewString()
		err = m.AddKey(context.TODO(), set, addKey)
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), set, suffix)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, priv, canonicalizeThumbprints(got.Keys))

		// Because MySQL
		time.Sleep(time.Second * 2)

		newKID := "new-key-id:" + suffix
		pub[0].KeyID = newKID
		pub[0].Use = "sig"
		err = m.AddKey(context.TODO(), set, First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), set, newKID)
		require.NoError(t, err)
		newKey := First(got.Keys)
		assert.EqualValues(t, "sig", newKey.Use)

		newKey.Use = "enc"
		err = m.UpdateKey(context.TODO(), set, newKey)
		require.NoError(t, err)
		updated, err := m.GetKey(context.TODO(), set, newKID)
		require.NoError(t, err)
		updatedKey := First(updated.Keys)
		assert.EqualValues(t, "enc", updatedKey.Use)

		keys, err = m.GetKeySet(context.TODO(), set)
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
		err = m.DeleteKey(context.TODO(), set, suffix)
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), set, suffix)
		require.Error(t, err)

		keys, err = m.GetKeySet(context.TODO(), set)
		require.NoError(t, err)
		assert.EqualValues(t, beforeDeleteKeysCount-1, len(keys.Keys))
	}
}

func TestHelperManagerKeySet(m Manager, algo string, keys *jose.JSONWebKeySet, suffix string, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		set := uuid.NewString()
		_, err := m.GetKeySet(context.TODO(), algo+set)
		require.Error(t, err)

		err = m.AddKeySet(context.TODO(), algo+set, keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), algo+set)
		require.NoError(t, err)
		assertx.EqualAsJSON(t, canonicalizeThumbprints(keys.Key(suffix)), canonicalizeThumbprints(got.Key(suffix)))
		assertx.EqualAsJSON(t, canonicalizeThumbprints(keys.Key(suffix)), canonicalizeThumbprints(got.Key(suffix)))

		for i := range got.Keys {
			got.Keys[i].Use = "enc"
		}
		err = m.UpdateKeySet(context.TODO(), algo+set, got)
		require.NoError(t, err)

		updated, err := m.GetKeySet(context.TODO(), algo+set)
		require.NoError(t, err)
		assert.EqualValues(t, "enc", updated.Key(suffix)[0].Public().Use)
		assert.EqualValues(t, "enc", updated.Key(suffix)[0].Use)

		err = m.DeleteKeySet(context.TODO(), algo+set)
		require.NoError(t, err)

		_, err = m.GetKeySet(context.TODO(), algo+set)
		require.Error(t, err)
	}
}

func TestHelperManagerGenerateAndPersistKeySet(m Manager, alg string, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		_, err := m.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		keys, err := m.GenerateAndPersistKeySet(context.TODO(), "foo", "bar", alg, "sig")
		require.NoError(t, err)
		genPub, err := FindPublicKey(keys)
		require.NoError(t, err)
		require.NotEmpty(t, genPub)
		genPriv, err := FindPrivateKey(keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), "foo")
		require.NoError(t, err)
		gotPub, err := FindPublicKey(got)
		require.NoError(t, err)
		require.NotEmpty(t, gotPub)
		gotPriv, err := FindPrivateKey(got)
		require.NoError(t, err)

		assertx.EqualAsJSON(t, canonicalizeKeyThumbprints(genPub), canonicalizeKeyThumbprints(gotPub))

		assert.EqualValues(t, genPriv.KeyID, gotPriv.KeyID)

		err = m.DeleteKeySet(context.TODO(), "foo")
		require.NoError(t, err)

		_, err = m.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)
	}
}

func TestHelperManagerNIDIsolationKeySet(t1 Manager, t2 Manager, alg string) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := t1.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)
		_, err = t2.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		_, err = t1.GenerateAndPersistKeySet(context.TODO(), "foo", "bar", alg, "sig")
		require.NoError(t, err)
		keys, err := t1.GetKeySet(context.TODO(), "foo")
		require.NoError(t, err)
		_, err = t2.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		err = t2.DeleteKeySet(context.TODO(), "foo")
		require.Error(t, err)
		err = t1.DeleteKeySet(context.TODO(), "foo")
		require.NoError(t, err)
		_, err = t1.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		err = t1.AddKeySet(context.TODO(), "foo", keys)
		require.NoError(t, err)
		err = t2.DeleteKeySet(context.TODO(), "foo")
		require.Error(t, err)

		for i := range keys.Keys {
			keys.Keys[i].Use = "enc"
		}
		err = t1.UpdateKeySet(context.TODO(), "foo", keys)
		require.Error(t, err)
		for i := range keys.Keys {
			keys.Keys[i].Use = "err"
		}
		err = t2.UpdateKeySet(context.TODO(), "foo", keys)
		require.Error(t, err)
		updated, err := t1.GetKeySet(context.TODO(), "foo")
		require.NoError(t, err)
		for i := range updated.Keys {
			assert.EqualValues(t, "enc", updated.Keys[i].Use)
		}

		err = t1.DeleteKeySet(context.TODO(), "foo")
		require.Error(t, err)
	}
}

func TestHelperNID(t1ValidNID Manager, t2InvalidNID Manager) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		jwks, err := GenerateJWK(ctx, jose.RS256, "2022-03-11-ks-1-kid", "test")
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

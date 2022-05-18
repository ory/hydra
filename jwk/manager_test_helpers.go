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
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"
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
	pub := canonicalizeThumbprints(keys.Key("public:" + suffix))
	priv := canonicalizeThumbprints(keys.Key("private:" + suffix))

	return func(t *testing.T) {
		_, err := m.GetKey(context.TODO(), algo+"faz", "baz")
		assert.NotNil(t, err)

		err = m.AddKey(context.TODO(), algo+"faz", First(priv))
		require.NoError(t, err)

		got, err := m.GetKey(context.TODO(), algo+"faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, canonicalizeThumbprints(got.Keys))

		err = m.AddKey(context.TODO(), algo+"faz", First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), algo+"faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, canonicalizeThumbprints(got.Keys))

		got, err = m.GetKey(context.TODO(), algo+"faz", "public:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, pub, canonicalizeThumbprints(got.Keys))

		// Because MySQL
		time.Sleep(time.Second * 2)

		First(pub).KeyID = "new-key-id:" + suffix
		First(pub).Use = "sig"
		err = m.AddKey(context.TODO(), algo+"faz", First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), algo+"faz", "new-key-id:"+suffix)
		require.NoError(t, err)
		newKey := First(got.Keys)
		assert.EqualValues(t, "sig", newKey.Use)

		newKey.Use = "enc"
		err = m.UpdateKey(context.TODO(), algo+"faz", newKey)
		require.NoError(t, err)
		updated, err := m.GetKey(context.TODO(), algo+"faz", "new-key-id:"+suffix)
		require.NoError(t, err)
		updatedKey := First(updated.Keys)
		assert.EqualValues(t, "enc", updatedKey.Use)

		keys, err = m.GetKeySet(context.TODO(), algo+"faz")
		require.NoError(t, err)
		assert.EqualValues(t, "new-key-id:"+suffix, First(keys.Keys).KeyID)

		beforeDeleteKeysCount := len(keys.Keys)
		err = m.DeleteKey(context.TODO(), algo+"faz", "public:"+suffix)
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), algo+"faz", "public:"+suffix)
		require.Error(t, err)

		keys, err = m.GetKeySet(context.TODO(), algo+"faz")
		require.NoError(t, err)
		assert.EqualValues(t, beforeDeleteKeysCount-1, len(keys.Keys))
	}
}

func TestHelperManagerKeySet(m Manager, algo string, keys *jose.JSONWebKeySet, suffix string, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		_, err := m.GetKeySet(context.TODO(), algo+"foo")
		require.Error(t, err)

		err = m.AddKeySet(context.TODO(), algo+"bar", keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), algo+"bar")
		require.NoError(t, err)
		assert.Equal(t, canonicalizeThumbprints(keys.Key("public:"+suffix)), canonicalizeThumbprints(got.Key("public:"+suffix)))
		assert.Equal(t, canonicalizeThumbprints(keys.Key("private:"+suffix)), canonicalizeThumbprints(got.Key("private:"+suffix)))

		for i, _ := range got.Keys {
			got.Keys[i].Use = "enc"
		}
		err = m.UpdateKeySet(context.TODO(), algo+"bar", got)
		require.NoError(t, err)
		updated, err := m.GetKeySet(context.TODO(), algo+"bar")
		require.NoError(t, err)
		assert.EqualValues(t, "enc", First(updated.Key("public:"+suffix)).Use)
		assert.EqualValues(t, "enc", First(updated.Key("private:"+suffix)).Use)

		err = m.DeleteKeySet(context.TODO(), algo+"bar")
		require.NoError(t, err)

		_, err = m.GetKeySet(context.TODO(), algo+"bar")
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
		genPriv, err := FindPrivateKey(keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), "foo")
		require.NoError(t, err)
		gotPub, err := FindPublicKey(got)
		require.NoError(t, err)
		gotPriv, err := FindPrivateKey(got)
		require.NoError(t, err)

		assert.Equal(t, canonicalizeKeyThumbprints(genPub), canonicalizeKeyThumbprints(gotPub))
		assert.Equal(t, canonicalizeKeyThumbprints(genPriv), canonicalizeKeyThumbprints(gotPriv))

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
		kg := RS256Generator{}
		jwks, err := kg.Generate("2022-03-11-ks-1-kid", "test")
		require.NoError(t, err)
		require.Error(t, t2InvalidNID.AddKey(ctx, "2022-03-11-k-1", &jwks.Keys[0]))
		require.NoError(t, t1ValidNID.AddKey(ctx, "2022-03-11-k-1", &jwks.Keys[0]))
		require.Error(t, t2InvalidNID.AddKeySet(ctx, "2022-03-11-ks-1", jwks))
		require.NoError(t, t1ValidNID.AddKeySet(ctx, "2022-03-11-ks-1", jwks))
		require.Error(t, t2InvalidNID.DeleteKey(ctx, "2022-03-11-ks-1", jwks.Keys[0].KeyID))
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
		gks2.Keys[1].Use = "enc"
		require.Error(t, t2InvalidNID.UpdateKeySet(ctx, "2022-03-11-ks-2", gks2))
		require.NoError(t, t1ValidNID.UpdateKeySet(ctx, "2022-03-11-ks-2", gks2))
		require.Error(t, t2InvalidNID.DeleteKeySet(ctx, "2022-03-11-ks-2"))
		require.NoError(t, t1ValidNID.DeleteKeySet(ctx, "2022-03-11-ks-2"))
	}
}

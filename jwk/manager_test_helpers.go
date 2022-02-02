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

func TestHelperManagerKeySet(m Manager, algo string, keys *jose.JSONWebKeySet, suffix string) func(t *testing.T) {
	return func(t *testing.T) {
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

func TestHelperManagerGenerateAndPersistKeySet(m Manager, alg string) func(t *testing.T) {
	return func(t *testing.T) {
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

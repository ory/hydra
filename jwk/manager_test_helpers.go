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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"
)

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func TestHelperManagerKey(m Manager, keys *jose.JSONWebKeySet, suffix string) func(t *testing.T) {
	pub := keys.Key("public:" + suffix)
	priv := keys.Key("private:" + suffix)

	return func(t *testing.T) {
		t.Parallel()
		_, err := m.GetKey(context.TODO(), "faz", "baz")
		assert.NotNil(t, err)

		err = m.AddKey(context.TODO(), "faz", First(priv))
		require.NoError(t, err)

		got, err := m.GetKey(context.TODO(), "faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, got.Keys)

		err = m.AddKey(context.TODO(), "faz", First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), "faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, got.Keys)

		got, err = m.GetKey(context.TODO(), "faz", "public:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, pub, got.Keys)

		// Because MySQL
		time.Sleep(time.Second * 2)

		First(pub).KeyID = "new-key-id:" + suffix
		err = m.AddKey(context.TODO(), "faz", First(pub))
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), "faz", "new-key-id:"+suffix)
		require.NoError(t, err)

		keys, err = m.GetKeySet(context.TODO(), "faz")
		require.NoError(t, err)
		assert.EqualValues(t, "new-key-id:"+suffix, First(keys.Keys).KeyID)

		err = m.DeleteKey(context.TODO(), "faz", "public:"+suffix)
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), "faz", "public:"+suffix)
		require.Error(t, err)
	}
}

func TestHelperManagerKeySet(m Manager, keys *jose.JSONWebKeySet, suffix string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		_, err := m.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		err = m.AddKeySet(context.TODO(), "bar", keys)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), "bar")
		require.NoError(t, err)
		assert.Equal(t, keys.Key("public:"+suffix), got.Key("public:"+suffix))
		assert.Equal(t, keys.Key("private:"+suffix), got.Key("private:"+suffix))

		err = m.DeleteKeySet(context.TODO(), "bar")
		require.NoError(t, err)

		_, err = m.GetKeySet(context.TODO(), "bar")
		require.Error(t, err)
	}
}

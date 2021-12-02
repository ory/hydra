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
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/pem"
	"strings"
	"testing"

	"gopkg.in/square/go-jose.v2/cryptosigner"

	"gopkg.in/square/go-jose.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIder(t *testing.T) {
	assert.True(t, len(Ider("public", "")) > len("public:"))
	assert.Equal(t, "public:foo", Ider("public", "foo"))
}

func TestHandlerFindPublicKey(t *testing.T) {
	var testRSGenerator = RS256Generator{}
	var testECDSAGenerator = ECDSA256Generator{}
	var testEdDSAGenerator = EdDSAGenerator{}

	t.Run("Test_Helper/Run_FindPublicKey_With_RSA", func(t *testing.T) {
		RSIDKS, _ := testRSGenerator.Generate("test-id-1", "sig")
		keys, err := FindPublicKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("public", "test-id-1"))
		assert.IsType(t, keys.Key, new(rsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, _ := testECDSAGenerator.Generate("test-id-2", "sig")
		keys, err := FindPublicKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("public", "test-id-2"))
		assert.IsType(t, keys.Key, new(ecdsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, _ := testEdDSAGenerator.Generate("test-id-3", "sig")
		keys, err := FindPublicKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("public", "test-id-3"))
		assert.IsType(t, keys.Key, ed25519.PublicKey{})
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_KeyNotFound", func(t *testing.T) {
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestHandlerFindPrivateKey(t *testing.T) {
	var testRSGenerator = RS256Generator{}
	var testECDSAGenerator = ECDSA256Generator{}
	var testEdDSAGenerator = EdDSAGenerator{}

	t.Run("Test_Helper/Run_FindPrivateKey_With_RSA", func(t *testing.T) {
		RSIDKS, _ := testRSGenerator.Generate("test-id-1", "sig")
		keys, err := FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("private", "test-id-1"))
		assert.IsType(t, keys.Key, new(rsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, _ := testECDSAGenerator.Generate("test-id-2", "sig")
		keys, err := FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("private", "test-id-2"))
		assert.IsType(t, keys.Key, new(ecdsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, _ := testEdDSAGenerator.Generate("test-id-3", "sig")
		keys, err := FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, Ider("private", "test-id-3"))
		assert.IsType(t, keys.Key, ed25519.PrivateKey{})
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_KeyNotFound", func(t *testing.T) {
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestPEMBlockForKey(t *testing.T) {
	var testRSGenerator = RS256Generator{}
	var testECDSAGenerator = ECDSA256Generator{}
	var testEdDSAGenerator = EdDSAGenerator{}

	t.Run("Test_Helper/Run_PEMBlockForKey_With_RSA", func(t *testing.T) {
		RSIDKS, _ := testRSGenerator.Generate("test-id-1", "sig")
		key, err := FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		pemBlock, err := PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "RSA PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, _ := testECDSAGenerator.Generate("test-id-2", "sig")
		key, err := FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "EC PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, _ := testEdDSAGenerator.Generate("test-id-3", "sig")
		key, err := FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_InvalidKeyType", func(t *testing.T) {
		key := dsa.PrivateKey{}
		_, err := PEMBlockForKey(key)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "Invalid key type"))
	})
}

func TestExcludeOpaquePrivateKeys(t *testing.T) {
	var testRSGenerator = RS256Generator{}

	opaqueKeys, err := testRSGenerator.Generate("test-id-1", "sig")
	assert.NoError(t, err)
	assert.Len(t, opaqueKeys.Keys, 2)
	opaqueKeys.Keys[0].Key = cryptosigner.Opaque(opaqueKeys.Keys[0].Key.(*rsa.PrivateKey))
	keys := ExcludeOpaquePrivateKeys(opaqueKeys)
	assert.Len(t, keys.Keys, 1)
	assert.IsType(t, new(rsa.PublicKey), keys.Keys[0].Key)
}

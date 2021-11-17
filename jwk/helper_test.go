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
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
}

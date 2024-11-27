// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk_test

import (
	"context"
	"crypto"
	"crypto/dsa" //lint:ignore SA1019 used for testing invalid key types
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"strings"
	"testing"

	"github.com/ory/hydra/v2/internal/testhelpers"

	hydra "github.com/ory/hydra-client-go/v2"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/cryptosigner"
	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
)

type fakeSigner struct {
	pk crypto.PublicKey
}

func (f *fakeSigner) Sign(_ io.Reader, _ []byte, _ crypto.SignerOpts) ([]byte, error) {
	return []byte("signature"), nil
}

func (f *fakeSigner) Public() crypto.PublicKey {
	return f.pk
}

func TestHandlerFindPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("Test_Helper/Run_FindPublicKey_With_RSA", func(t *testing.T) {
		t.Parallel()
		RSIDKS, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_Opaque", func(t *testing.T) {
		t.Parallel()
		key, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		RSIDKS := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
			Algorithm:                   "RS256",
			Use:                         "sig",
			Key:                         cryptosigner.Opaque(&fakeSigner{pk: key.Keys[0].Public().Key}),
			KeyID:                       "test-id-1",
			Certificates:                []*x509.Certificate{},
			CertificateThumbprintSHA1:   []uint8{},
			CertificateThumbprintSHA256: []uint8{},
		}, {
			Algorithm:                   "RS256",
			Use:                         "sig",
			Key:                         key.Keys[0].Public().Key,
			KeyID:                       "test-id-1",
			Certificates:                []*x509.Certificate{},
			CertificateThumbprintSHA1:   []uint8{},
			CertificateThumbprintSHA256: []uint8{},
		}}}
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, "test-id-1", keys.KeyID)
		assert.IsType(t, new(rsa.PublicKey), keys.Key)
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_ECDSA", func(t *testing.T) {
		t.Parallel()
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_EdDSA", func(t *testing.T) {
		t.Parallel()
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-3")
		assert.IsType(t, keys.Key, ed25519.PublicKey{})
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_KeyNotFound", func(t *testing.T) {
		t.Parallel()
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := jwk.FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestHandlerFindPrivateKey(t *testing.T) {
	t.Parallel()
	t.Run("Test_Helper/Run_FindPrivateKey_With_RSA", func(t *testing.T) {
		RSIDKS, _ := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		keys, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-3")
		assert.IsType(t, keys.Key, ed25519.PrivateKey{})
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_KeyNotFound", func(t *testing.T) {
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := jwk.FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestPEMBlockForKey(t *testing.T) {
	t.Parallel()
	t.Run("Test_Helper/Run_PEMBlockForKey_With_RSA", func(t *testing.T) {
		RSIDKS, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "RSA PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "EC PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_InvalidKeyType", func(t *testing.T) {
		key := dsa.PrivateKey{}
		_, err := jwk.PEMBlockForKey(key)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "Invalid key type"))
	})
}

func TestExcludeOpaquePrivateKeys(t *testing.T) {
	t.Parallel()
	opaqueKeys, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
	assert.NoError(t, err)
	require.Len(t, opaqueKeys.Keys, 1)
	opaqueKeys.Keys[0].Key = cryptosigner.Opaque(opaqueKeys.Keys[0].Key.(*rsa.PrivateKey))

	keys := jwk.ExcludeOpaquePrivateKeys(opaqueKeys)

	require.Len(t, keys.Keys, 1)
	k := keys.Keys[0]
	_, isPublic := k.Key.(*rsa.PublicKey)
	assert.True(t, isPublic)
}

func TestGetOrGenerateKeys(t *testing.T) {
	t.Parallel()
	reg := testhelpers.NewMockedRegistry(t, &contextx.Default{})

	setId := uuid.NewUUID().String()
	keyId := uuid.NewUUID().String()

	keySet, _ := jwk.GenerateJWK(context.Background(), jose.RS256, keyId, "sig")
	keySetWithoutPrivateKey := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{keySet.Keys[0].Public()},
	}

	km := func(t *testing.T) *MockManager {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		return NewMockManager(ctrl)
	}

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GetKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(nil, errors.Wrap(x.ErrNotFound, ""))
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GetKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySet, nil)
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.NoError(t, err)
		assert.Equal(t, privKey, &keySet.Keys[0])
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySetWithoutPrivateKey, nil).Times(1)
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "key not found")
	})
}

func TestOnlyPublicSDKKeys(t *testing.T) {
	set, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
	require.NoError(t, err)

	out, err := json.Marshal(set.Keys)
	require.NoError(t, err)

	var sdkSet []hydra.JsonWebKey
	require.NoError(t, json.Unmarshal(out, &sdkSet))

	assert.NotEmpty(t, sdkSet[0].P)
	result, err := jwk.OnlyPublicSDKKeys(sdkSet)
	require.NoError(t, err)

	assert.Empty(t, result[0].P)
}

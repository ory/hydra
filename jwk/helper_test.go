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
	"sync"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/cryptosigner"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
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
		RSIDKS, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_Opaque", func(t *testing.T) {
		t.Parallel()
		key, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
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
		ECDSAIDKS, err := jwk.GenerateJWK(jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_EdDSA", func(t *testing.T) {
		t.Parallel()
		EdDSAIDKS, err := jwk.GenerateJWK(jose.EdDSA, "test-id-3", "sig")
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
		RSIDKS, _ := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
		keys, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(jose.EdDSA, "test-id-3", "sig")
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
		RSIDKS, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "RSA PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "EC PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(jose.EdDSA, "test-id-3", "sig")
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
	opaqueKeys, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
	assert.NoError(t, err)
	require.Len(t, opaqueKeys.Keys, 1)
	opaqueKeys.Keys[0].Key = cryptosigner.Opaque(opaqueKeys.Keys[0].Key.(*rsa.PrivateKey))

	keys := jwk.ExcludeOpaquePrivateKeys(opaqueKeys)

	require.Len(t, keys.Keys, 1)
	k := keys.Keys[0]
	_, isPublic := k.Key.(*rsa.PublicKey)
	assert.True(t, isPublic)
}

type regWithManager struct {
	*driver.RegistrySQL
	km jwk.Manager
}

func (r regWithManager) KeyManager() jwk.Manager { return r.km }

// slowKeyManager delegates to a real manager but delays reads, simulating a
// database that responds slowly.
type slowKeyManager struct {
	jwk.Manager
	delay time.Duration
}

func (s *slowKeyManager) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	time.Sleep(s.delay)
	return s.Manager.GetKeySet(ctx, set)
}

func TestGetOrGenerateKeysConcurrency(t *testing.T) {
	t.Parallel()

	t.Run("concurrent misses generate exactly one key set", func(t *testing.T) {
		t.Parallel()
		reg := testhelpers.NewRegistryMemory(t)
		setID := uuid.NewUUID().String()

		start := make(chan struct{})
		keys := make([]*jose.JSONWebKey, 20)
		var wg sync.WaitGroup
		for i := range keys {
			wg.Go(func() {
				<-start
				key, err := jwk.GetOrGenerateKeys(t.Context(), reg, setID, "RS256")
				if assert.NoError(t, err) {
					keys[i] = key
				}
			})
		}
		close(start)
		wg.Wait()

		set, err := reg.KeyManager().GetKeySet(t.Context(), setID)
		require.NoError(t, err)
		require.Len(t, set.Keys, 1)
		for _, key := range keys {
			require.NotNil(t, key)
			assert.Equal(t, set.Keys[0].KeyID, key.KeyID)
		}
	})

	t.Run("generation inside a transaction does not deadlock", func(t *testing.T) {
		t.Parallel()
		reg := testhelpers.NewRegistryMemory(t)
		setID := uuid.NewUUID().String()

		// Generation from within a transaction (the token endpoint wraps
		// NewAccessResponse in one) must run on that transaction: a separate
		// connection would deadlock against the transaction's locks.
		begin := time.Now()
		var generated *jose.JSONWebKey
		err := reg.Transaction(t.Context(), func(ctx context.Context) error {
			var err error
			generated, err = jwk.GetOrGenerateKeys(ctx, reg, setID, "RS256")
			if err != nil {
				return err
			}
			return errors.New("force rollback")
		})
		require.ErrorContains(t, err, "force rollback")
		require.NotNil(t, generated)
		assert.Less(t, time.Since(begin), 30*time.Second)

		key, err := jwk.GetOrGenerateKeys(t.Context(), reg, setID, "RS256")
		require.NoError(t, err)
		if reg.Config().HSMEnabled() {
			// HSM key storage is not transactional: the keys survive the
			// rollback and the next caller reuses them.
			assert.Equal(t, generated.KeyID, key.KeyID)
		} else {
			// The keys were generated within the caller's transaction and
			// rolled back with it; the next caller generates a fresh set.
			assert.NotEqual(t, generated.KeyID, key.KeyID)
		}
	})

	t.Run("reads of an existing key set are not serialized", func(t *testing.T) {
		t.Parallel()
		reg := testhelpers.NewRegistryMemory(t)
		if reg.Config().HSMEnabled() {
			// The HSM key manager serializes all operations on a process-wide
			// per-token lock, so parallel tests generating keys distort any
			// read-latency measurement.
			t.Skip("read concurrency cannot be measured with an HSM-backed key manager")
		}
		setID := uuid.NewUUID().String()

		_, err := jwk.GetOrGenerateKeys(t.Context(), reg, setID, "RS256")
		require.NoError(t, err)

		const delay = 250 * time.Millisecond
		const n = 8
		slowReg := regWithManager{
			RegistrySQL: reg,
			km:          &slowKeyManager{Manager: reg.KeyManager(), delay: delay},
		}

		begin := time.Now()
		var wg sync.WaitGroup
		for range n {
			wg.Go(func() {
				_, err := jwk.GetOrGenerateKeys(t.Context(), slowReg, setID, "RS256")
				assert.NoError(t, err)
			})
		}
		wg.Wait()

		// Serialized reads take at least n*delay; concurrent reads finish in
		// roughly one delay.
		assert.Less(t, time.Since(begin), n*delay/2)
	})
}

func TestGetOrGenerateKeys(t *testing.T) {
	t.Parallel()
	reg := testhelpers.NewRegistryMemory(t)

	setID := uuid.NewUUID().String()
	keyID := uuid.NewUUID().String()

	keySet, err := jwk.GenerateJWK(jose.RS256, keyID, "sig")
	require.NoError(t, err)
	require.Len(t, keySet.Keys, 1)
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
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setID)).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(t.Context(), regWithManager{RegistrySQL: reg, km: keyManager}, setID, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setID)).Return(nil, errors.Wrap(x.ErrNotFound, "")).Times(2)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setID), gomock.Eq(""), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(t.Context(), regWithManager{RegistrySQL: reg, km: keyManager}, setID, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setID)).Return(keySetWithoutPrivateKey, nil).Times(2)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setID), gomock.Eq(""), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(t.Context(), regWithManager{RegistrySQL: reg, km: keyManager}, setID, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GetKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setID)).Return(keySetWithoutPrivateKey, nil).Times(2)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setID), gomock.Eq(""), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySet, nil)
		privKey, err := jwk.GetOrGenerateKeys(t.Context(), regWithManager{RegistrySQL: reg, km: keyManager}, setID, "RS256")
		assert.NoError(t, err)
		assert.Equal(t, privKey, &keySet.Keys[0])
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setID)).Return(keySetWithoutPrivateKey, nil).Times(2)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setID), gomock.Eq(""), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySetWithoutPrivateKey, nil).Times(1)
		privKey, err := jwk.GetOrGenerateKeys(t.Context(), regWithManager{RegistrySQL: reg, km: keyManager}, setID, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "key not found")
	})
}

func TestOnlyPublicSDKKeys(t *testing.T) {
	set, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
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

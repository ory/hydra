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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk_test

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/cryptosigner"

	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/ory/fosite/token/jwt"

	"github.com/ory/fosite/token/jwt"
	. "github.com/ory/hydra/jwk"
)

func TestRS256JWTStrategy(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)
	m := reg.KeyManager()

	_, err := m.GenerateAndPersistKeySet(context.TODO(), "foo-set", "foo", "RS256", "sig")
	require.NoError(t, err)

	s, err := NewRS256JWTStrategy(*conf, reg, func() string {
		return "foo-set"
	})

	require.NoError(t, err)
	a, b, err := s.Generate(context.TODO(), jwt2.MapClaims{"foo": "bar"}, &jwt.Headers{})
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.NotEmpty(t, b)

	_, err = s.Validate(context.TODO(), a)
	require.NoError(t, err)

	kidFoo, err := s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)

	_, err = m.GenerateAndPersistKeySet(context.TODO(), "foo-set", "bar", "RS256", "sig")
	require.NoError(t, err)

	a, b, err = s.Generate(context.TODO(), jwt2.MapClaims{"foo": "bar"}, &jwt.Headers{})
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.NotEmpty(t, b)

	_, err = s.Validate(context.TODO(), a)
	require.NoError(t, err)

	kidBar, err := s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)

	if conf.HsmEnabled() {
		assert.Equal(t, "foo", kidFoo)
		assert.Equal(t, "bar", kidBar)
	} else {
		assert.Equal(t, "public:foo", kidFoo)
		assert.Equal(t, "public:bar", kidBar)
	}
}

func TestRS256JWTStrategy_Refresh(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	ctrl := gomock.NewController(t)
	keyManager := NewMockManager(ctrl)
	reg := NewMockInternalRegistry(ctrl)
	defer ctrl.Finish()

	reg.EXPECT().KeyManager().Return(keyManager).AnyTimes()

	setId := uuid.NewUUID().String()
	keyId := uuid.NewUUID().String()

	rsaGenerator := &RS256Generator{KeyLength: 1024}
	rsaKeySet, err := rsaGenerator.Generate(keyId, "sig")
	require.NoError(t, err)
	edsaGenerator := &ECDSA256Generator{}
	edsaKeySet, err := edsaGenerator.Generate(keyId, "sig")
	require.NoError(t, err)

	t.Run("With_RsaKeyPair", func(t *testing.T) {
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(rsaKeySet, nil)
		strategy, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.NoError(t, err)
		require.IsType(t, new(rsa.PrivateKey), strategy.RS256JWTStrategy.PrivateKey)
	})

	t.Run("With_OpaqueKeyPair", func(t *testing.T) {
		opaquePrivateKey := cryptosigner.Opaque(rsaKeySet.Keys[0].Key.(*rsa.PrivateKey))
		keySetWithOpaquePrivateKey := &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{{
				Algorithm:                   "RS256",
				Use:                         "sig",
				Key:                         opaquePrivateKey,
				KeyID:                       keyId,
				Certificates:                []*x509.Certificate{},
				CertificateThumbprintSHA1:   []uint8{},
				CertificateThumbprintSHA256: []uint8{},
			}, rsaKeySet.Keys[1]},
		}

		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithOpaquePrivateKey, nil)
		strategy, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.NoError(t, err)
		require.IsType(t, opaquePrivateKey, strategy.RS256JWTStrategy.PrivateKey)
	})

	t.Run("With_GetKeySetError", func(t *testing.T) {
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(nil, errors.New("GetKeySetError"))
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "GetKeySetError")
	})

	t.Run("With_FindPublicKeyError", func(t *testing.T) {
		keySetWithoutPublicKey := &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{rsaKeySet.Keys[0]},
		}
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPublicKey, nil)
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "key not found")
	})

	t.Run("With_FindPrivateKeyError", func(t *testing.T) {
		keySetWithoutPrivateKey := &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{rsaKeySet.Keys[1]},
		}
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "key not found")
	})

	t.Run("With_PublicKeyTypeError", func(t *testing.T) {
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(edsaKeySet, nil)
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "unable to type assert key to *rsa.PublicKey")
	})

	t.Run("With_PrivateKeyTypeError", func(t *testing.T) {
		keyInvalidPrivateKeyType := &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{edsaKeySet.Keys[0], rsaKeySet.Keys[1]},
		}
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keyInvalidPrivateKeyType, nil)
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "unknown private key type")
	})

	t.Run("With_KeyPairIdsNotMatchError", func(t *testing.T) {
		rsaKeySet.Keys[0].KeyID = uuid.NewUUID().String()
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(rsaKeySet, nil)
		_, err := NewRS256JWTStrategy(*conf, reg, func() string {
			return setId
		})
		require.EqualError(t, err, "public and private key pair kids do not match")
		rsaKeySet.Keys[0].KeyID = rsaKeySet.Keys[1].KeyID
	})
}

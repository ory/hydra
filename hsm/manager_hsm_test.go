// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build hsm
// +build hsm

package hsm_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"
	"testing"

	"github.com/ThalesGroup/crypto11"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/cryptosigner"
	"github.com/golang/mock/gomock"
	"github.com/miekg/pkcs11"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/hsm"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
)

func TestDefaultKeyManager_HSMEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	reg, err := driver.New(t.Context(),
		driver.WithConfigOptions(configx.WithValues(map[string]any{
			config.KeyDSN:     "memory",
			config.HSMEnabled: true,
		})),
		driver.WithHSMContext(mockHsmContext),
	)
	require.NoError(t, err)
	assert.IsType(t, &jwk.ManagerStrategy{}, reg.KeyManager())
}

func TestKeyManager_HsmKeySetPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	keySetPrefix := "application_specific_prefix."
	c.MustSet(context.Background(), config.HSMKeySetPrefix, keySetPrefix)
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKey3072, err := rsa.GenerateKey(rand.Reader, 3072)
	require.NoError(t, err)
	rsaKey4096, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)

	rsaKeyPair3072 := NewMockSignerDecrypter(ctrl)
	rsaKeyPair3072.EXPECT().Public().Return(&rsaKey3072.PublicKey).AnyTimes()

	rsaKeyPair4096 := NewMockSignerDecrypter(ctrl)
	rsaKeyPair4096.EXPECT().Public().Return(&rsaKey4096.PublicKey).AnyTimes()

	ecdsaKeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaKeyPair.EXPECT().Public().Return(&ecdsaKey.PublicKey).AnyTimes()

	var kid = uuid.New()

	expectedPrefixedOpenIDConnectKeyName := fmt.Sprintf("%s%s", keySetPrefix, x.OpenIDConnectKeyName)

	t.Run("case=GenerateAndPersistKeySet", func(t *testing.T) {
		privateAttrSet, publicAttrSet := expectedKeyAttributes(t, expectedPrefixedOpenIDConnectKeyName, kid)
		hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return(nil, nil)
		hsmContext.EXPECT().GenerateRSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(4096)).Return(rsaKeyPair4096, nil)

		got, err := m.GenerateAndPersistKeySet(context.TODO(), x.OpenIDConnectKeyName, kid, "RS256", "sig")

		assert.NoError(t, err)
		expectedKeySet := expectedKeySet(rsaKeyPair4096, kid, "RS256", "sig")
		if !reflect.DeepEqual(got, expectedKeySet) {
			t.Errorf("GenerateAndPersistKeySet() got = %v, want %v", got, expectedKeySet)
		}
	})
	t.Run("case=GetKey", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return(rsaKeyPair4096, nil)
		hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair4096), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)

		got, err := m.GetKey(context.TODO(), x.OpenIDConnectKeyName, kid)

		assert.NoError(t, err)
		expectedKeySet := expectedKeySet(rsaKeyPair4096, kid, "RS256", "sig")
		if !reflect.DeepEqual(got, expectedKeySet) {
			t.Errorf("GetKey() got = %v, want %v", got, expectedKeySet)
		}
	})
	t.Run("case=GetKeyMinimalRsaKeyLengthError", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return(rsaKeyPair3072, nil)

		_, err := m.GetKey(context.TODO(), x.OpenIDConnectKeyName, kid)

		assert.ErrorIs(t, err, jwk.ErrMinimalRsaKeyLength)
	})
	t.Run("case=GetKeySet", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return([]crypto11.Signer{rsaKeyPair4096}, nil)
		hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair4096), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(kid)), nil)
		hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair4096), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)

		got, err := m.GetKeySet(context.TODO(), x.OpenIDConnectKeyName)

		assert.NoError(t, err)
		expectedKeySet := expectedKeySet(rsaKeyPair4096, kid, "RS256", "sig")
		if !reflect.DeepEqual(got, expectedKeySet) {
			t.Errorf("GetKey() got = %v, want %v", got, expectedKeySet)
		}
	})
	t.Run("case=GetKeySetMinimalRsaKeyLengthError", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return([]crypto11.Signer{rsaKeyPair3072}, nil)
		hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair3072), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(kid)), nil)

		_, err := m.GetKeySet(context.TODO(), x.OpenIDConnectKeyName)

		assert.ErrorIs(t, err, jwk.ErrMinimalRsaKeyLength)
	})
	t.Run("case=DeleteKey", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return(rsaKeyPair4096, nil)
		rsaKeyPair4096.EXPECT().Delete().Return(nil)

		err := m.DeleteKey(context.TODO(), x.OpenIDConnectKeyName, kid)

		assert.NoError(t, err)
	})
	t.Run("case=DeleteKeySet", func(t *testing.T) {
		hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(expectedPrefixedOpenIDConnectKeyName))).Return([]crypto11.Signer{rsaKeyPair4096}, nil)
		rsaKeyPair4096.EXPECT().Delete().Return(nil)

		err := m.DeleteKeySet(context.TODO(), x.OpenIDConnectKeyName)

		assert.NoError(t, err)
	})
}

func TestKeyManager_GenerateAndPersistKeySet(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)

	rsaKeyPair := NewMockSignerDecrypter(ctrl)
	rsaKeyPair.EXPECT().Public().Return(&rsaKey.PublicKey).AnyTimes()

	ecdsaKeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaKeyPair.EXPECT().Public().Return(&ecdsaKey.PublicKey).AnyTimes()

	var kid = uuid.New()

	type args struct {
		ctx context.Context
		set string
		kid string
		alg string
		use string
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		args       args
		want       *jose.JSONWebKeySet
		wantErrMsg string
		wantErr    error
	}{
		{
			name: "Generate RS256",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "RS256",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateRSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(4096)).Return(rsaKeyPair, nil)
			},
			want: expectedKeySet(rsaKeyPair, kid, "RS256", "sig"),
		},
		{
			name: "Generate RS256 with GenerateRSAKeyPairWithAttributes Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "RS256",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateRSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(4096)).Return(nil, errors.New("GenerateRSAKeyPairWithAttributesError"))
			},
			wantErrMsg: "GenerateRSAKeyPairWithAttributesError",
		},
		{
			name: "Generate ES256",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "ES256",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateECDSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(elliptic.P256())).Return(ecdsaKeyPair, nil)
			},
			want: expectedKeySet(ecdsaKeyPair, kid, "ES256", "sig"),
		},
		{
			name: "Generate ES256 with GenerateECDSAKeyPairWithAttributes Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "ES256",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateECDSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(elliptic.P256())).Return(nil, errors.New("GenerateECDSAKeyPairWithAttributesError"))
			},
			wantErrMsg: "GenerateECDSAKeyPairWithAttributesError",
		},
		{
			name: "Generate ES512",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "ES512",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateECDSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(elliptic.P521())).Return(ecdsaKeyPair, nil)
			},
			want: expectedKeySet(ecdsaKeyPair, kid, "ES512", "sig"),
		},
		{
			name: "Generate ES512 GenerateECDSAKeyPairWithAttributes Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "ES512",
				use: "sig",
			},
			setup: func(t *testing.T) {
				privateAttrSet, publicAttrSet := expectedKeyAttributes(t, x.OpenIDConnectKeyName, kid)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
				hsmContext.EXPECT().GenerateECDSAKeyPairWithAttributes(gomock.Eq(publicAttrSet), gomock.Eq(privateAttrSet), gomock.Eq(elliptic.P521())).Return(nil, errors.New("GenerateECDSAKeyPairWithAttributesError"))
			},
			wantErrMsg: "GenerateECDSAKeyPairWithAttributesError",
		},
		{
			name: "Generate unsupported",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "ES384",
				use: "sig",
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
			},
			wantErr: errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm),
		},
		{
			name: "Generate with FindKeyPair Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
				alg: "RS256",
				use: "sig",
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, errors.New("FindKeyPairError"))
			},
			wantErrMsg: "FindKeyPairError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got, err := m.GenerateAndPersistKeySet(tt.args.ctx, tt.args.set, tt.args.kid, tt.args.alg, tt.args.use)
			if tt.wantErr != nil {
				require.Nil(t, got)
				require.IsType(t, tt.wantErr, err)
			} else if len(tt.wantErrMsg) != 0 {
				require.Nil(t, got)
				require.EqualError(t, err, tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateAndPersistKeySet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyManager_GetKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)
	rsaKeyPair := NewMockSignerDecrypter(ctrl)
	rsaKeyPair.EXPECT().Public().Return(&rsaKey.PublicKey).AnyTimes()

	ecdsaP256Key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	ecdsaP256KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP256KeyPair.EXPECT().Public().Return(&ecdsaP256Key.PublicKey).AnyTimes()

	ecdsaP521Key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)
	ecdsaP521KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP521KeyPair.EXPECT().Public().Return(&ecdsaP521Key.PublicKey).AnyTimes()

	ecdsaP224Key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	require.NoError(t, err)
	ecdsaP224KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP224KeyPair.EXPECT().Public().Return(&ecdsaP224Key.PublicKey).AnyTimes()

	var kid = uuid.New()

	type args struct {
		ctx context.Context
		set string
		kid string
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		args       args
		want       *jose.JSONWebKeySet
		wantErrMsg string
		wantErr    error
	}{
		{
			name: "Get RS256 sig",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(rsaKeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
			},
			want: expectedKeySet(rsaKeyPair, kid, "RS256", "sig"),
		},
		{
			name: "Get RS256 enc",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(rsaKeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true), nil)
			},
			want: expectedKeySet(rsaKeyPair, kid, "RS256", "enc"),
		},
		{
			name: "Key usage attribute error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(rsaKeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, errors.New("GetAttributeError"))
			},
			want: expectedKeySet(rsaKeyPair, kid, "RS256", "sig"),
		},
		{
			name: "Get ES256 sig",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(ecdsaP256KeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP256KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
			},
			want: expectedKeySet(ecdsaP256KeyPair, kid, "ES256", "sig"),
		},
		{
			name: "Get ES256 enc",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(ecdsaP256KeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP256KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true), nil)
			},
			want: expectedKeySet(ecdsaP256KeyPair, kid, "ES256", "enc"),
		},
		{
			name: "Get ES512 sig",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(ecdsaP521KeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP521KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
			},
			want: expectedKeySet(ecdsaP521KeyPair, kid, "ES512", "sig"),
		},
		{
			name: "Get ES512 enc",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(ecdsaP521KeyPair, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP521KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true), nil)
			},
			want: expectedKeySet(ecdsaP521KeyPair, kid, "ES512", "enc"),
		},
		{
			name: "Key not found",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
			},
			wantErrMsg: "Not Found",
		},
		{
			name: "FindKeyPair Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, errors.New("FindKeyPairError"))
			},
			wantErrMsg: "FindKeyPairError",
		},
		{
			name: "Unsupported elliptic curve",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(ecdsaP224KeyPair, nil)
			},
			wantErr: errors.WithStack(jwk.ErrUnsupportedEllipticCurve),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got, err := m.GetKey(tt.args.ctx, tt.args.set, tt.args.kid)
			if tt.wantErr != nil {
				require.Nil(t, got)
				require.IsType(t, tt.wantErr, err)
			} else if len(tt.wantErrMsg) != 0 {
				require.Nil(t, got)
				require.EqualError(t, err, tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyManager_GetKeySet(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)
	rsaKid := uuid.New()
	rsaKeyPair := NewMockSignerDecrypter(ctrl)
	rsaKeyPair.EXPECT().Public().Return(&rsaKey.PublicKey).AnyTimes()

	ecdsaP256Key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	ecdsaP256Kid := uuid.New()
	ecdsaP256KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP256KeyPair.EXPECT().Public().Return(&ecdsaP256Key.PublicKey).AnyTimes()

	ecdsaP521Key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)
	ecdsaP521Kid := uuid.New()
	ecdsaP521KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP521KeyPair.EXPECT().Public().Return(&ecdsaP521Key.PublicKey).AnyTimes()

	ecdsaP224Key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	require.NoError(t, err)
	ecdsaP224Kid := uuid.New()
	ecdsaP224KeyPair := NewMockSignerDecrypter(ctrl)
	ecdsaP224KeyPair.EXPECT().Public().Return(&ecdsaP224Key.PublicKey).AnyTimes()

	allKeys := []crypto11.Signer{rsaKeyPair, ecdsaP256KeyPair, ecdsaP521KeyPair}

	var keys []jose.JSONWebKey
	keys = append(keys, createJSONWebKeys(rsaKeyPair, rsaKid, "RS256", "sig")...)
	keys = append(keys, createJSONWebKeys(ecdsaP256KeyPair, ecdsaP256Kid, "ES256", "sig")...)
	keys = append(keys, createJSONWebKeys(ecdsaP521KeyPair, ecdsaP521Kid, "ES512", "sig")...)

	type args struct {
		ctx context.Context
		set string
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		args       args
		want       *jose.JSONWebKeySet
		wantErrMsg string
		wantErr    error
	}{
		{
			name: "With multiple keys per set",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(allKeys, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(rsaKid)), nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP256KeyPair), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(ecdsaP256Kid)), nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP256KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP521KeyPair), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(ecdsaP521Kid)), nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP521KeyPair), gomock.Eq(crypto11.CkaDecrypt)).Return(nil, nil)
			},
			want: &jose.JSONWebKeySet{Keys: keys},
		},
		{
			name: "GetCkaIdAttributeError Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(allKeys, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(rsaKeyPair), gomock.Eq(crypto11.CkaId)).Return(nil, errors.New("GetCkaIdAttributeError"))
			},
			wantErrMsg: "GetCkaIdAttributeError",
		},
		{
			name: "Key set not found",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
			},
			wantErrMsg: "Not Found",
		},
		{
			name: "FindKeyPairs Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, errors.New("FindKeyPairsError"))
			},
			wantErrMsg: "FindKeyPairsError",
		},
		{
			name: "Unsupported elliptic curve",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return([]crypto11.Signer{ecdsaP224KeyPair}, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(ecdsaP224KeyPair), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(ecdsaP224Kid)), nil)
			},
			wantErr: errors.WithStack(jwk.ErrUnsupportedEllipticCurve),
		},
		{
			name: "Invalid key type Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				keyPair := NewMockSignerDecrypter(ctrl)
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return([]crypto11.Signer{keyPair}, nil)
				hsmContext.EXPECT().GetAttribute(gomock.Eq(keyPair), gomock.Eq(crypto11.CkaId)).Return(pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(rsaKid)), nil)
				keyPair.EXPECT().Public().Return(nil).Times(1)
			},
			wantErr: errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got, err := m.GetKeySet(tt.args.ctx, tt.args.set)
			if tt.wantErr != nil {
				require.Nil(t, got)
				require.IsType(t, tt.wantErr, err)
			} else if len(tt.wantErrMsg) != 0 {
				require.Nil(t, got)
				require.EqualError(t, err, tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyManager_DeleteKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKeyPair := NewMockSignerDecrypter(ctrl)

	kid := uuid.New()

	type args struct {
		ctx context.Context
		set string
		kid string
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		args       args
		wantErrMsg string
	}{
		{
			name: "Existing key",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(rsaKeyPair, nil)
				rsaKeyPair.EXPECT().Delete().Return(nil)
			},
		},
		{
			name: "Key not found",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
			},
			wantErrMsg: "Not Found",
		},
		{
			name: "FindKeyPair Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, errors.New("FindKeyPairError"))
			},
			wantErrMsg: "FindKeyPairError",
		},
		{
			name: "Delete Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
				kid: kid,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPair(gomock.Eq([]byte(kid)), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(rsaKeyPair, nil)
				rsaKeyPair.EXPECT().Delete().Return(errors.New("DeleteError"))
			},
			wantErrMsg: "DeleteError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			if err := m.DeleteKey(tt.args.ctx, tt.args.set, tt.args.kid); len(tt.wantErrMsg) != 0 {
				require.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

func TestKeyManager_DeleteKeySet(t *testing.T) {
	ctrl := gomock.NewController(t)
	hsmContext := NewMockContext(ctrl)
	defer ctrl.Finish()
	l := logrusx.New("", "")
	c := config.MustNew(t, l, configx.SkipValidation())
	m := hsm.NewKeyManager(hsmContext, c)

	rsaKeyPair1 := NewMockSignerDecrypter(ctrl)
	rsaKeyPair2 := NewMockSignerDecrypter(ctrl)
	allKeys := []crypto11.Signer{rsaKeyPair1, rsaKeyPair2}

	type args struct {
		ctx context.Context
		set string
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		args       args
		wantErrMsg string
	}{
		{
			name: "Existing key",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(allKeys, nil)
				rsaKeyPair1.EXPECT().Delete().Return(nil)
				rsaKeyPair2.EXPECT().Delete().Return(nil)
			},
		},
		{
			name: "Key not found",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, nil)
			},
			wantErrMsg: "Not Found",
		},
		{
			name: "FindKeyPairs Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(nil, errors.New("FindKeyPairsError"))
			},
			wantErrMsg: "FindKeyPairsError",
		},
		{
			name: "Delete Error",
			args: args{
				ctx: context.TODO(),
				set: x.OpenIDConnectKeyName,
			},
			setup: func(t *testing.T) {
				hsmContext.EXPECT().FindKeyPairs(gomock.Nil(), gomock.Eq([]byte(x.OpenIDConnectKeyName))).Return(allKeys, nil)
				rsaKeyPair1.EXPECT().Delete().Return(errors.New("DeleteError"))
			},
			wantErrMsg: "DeleteError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			if err := m.DeleteKeySet(tt.args.ctx, tt.args.set); len(tt.wantErrMsg) != 0 {
				require.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

func TestKeyManager_AddKey(t *testing.T) {
	m := &hsm.KeyManager{
		Context: nil,
	}
	err := m.AddKey(context.TODO(), x.OpenIDConnectKeyName, &jose.JSONWebKey{})
	assert.ErrorIs(t, err, hsm.ErrPreGeneratedKeys)
}

func TestKeyManager_AddKeySet(t *testing.T) {
	m := &hsm.KeyManager{
		Context: nil,
	}
	err := m.AddKeySet(context.TODO(), x.OpenIDConnectKeyName, &jose.JSONWebKeySet{})
	assert.ErrorIs(t, err, hsm.ErrPreGeneratedKeys)
}

func TestKeyManager_UpdateKey(t *testing.T) {
	m := &hsm.KeyManager{
		Context: nil,
	}
	err := m.UpdateKey(context.TODO(), x.OpenIDConnectKeyName, &jose.JSONWebKey{})
	assert.ErrorIs(t, err, hsm.ErrPreGeneratedKeys)
}

func TestKeyManager_UpdateKeySet(t *testing.T) {
	m := &hsm.KeyManager{
		Context: nil,
	}
	err := m.UpdateKeySet(context.TODO(), x.OpenIDConnectKeyName, &jose.JSONWebKeySet{})
	assert.ErrorIs(t, err, hsm.ErrPreGeneratedKeys)
}

func expectedKeyAttributes(t *testing.T, set, kid string) (crypto11.AttributeSet, crypto11.AttributeSet) {
	privateAttrSet, err := crypto11.NewAttributeSetWithIDAndLabel([]byte(kid), []byte(set))
	require.NoError(t, err)
	publicAttrSet, err := crypto11.NewAttributeSetWithIDAndLabel([]byte(kid), []byte(set))
	require.NoError(t, err)
	publicAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, false),
	})
	privateAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, false),
	})
	return privateAttrSet, publicAttrSet
}

func expectedKeySet(keyPair *MockSignerDecrypter, kid, alg, use string) *jose.JSONWebKeySet {
	return &jose.JSONWebKeySet{Keys: createJSONWebKeys(keyPair, kid, alg, use)}
}

func createJSONWebKeys(keyPair *MockSignerDecrypter, kid string, alg string, use string) []jose.JSONWebKey {
	return []jose.JSONWebKey{{
		Algorithm:                   alg,
		Use:                         use,
		Key:                         cryptosigner.Opaque(keyPair),
		KeyID:                       kid,
		Certificates:                []*x509.Certificate{},
		CertificateThumbprintSHA1:   []uint8{},
		CertificateThumbprintSHA256: []uint8{},
	}}
}

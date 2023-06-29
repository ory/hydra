// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk_test

import (
	"context"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
)

func TestKeyManagerStrategy(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	softwareKeyManager := NewMockManager(ctrl)
	hardwareKeyManager := NewMockManager(ctrl)
	keyManager := jwk.NewManagerStrategy(hardwareKeyManager, softwareKeyManager)
	defer ctrl.Finish()
	hwKeySet := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			KeyID: "hwKeyID",
		}},
	}
	swKeySet := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			KeyID: "swKeyID",
		}},
	}

	t.Run("GenerateAndPersistKeySet_WithResult", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1"), gomock.Any(), gomock.Any()).Return(hwKeySet, nil)
		resultKeySet, err := keyManager.GenerateAndPersistKeySet(context.TODO(), "set1", "kid1", "RS256", "sig")
		assert.NoError(t, err)
		assert.Equal(t, hwKeySet, resultKeySet)
	})

	t.Run("GenerateAndPersistKeySet_WithError", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1"), gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
		resultKeySet, err := keyManager.GenerateAndPersistKeySet(context.TODO(), "set1", "kid1", "RS256", "sig")
		assert.Error(t, err, "test")
		assert.Nil(t, resultKeySet)
	})

	t.Run("AddKey", func(t *testing.T) {
		softwareKeyManager.EXPECT().AddKey(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(nil)
		err := keyManager.AddKey(context.TODO(), "set1", nil)
		assert.NoError(t, err)
	})

	t.Run("AddKey_WithError", func(t *testing.T) {
		softwareKeyManager.EXPECT().AddKey(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(errors.New("test"))
		err := keyManager.AddKey(context.TODO(), "set1", nil)
		assert.Error(t, err, "test")
	})

	t.Run("AddKeySet", func(t *testing.T) {
		softwareKeyManager.EXPECT().AddKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(nil)
		err := keyManager.AddKeySet(context.TODO(), "set1", nil)
		assert.NoError(t, err)
	})

	t.Run("AddKeySet_WithError", func(t *testing.T) {
		softwareKeyManager.EXPECT().AddKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(errors.New("test"))
		err := keyManager.AddKeySet(context.TODO(), "set1", nil)
		assert.Error(t, err, "test")
	})

	t.Run("UpdateKey", func(t *testing.T) {
		softwareKeyManager.EXPECT().UpdateKey(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(nil)
		err := keyManager.UpdateKey(context.TODO(), "set1", nil)
		assert.NoError(t, err)
	})

	t.Run("UpdateKey_WithError", func(t *testing.T) {
		softwareKeyManager.EXPECT().UpdateKey(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(errors.New("test"))
		err := keyManager.UpdateKey(context.TODO(), "set1", nil)
		assert.Error(t, err, "test")
	})

	t.Run("UpdateKeySet", func(t *testing.T) {
		softwareKeyManager.EXPECT().UpdateKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(nil)
		err := keyManager.UpdateKeySet(context.TODO(), "set1", nil)
		assert.NoError(t, err)
	})

	t.Run("UpdateKeySet_WithError", func(t *testing.T) {
		softwareKeyManager.EXPECT().UpdateKeySet(gomock.Any(), gomock.Eq("set1"), gomock.Any()).Return(errors.New("test"))
		err := keyManager.UpdateKeySet(context.TODO(), "set1", nil)
		assert.Error(t, err, "test")
	})

	t.Run("GetKey_WithResultFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(hwKeySet, nil)
		resultKeySet, err := keyManager.GetKey(context.TODO(), "set1", "kid1")
		assert.NoError(t, err)
		assert.Equal(t, hwKeySet, resultKeySet)
	})

	t.Run("GetKey_WithErrorFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil, errors.New("test"))
		resultKeySet, err := keyManager.GetKey(context.TODO(), "set1", "kid1")
		assert.Error(t, err, "test")
		assert.Nil(t, resultKeySet)
	})

	t.Run("GetKey_WithErrNotFoundFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil, errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(swKeySet, nil)
		resultKeySet, err := keyManager.GetKey(context.TODO(), "set1", "kid1")
		assert.NoError(t, err)
		assert.Equal(t, swKeySet, resultKeySet)
	})

	t.Run("GetKey_WithErrNotFoundFromSoftwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil, errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().GetKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil, errors.WithStack(x.ErrNotFound))
		resultKeySet, err := keyManager.GetKey(context.TODO(), "set1", "kid1")
		assert.Error(t, err, "Not Found")
		assert.Nil(t, resultKeySet)
	})

	t.Run("GetKeySet_WithResultFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(hwKeySet, nil)
		resultKeySet, err := keyManager.GetKeySet(context.TODO(), "set1")
		assert.NoError(t, err)
		assert.Equal(t, hwKeySet, resultKeySet)
	})

	t.Run("GetKeySet_WithErrorFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil, errors.New("test"))
		resultKeySet, err := keyManager.GetKeySet(context.TODO(), "set1")
		assert.Error(t, err, "test")
		assert.Nil(t, resultKeySet)
	})

	t.Run("GetKeySet_WithErrNotFoundFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil, errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(swKeySet, nil)
		resultKeySet, err := keyManager.GetKeySet(context.TODO(), "set1")
		assert.NoError(t, err)
		assert.Equal(t, swKeySet, resultKeySet)
	})

	t.Run("GetKeySet_WithErrNotFoundFromSoftwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil, errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil, errors.WithStack(x.ErrNotFound))
		resultKeySet, err := keyManager.GetKeySet(context.TODO(), "set1")
		assert.Error(t, err, "Not Found")
		assert.Nil(t, resultKeySet)
	})

	t.Run("DeleteKey_FromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil)
		err := keyManager.DeleteKey(context.TODO(), "set1", "kid1")
		assert.NoError(t, err)
	})

	t.Run("DeleteKey_WithErrorFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(errors.New("test"))
		err := keyManager.DeleteKey(context.TODO(), "set1", "kid1")
		assert.Error(t, err, "test")
	})

	t.Run("DeleteKey_WithErrNotFoundFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(nil)
		err := keyManager.DeleteKey(context.TODO(), "set1", "kid1")
		assert.NoError(t, err)
	})

	t.Run("DeleteKey_WithErrNotFoundFromSoftwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().DeleteKey(gomock.Any(), gomock.Eq("set1"), gomock.Eq("kid1")).Return(errors.WithStack(x.ErrNotFound))
		err := keyManager.DeleteKey(context.TODO(), "set1", "kid1")
		assert.Error(t, err, "Not Found")
	})

	t.Run("DeleteKeySet_FromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil)
		err := keyManager.DeleteKeySet(context.TODO(), "set1")
		assert.NoError(t, err)
	})

	t.Run("DeleteKeySet_WithErrorFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(errors.New("test"))
		err := keyManager.DeleteKeySet(context.TODO(), "set1")
		assert.Error(t, err, "test")
	})

	t.Run("DeleteKeySet_WithErrNotFoundFromHardwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(nil)
		err := keyManager.DeleteKeySet(context.TODO(), "set1")
		assert.NoError(t, err)
	})

	t.Run("DeleteKeySet_WithErrNotFoundFromSoftwareKeyManager", func(t *testing.T) {
		hardwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(errors.WithStack(x.ErrNotFound))
		softwareKeyManager.EXPECT().DeleteKeySet(gomock.Any(), gomock.Eq("set1")).Return(errors.WithStack(x.ErrNotFound))
		err := keyManager.DeleteKeySet(context.TODO(), "set1")
		assert.Error(t, err, "Not Found")
	})
}

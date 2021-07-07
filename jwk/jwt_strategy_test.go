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
	"fmt"
	"testing"

	"github.com/pborman/uuid"

	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/ory/fosite/token/jwt"

	"github.com/ory/fosite/token/jwt"
	. "github.com/ory/hydra/jwk"
)

func TestRS256JWTStrategy_withSoftwareKeyStore(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)

	if conf.HsmEnabled() {
		t.Skip("Hardware Security Module enabled. Skipping test.")
	}

	testGenerator := &RS256Generator{}

	m := reg.KeyManager()
	ks, err := testGenerator.Generate("foo", "sig")
	require.NoError(t, err)
	require.NoError(t, m.AddKeySet(context.TODO(), "foo-set", ks))

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

	kid, err := s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, "public:foo", kid)

	ks, err = testGenerator.Generate("bar", "sig")
	require.NoError(t, err)
	require.NoError(t, m.AddKeySet(context.TODO(), "foo-set", ks))

	a, b, err = s.Generate(context.TODO(), jwt2.MapClaims{"foo": "bar"}, &jwt.Headers{})
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.NotEmpty(t, b)

	_, err = s.Validate(context.TODO(), a)
	require.NoError(t, err)

	kid, err = s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, "public:bar", kid)
}

func TestRS256JWTStrategy_withHardwareKeyStore(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)

	if !conf.HsmEnabled() {
		t.Skip("Hardware Security Module not enabled. Skipping test.")
	}

	m := reg.KeyManager()

	var kid1 = uuid.New()
	_, err := m.GenerateKeySet(context.TODO(), "foo-set", kid1, "RS256", "sig")
	require.NoError(t, err)

	s, err := NewRS256JWTStrategy(*conf, reg, func() string {
		return "foo-set"
	})

	require.NoError(t, err)
	token, sig, err := s.Generate(context.TODO(), jwt2.MapClaims{"foo": "bar"}, &jwt.Headers{})
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, sig)

	_, err = s.Validate(context.TODO(), token)
	require.NoError(t, err)

	kid, err := s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("public:%s", kid1), kid)

	err = m.DeleteKeySet(context.TODO(), "foo-set")
	require.NoError(t, err)

	var kid2 = uuid.New()
	_, err = m.GenerateKeySet(context.TODO(), "foo-set", kid2, "RS256", "sig")

	token, sig, err = s.Generate(context.TODO(), jwt2.MapClaims{"foo": "bar"}, &jwt.Headers{})
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, sig)

	_, err = s.Validate(context.TODO(), token)
	require.NoError(t, err)

	kid, err = s.GetPublicKeyID(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("public:%s", kid2), kid)
}

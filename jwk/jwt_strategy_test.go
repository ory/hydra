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
	"testing"

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

	_, err := m.GenerateKeySet(context.TODO(), "foo-set", "foo", "RS256", "sig")
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

	_, err = m.GenerateKeySet(context.TODO(), "foo-set", "bar", "RS256", "sig")
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

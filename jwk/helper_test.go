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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"
)

func TestFindKeyByPrefix(t *testing.T) {
	jwks := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: "public:foo"},
		{KeyID: "private:foo"},
	}}

	key, err := FindKeyByPrefix(jwks, "public")
	require.NoError(t, err)
	assert.Equal(t, "public:foo", key.KeyID)

	key, err = FindKeyByPrefix(jwks, "private")
	require.NoError(t, err)
	assert.Equal(t, "private:foo", key.KeyID)

	_, err = FindKeyByPrefix(jwks, "asdf")
	require.Error(t, err)

	jwks = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: "public:"},
		{KeyID: "private:"},
	}}

	key, err = FindKeyByPrefix(jwks, "public")
	require.NoError(t, err)
	assert.Equal(t, "public:", key.KeyID)

	key, err = FindKeyByPrefix(jwks, "private")
	require.NoError(t, err)
	assert.Equal(t, "private:", key.KeyID)

	_, err = FindKeyByPrefix(jwks, "asdf")
	require.Error(t, err)

	jwks = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: ""},
	}}
	require.Error(t, err)
}

func TestIder(t *testing.T) {
	assert.True(t, len(ider("public", "")) > len("public:"))
	assert.Equal(t, "public:foo", ider("public", "foo"))
}

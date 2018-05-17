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

package oauth2

import (
	"testing"

	"github.com/oleiade/reflections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.Equal(t, c, d, "%s", k)
	}
}

func AssertObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.NotEqual(t, c, d, "%s", k)
	}
}

func RequireObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.Equal(t, c, d, "%s", k)
	}
}
func RequireObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.NotEqual(t, c, d, "%s", k)
	}
}

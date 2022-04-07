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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"
)

func TestGenerator(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	for k, c := range []struct {
		g     KeyGenerator
		use   string
		check func(*jose.JSONWebKeySet)
	}{
		{
			g:   &RS256Generator{},
			use: "sig",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
				assert.Equal(t, "sig", ks.Keys[0].Use)
				assert.Equal(t, "sig", ks.Keys[1].Use)
			},
		},
		{
			g:   &ECDSA512Generator{},
			use: "enc",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
				assert.Equal(t, "enc", ks.Keys[0].Use)
				assert.Equal(t, "enc", ks.Keys[1].Use)
			},
		},
		{
			g:   &ECDSA256Generator{},
			use: "sig",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
				assert.Equal(t, "sig", ks.Keys[0].Use)
				assert.Equal(t, "sig", ks.Keys[1].Use)
			},
		},
		{
			g:   &HS256Generator{},
			use: "sig",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 1)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.Equal(t, "sig", ks.Keys[0].Use)
			},
		},
		{
			g:   &HS512Generator{},
			use: "enc",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 1)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.Equal(t, "enc", ks.Keys[0].Use)
			},
		},
		{
			g:   &EdDSAGenerator{},
			use: "sig",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
				assert.Equal(t, "sig", ks.Keys[0].Use)
				assert.Equal(t, "sig", ks.Keys[1].Use)
			},
		},
		{
			g:   &EdDSAGenerator{},
			use: "enc",
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks.Keys, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
				assert.Equal(t, "enc", ks.Keys[0].Use)
				assert.Equal(t, "enc", ks.Keys[1].Use)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			keys, err := c.g.Generate("foo", c.use)
			require.NoError(t, err)
			if err == nil {
				c.check(keys)
			}
		})
	}
}

// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwk

import (
	"testing"

	"fmt"

	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	for k, c := range []struct {
		g     KeyGenerator
		check func(*jose.JSONWebKeySet)
	}{
		{
			g: &RS256Generator{},
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
			},
		},
		{
			g: &ECDSA512Generator{},
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
			},
		},
		{
			g: &ECDSA256Generator{},
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks, 2)
				assert.NotEmpty(t, ks.Keys[0].Key)
				assert.NotEmpty(t, ks.Keys[1].Key)
			},
		},
		{
			g: &HS256Generator{},
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks, 1)
				assert.NotEmpty(t, ks.Keys[0].Key)
			},
		},
		{
			g: &HS512Generator{},
			check: func(ks *jose.JSONWebKeySet) {
				assert.Len(t, ks, 1)
				assert.NotEmpty(t, ks.Keys[0].Key)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			keys, err := c.g.Generate("foo")
			require.NoError(t, err)
			if err != nil {
				c.check(keys)
			}
		})
	}
}

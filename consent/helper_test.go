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

package consent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
)

func TestSanitizeClient(t *testing.T) {
	c := &client.Client{
		Secret: "some-secret",
	}
	ar := &fosite.AuthorizeRequest{
		Request: fosite.Request{
			Client: c,
		},
	}
	got := sanitizeClientFromRequest(ar)
	assert.Empty(t, got.Secret)
	assert.NotEmpty(t, c.Secret)
}

func TestMatchScopes(t *testing.T) {
	for k, tc := range []struct {
		granted         []HandledConsentRequest
		requested       []string
		expectChallenge string
	}{
		{
			granted:         []HandledConsentRequest{{Challenge: "1", GrantedScope: []string{"foo", "bar"}}},
			requested:       []string{"foo", "bar"},
			expectChallenge: "1",
		},
		{
			granted:         []HandledConsentRequest{{Challenge: "1", GrantedScope: []string{"foo", "bar"}}},
			requested:       []string{"foo", "bar", "baz"},
			expectChallenge: "",
		},
		{
			granted: []HandledConsentRequest{
				{Challenge: "1", GrantedScope: []string{"foo", "bar"}},
				{Challenge: "2", GrantedScope: []string{"foo", "bar"}},
			},
			requested:       []string{"foo", "bar"},
			expectChallenge: "1",
		},
		{
			granted: []HandledConsentRequest{
				{Challenge: "1", GrantedScope: []string{"foo", "bar"}},
				{Challenge: "2", GrantedScope: []string{"foo", "bar", "baz"}},
			},
			requested:       []string{"foo", "bar", "baz"},
			expectChallenge: "2",
		},
		{
			granted: []HandledConsentRequest{
				{Challenge: "1", GrantedScope: []string{"foo", "bar"}},
				{Challenge: "2", GrantedScope: []string{"foo", "bar", "baz"}},
			},
			requested:       []string{"zab"},
			expectChallenge: "",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			got := matchScopes(fosite.ExactScopeStrategy, tc.granted, tc.requested)
			if tc.expectChallenge == "" {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tc.expectChallenge, got.Challenge)
		})
	}
}

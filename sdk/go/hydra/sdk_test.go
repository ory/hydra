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

package hydra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	var sdk SDK
	var err error
	sdk, err = NewSDK(&Configuration{
		EndpointURL:  "http://localhost:4444/",
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	assert.NoError(t, err)
	assert.NotNil(t, sdk)
}

func TestErrorHandlers(t *testing.T) {
	for k, c := range []Configuration{
		{
			EndpointURL:  "http://localhost:4444/",
			ClientSecret: "bar",
			Scopes:       []string{"foo"},
		},
		{
			EndpointURL: "http://localhost:4444/",
			ClientID:    "bar",
			Scopes:      []string{"foo"},
		},
		{
			ClientID:     "foo",
			ClientSecret: "bar",
			Scopes:       []string{"foo"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			sdk, err := NewSDK(&c)
			assert.Error(t, err)
			assert.Nil(t, sdk)
		})
	}
}

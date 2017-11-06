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

package warden_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionAllowed(t *testing.T) {
	for n, w := range wardens {
		t.Run("warden="+n, func(t *testing.T) {
			for k, c := range accessRequestTokenTestCases {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					ctx, err := w.TokenAllowed(context.Background(), c.token, c.req, c.scopes...)
					if c.expectErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					if err == nil && c.assert != nil {
						c.assert(t, ctx)
					}
				})
			}
		})
	}
}

func TestAllowed(t *testing.T) {
	for n, w := range wardens {
		t.Run("warden="+n, func(t *testing.T) {
			for k, c := range accessRequestTestCases {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					err := w.IsAllowed(context.Background(), c.req)
					if c.expectErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}
				})
			}
		})
	}
}

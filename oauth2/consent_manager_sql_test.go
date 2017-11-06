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

package oauth2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsentRequestSqlDataTransforms(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		d string
		r *ConsentRequest
	}{
		{
			d: "fully hydrated request object",
			r: &ConsentRequest{
				ID:               "id",
				ClientID:         "client-id",
				RequestedScopes:  []string{"foo", "bar"},
				GrantedScopes:    []string{"baz", "bar"},
				CSRF:             "some-csrf",
				ExpiresAt:        time.Now().Round(time.Second),
				Consent:          ConsentRequestAccepted,
				DenyReason:       "some reason",
				AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
				IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
				RedirectURL:      "https://redirect-me/foo",
				Subject:          "Peter",
			},
		},
	} {
		t.Run(tc.d, func(t *testing.T) {
			s, err := newConsentRequestSqlData(tc.r)
			require.Nil(t, err)

			o, err := s.toConsentRequest()
			require.NoError(t, err)

			assert.EqualValues(t, tc.r, o)
		})
	}
}

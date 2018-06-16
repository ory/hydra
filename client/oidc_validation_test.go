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

package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSectorIdentifierURL(t *testing.T) {
	var payload string

	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	}
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	v := &DynamicValidator{
		c: ts.Client(),
	}

	for k, tc := range []struct {
		p         string
		r         []string
		u         string
		expectErr bool
	}{
		{
			u:         "",
			expectErr: true,
		},
		{
			u:         "http://foo/bar",
			expectErr: true,
		},
		{
			u:         ts.URL,
			expectErr: true,
		},
		{
			p:         `["http://foo"]`,
			u:         ts.URL,
			expectErr: false,
			r:         []string{"http://foo"},
		},
		{
			p:         `["http://foo"]`,
			u:         ts.URL,
			expectErr: true,
			r:         []string{"http://foo", "http://not-foo"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			payload = tc.p
			err := v.validateSectorIdentifierURL(tc.u, tc.r)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

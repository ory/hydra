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

package x

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/x/sqlcon"
)

func TestErrorEnhancer(t *testing.T) {
	for k, tc := range []struct {
		in  error
		out string
	}{
		{
			in:  sqlcon.ErrNoRows,
			out: "{\"error\":\"Unable to locate the resource\",\"error_description\":\"\"}",
		},
		{
			in:  errorsx.WithStack(sqlcon.ErrNoRows),
			out: "{\"error\":\"Unable to locate the resource\",\"error_description\":\"\"}",
		},
		{
			in:  errors.New("bla"),
			out: "{\"error\":\"error\",\"error_description\":\"The error is unrecognizable\"}",
		},
		{
			in:  errors2.New("bla"),
			out: "{\"error\":\"error\",\"error_description\":\"The error is unrecognizable\"}",
		},
		{
			in:  fosite.ErrInvalidRequest,
			out: "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Make sure that the various parameters are correct, be aware of case sensitivity and trim your parameters. Make sure that the client you are using has exactly whitelisted the redirect_uri you specified.\"}",
		},
		{
			in:  errorsx.WithStack(fosite.ErrInvalidRequest),
			out: "{\"error\":\"invalid_request\",\"error_description\":\"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Make sure that the various parameters are correct, be aware of case sensitivity and trim your parameters. Make sure that the client you are using has exactly whitelisted the redirect_uri you specified.\"}",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			err := ErrorEnhancer(new(http.Request), tc.in)
			out, err2 := json.Marshal(err)
			require.NoError(t, err2)
			assert.Equal(t, tc.out, string(out))
		})
	}
}

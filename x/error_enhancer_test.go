// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"net/http"
	"testing"

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
			in:  errors.WithStack(sqlcon.ErrNoRows),
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
			in:  errors.WithStack(fosite.ErrInvalidRequest),
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

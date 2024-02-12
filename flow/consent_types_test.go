// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
)

func TestToRFCError(t *testing.T) {
	for k, tc := range []struct {
		input  *RequestDeniedError
		expect *fosite.RFC6749Error
	}{
		{
			input: &RequestDeniedError{
				Name:  "not empty",
				Valid: true,
			},
			expect: &fosite.RFC6749Error{
				ErrorField:       "not empty",
				DescriptionField: "",
				CodeField:        fosite.ErrInvalidRequest.CodeField,
				DebugField:       "",
			},
		},
		{
			input: &RequestDeniedError{
				Name:        "",
				Description: "not empty",
				Valid:       true,
			},
			expect: &fosite.RFC6749Error{
				ErrorField:       "request_denied",
				DescriptionField: "not empty",
				CodeField:        fosite.ErrInvalidRequest.CodeField,
				DebugField:       "",
			},
		},
		{
			input: &RequestDeniedError{Valid: true},
			expect: &fosite.RFC6749Error{
				ErrorField:       "request_denied",
				DescriptionField: "",
				HintField:        "",
				CodeField:        fosite.ErrInvalidRequest.CodeField,
				DebugField:       "",
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			require.EqualValues(t, tc.input.ToRFCError(), tc.expect)
		})
	}
}

func TestRequestDeniedError(t *testing.T) {
	var e *RequestDeniedError
	v, err := e.Value()
	require.NoError(t, err)
	assert.EqualValues(t, "{}", fmt.Sprintf("%v", v))
}

func TestAcceptOAuth2ConsentRequest_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(AcceptOAuth2ConsentRequest))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestOAuth2ConsentSession_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(OAuth2ConsentSession))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestHandledLoginRequest_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(HandledLoginRequest))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestOAuth2ConsentRequestOpenIDConnectContext_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(OAuth2ConsentRequestOpenIDConnectContext))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestLogoutRequest_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(LogoutRequest))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestLoginRequest_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(LoginRequest))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestOAuth2ConsentRequest_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(OAuth2ConsentRequest))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

func TestAcceptOAuth2ConsentRequestSession_MarshalJSON(t *testing.T) {
	out, err := json.Marshal(new(AcceptOAuth2ConsentRequestSession))
	require.NoError(t, err)
	snapshotx.SnapshotT(t, string(out))
}

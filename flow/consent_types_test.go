// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
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
	assert.JSONEq(t, "{}", fmt.Sprintf("%v", v))
}

func TestAcceptOAuth2ConsentRequest_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(AcceptOAuth2ConsentRequest))
}

func TestOAuth2ConsentSession_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(OAuth2ConsentSession))
}

func TestHandledLoginRequest_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(HandledLoginRequest))
}

func TestOAuth2ConsentRequestOpenIDConnectContext_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(OAuth2ConsentRequestOpenIDConnectContext))
}

func TestLogoutRequest_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(LogoutRequest))
}

func TestLoginRequest_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(LoginRequest))
}

func TestOAuth2ConsentRequest_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(OAuth2ConsentRequest))
}

func TestAcceptOAuth2ConsentRequestSession_MarshalJSON(t *testing.T) {
	snapshotx.SnapshotT(t, new(AcceptOAuth2ConsentRequestSession))
}

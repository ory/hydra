// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPurposeAADWireValuesAreStable locks the additional-authenticated-data
// bytes produced for each codec purpose. The purpose is serialized into the
// AAD ({"p":N}, omitted when zero) and bound into every flow ciphertext —
// including AEAD authorization codes. The constants are iota-derived;
// reordering or inserting one would shift these values and silently
// invalidate every in-flight code across a deploy. This test guards the
// wire values so such a change fails loudly instead.
func TestPurposeAADWireValuesAreStable(t *testing.T) {
	cases := []struct {
		name string
		opt  CodecOption
		want string
	}{
		{"login challenge", AsLoginChallenge, `{}`}, // 0 is omitted by omitempty.
		{"login verifier", AsLoginVerifier, `{"p":1}`},
		{"device challenge", AsDeviceChallenge, `{"p":2}`},
		{"device verifier", AsDeviceVerifier, `{"p":3}`},
		{"consent challenge", AsConsentChallenge, `{"p":4}`},
		{"consent verifier", AsConsentVerifier, `{"p":5}`},
		{"authorize code", AsAuthorizeCode, `{"p":6}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, string(additionalDataFromOpts(tc.opt)))
		})
	}
}

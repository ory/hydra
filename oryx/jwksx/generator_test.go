// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwksx

import (
	"fmt"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/require"
)

func TestGenerateSigningKeys(t *testing.T) {
	for _, alg := range GenerateSigningKeysAvailableAlgorithms() {
		t.Run(fmt.Sprintf("alg=%s", alg), func(t *testing.T) {
			key, err := GenerateSigningKeys("", alg, 0)
			require.NoError(t, err)
			t.Logf("%+v", key)
		})
	}

	for _, tc := range []struct {
		alg  jose.SignatureAlgorithm
		bits int
	}{
		{alg: jose.HS256, bits: 128}, // should fail because minimum 256 bit
		{alg: jose.HS384, bits: 256}, // should fail because minimum 384 bit
		{alg: jose.HS512, bits: 384}, // should fail because minimum 512 bit
		{alg: jose.HS512, bits: 555}, // should fail because not modulo 8
	} {
		t.Run(fmt.Sprintf("alg=%s/bit=%d", tc.alg, tc.bits), func(t *testing.T) {
			_, err := GenerateSigningKeys("", string(tc.alg), tc.bits)
			require.Error(t, err)
		})
	}
}

// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ipx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAssociatedIPAllowed(t *testing.T) {
	for _, disallowed := range []string{
		"localhost",
		"https://localhost/foo?bar=baz#zab",
		"127.0.0.0",
		"127.255.255.255",
		"172.16.0.0",
		"172.31.255.255",
		"192.168.0.0",
		"192.168.255.255",
		"10.0.0.0",
		"10.255.255.255",
		"::1",
	} {
		t.Run("case="+disallowed, func(t *testing.T) {
			require.Error(t, IsAssociatedIPAllowed(disallowed))
		})
	}

	// Do not error if invalid data is used
	require.NoError(t, IsAssociatedIPAllowed("idonotexist"))
	require.NoError(t, IsAssociatedIPAllowedWhenSet(""))
	require.NoError(t, AreAllAssociatedIPsAllowed(map[string]string{
		"foo": "https://google.com",
		"bar": "microsoft.com",
	}))
	require.Error(t, AreAllAssociatedIPsAllowed(map[string]string{
		"foo": "https://google.com",
		"bar": "microsoft.com",
		"baz": "localhost",
	}))
}

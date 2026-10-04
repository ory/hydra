// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/x"
)

// Test clients must not share http.DefaultTransport: httptest.Server.Close
// calls http.DefaultTransport.CloseIdleConnections, so a parallel subtest
// closing its server would break another subtest's in-flight request with
// "http: CloseIdleConnections called".
func TestTestClientsOwnTheirTransport(t *testing.T) {
	for name, hc := range map[string]*http.Client{
		"NewEmptyJarClient": x.NewEmptyJarClient(t),
		"NewTestClient":     x.NewTestClient(t),
	} {
		t.Run(name, func(t *testing.T) {
			require.NotNil(t, hc.Transport)
			require.NotSame(t, http.DefaultTransport, hc.Transport)
		})
	}
}

// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockrc struct{ dm bool }

func (m *mockrc) IsDevelopmentMode(ctx context.Context) bool {
	return m.dm
}

func TestIsRedirectURISecure(t *testing.T) {
	for d, c := range []struct {
		u   string
		err bool
		dev bool
	}{
		{u: "http://google.com", err: true},
		{u: "https://google.com", err: false},
		{u: "http://localhost", err: false},
		{u: "http://test.localhost", err: false},
		{u: "wta://auth", err: false},
		{u: "http://foo.com/bar", err: false, dev: true},
		{u: "http://baz.com/bar", err: false, dev: true},
	} {
		uu, err := url.Parse(c.u)
		require.NoError(t, err)
		assert.Equal(t, !c.err, IsRedirectURISecure(&mockrc{dm: c.dev})(context.Background(), uu), "case %d", d)
	}
}

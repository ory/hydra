// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRegistrationClientURI(t *testing.T) {
	base, err := url.Parse("http://localhost:4444/")
	require.NoError(t, err)

	for id, expected := range map[string]string{
		"4b2c1e7a-9d3f-4c1a-8e55-0a1b2c3d4e5f": "http://localhost:4444/oauth2/register/4b2c1e7a-9d3f-4c1a-8e55-0a1b2c3d4e5f",
		"team client":                          "http://localhost:4444/oauth2/register/team%20client",
		"a/b":                                  "http://localhost:4444/oauth2/register/a%2Fb",
		"100%":                                 "http://localhost:4444/oauth2/register/100%25",
	} {
		t.Run("id="+id, func(t *testing.T) {
			assert.Equal(t, expected, buildRegistrationClientURI(base, id))
		})
	}
}

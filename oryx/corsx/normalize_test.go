// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"
)

func TestNormalizeOrigins(t *testing.T) {
	assert.EqualValues(t,
		[]string{"https://example.org:1234"},
		NormalizeOrigins([]url.URL{*urlx.ParseOrPanic("https://example.org:1234/asdf")}))
}

func TestNormalizeOriginStrings(t *testing.T) {
	actual, err := NormalizeOriginStrings([]string{"https://example.org:1234/asdf"})
	require.NoError(t, err)
	assert.EqualValues(t, []string{"https://example.org:1234"}, actual)
}

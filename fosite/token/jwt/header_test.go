// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderToMap(t *testing.T) {
	header := &Headers{}
	header.Add("foo", "bar")
	assert.Equal(t, "bar", header.Get("foo"))
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
	}, header.ToMap())
}

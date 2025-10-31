// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/ory/hydra/v2/fosite"
)

func TestAccessResponse(t *testing.T) {
	ar := NewAccessResponse()
	ar.SetAccessToken("access")
	ar.SetTokenType("bearer")
	ar.SetExtra("access_token", "invalid")
	ar.SetExtra("foo", "bar")
	assert.Equal(t, "access", ar.GetAccessToken())
	assert.Equal(t, "bearer", ar.GetTokenType())
	assert.Equal(t, "bar", ar.GetExtra("foo"))
	assert.Equal(t, map[string]interface{}{
		"access_token": "access",
		"token_type":   "bearer",
		"foo":          "bar",
	}, ar.ToMap())
}

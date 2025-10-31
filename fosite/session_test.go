// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	var s *DefaultSession
	assert.Empty(t, s.GetSubject())
	assert.Empty(t, s.GetUsername())
	assert.Nil(t, s.Clone())
}

func TestZeroSession(t *testing.T) {
	var s *DefaultSession = new(DefaultSession)
	assert.Empty(t, s.GetSubject())
	assert.Empty(t, s.GetUsername())
	assert.Empty(t, s.Clone())
	assert.Empty(t, s.GetExpiresAt(AccessToken))
}

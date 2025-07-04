// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPointer(t *testing.T) {
	s := "TestString"
	assert.Equal(t, &s, GetPointer(s))
}

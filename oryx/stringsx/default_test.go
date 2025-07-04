// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultIfEmpty(t *testing.T) {
	assert.Equal(t, DefaultIfEmpty("", "default"), "default")
	assert.Equal(t, DefaultIfEmpty("custom", "default"), "custom")
}

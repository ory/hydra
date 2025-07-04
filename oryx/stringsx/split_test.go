// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitNonEmpty(t *testing.T) {
	// assert.Len(t, strings.Split("", " "), 1)
	assert.Len(t, Splitx("", " "), 0)
}

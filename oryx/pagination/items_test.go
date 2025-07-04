// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxItemsPerPage(t *testing.T) {
	assert.Equal(t, 0, MaxItemsPerPage(100, 0))
	assert.Equal(t, 10, MaxItemsPerPage(100, 10))
	assert.Equal(t, 100, MaxItemsPerPage(100, 110))
}

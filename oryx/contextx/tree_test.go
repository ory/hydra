// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeContext(t *testing.T) {
	assert.True(t, IsRootContext(RootContext))
	assert.True(t, IsRootContext(context.WithValue(RootContext, "foo", "bar"))) //lint:ignore SA1029 builtin type for context is OK in test
	assert.False(t, IsRootContext(context.Background()))
}

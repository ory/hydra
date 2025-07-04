// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package errorsx

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWithStack(t *testing.T) {
	t.Run("case=wrap", func(t *testing.T) {
		orig := errors.New("hi")
		wrap := WithStack(orig)

		assert.EqualValues(t, orig.(StackTracer).StackTrace(), wrap.(StackTracer).StackTrace())
		assert.EqualValues(t, orig.(StackTracer).StackTrace(), WithStack(wrap).(StackTracer).StackTrace())
	})
}

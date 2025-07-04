// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpMessage(t *testing.T) {
	assert.NotEmpty(t, HelpMessage())
}

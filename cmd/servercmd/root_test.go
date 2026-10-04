// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package servercmd

import (
	"testing"

	"github.com/ory/x/cmdx"
)

func TestUsageStrings(t *testing.T) {
	cmdx.AssertUsageTemplates(t, NewRootCmd(nil, nil))
}
